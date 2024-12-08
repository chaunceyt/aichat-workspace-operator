# AIChat Workspace Operator

Kubernetes Operator that deploys Open WebUI along with Ollama for the OpenAI-compatible API, creating a multi-tenant "hosting" platform for private AIChat environments

The operator will 
- create each of the components here: https://github.com/open-webui/open-webui/tree/main/kubernetes/manifest/base
- plus serviceaccounts for each workload
- init the download of the requested models
- reconcile changes to replica, and resourcequota

Each pod will have securitycontext defined making the workload secure
Network policies to prevent traffic outside the namespace from interacting with endpoint

Database backend for AIChat workspaces:
- default sqlite3
- mysql
- postgres (should this be cnpg or custom)

Focus on xsmall llms
sLLMs
- qwen2.5-coder:3b
- smollm2
- phi3.5
- codegemma:2b
