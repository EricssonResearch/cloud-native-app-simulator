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
#from src.service.backend import client
import json

dtformat = '%Y-%m-%d %H:%M:%S.%f'
def gen_dict_extract(key, var):
    if hasattr(var,'items'):
        for k, v in var.items():
            if k == key:
                yield v
            if isinstance(v, dict):
                for result in gen_dict_extract(key, v):
                    yield result
            elif isinstance(v, list):
                for d in v:
                    for result in gen_dict_extract(key, d):
                        yield result

def create_point(username, measurement, value, time):
    json_body = {
        "measurement": measurement,
        "tags": {
            "user": username
        },
        "time": time,
        "fields": {
            "value": value
        }
    }

    return json_body

def create_and_save_network_data(json_data):
    json_body = []
    time = datetime.datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%SZ')
    i = list(gen_dict_extract("initial_time", json_data))
    f = list(gen_dict_extract("final_time", json_data))
    t = list(gen_dict_extract("task_type", json_data))
    num_cpu = t.count("cpu")
    num_mem = t.count("memory")
    num_sleep = t.count("sleep")
    latency = datetime.datetime.strptime(f[0], dtformat) - datetime.datetime.strptime(i[0], dtformat)
    delay = latency.seconds + (latency.microseconds / 1000000.0)
    mec_number = len(f)
    chain_no = int(json_data["chain_no"])
    data_points = {
        "network_latency": delay,
        "mec_number": mec_number,
        "chain_no": chain_no,
        "cpu_tasks": num_cpu,
        "mem_tasks": num_mem,
        "sleep_tasks": num_sleep,
    }
    for key, value in data_points.items():
        json_body.append((create_point("network", key, value, time)))
 #       client.write_points(json_body)
