#!/bin/bash

# optional argument handling
if [[ "$1" == "version" ]]
then
    echo "1.0.0"
    exit 0
fi

# optional argument handling
if [[ "$1" == "config" ]]
then
    echo "$KUBECONFIG"
    exit 0
fi

if [[ "$1" == "get"  ]] && [[ "$2" == "ob" ]]
then
    kubectl exec -it -n kube-system kube-obituary -- sqlite3 /opt/kube-obituaries/obituaries.db 'select * from pods where name = $3'
    exit 0
fi
