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
import path
from aiohttp import ClientSession
import asyncio
import uuid

# TODO: So far, we only support a hard-coded namespace. For more flexible support of namespaces we will need to pass that info as part of the config map
# TODO: So far, we only support http client
FORMATTED_REMOTE_URL = "http://{0}:{1}/{2}"


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


def run_task(service_endpoint):
    headers = getForwardHeaders(request)

    if service_endpoint["forward_requests"] == "asynchronous":
        asyncio.set_event_loop(asyncio.new_event_loop())
        loop = asyncio.get_event_loop()
        response = loop.run_until_complete(async_tasks(service_endpoint, headers))
        return response
    else: # "synchronous"
        asyncio.set_event_loop(asyncio.new_event_loop())
        loop = asyncio.get_event_loop()
        response = loop.run_until_complete(sync_tasks(service_endpoint, headers))

        return response


async def async_tasks(service_endpoint, headers):
    async with ClientSession() as session:
        # TODO: CPU-bounded tasks not supported yet
        io_tasks = []

        if request.get_data() == bytes("", "utf-8"):
            json_data = {}
        else:
            json_data = request.json
        if len(service_endpoint["called_services"]) > 0:
            for svc in service_endpoint["called_services"]:
                io_task = asyncio.create_task(execute_io_bounded_task(session=session, target_service=svc,
                                                                      json_data=json_data, forward_headers=headers))
                io_tasks.append(io_task)
            services = await asyncio.gather(*io_tasks)

        # Concatenate json responses
    response = {}
    response["services"] = []
    response["statuses"] = []

    if len(service_endpoint["called_services"]) > 0:
        for svc in services:
            response["services"] += svc["services"]
            response["statuses"] += svc["statuses"]
    return response


async def sync_tasks(service_endpoint, headers):
    async with ClientSession() as session:
        # TODO: CPU-bounded tasks not supported yet
        response = {}
        response["services"] = []
        response["statuses"] = []
        if request.get_data() == bytes("", "utf-8"):
            json_data = {}
        else:
            json_data = request.json
        if len(service_endpoint["called_services"]) > 0:
            for svc in service_endpoint["called_services"]:
                res = execute_io_bounded_task(session=session, target_service=svc, json_data=json_data, forward_headers=headers)
                response["services"] += res["services"]
                response["statuses"] += res["statuses"]

    return response


def execute_cpu_bounded_task(origin_service_name, target_service, headers):
    task_id = str(uuid.uuid4())
    task_config = {}
    task_config["task_id"] = task_id
    task_config["cpu_consumption"] = target_service["cpu_consumption"]
    task_config["network_consumption"] = target_service["network_consumption"]
    task_config["memory_consumption"] = target_service["memory_consumption"]

    # TODO: Implement resource stress emulation...

    response_object = {
        "status": "CPU-bounded task executed",
        "data": {
            "svc_name": origin_service_name,
            "task_id": task_config["task_id"]
        }
    }
    return response_object


async def execute_io_bounded_task(session, target_service, json_data, forward_headers={}):
    # TODO: Request forwarding to a service on a particular cluster is not supported yet
    # TODO: Requests for other protocols than html are not supported yet

    dst = FORMATTED_REMOTE_URL.format(target_service["service"], target_service["port"], target_service["endpoint"])
    forward_headers.update({'Content-type' : 'application/json'})

    # TODO: traffic_forward_ratio not supported yet
    traffic_forward_ratio = target_service["traffic_forward_ratio"]

    res = await session.post(dst, data=json_data, headers=forward_headers)
    res_payload = await res.json()

    response = {}
    response["services"] = res_payload["services"]
    response['services'].append(str(res.url))
    response["statuses"] = res_payload["statuses"]
    response['statuses'].append(res.status)
    return response
