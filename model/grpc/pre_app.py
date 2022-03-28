from jinja2 import Template
import json
import os

with open("../config/conf.json", "r") as f:
    conf = json.load(f)

with open("template/app.jinja", "r") as j:
    temp = j.read()


proto_temp = Template(temp)
filled_temp = proto_temp.render({"endpoints": conf['endpoints'], "service_name": os.environ['SERVICE_NAME']})

with open("app.py", "w") as output:
    output.write(filled_temp)
output.close()

