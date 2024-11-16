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

In another terminal, allowing the logs to stream.

```
make run
```

## Create CR

```
kubectl create ns a-team
kustomize build config/samples/ | kubectl -n a-team apply -f -
```