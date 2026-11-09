package vmapi

import "time"

type BootstrapEventKind string

const (
	EventBaseReady         BootstrapEventKind = "base-ready"
	EventTemplateReady     BootstrapEventKind = "template-ready"
	EventDockerEngineReady BootstrapEventKind = "docker-engine-ready"
	EventDeliverySucceeded BootstrapEventKind = "delivery-succeeded"
	EventDeliveryFailed    BootstrapEventKind = "delivery-failed"
)

type TemplateEvidence struct{ MAC, BootIdentity string }
type DockerEvidence struct{ EngineIdentity string }

// Events carry only permitted typed evidence, never free-form output or secrets.
// TemplateReady requires Template; DockerEngineReady requires Docker. Invalid
// evidence is never success. DeliverySucceeded corresponds to PULSE_HANDOFF_OK,
// independently from base/template/application readiness and power transitions.
type BootstrapEvent struct {
	Cursor            Cursor
	At                time.Time
	Instance          InstanceRef
	IntentFingerprint string
	Kind              BootstrapEventKind
	Template          *TemplateEvidence
	Docker            *DockerEvidence
}

// Positive Limit is bounded by Info; After is backend-domain/exact-instance/stream scoped.
// Expiry returns ErrCursorExpired. Reads are bounded, nonblocking, ordered and
// observational. Adapters retain equivalent sanitized evidence across restart.
type EventRequest struct {
	Instance InstanceRef
	After    Cursor
	Limit    uint32
}

// CaughtUp means no further retained events now, not readiness or stream completion.
type EventPage struct {
	Items    []BootstrapEvent
	Next     Cursor
	CaughtUp bool
}
