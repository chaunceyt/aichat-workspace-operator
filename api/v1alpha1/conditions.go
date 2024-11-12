package v1alpha1

const (
	// ConditionTypeReady represents the fact that the reconciliation of
	// the resource has succeeded.
	ConditionTypeReady string = "Ready"

	// ConditionTypeConfigMapReady represents the fact that the reconciliation of
	// the ConfigMap has succeeded.
	ConditionTypeConfigMapReady string = "ConfigMapReady"

	// ReconciliationSucceededReason represents the fact that reconciliation has succeeded.
	ReconciliationSucceededReason string = "ReconciliationSucceeded"

	// ReconciliationFailedReason represents the fact that reconciliation has failed.
	ReconciliationFailedReason string = "ReconciliationFailed"

	// ProgressingReason represents the fact that the reconciliation of the
	// resource is underway.
	ProgressingReason string = "Progressing"
)
