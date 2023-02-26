#!/bin/bash
set -e

NAMESPACE=${NAMESPACE:-"nezha-demo2"}

start() {
    kubectl create ns ${NAMESPACE} || true
    # cr