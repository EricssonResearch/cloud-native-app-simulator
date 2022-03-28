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


