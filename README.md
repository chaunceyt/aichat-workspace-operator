# AIChat Workspace Operator

**Disclaimer**: This project is currently under development and may change rapidly, including breaking changes. Use with caution in production environments.

Create AIChat Workspaces powered by [Open WebUI](https://openwebui.com/) and [Ollama](https://ollama.com/)

## Objectives

Create a **LLM as a Service** using Ollama to provide the API for interacting with LLMs. 

Create a Kubernetes Operator that creates the `control-plane` for the **LLM as a Service**. The operator should manage the lifecycle of each Kubernetes resource needed to run the selected service in a namespace. By dynamically creating, and managing the resources needed to run the AIChat Workspace. The resources for each AIChat Workspace should be separated, reducing interference between tenants and optimizing resource utilization.

Create a Model from a modelfile that sets the `SYSTEM` prompt for the model using one-to-many [fabric/patterns](https://github.com/danielmiessler/fabric/tree/main/patterns) `SYSTEM` prompts. 


## Design and Implementation

### Design


The purpose of this design is to support running AIChat Workspaces within a Namespace as a Service offering.


![image](diagrams/aichatworkspace-operator-components.png)

Personas that will interact with the environment. Platform Engineers, and AIChatWorkspace owners

![image](diagrams/simple-persona-workflow.png)

## Implementation

* Used go version `1.23`
* Used Kubebuilder `v4.3.1` to generate most of the code.
* Used Kind to create the K8s cluster running `v1.31.0`.
* KEDA's `HTTPScaledObject` is being used to address the **scale-to-zero** requirement.
* Defaulted to xSmall LLMs i.e. `gemma2:2b`, `llama3.2:1b`, or `qwen2.5-coder:1.5b` in custom resource manifests. 

Legend: ❌ roadmap ✅ completed initial implementation

* ✅ Controller for `AIChatWorkspace` resources
* ✅ The `AIChatWorkspace` represents the resources needed to run [Open WebUI](https://openwebui.com/) and [Ollama](https://ollama.com/) in a namespace.
* ✅ When a new `AIChatWorkspace` is created. Create each Kubernetes resource located [here]](https://github.com/open-webui/open-webui/tree/main/kubernetes/manifest/base) as the base.
* ✅ Handle updates to `AIChatWorkspace` and owned resources
* ✅ Handle deletions of `AIChatWorkspace` ensuring the associated resources are also deleted.
* ✅ Handle pulling in requested models
* ✅ Create model from modelfile using a SYSTEM prompts from [fabric/patterns](https://github.com/danielmiessler/fabric/tree/main/patterns)
* ✅ API endpoint for register and login and calling a protected endpoint. (use: curl, postman, etc)
* Manage the lifecycle of each application (Open WebUI and Ollama)
* ✅ e2e testing (using Kyverno Chainsaw)
* ❌ Scale-to-Zero after no request are received for a period of time. (scale up on new requests) [testing](hack/scale-to-zero/)
* ❌ Add auth to the Ollama endpoint, consider envoy sidecar proxy providing auth or use basic-auth for ingress-nginx. (basic-auth, jwt) [testing](hack/envoy-sidecar/)
* ❌ List of resource in the describe of the aichatworkspace object. (pods, pvc, svc, models running, etc)
* ✅ Support for most of the `system.md` located under [fabric/patterns](https://github.com/danielmiessler/fabric/tree/main/patterns)
* ❌ API endpoint to manage AIChat workspace. (use: curl, postman, etc)
* ❌ Built in SRE that monitors the events of AIChat Workspace workloads and interact with LLM to identify a solution.
* ❌ Observability (i.e, grafana, loki, prometheus, promtail, and tempo)
* ❌ Helm chart (with unittest)
* ❌ Web Frontend

### Components that will be dynamically created and managed:

* ✅ Namespace - for isolation
* ✅ ResourceQuota - optimizing resource usage
* ✅ ServiceAccount - one per service (to ensure the workloads are not using the default)
* ✅ Statefulset - running Ollama
* ✅ Deployment - running Open WebUI
* ✅ PVC to store the downloaded LLMs
* ✅ Kubernetes Service for Ollama
* ✅ Kubernetes Service for Open WebUI
* ✅ Ingress object for Ollama
* ✅ Ingress object for Open WebUI
* ❌ KEDA HTTPScaledObject to scale the Open WebUI to zero after no requests are received based on `scaledownPeriod`.
* ❌ K8s ExternalService for open-webui scale-to-zero functionality
* ❌ NetworkPolicy allow traffic from ingress controller namespace to Open WebUI and Ollama

### Dependencies

This project has a number of external dependencies that contribute to the **LLM as a Service** solution

* ✅ Database for the AIChat Workspace API. User information
* ✅ Ingress controller
* ❌ CNI that supports network policies (optional)
* ❌ Statefulset running Ollama (will be used for SRE features)
* ❌ Vault to manage secrets for the API database
* ✅ KEDA for scale-to-zero of AIChat Workspaces

### Features before considered feature complete

# Getting Started

### Prerequisites
* go version v1.23.0+
* kubectl version v1.31.0+.
* A KinD Kubernetes cluster. (Developed/Tested using Kind v1.31.0+)


### To Test this project

* To run unit tests.

```sh
make test
```

* To run chainsaw e2e tests.

```sh
make chainsaw
```

### Run the project

Running the script below will create install everything needed to run this locally using a kind create cluster.

```sh
# Setup KinD cluster and install needed dependencies.
./from-scratch.sh

# Run the sample
kustomize build config/samples/ | kubectl -n aichat-workspace-operator-system apply -f -
kubectl rollout status deploy team-a-aichat-openwebui -n team-a-aichat
kubectl rollout status sts team-a-aichat-ollama -n team-a-aichat

# 
kubectl get all,ing,pvc,resourcequota -n team-a-aichat
```