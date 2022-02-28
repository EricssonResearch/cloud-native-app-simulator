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

from src.service import path

def init_influxdb(host=path.INFLUXDB_HOST, db=path.INFLUXDB_DATABASE, user=path.INFLUXDB_USERNAME, password=path.INFLUXDB_PASSWORD):
    client = InfluxDBClient(host, 8086, user, password, db)
    return (client)

def create_influxdb(client, db_name=path.INFLUXDB_DATABASE):
    dbs = client.get_list_database()
    if not any(db['name'] == db_name for db in dbs):
        client.create_database(db_name)
    client.switch_database(db_name)

#client = init_influxdb()
#create_influxdb(client)