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
import os

# script_dir = os.path.dirname(__file__)
configFilePath = os.environ['CONF']
# configFilePath = os.path.join(script_dir, 'conf.json')
logger = logging.getLogger(__name__)

def process_configfile():
    with open(configFilePath) as f:
        d = json.load(f)
    return d

config_data = process_configfile()

HOP = config_data["Hop"]
# TYPE = config_data["Type"]
INFLUXDB_DATABASE = "latency"
INFLUXDB_HOST = "influxdb"
INFLUXDB_PASSWORD = "root"
INFLUXDB_USERNAME = "root"