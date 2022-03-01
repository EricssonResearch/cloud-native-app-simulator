# Useful Commands

```bash
kubectl create configmap grafana-config \
    --from-file=grafana-dashboard-provider.yaml=grafana-dashboard-provider.yaml
```

```bash
kubectl create secret generic grafana-creds \
  --from-literal=GF_SECURITY_ADMIN_USER=admin \
  --from-literal=GF_SECURITY_ADMIN_PASSWORD=graphsRcool
```
