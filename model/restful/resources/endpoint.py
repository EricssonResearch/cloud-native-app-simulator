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

from flask_restful import Resource, abort
from path import SERVICE_CONFIG
from restful.utils import task


def not_found(endpoint):
    abort(404, message="Endpoint {} doesn't exist".format(endpoint))


class Endpoint(Resource):
    def get(self, endpoint=None):
        if endpoint is None:
            message = {"status": "ok"}
            return message
        not_found(endpoint)
    
    def post(self, endpoint=None):
        if endpoint is None:
            message = {"status": "ok"}
            return message
        else:
            for ep in SERVICE_CONFIG['endpoints']:
                if ep['name'] == endpoint:
                    res = task.run_task(service_endpoint=ep)
                    return res
        not_found(endpoint)


