package vmapi

import "context"

// Backend is the concurrent, configured-domain VM contract. Results are detached
// snapshots. Optional methods remain present and report ErrUnsupported without
// dispatch. Close is idempotent: cancel and join workers, release resources, but
// do not stop independent VMs or discard durable bootstrap evidence.
type Backend interface {
	Info(context.Context) (BackendInfo, error)
	List(context.Context, ListRequest) (InstancePage, error)
	Get(context.Context, GetRequest) (Instance, error)
	Create(context.Context, CreateRequest) (CreateResult, error)
	Suspend(context.Context, InstanceMutation) (MutationResult, error)
	Delete(context.Context, InstanceMutation) (MutationResult, error)
	Events(context.Context, EventRequest) (EventPage, error)
	Close() error
	Stop(context.Context, StopRequest) (MutationResult, error)
	Exec(context.Context, ExecRequest) (ExecResult, error)
	EnsureEndpoint(context.Context, EnsureEndpointRequest) (Endpoint, error)
	GetEndpoint(context.Context, EndpointRequest) (Endpoint, error)
	AdoptSource(context.Context, AdoptSourceRequest) (Source, error)
	RetainSource(context.Context, RetainSourceRequest) error
	Recover(context.Context) (RecoveryResult, error)
	RemoveEndpoint(context.Context, RemoveEndpointRequest) (RemovalResult, error)
	EndpointDiagnostics(context.Context, EndpointDiagnosticRequest) (EndpointDiagnosticPage, error)
}
