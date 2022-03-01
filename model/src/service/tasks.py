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

import datetime
import json
import logging
import time
import uuid
import random
import requests
from src.service.util import gen_dict_extract
from src.service import path

DEFAULT_STATUS_ADDRESS = "http://{0}.edge-namespace:80/status"
DEFAULT_ROOT_ADDRESS = "http://{0}.edge-namespace:80/"

hop = path.HOP
# task_type =  path.TYPE
logger = logging.getLogger(__name__)

def range_prod(lo,hi):
    if lo+1 < hi:
        mid = (hi+lo)//2
        return range_prod(lo,mid) * range_prod(mid+1,hi)
    if lo == hi:
        return lo
    return lo*hi

def treefactorial(n):
    if n < 2:
        return 1
    return range_prod(1, n)

def eat_cpu():
    num = random.randint(100, 500)
    treefactorial(num)

def eat_memory():
    num = random.randint(1,6)
    a = [1] * (10 ** num)
    b = [2] * (2 * 10 ** num)
    del b
    return a

def do_sleep():
    num = random.uniform(0.01, 0.02)
    time.sleep(num)
    return True

def make_request(task_id, chain_no, address, request_body, task, request_to_next_service, forward_headers={}):
    task_type = task
    if task_type == "cpu":
        eat_cpu()
    elif task_type == "memory":
        eat_memory()
    elif task_type == "sleep":
        do_sleep()
    v = list(gen_dict_extract(chain_no, hop))[0]
    logger.info(v)
    if v:
        for index,item in enumerate(v):
            dst = DEFAULT_ROOT_ADDRESS.format(item)
            d = {
                "previous"  : request_body,
                "chain_no" : chain_no,
                "task_id" : task_id,
                "task_type" : task_type,
                "address" : address,
                "request_type": request_to_next_service,
            }
            forward_headers.update({'Content-type' : 'application/json'})
            res = requests.post(dst, data=json.dumps(d, indent=4, default=str), headers=forward_headers)
    if v is None:
       req_id = str(uuid.uuid4())
       d = {
            "previous":request_body,
            "req_id": req_id,
            "chain_no": chain_no,
            "task_id" : task_id,
            "task_type" : task_type,
            "address" : address,
        }
       address = list(gen_dict_extract("initial", request_body))[0]
       dst = DEFAULT_STATUS_ADDRESS.format(address)
       forward_headers.update({'Content-type' : 'application/json'})
#       res = requests.post(dst, data=json.dumps(d, indent=4, default=str), headers=forward_headers)

def execute_task(json_data, address, headers):
    task_id = str(uuid.uuid4())
    chain_no = json_data["chain_no"]
    request_type = json_data["request_type"]
    request_to_next_service = str(int(request_type) + 1)
    request_task_type = list(gen_dict_extract("request_task_type", json_data))[0]
    if int(request_type) > len(request_task_type):
        task = random.choice(["cpu", "memory", "sleep"])
    else:
        task = list(gen_dict_extract(request_type, request_task_type))[0]
    make_request(task_id, chain_no, address, json_data, task, request_to_next_service, headers)
    response_object = {
        "status": "Task created",
        "data": {
            "task_id": task_id
        }
    }
    return response_object