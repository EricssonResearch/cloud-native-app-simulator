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

from wsgiref import headers
from flask import Blueprint, jsonify, request
from aiohttp import ClientSession
import asyncio
import subprocess
import sys
from concurrent.futures import ThreadPoolExecutor, as_completed
import requests

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

    res_payload = create_payload(service_endpoint["network_complexity"]["response_payload_size"])

    source_svc = {}
    source_svc["service"] = service_name
    source_svc["endpoint"] = service_endpoint["name"]

    response = create_response()

    execution_mode = service_endpoint["execution_mode"]
    if execution_mode == "sequential":
        # Network task
        if ("network_complexity" in service_endpoint) and (len(service_endpoint["network_complexity"]["called_services"]) > 0):
            nw_response, _ = run_network_task(source_svc, service_endpoint, headers, res_payload)
            response = concatenate_response_simple(response, nw_response)

        # CPU task
        if ("cpu_complexity" in service_endpoint) and len(service_endpoint["cpu_complexity"]["execution_time"]) > 0:
            cpu_response, _ = execute_cpu_bounded_task(conf=service_endpoint["cpu_complexity"])
            response["cpu_task"]["services"].append(source_svc["service"]+"/"+source_svc["endpoint"])
            response["cpu_task"]["statuses"].append(cpu_response)

        # Memory task
        if ("memory_complexity" in service_endpoint) and len(service_endpoint["memory_complexity"]["execution_time"]) > 0:
            mem_response, _ = execute_memory_bounded_task(conf=service_endpoint["memory_complexity"])
            response["memory_task"]["services"].append(source_svc["service"]+"/"+source_svc["endpoint"])
            response["memory_task"]["statuses"].append(mem_response)

    else: # "parallel"
        executor = ThreadPoolExecutor(max_workers=3)
        task_futures = []

        # Network task
        if ("network_complexity" in service_endpoint) and (len(service_endpoint["network_complexity"]["called_services"]) > 0):
            nw_future = executor.submit(run_network_task, source_svc, service_endpoint, headers, res_payload)
            task_futures.append(nw_future)

        # CPU task
        if ("cpu_complexity" in service_endpoint) and len(service_endpoint["cpu_complexity"]["execution_time"]) > 0:
            cpu_future = executor.submit(execute_cpu_bounded_task, service_endpoint["cpu_complexity"])
            task_futures.append(cpu_future)

        # Memory task
        if ("memory_complexity" in service_endpoint) and len(service_endpoint["memory_complexity"]["execution_time"]) > 0:
            mem_future = executor.submit(execute_memory_bounded_task, service_endpoint["memory_complexity"])
            task_futures.append(mem_future)

        # Wait until all threads are done with their tasks
        for future in as_completed(task_futures):
            r, task_type = future.result()
            if task_type == "network":
                response = concatenate_response_simple(response, r)
            elif task_type == "cpu":
                response["cpu_task"]["services"].append(source_svc["service"]+"/"+source_svc["endpoint"])
                response["cpu_task"]["statuses"].append(r)
            elif task_type == "memory":
                response["memory_task"]["services"].append(source_svc["service"]+"/"+source_svc["endpoint"])
                response["memory_task"]["statuses"].append(r)

        executor.shutdown()
  
    return response


def execute_cpu_bounded_task(conf):
    if len(conf["cpu_affinity"]) > 0:
        res = subprocess.run(['stress-ng --class cpu --cpu %s --cpu-method %s --taskset %s --cpu-load %s --timeout %s --metrics-brief' % (conf["workers"], conf["method"], ",".join(str(cpu_id) for cpu_id in conf["cpu_affinity"]), conf["cpu_load"], conf["execution_time"])], capture_output=True, shell=True)
    else:
        res = subprocess.run(['stress-ng --class cpu --cpu %s --cpu-method %s --cpu-load %s --timeout %s --metrics-brief' % (conf["workers"], conf["method"], conf["cpu_load"], conf["execution_time"])], capture_output=True, shell=True)

    return res.stderr.decode("utf-8"), "cpu"


def execute_memory_bounded_task(conf):
    res = subprocess.run(['stress-ng --class memory --vm %s --vm-method %s --vm-bytes %s --timeout %s --metrics-brief' % (conf["workers"], conf["method"], conf["bytes_load"], conf["execution_time"])], capture_output=True, shell=True)

    return res.stderr.decode("utf-8"), "memory"


def run_network_task(source_svc, service_endpoint, headers, res_payload):
    
    if service_endpoint["network_complexity"]["forward_requests"] == "asynchronous":
        asyncio.set_event_loop(asyncio.new_event_loop())
        loop = asyncio.get_event_loop()
        nw_response = loop.run_until_complete(async_network_task(source_svc, service_endpoint, headers, res_payload))
    else: # "synchronous"
        nw_response = sync_network_task(source_svc, service_endpoint, headers, res_payload)

    return nw_response, "network"


