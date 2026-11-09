package vmapi

import "time"

type DeliveryPhase string

const (
	DeliveryUnknown    DeliveryPhase = ""
	DeliveryAccepted   DeliveryPhase = "accepted"
	DeliveryDispatched DeliveryPhase = "dispatched"
	DeliveryDelivered  DeliveryPhase = "delivered"
	DeliveryFailed     DeliveryPhase = "failed"
)

// Exec preserves argv and exact stdin. Deadline is required/absolute; output bound
// is positive and aggregate. Context bounds this call; Deadline bounds guest work
// independently where supported. Adapter configuration owns transport, credentials
// and inherited host-key verification. No automatic retry or Recover of arbitrary
// Exec is permitted: delivery is not exactly once.
type ExecRequest struct {
	InstanceMutation
	Argv             []string
	Stdin            []byte
	Environment      map[string]string
	WorkingDirectory string
	Deadline         time.Time
	OutputLimitBytes uint32
}

// Nonzero exit is a result, transport failure an error, possibly with uncertain
// delivery. Nil ExitCode is not zero. Output is sensitive: only this explicit
// result exposes it, never events, errors, observations or ordinary diagnostics.
type ExecResult struct {
	Instance                         InstanceRef
	Delivery                         DeliveryPhase
	ExitCode                         *int
	Stdout, Stderr                   []byte
	StdoutTruncated, StderrTruncated bool
}
