# AIChat Workspace Operator

Create AIChat Workspaces, using [Open WebUI](https://openwebui.com/) as the front-end to [Ollama](https://ollama.com/).

## Description

Many companies require an assignment, so they can see how one would solve a probem they have defined, during the interview process.

This is a pre self-assignment:
- Write an Kubernetes Operator that manages an application.
- Write an API that handles the CRUD for the CR built.

The operator will create each component needed to run what an AIChat Workspace. Consisting of components defined in this [git repo](https://github.com/open-webui/open-webui/tree/main/kubernetes/manifest/base) along with a couple of more objects (i.e. serviceaccounts)

The API will be a front-end to the operator. Removing the need for direct interaction with Kubernetes, providing all of the CRUD actions needed to manage an `kind: AIChatWorkspace`

When a `kind: AIChatWorkspace` is submitted, the following objects will be created.

### ResourceQuota
Since several AIChatWorkspaces may share a cluster with a fixed number of nodes, there could be a concern that one workspace could use more than its fair share of cluster resources.

Resource quotas will be used to address the resources usage per tenant (namespace).

### ServiceAccounts
To ensure each workload created doesn't get attached to the `default` service account in the namespace created for the workspace, for each workload the operator will create a serviceaccont specific for the pod.

- serviceaccount for ollama api
- serviceaccount for openwebui

### StatefulSet (Ollama API)
A StatefulSet will run the Ollama API along with a Kubernetes Service to route traffic to the pod.

- statefulset for Ollama API
- service for Ollama API

### Deployment (Open WebUI)
A Deployment will run the Open WebUI workload along with a Kubernetes Ingress and Service to route traffic to the pod.

- deployment for Open WebUI
- service for Open WebUI
- ingress for Open WebUI

## Minimum viable product

The goal of this self-assignment is to show an ability to devlop solution around the stated [git repo](https://github.com/open-webui/open-webui/tree/main/kubernetes/manifest/base) , implement a pattern for testing the solution e2e, implement automation to deploy to "production".

MVP = the ability to create an `kind: AIChatWorkspace` and connect to the workspace via the ingress endpoint or port-forward with `ENV=dev` set using `kustomize build config/samples/ | kubectl -n aichat-workspace-operator-system apply -f -`.

##  After getting started

Initial work on the self-assignment...

- Set a goal to leverage some of the stated benefits of a Domain Driven, Data Oriented Architecture building this operator.
- Identified initial Kubebuilder [workflow](notes/kubebuilder-workflow.md)
- Implemented the [AIChatWorkspaceSpec](api/v1alpha1/aichatworkspace_types.go)
- Implemented the [Reconcile](internal/controller/) logic
- At time of writing this line, the MVP is achieved, submitting a CR manifest results in stated items being created. One can port-forward and download a Model and submit a prompt, which will result in a reply from the LLM.