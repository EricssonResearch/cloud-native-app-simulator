# Useful Commands

```bash
kubectl create configmap grafana-config \
    --from-file=influxdb-datasource.yaml=influxdb-datasource.yaml \
    --from-file=grafana-dashboard-provider.yaml=grafana-dashboard-provider.yaml \
    --from-file=latency-dashboard.json=latency-dashboard.json
```

```bash
kubectl create secret generic influxdb-creds \
  --from-literal=INFLUXDB_DATABASE=latency \
  --from-literal=INFLUXDB_USERNAME=root \
  --from-literal=INFLUXDB_PASSWORD=root \
  --from-literal=INFLUXDB_HOST=influxdb
```

```bash
kubectl create secret generic grafana-creds \
  --from-literal=GF_SECURITY_ADMIN_USER=admin \
  --from-literal=GF_SECURITY_ADMIN_PASSWORD=graphsRcool
```
