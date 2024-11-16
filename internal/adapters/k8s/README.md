This package is responsible for reconciling the k8s cluster objects

		/**
		  Create the following:
		  some derived from https://github.com/open-webui/open-webui/tree/main/kubernetes/manifest/base
		  - namespace
		  - resourcequota
		  - serviceaccount for ollama api
		  - serviceaccount for openwebui
		  - statefulset for Ollama API
		  - service for Ollama API
		  - deployment for Open WebUI
		  - service for Open WebUI
		  - ingress for Open WebUI
		**/


ensure<object>

createServiceAccount
createDeployment
createStatefulSet
