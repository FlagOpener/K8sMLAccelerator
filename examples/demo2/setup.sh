#!/bin/bash
set -e

NAMESPACE=${NAMESPACE:-"nezha-demo2"}

start() {
    kubectl create ns ${NAMESPACE} || true
    # create s3-cache config as a configmap
    kubectl create -n ${NAMESPACE} configmap s3-cache-cfg --from-file=s3-cache.conf || true
    # create s3-cache and svc
    kubectl apply -n ${NAMESPACE} -f s3.yaml
    # create webhook svc
    cat ../../deploy/mutatingwebhook.yaml | sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g" | kubectl apply -f -
    # patch host aliases
    SVC=$(kubectl get svc -n ${NAMESPACE} s3-cache -o jsonpath={.spec.clusterIP})    
    SERVERS=$(grep server_name s3-cache.conf |tr -d ';' |awk '{print $2}')
    file=$(mktemp temp.XXX.yaml)
    cat > ${file} <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: hostaliases-config
data:
  config: |
      - 