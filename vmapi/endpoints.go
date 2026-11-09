package vmapi

import "net/netip"

// Ensure converges only for identical allocation, instance, binding and target.
// Validate exact immutable ownership for all endpoint operations; mismatches and
// occupied/substituted ports conflict. Restore across guest address changes or
// report uncertainty. Listening does not establish authenticated Docker readiness.
type EnsureEndpointRequest struct {
	InstanceMutation
	Allocation AllocationID
	Bind       netip.AddrPort
	GuestPort  uint16
}
type EndpointState string

const (
	EndpointUnknown   EndpointState = ""
	EndpointPending   EndpointState = "pending"
	EndpointListening EndpointState = "listening"
	EndpointAbsent    EndpointState = "absent"
)

// Bind is the actual exact worker TCP binding, valid IP and nonzero port.
// GuestPort is the nonzero guest TCP target, never vsock.
type Endpoint struct {
	Ref         EndpointRef
	Instance    InstanceRef
	Bind        netip.AddrPort
	GuestPort   uint16
	State       EndpointState
	Observation Observation
}
type EndpointRequest struct {
	Instance InstanceRef
	Endpoint EndpointRef
}
type RemoveEndpointRequest struct {
	InstanceMutation
	Endpoint EndpointRef
}

// Completion confirms exact owned native resources absent, not merely acceptance.
type RemovalResult struct {
	Endpoint                   EndpointRef
	Instance                   InstanceRef
	Dispatch                   Dispatch
	Complete, AbsenceConfirmed bool
}

// Diagnostics is optional, observational and exact-ownership scoped. Positive
// byte bound and cursor bind backend-domain/query/endpoint/snapshot; expiry is ErrCursorExpired.
// Records are sanitized helper evidence, never raw config, secrets or Exec output.
// Unsupported diagnostics cannot imply health.
type EndpointDiagnosticRequest struct {
	EndpointRequest
	Cursor     Cursor
	LimitBytes uint32
}
type EndpointDiagnostic struct {
	Code    ErrorCode
	Message string
}
type EndpointDiagnosticPage struct {
	Items     []EndpointDiagnostic
	Next      Cursor
	Truncated bool
}
