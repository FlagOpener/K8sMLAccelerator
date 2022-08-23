
#!/bin/bash
set -e

NAMESPACE=${NAMESPACE:-"nezha-demo"}

start() {
    kubectl create ns ${NAMESPACE} || true
    # create nginx config as a configmap
    kubectl create -n ${NAMESPACE} configmap nginx-proxy --from-file=nginx.conf || true
    # create nginx and svc
    kubectl apply -n ${NAMESPACE} -f nginx.yaml
    # create csr
    ../../deploy/create-signed-crt.sh
    # create webhook svc
    CA_BUNDLE=$(kubectl get configmap -n kube-system extension-apiserver-authentication -o=jsonpath='{.data.client-ca-file}' | base64 | tr -d '\n')
    cat ../../deploy/mutatingwebhook.yaml | sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g" | kubectl apply -f -
    # patch host aliases
    SVC=$(kubectl get svc -n ${NAMESPACE} proxy-cache -o jsonpath={.spec.clusterIP})    
    SERVERS=$(grep server_name nginx.conf |tr -d ';' |awk '{print $2}')
    file=$(mktemp temp.XXX.yaml)
    cat > ${file} <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: hostaliases-config
data: