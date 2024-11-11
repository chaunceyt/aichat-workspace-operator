# Notes

aichat-workspace-operator - will create "AI Chat" environments with Open WebUI as the front-end to Ollama.

- https://openwebui.com/
- https://ollama.com/

Why?

* Improve ones understanding of using kubebuilder and the workflow centered around building Kubernetes operators.
* Improve ones understanding of controller-runtime.
* Learn

What?

Create an operator the creates AI Chat workspaces. That consist of Open WebUI configured to use an Ollama backend. The idea is to get a couple of values from an API request that would result in the creation of an AIChatWorkspace Custom resource.

When a AIChatWorkspace CR is submitted the following objects will be created:

* Namespace - for isolation
* ServiceAccount - one per service (just to ensure the workloads are not using the default)
* Pod - running ollama
* Deployment - running Open WebUI
* PVC to store the downloaded LLMs
* Kubernetes Service for Ollama 11434
* Kubernetes Service for Open WebUI 3000
* Ingress object for Open WebUI


We need very few inputs to create the workspace.

Name - name of the workspace. This will be used for the namespace, DNS address for workspace.
Models - list of models. This will be used to init the workspace by downloading the models.

Once submitted and the environment created, the DNS address and instructions on how to login are delivered to the requestor.

The API

I would like to not have one interact with Kubernetes and was thinking of creating an API to create, update, or delete the AIChatWorkspace CR. Something simple, maybe a cli.
