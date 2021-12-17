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

import sys
import os
import json
from redis import Redis, RedisError
from rq import Queue, Connection, SimpleWorker
from src.service.logger import initialize_logging

configFilePath = os.environ["CONF"]
with open(configFilePath) as f:
    d = json.load(f)

redis = Redis(host="localhost", db=0, socket_connect_timeout=2)
listen = ['high', 'default']


if __name__ == "__main__":
    initialize_logging()
    with Connection(redis):
        worker = SimpleWorker([Queue(name=name, is_async=False) for name in listen])
        worker.work()
