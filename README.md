# kube-obituaries
Archives pod manifests at termination.

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
