package vmapi

import "time"

// Reads are observational: never create directories, repair or restart delivery.
// List includes all discoverable managed, stopped/source and foreign capacity
// consumers in the configured backend domain. Unassignable identities appear
// in Issues and capacity.
// Incomplete observation is unavailable, never evidence of absence/free capacity.
type ListRequest struct {
	Page PageRequest
}

// Get permits name discovery after uncertain Create. ExpectedUID mismatch conflicts.
type GetRequest struct {
	Name        string
	ExpectedUID InstanceUID
}

// ID is informational and cannot prove domain continuity or authorize routing.
type BackendInfo struct {
	ID           string
	Capabilities Capabilities
	Limits       Limits
	Capacity     Capacity
	Observation  Observation
}

// Baseline methods are required. Unknown capability values mean unsupported.
// Docker requires Exec, Stop and TCPGuestEndpoints before any Create.
type Capabilities struct {
	GuestOS                                                       []GuestOS
	CloneModes                                                    []CloneMode
	Exec, Stop, TCPGuestEndpoints                                 bool
	ProbeLocalSources, AdoptSource, RetainSource                  bool
	Recover, RestartDelivery, RemoveEndpoint, EndpointDiagnostics bool
	LifecycleEndpointCleanup                                      bool
}

// Available is advisory even when Complete: core reserves unobserved durable intents.
type Capacity struct {
	Total, ManagedUsed, ForeignUsed, Available Hardware
	Complete                                   bool
}

// All bounds and retention durations are positive advertised maxima.
type Limits struct {
	PageItems, EventItems, OutputBytes, DiagnosticBytes, RecoveryItems, ConditionItems uint32
	CommandInputBytes, BootstrapInputBytes                                             uint64
	EventRetention                                                                     time.Duration
}

// Limit must be positive and within Limits. Cursors bind query, backend domain and
// stable snapshot; expiry returns ErrCursorExpired. Restart enumeration on expiry.
type PageRequest struct {
	Limit  uint32
	Cursor Cursor
}
type InventoryIssue struct {
	Resource string
	Code     ErrorCode
	Message  string
}

// Empty Next ends enumeration. Complete describes the whole stable snapshot;
// partial retrieval never proves absence. Issues are bounded and sanitized.
type InstancePage struct {
	Items       []Instance
	Next        Cursor
	Snapshot    SnapshotID
	Observation Observation
	Complete    bool
	Issues      []InventoryIssue
}
