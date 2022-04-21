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


def run_task(service_name, service_endpoint):
    headers = getForwardHeaders(request)

    response_payload_size = service_endpoint["network_complexity"]["response_payload_size"]
    response_payload = subprocess.run(['cat /dev/urandom | tr -dc "[:alnum:]" | head -c${1:-%s}' % response_payload_size], capture_output=True, shell=True)
    res_payload = response_payload.stdout.decode("utf-8")

    source_svc = {}
    source_svc["service"] = service_name
    source_svc["endpoint"] = service_endpoint["name"]

    # CPU task
    if service_endpoint["cpu_complexity"]:
        cpu_response = execute_cpu_bounded_task(conf=service_endpoint["cpu_complexity"])

    # TODO: Memory task
    mem_response = {}

    # Network task
    if service_endpoint["network_complexity"]["forward_requests"] == "asynchronous":
        asyncio.set_event_loop(asyncio.new_event_loop())
        loop = asyncio.get_event_loop()
        nw_response = loop.run_until_complete(async_network_task(source_svc, service_endpoint, headers, res_payload))
    else: # "synchronous"
        asyncio.set_event_loop(asyncio.new_event_loop())
        loop = asyncio.get_event_loop()
        nw_response = loop.run_until_complete(sync_network_task(source_svc, service_endpoint, headers, res_payload))

    nw_response["cpu_task"]["statuses"].append(cpu_response["status"])
  
    return nw_response


async def async_network_task(source_svc, service_endpoint, headers, res_payload):
    async with ClientSession() as session:
        io_tasks = []

        if len(service_endpoint["network_complexity"]["called_services"]) > 0:
            for target_svc in service_endpoint["network_complexity"]["called_services"]:
                io_task = asyncio.create_task(execute_io_bounded_task(session=session, source_service=source_svc, target_service=target_svc, sync=False, forward_headers=headers))
                io_tasks.append(io_task)
            services = await asyncio.gather(*io_tasks)

    # Concatenate json responses
    response = {}
    response["cpu_task"] = {}
    response["cpu_task"]["services"] = []
    response["cpu_task"]["statuses"] = []

    response["memory_task"] = {}

    response["network_task"] = {}
    response["network_task"]["services"] = []
    response["network_task"]["statuses"] = []
    response["network_task"]["payload"] = res_payload

    if len(service_endpoint["network_complexity"]["called_services"]) > 0:
        for svc in services:
            response["cpu_task"]["services"] += svc["cpu_task"]["services"]
            response["cpu_task"]["statuses"] += svc["cpu_task"]["statuses"]

            response["network_task"]["services"] += svc["network_task"]["services"]
            response["network_task"]["statuses"] += svc["network_task"]["statuses"]

    return response


async def sync_network_task(source_svc, service_endpoint, headers, res_payload):
    async with ClientSession() as session:
        response = {}
        response["cpu_task"] = {}
        response["cpu_task"]["services"] = []
        response["cpu_task"]["statuses"] = []

        response["memory_task"] = {}

        response["network_task"] = {}
        response["network_task"]["services"] = []
        response["network_task"]["statuses"] = []
        response["network_task"]["payload"] = res_payload
        
        if len(service_endpoint["network_complexity"]["called_services"]) > 0:
            for target_svc in service_endpoint["network_complexity"]["called_services"]:
                res = await execute_io_bounded_task(session=session, source_service=source_svc, target_service=target_svc, sync=True, forward_headers=headers)
                
                # Concatenate json responses
                response["cpu_task"]["services"] += res["cpu_task"]["services"]
                response["cpu_task"]["statuses"] += res["cpu_task"]["statuses"]

                response["network_task"]["services"] += res["network_task"]["services"]
                response["network_task"]["statuses"] += res["network_task"]["statuses"]

    return response


def execute_cpu_bounded_task(conf):
    # TODO: Implement resource stress emulation...
    
    res = subprocess.run(['stress-ng --class cpu --cpu %s --cpu-method %s --taskset %s --cpu-load %s --timeout %f --metrics-brief' % (conf["workers"], conf["method"], ",".join(str(cpu_id) for cpu_id in conf["cpu_affinity"]), conf["cpu_load"], conf["execution_time"])], capture_output=True, shell=True)

    response = {
        "status": res.stderr.decode("utf-8") 
    }
    return response


async def execute_io_bounded_task(session, source_service, target_service, sync, forward_headers={}):

    dst = FORMATTED_REMOTE_URL.format(target_service["service"], target_service["port"], target_service["endpoint"])
    forward_headers.update({'Content-type' : 'application/json'})

    json_data = {}
    json_data["payload"] = ""

    forward_ratio = target_service["traffic_forward_ratio"]
    request_payload_size = target_service["request_payload_size"]

    response = {}
    response["cpu_task"] = {}
    response["cpu_task"]["services"] = []
    response["cpu_task"]["statuses"] = []

    response["memory_task"] = {}

    response["network_task"] = {}
    response["network_task"]["services"] = []
    response["network_task"]["statuses"] = []


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

                    response["cpu_task"]["services"] += res_payload["cpu_task"]["services"]
                    response["cpu_task"]['services'].append("("+source_service["service"]+"/"+source_service["endpoint"]+", "+target_service["service"]+"/"+target_service["endpoint"]+")")
                    response["cpu_task"]["statuses"] += res_payload["cpu_task"]["statuses"]
                    
                    response["network_task"]["services"] += res_payload["network_task"]["services"]
                    response["network_task"]['services'].append("("+source_service["service"]+"/"+source_service["endpoint"]+", "+target_service["service"]+"/"+target_service["endpoint"]+")")
                    response["network_task"]["statuses"] += res_payload["network_task"]["statuses"]
                    response["network_task"]['statuses'].append(res.status)

            else: # "synchronous"
                async with ClientSession() as session:
                    for i in range(forward_ratio):
                        res = await session.post(dst, data=json_data, headers=forward_headers)
                        res_payload = await res.json()

                        # Concatenate json responses
                        response["cpu_task"]["services"] += res_payload["cpu_task"]["services"]
                        response["cpu_task"]['services'].append("("+source_service["service"]+"/"+source_service["endpoint"]+", "+target_service["service"]+"/"+target_service["endpoint"]+")")
                        response["cpu_task"]["statuses"] += res_payload["cpu_task"]["statuses"]

                        response["network_task"]["services"] += res_payload["network_task"]["services"]
                        response["network_task"]['services'].append("("+source_service["service"]+"/"+source_service["endpoint"]+", "+target_service["service"]+"/"+target_service["endpoint"]+")")
                        response["network_task"]["statuses"] += res_payload["network_task"]["statuses"]
                        response["network_task"]['statuses'].append(res.status)
                
    return response