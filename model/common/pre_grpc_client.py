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

from jinja2 import Template
import json

with open("../config/conf.json", "r") as f:
    conf = json.load(f)

with open("template/grpc_client.jinja", "r") as j:
    temp = j.read()

called_svc = []

for endpoint in conf['endpoints']:
    for svc in endpoint['calledServices']:
        if svc['protocol'] == "grpc":
            called_svc.append(svc)

grpc_temp = Template(temp)
filled_temp = grpc_temp.render({"called_svc": called_svc})

with open("grpc_client.py", "w") as output:
    output.write(filled_temp)
output.close()


