#!/bin/bash -e

set -o pipefail || true  # trace ERR through pipes
set -o errtrace || true # trace ERR through commands and functions
set -o errexit || true  # exit the script if any statement returns a non-true return value

on_error(){
    CODE="${?}" && \
    set +x && \
    printf "[ERROR] Error(code=%s) occurred at %s:%s command: %s\n" \
        "${CODE}" "${BASH_SOURCE}" "${LINENO}" "${BASH_COMMAND}"
}
trap on_error ERR

_CLUSTER_NAME=aichatworkspace

kind delete cluster --name "${_CLUSTER_NAME}"
sleep 1
kind create cluster --name "${_CLUSTER_NAME}" --config hack/kind-cluster.yaml
sleep 1

# Install Ingress controller
# setting domains to *.localtest.me
kubectl apply -f hack/ingress-deploy.yaml

# Install KEDA
# relying on a couple of scalers (http-add-on and the kubernetes workload)
helm install keda kedacore/keda --create-namespace --namespace keda

# Install http-add-on to handle scaling to zero
helm install http-add-on kedacore/keda-add-ons-http  \
  --create-namespace --namespace keda \
  --set interceptor.responseHeaderTimeout=120s

# Setup AIChat Workspace Operator
make generate
make manifests
make install
IMG=aichatworkspace:v1 make docker-build
kind load docker-image aichatworkspace:v1 --name "${_CLUSTER_NAME}"
IMG=aichatworkspace:v1 make deploy
kubectl rollout restart deploy aichat-workspace-operator-controller-manager -n aichat-workspace-operator-system
kubectl rollout status deploy aichat-workspace-operator-controller-manager -n aichat-workspace-operator-system
kubectl get po -n aichat-workspace-operator-system
kubectl get node -o wide
