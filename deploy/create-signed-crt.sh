#!/bin/bash

set -e

usage() {
    cat <<EOF
Generate certificate suitable for use with an hostaliases-injector webhook service.
This script uses k8s' CertificateSigningRequest API t