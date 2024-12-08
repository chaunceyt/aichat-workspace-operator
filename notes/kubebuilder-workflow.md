# Kubebuilder workflow

## Init project

Create an empty folder and run the following command within that folder.

```
kubebuilder init --domain aichatworkspaces.io --repo github.com/chaunceyt/aichat-workspace-operator
```

## Create API

```
kubebuilder create api --group apps --version v1alpha1 --kind AIChatWorkspace
```

## Define the Spec for the API

For this operator edit the following file.

```
vi api/v1alpha1/aichatworkspace_types.go
```

## Create/Update the DeepCopy implementations

Run this command after making changes to `api/v1alpha1/aichatworkspace_types.go` file

```
make generate
```

## Create/Update the CRD manifests

```
make manifests
```

## Install the crds

```
make install
```

## Running locally

Run `./from-scratch.sh` and in another terminal monitor the logs of the operator.

## Commands Excuted during development

### After running from-scratch.sh

Setup initial AIChat Workspace

```
kustomize build config/samples/ | kubectl -n aichat-workspace-operator-system apply -f -
kubectl rollout status deploy team-a-aichat-openwebui -n team-a-aichat
kubectl rollout status sts team-a-aichat-ollama -n team-a-aichat
```

### scale-to-zero setup

```
kubectl apply -f hack/externalservice-openwebui.yaml -n team-a-aichat
kubectl apply -f hack/externalservice-ollama.yaml -n team-a-aichat
kubectl apply -f hack/open-webui-ingress.yaml -n team-a-aichat
kubectl apply -f hack/ollama-ingress.yaml -n team-a-aichat
kubectl apply -f hack/keda-httpscaledobject-openwebui.yaml -n team-a-aichat
kubectl apply -f hack/keda-k8s-workload-scaledobject.yaml -n team-a-aichat
```

### rebuild container image

```
IMG=aichatworkspace:v1 make docker-build
kind load docker-image aichatworkspace:v1
kubectl rollout restart deploy aichat-workspace-operator-controller-manager -n aichat-workspace-operator-system
kubectl rollout status deploy aichat-workspace-operator-controller-manager -n aichat-workspace-operator-system
kubectl get po -n aichat-workspace-operator-system
```
