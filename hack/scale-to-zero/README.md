# Scaling to Zero
Working to optimizing resource utilization, scaling to Zero is an important feature. This configuration is setup to scale Open WebUI to zero after a period of time. This functionality will be integrated into the operator once everything works as expected.

## Current state

* Since the default state is scale-to-zero, need to address the timeout pulling container image when pull policy is Always vs IfNotPresent.

### Resources
* https://github.com/kedacore/http-add-on/blob/main/docs/install.md
* https://cloud.google.com/kubernetes-engine/docs/tutorials/scale-to-zero-using-keda