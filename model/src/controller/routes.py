"""
Copyright 2021 Ericsson AB

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import logging
from flask import Blueprint, jsonify, request
from src.service.tasks import execute_io_bounded_task
from src.service import path
from aiohttp import ClientSession
import asyncio

simple_page = Blueprint("simple_page", __name__,)
logger = logging.getLogger(__name__)

# TODO: So far, we only support one endpoint per service...
service_endpoint = path.SERVICE_CONFIG["endpoints"][0]
service_processes = path.SERVICE_CONFIG["processes"]

def getForwardHeaders(request):
    '''
    function to propagate header from inbound to outbound
    '''

    #incoming_headers = ['user-agent', 'x-request-id', 'x-datadog-trace-id', 'x-datadog-parent-id', 'x-datadog-sampled']
    incoming_headers = ['user-agent', 'end-user', 'x-request-id', 'x-b3-traceid', 'x-b3-spanid', 'x-b3-parentspanid', 'x-b3-sampled', 'x-b3-flags']

    # propagate headers manually
    headers = {}
    for ihdr in incoming_headers:
        val = request.headers.get(ihdr)
        if val is not None:
            headers[ihdr] = val
    return headers

@simple_page.route(service_endpoint["name"], methods=["POST"])
async def run_task():
    headers = getForwardHeaders(request)

    if service_endpoint["forward_requests"] == "asynchronous":
        async with ClientSession() as session:
            # TODO: CPU-bounded tasks not supported yet
            io_tasks = []
            for svc in service_endpoint["called_services"]:
                io_task = asyncio.create_task(execute_io_bounded_task(session=session, target_service=svc, forward_headers=headers))
                io_tasks.append(io_task)
            services = await asyncio.gather(*io_tasks)

        # Concatenate json responses
        response = {}
        response["services"] = []
        response["statuses"] = []

        for svc in services:
            response["services"] += svc["services"]
            response["statuses"] += svc["statuses"]
        return response

    else: # "synchronous"
        async with ClientSession() as session:
            # TODO: CPU-bounded tasks not supported yet
            response = '<h1>Called services:</h1>'
            for svc in service_endpoint["called_services"]:
                res = execute_io_bounded_task(session=session, target_service=svc, forward_headers=headers)
                response += f"<p>Service: {res['service']} --- Status Code: {res['status']}</p>"

        return response

@simple_page.route("/", methods=["GET"])
def task_status():
    return "OK", 200