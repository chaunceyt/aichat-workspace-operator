# Design

Here I will explain the design philosophy that I am applying to the AIChat Workspace Operator and API to govern how I would extend and maintain this application. Initially, I focused on making things work and then focus on cleaning up.

Not being an expert, this is an evolutionary architecture design and will evolve over time as new requirements, tech debt, and new technologies emerge. Also, as I'm reading Domain-Driven Design with Golang, I'm hoping to implement that pattern and ensure the operator works with each version of Golang and Kubernetes.

# Directory Structure

The application structure is a result of running the required kubebuilder commands to initialize the operator.

### api/v1alpha1

This was created during the initialization of the project using kubebuilder. It contains the Spec definition for the AIChatWorkspace

### docs/

This directory contains docs about the project.

### hack/
This was created during the initialization of the project using kubebuilder. I used it for staging files use during the provisioning of the local development environment. It will contain configuration that will be provisioned by the operator at some point in the future.


### internal/

Here I add packages that are internal to this project only.

### internal/adapters/
This directory contains all the packages that provide logic for communicating with external resources a.k.a infrastructure. In other words these packages are a translation layer between the domain and a specific external technology. e.g. ollama, mysql, etc.

### internal/adapters/auth
This directory contains auth login, generating, refreshing, and validating JWT tokens.

### internal/adapters/database
This directory contains database routines to be performed by the API 

### internal/adapters/k8s
This directory contains functions that represent K8s objections being created in the cluster. i.e. namespace, serviceaccount, persistantvolumeclaim, service, etc.

### internal/adapters/middlewares
This directory contains authz for interacting with protected endpoints.

### internal/adapters/models
This directory contains data schema representation.

### internal/adapters/utils
This directory contains utilities use within the operator and api

### internal/controller
This directory contains the reconcile logic for the operator.

### internal/webapi/
This directory contains the API endpoint for managing an AIChat Workspace (wip). Will support user registration, login and managining AIChat workspace.

### notes/

This directory contains notes used as this project is being developed.

### test/
This was created during the initialization of the project using kubebuilder. The kyverno chainsaw test was added to this directory

