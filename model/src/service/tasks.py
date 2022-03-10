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
import uuid

# TODO: So far, we only support a hard-coded namespace. For more flexible support of namespaces we will need to pass that info as part of the config map
# TODO: So far, we only support http client
FORMATTED_REMOTE_URL = "http://{0}:{1}{2}"

logger = logging.getLogger(__name__)

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

async def execute_io_bounded_task(session, target_service, forward_headers={}):
    # TODO: Request forwarding to a service on a particular cluster is not supported yet
    # TODO: Requests for other protocols than html are not supported yet
    logger.info(target_service)

    dst = FORMATTED_REMOTE_URL.format(target_service["service"], target_service["port"], target_service["endpoint"])
    forward_headers.update({'Content-type' : 'application/json'})

    # TODO: traffic_forward_ratio not supported yet
    traffic_forward_ratio = target_service["traffic_forward_ratio"]

    res = await session.post(dst, headers=forward_headers)
    return {'service': res.url, 'status': res.status}