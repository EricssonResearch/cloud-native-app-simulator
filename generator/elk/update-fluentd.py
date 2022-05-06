#
# Copyright 2021 Ericsson AB
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
import yaml, sys
from kubernetes import client, config

with open("fluentd.yaml") as f:
    y = yaml.safe_load_all(f)
    contexts = config.list_kube_config_contexts()
    for context in contexts[0]:
        if context['context']['cluster'] == sys.argv[3]:
            config.load_config(context=context['name'])
            for part in y:
                if part['kind'] == "ServiceAccount":
                    k8s_core_api_v1 = client.CoreV1Api()
                    resp = k8s_core_api_v1.create_namespaced_service_account(namespace="default", body=part)
                    print("Fluentd ServiceAccount created in " + sys.argv[3])
                elif part['kind'] == "ClusterRole":
                    k8s_RBAC_authorization_api = client.RbacAuthorizationV1Api()
                    resp = k8s_RBAC_authorization_api.create_cluster_role(body=part)
                    print("Fluentd ClusterRole created in " + sys.argv[3])
                elif part['kind'] == "ClusterRoleBinding":
                    k8s_RBAC_authorization_api = client.RbacAuthorizationV1Api()
                    resp = k8s_RBAC_authorization_api.create_cluster_role_binding(body=part)
                    print("Fluentd ClusterRoleBinding created in " + sys.argv[3])
                elif part['kind'] == "DaemonSet":
                    for env in part['spec']['template']['spec']['containers'][0]['env']:
                        if env['name'] == "FLUENT_ELASTICSEARCH_HOST":
                            env['value'] = sys.argv[1]
                        if env['name'] == "FLUENT_ELASTICSEARCH_PORT":
                            env['value'] = sys.argv[2]
                    k8s_apps_api = client.AppsV1Api()
                    resp = k8s_apps_api.create_namespaced_daemon_set(namespace="default", body=part)
                    print("Fluentd DeamonSet created in " + sys.argv[3])




