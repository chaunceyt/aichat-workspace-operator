#!/bin/bash

kind delete cluster
sleep 1
kind create cluster --config hack/kind-cluster.yaml
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
kind load docker-image aichatworkspace:v1
IMG=aichatworkspace:v1 make deploy
kubectl rollout restart deploy aichat-workspace-operator-controller-manager -n aichat-workspace-operator-system
kubectl rollout status deploy aichat-workspace-operator-controller-manager -n aichat-workspace-operator-system
kubectl get po -n aichat-workspace-operator-system

sleep 1
kustomize build config/samples/ | kubectl -n aichat-workspace-operator-system apply -f -