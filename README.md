# kube-obituaries
Archives pod manifests at termination.

![alt text](https://github.com/lannparty/kube-obituaries/blob/main/kube-obituaries.png?raw=true)

Cluster Installation
```
kubectl apply -f manifests/kube-obituaries.yaml
```

Kubectl Plugin Installation
```
chmod +x kubectl-plugin/kubectl-ob
cp kubectl-plugin/kubectl-ob /usr/local/bin
```

Usage:
```
kubectl ob <DELETED_POD_NAME>
```
Feed into jq for pretty print.
```
kubectl ob <DELETED_POD_NAME> |jq -r
```
