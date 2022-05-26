# Kubernetes Obituaries
Archives pod manifests at termination.

![alt text](https://github.com/lannparty/kob/blob/main/kob-arch.png?raw=true)

## Controller installation
```
kubectl apply -f manifests/kob.yaml
```

## Kubectl plugin installation
```
chmod +x kubectl-plugin/kubectl-ob
cp kubectl-plugin/kubectl-ob /usr/local/bin
```

## CLI Usage:
```
kubectl ob <DELETED_POD_NAME>
```
Feed into jq for pretty print
```
kubectl ob <DELETED_POD_NAME> |jq -r
```
Supports wildcard
```
kubectl ob calico% |jq -r
```
