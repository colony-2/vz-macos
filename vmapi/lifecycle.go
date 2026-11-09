package vmapi

import "time"

type CloneMode string

const (
	CloneCold      CloneMode = "cold"
	CloneSuspended CloneMode = "suspended"
)

type StopMode string

const (
	StopGraceful StopMode = "graceful"
	StopForce    StopMode = "force"
)

type Dispatch string

const (
	DispatchUnknown  Dispatch = ""
	DispatchAccepted Dispatch = "accepted"
	DispatchRejected Dispatch = "rejected"
)

// Bootstrap input is sensitive, privately persisted subject to advertised bounds.
// Deadline is absolute; recipe identity identifies eligible idempotent recovery.
type BootstrapIntent struct {
	Recipe   string
	Stdin    []byte
	Deadline time.Time
}
type EffectiveIntent struct {
	Allocation  AllocationID
	Owner       string
	Source      SourceSelector
	Resolved    ResolvedSource
	Mode        CloneMode
	Hardware    *Hardware
	Fingerprint string
}

// Create provisions, launches/resumes and starts bootstrap without waiting for
// readiness. Mode is explicit. Cold clones allow hardware changes/randomized MAC.
// Suspended clones preserve memory/MAC and require compatibility and exclusive
// collision checks including foreign objects; never silently fall back to cold.
// All mutations atomically validate immutable UID within their configured
// domain; replacements conflict. Persist native intent/delivery evidence before irreversible work.
// Definite rejection dispatches nothing. After possible dispatch return
// ErrUnknownOutcome with RequestID, even on cancellation (optionally joined).
// Reconcile through journal-backed name observation and effective-intent validation;
// never blindly replay Create when a retained record is missing.
type CreateRequest struct {
	Mutation
	Name              string
	Allocation        AllocationID
	Owner             string
	Source            SourceSelector
	Mode              CloneMode
	Hardware          *Hardware
	IntentFingerprint string
	Bootstrap         BootstrapIntent
}

// Successful acceptance includes stable UID and normalized nonsecret intent.
// Dispatch, power, delivery and core readiness are independent.
type CreateResult struct {
	Instance Instance
	Dispatch Dispatch
}

// Suspend completes only with observed nonrunning, desired suspended, known
// reached convergence and no active suspension. Delete completes only after exact
// UID and owned resources are confirmed absent, never a same-name replacement.
type MutationResult struct {
	Instance                   Instance
	Dispatch                   Dispatch
	Complete, AbsenceConfirmed bool
}
type StopRequest struct {
	InstanceMutation
	Mode StopMode
}
