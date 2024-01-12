package v1alpha1

const (
	PendingState = "PENDING"
	ReadyState   = "READY"
	TaintedState = "TAINTED"
)

type ValueSource struct {
	SecretKeyRef *SecretReference `json:"secretKeyRef,omitempty"`
}

type SecretReference struct {
	// Name of the secret
	Name string `json:"name"`

	// Key of the secret
	Key string `json:"key"`
}
