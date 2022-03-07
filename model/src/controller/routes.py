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
from src.service.tasks import attend_request
from src.service import path

simple_page = Blueprint("simple_page", __name__,)
logger = logging.getLogger(__name__)

service_config = path.SERVICE_CONFIG

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

@simple_page.route(service_config["name"], methods=["POST", "GET"])
def run_task():
    if request.method == "POST":
        headers = getForwardHeaders(request)
        response_object = attend_request(service_config, headers)
        return jsonify(response_object), 200
    else:
        return "request received"

@simple_page.route("/", methods=["POST", "GET"])
def task_status():
    return "OK", 200

@simple_page.route("/loadtest", methods=["GET"])
def load_test():
    return "request received"