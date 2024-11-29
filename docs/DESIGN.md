# Design

Here I will explain the design philosophy that I am applying to the AIChat Workspace Operator and API to govern how I would extend and maintain this application.

This is an evolutionary architecture design and it aims to evolve over time as new requirements, tech debt, and new technologies emerge.

# Directory Structure

The application structure is a result of running the required kubebuilder commands to initialize the operator.

`internal/`. Here I add packages that are internal to this project only.

`internal/adapters/`
This directory contains all the packages that provide logic for communicating with external resources a.k.a infrastructure. In other words these packages are a translation layer between the domain and a specific external technology. e.g. ollama, mysql, etc.

`internal/webapi/`
This directory contains the API endpoint for managing an AIChat Workspace (wip)