async def async_network_task(source_svc, service_endpoint, headers, res_payload):

    async with ClientSession() as session:
        response = create_response()
        response["network_task"]["payload"] = res_payload

        io_tasks = []

        if len(service_endpoint["network_complexity"]["called_services"]) > 0:
            for target_svc in service_endpoint["network_complexity"]["called_services"]:
                io_task = asyncio.create_task(async_execute_io_bounded_task(session=session, source_service=source_svc, target_service=target_svc, forward_headers=headers))
                io_tasks.append(io_task)
            services = await asyncio.gather(*io_tasks)

        if len(service_endpoint["network_complexity"]["called_services"]) > 0:
            for svc in services:
                response = concatenate_response_simple(response, svc)

    return response


def sync_network_task(source_svc, service_endpoint, headers, res_payload):

    response = create_response()
    response["network_task"]["payload"] = res_payload
    
    if len(service_endpoint["network_complexity"]["called_services"]) > 0:
        for target_svc in service_endpoint["network_complexity"]["called_services"]:
            res = sync_execute_io_bounded_task(source_service=source_svc, target_service=target_svc, forward_headers=headers)
            response = concatenate_response_simple(response, res)

    return response


async def async_execute_io_bounded_task(session, source_service, target_service, forward_headers={}):

    dst = FORMATTED_REMOTE_URL.format(target_service["service"], target_service["port"], target_service["endpoint"])
    forward_headers.update({'Content-type' : 'application/json'})

    response = create_response()

    json_data = {}
    json_data["payload"] = create_payload(target_service["request_payload_size"])

    forward_ratio = target_service["traffic_forward_ratio"]
    if forward_ratio > 0:
        async with ClientSession() as session:
            io_tasks = []

            for i in range(forward_ratio):
                io_task = asyncio.create_task(session.post(dst, data=json_data, headers=forward_headers))
                io_tasks.append(io_task)
            calls = await asyncio.gather(*io_tasks)

        for res in calls:
            res_payload = await res.json()
            response = concatenate_response(response, res_payload, source_service, target_service)
            response["network_task"]['statuses'].append(res.status)

    return response


def sync_execute_io_bounded_task(source_service, target_service, forward_headers={}):

    dst = FORMATTED_REMOTE_URL.format(target_service["service"], target_service["port"], target_service["endpoint"])
    forward_headers.update({'Content-type' : 'application/json'})

    response = create_response()

    json_data = {}
    json_data["payload"] = create_payload(target_service["request_payload_size"])

    forward_ratio = target_service["traffic_forward_ratio"]       
    if forward_ratio > 0:
            for i in range(forward_ratio):
                res = requests.post(dst, data=json_data, headers=forward_headers)
                response = concatenate_response(response, res.json(), source_service, target_service)
                response["network_task"]['statuses'].append(res.status_code)
                
    return response


def create_payload(payload_size):

    request_payload = subprocess.run(['cat /dev/urandom | tr -dc "[:alnum:]" | head -c${1:-%s}' % payload_size], capture_output=True, shell=True)
    
    return request_payload.stdout.decode("utf-8")


def create_response():

    response = {}
    response["cpu_task"] = {}
    response["cpu_task"]["services"] = []
    response["cpu_task"]["statuses"] = []

    response["memory_task"] = {}
    response["memory_task"]["services"] = []
    response["memory_task"]["statuses"] = []

    response["network_task"] = {}
    response["network_task"]["services"] = []
    response["network_task"]["statuses"] = []
    response["network_task"]["payload"] = ""

    return response


def concatenate_response_simple(response, res):

    response["cpu_task"]["services"] += res["cpu_task"]["services"]
    response["cpu_task"]["statuses"] += res["cpu_task"]["statuses"]

    response["memory_task"]["services"] += res["memory_task"]["services"]
    response["memory_task"]["statuses"] += res["memory_task"]["statuses"]

    response["network_task"]["services"] += res["network_task"]["services"]
    response["network_task"]["statuses"] += res["network_task"]["statuses"]
    response["network_task"]["payload"] += res["network_task"]["payload"]

    return response


def concatenate_response(response, res_payload, source_service, target_service):

    response["cpu_task"]["services"] += res_payload["cpu_task"]["services"]
    response["cpu_task"]["statuses"] += res_payload["cpu_task"]["statuses"]

    response["memory_task"]["services"] += res_payload["memory_task"]["services"]
    response["memory_task"]["statuses"] += res_payload["memory_task"]["statuses"]

    response["network_task"]["services"] += res_payload["network_task"]["services"]
    response["network_task"]['services'].append("("+source_service["service"]+"/"+source_service["endpoint"]+", "+target_service["service"]+"/"+target_service["endpoint"]+")")
    response["network_task"]["statuses"] += res_payload["network_task"]["statuses"]
    response["network_task"]["payload"] = res_payload["network_task"]["payload"]

    return response