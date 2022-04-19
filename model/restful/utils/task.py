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
from wsgiref import headers
from flask import Blueprint, jsonify, request
import path
from aiohttp import ClientSession
import asyncio
import uuid
import subprocess
import sys

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

    response_payload_size = service_endpoint["response_payload_size"]
    response_payload = subprocess.run(['cat /dev/urandom | tr -dc "[:alnum:]" | head -c${1:-%s}' % response_payload_size], capture_output=True, shell=True)
    res_payload = response_payload.stdout.decode("utf-8")

    if service_endpoint["forward_requests"] == "asynchronous":
        asyncio.set_event_loop(asyncio.new_event_loop())
        loop = asyncio.get_event_loop()
        response = loop.run_until_complete(async_tasks(service_endpoint, headers, res_payload))
        return response
    else: # "synchronous"
        asyncio.set_event_loop(asyncio.new_event_loop())
        loop = asyncio.get_event_loop()
        response = loop.run_until_complete(sync_tasks(service_endpoint, headers, res_payload))

        return response


async def async_tasks(service_endpoint, headers, res_payload):
    async with ClientSession() as session:
        io_tasks = []

        if len(service_endpoint["called_services"]) > 0:
            for svc in service_endpoint["called_services"]:
                io_task = asyncio.create_task(execute_io_bounded_task(session=session, target_service=svc, sync=False, forward_headers=headers))
                io_tasks.append(io_task)
            services = await asyncio.gather(*io_tasks)

    # Concatenate json responses
    response = {}
    response["services"] = []
    response["statuses"] = []
    response["payload"] = res_payload

    if len(service_endpoint["called_services"]) > 0:
        for svc in services:
            response["services"] += svc["services"]
            response["statuses"] += svc["statuses"]
    return response


async def sync_tasks(service_endpoint, headers, res_payload):
    async with ClientSession() as session:
        response = {}
        response["services"] = []
        response["statuses"] = []
        response["payload"] = res_payload
        
        if len(service_endpoint["called_services"]) > 0:
            for svc in service_endpoint["called_services"]:
                res = await execute_io_bounded_task(session=session, target_service=svc, sync=True, forward_headers=headers)
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

    # TODO: CPU-bounded tasks not supported yet
    # TODO: Implement resource stress emulation...

    response_object = {
        "status": "CPU-bounded task executed",
        "data": {
            "svc_name": origin_service_name,
            "task_id": task_config["task_id"]
        }
    }
    return response_object


async def execute_io_bounded_task(session, target_service, sync, forward_headers={}):

    dst = FORMATTED_REMOTE_URL.format(target_service["service"], target_service["port"], target_service["endpoint"])
    forward_headers.update({'Content-type' : 'application/json'})

    json_data = {}
    json_data["payload"] = ""

    forward_ratio = target_service["traffic_forward_ratio"]
    request_payload_size = target_service["request_payload_size"]

    responses = {}
    responses["services"] = []
    responses["statuses"] = []

    if request_payload_size:
        request_payload = subprocess.run(['cat /dev/urandom | tr -dc "[:alnum:]" | head -c${1:-%s}' % request_payload_size], capture_output=True, shell=True)
        json_data["payload"] = request_payload.stdout.decode("utf-8")

        if forward_ratio > 0:
            if not sync:
                async with ClientSession() as session:
                    io_tasks = []

                    for i in range(forward_ratio):
                        io_task = asyncio.create_task(session.post(dst, data=json_data, headers=forward_headers))
                        io_tasks.append(io_task)
                    calls = await asyncio.gather(*io_tasks)

                # Concatenate json responses
                for res in calls:
                    res_payload = await res.json()

                    responses["services"] += res_payload["services"]
                    responses['services'].append(str(res.url))
                    responses["statuses"] += res_payload["statuses"]
                    responses['statuses'].append(res.status)

            else: # "synchronous"
                async with ClientSession() as session:
                    for i in range(forward_ratio):
                        res = await session.post(dst, data=json_data, headers=forward_headers)
                        res_payload = await res.json()

                        responses["services"] += res_payload["services"]
                        responses['services'].append(str(res.url))
                        responses["statuses"] += res_payload["statuses"]
                        responses['statuses'].append(res.status)
                
    return responses