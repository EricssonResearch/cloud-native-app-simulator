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

import json
import logging
import uuid
import requests

FORMATTED_REMOTE_URL = "http://{0}.edge-namespace:80{1}"

logger = logging.getLogger(__name__)

def attend_request(service_config, headers):
    # TODO: Asynchronous attendance of requests not supported yet
    task_id = str(uuid.uuid4())
    task_config = {}
    task_config["task_id"] = task_id
    task_config["cpu_consumption"] = service_config["cpuConsumption"]
    task_config["network_consumption"] = service_config["networkConsumption"]
    task_config["memory_consumption"] = service_config["memoryConsumption"]
    response_object = execute_task(task_config)

    called_services = service_config["calledServices"]
    for remote_svc in called_services:
        make_request(remote_svc, headers)

    return response_object

def execute_task(task_config):
    # TODO: Implement resource stress emulation...
    response_object = {
        "status": "Task executed",
        "data": {
            "task_id": task_config["task_id"]
        }
    }
    return response_object

def make_request(remote_service, forward_headers={}):
    # TODO: Request forwarding to a service on a particular cluster is not supported yet
    # TODO: Requests for other protocols than html are not supported yet
    logger.info(remote_service)

    dst = FORMATTED_REMOTE_URL.format(remote_service["service"], remote_service["endpoint"])
    traffic_forward_ratio = remote_service["traffic_forward_ratio"]
    request_type = remote_service["requests"]

    # TODO: Asynchronous forwarding of traffic not supported yet
    for(_ in range(traffic_forward_ratio):
        forward_headers.update({'Content-type' : 'application/json'})
        res = requests.post(dst, headers=forward_headers)