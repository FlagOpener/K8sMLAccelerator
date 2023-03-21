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
    cat ../../deploy/mutatingwebhook.yaml | sed -e