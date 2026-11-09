package vmapi_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/tmc/vz-macos/vmapi"
)

type tartBackend struct{}

var _ vmapi.Backend = tartBackend{}

func (tartBackend) Info(context.Context) (vmapi.BackendInfo, error) {
	return vmapi.BackendInfo{ID: "tartBackend", Capabilities: vmapi.Capabilities{RestartDelivery: true, LifecycleEndpointCleanup: true}}, nil
}
func (tartBackend) List(context.Context, vmapi.ListRequest) (vmapi.InstancePage, error) {
	return vmapi.InstancePage{Complete: true}, nil
}
func (tartBackend) Get(_ context.Context, r vmapi.GetRequest) (vmapi.Instance, error) {
	return vmapi.Instance{Ref: vmapi.InstanceRef{Name: r.Name, UID: "local-uid"}, EffectiveIntent: journalFor().Expected}, nil
}
func (tartBackend) Create(_ context.Context, r vmapi.CreateRequest) (vmapi.CreateResult, error) {
	return vmapi.CreateResult{Dispatch: vmapi.DispatchUnknown}, &vmapi.Error{Code: vmapi.CodeUnknownOutcome, RequestID: r.RequestID}
}
func (tartBackend) Suspend(_ context.Context, r vmapi.InstanceMutation) (vmapi.MutationResult, error) {
	return vmapi.MutationResult{Instance: vmapi.Instance{Ref: r.Instance}, Dispatch: vmapi.DispatchAccepted}, nil
}
func (tartBackend) Delete(_ context.Context, r vmapi.InstanceMutation) (vmapi.MutationResult, error) {
	return vmapi.MutationResult{Instance: vmapi.Instance{Ref: r.Instance}, Dispatch: vmapi.DispatchAccepted, Complete: true, AbsenceConfirmed: true}, nil
}
func (tartBackend) Events(_ context.Context, r vmapi.EventRequest) (vmapi.EventPage, error) {
	return vmapi.EventPage{Items: []vmapi.BootstrapEvent{{Instance: r.Instance, Kind: vmapi.EventTemplateReady, Template: &vmapi.TemplateEvidence{MAC: "02:00:00:00:00:01", BootIdentity: "boot"}}, {Instance: r.Instance, Kind: vmapi.EventDeliverySucceeded}}, CaughtUp: true}, nil
}
func (tartBackend) Close() error { return nil }
func (tartBackend) Stop(context.Context, vmapi.StopRequest) (vmapi.MutationResult, error) {
	return vmapi.MutationResult{}, vmapi.ErrUnsupported
}
func (tartBackend) Exec(context.Context, vmapi.ExecRequest) (vmapi.ExecResult, error) {
	return vmapi.ExecResult{}, vmapi.ErrUnsupported
}
func (tartBackend) EnsureEndpoint(context.Context, vmapi.EnsureEndpointRequest) (vmapi.Endpoint, error) {
	return vmapi.Endpoint{}, vmapi.ErrUnsupported
}
func (tartBackend) GetEndpoint(context.Context, vmapi.EndpointRequest) (vmapi.Endpoint, error) {
	return vmapi.Endpoint{}, vmapi.ErrUnsupported
}
func (tartBackend) AdoptSource(context.Context, vmapi.AdoptSourceRequest) (vmapi.Source, error) {
	return vmapi.Source{}, vmapi.ErrUnsupported
}
func (tartBackend) RemoveEndpoint(context.Context, vmapi.RemoveEndpointRequest) (vmapi.RemovalResult, error) {
	return vmapi.RemovalResult{}, vmapi.ErrUnsupported
}
func (tartBackend) EndpointDiagnostics(context.Context, vmapi.EndpointDiagnosticRequest) (vmapi.EndpointDiagnosticPage, error) {
	return vmapi.EndpointDiagnosticPage{}, vmapi.ErrUnsupported
}
func (tartBackend) RetainSource(context.Context, vmapi.RetainSourceRequest) error {
	return vmapi.ErrUnsupported
}
func (tartBackend) Recover(context.Context) (vmapi.RecoveryResult, error) {
	return vmapi.RecoveryResult{}, vmapi.ErrUnsupported
}

type orchardBackend struct{}

var _ vmapi.Backend = orchardBackend{}

func (orchardBackend) Info(context.Context) (vmapi.BackendInfo, error) {
	return vmapi.BackendInfo{ID: "orchardBackend", Capabilities: vmapi.Capabilities{RestartDelivery: true, LifecycleEndpointCleanup: true}}, nil
}
func (orchardBackend) List(context.Context, vmapi.ListRequest) (vmapi.InstancePage, error) {
	return vmapi.InstancePage{Complete: true}, nil
}
func (orchardBackend) Get(_ context.Context, r vmapi.GetRequest) (vmapi.Instance, error) {
	return vmapi.Instance{Ref: vmapi.InstanceRef{Name: r.Name, UID: "remote-uid"}, EffectiveIntent: journalFor().Expected}, nil
}
func (orchardBackend) Create(_ context.Context, r vmapi.CreateRequest) (vmapi.CreateResult, error) {
	return vmapi.CreateResult{Dispatch: vmapi.DispatchUnknown}, &vmapi.Error{Code: vmapi.CodeUnknownOutcome, RequestID: r.RequestID}
}
func (orchardBackend) Suspend(_ context.Context, r vmapi.InstanceMutation) (vmapi.MutationResult, error) {
	return vmapi.MutationResult{Instance: vmapi.Instance{Ref: r.Instance}, Dispatch: vmapi.DispatchAccepted}, nil
}
func (orchardBackend) Delete(_ context.Context, r vmapi.InstanceMutation) (vmapi.MutationResult, error) {
	return vmapi.MutationResult{Instance: vmapi.Instance{Ref: r.Instance}, Dispatch: vmapi.DispatchAccepted, Complete: true, AbsenceConfirmed: true}, nil
}
func (orchardBackend) Events(_ context.Context, r vmapi.EventRequest) (vmapi.EventPage, error) {
	return vmapi.EventPage{Items: []vmapi.BootstrapEvent{{Instance: r.Instance, Kind: vmapi.EventTemplateReady, Template: &vmapi.TemplateEvidence{MAC: "02:00:00:00:00:01", BootIdentity: "boot"}}, {Instance: r.Instance, Kind: vmapi.EventDeliverySucceeded}}, CaughtUp: true}, nil
}
func (orchardBackend) Close() error { return nil }
func (orchardBackend) Stop(context.Context, vmapi.StopRequest) (vmapi.MutationResult, error) {
	return vmapi.MutationResult{}, vmapi.ErrUnsupported
}
func (orchardBackend) Exec(context.Context, vmapi.ExecRequest) (vmapi.ExecResult, error) {
	return vmapi.ExecResult{}, vmapi.ErrUnsupported
}
func (orchardBackend) EnsureEndpoint(context.Context, vmapi.EnsureEndpointRequest) (vmapi.Endpoint, error) {
	return vmapi.Endpoint{}, vmapi.ErrUnsupported
}
func (orchardBackend) GetEndpoint(context.Context, vmapi.EndpointRequest) (vmapi.Endpoint, error) {
	return vmapi.Endpoint{}, vmapi.ErrUnsupported
}
func (orchardBackend) AdoptSource(context.Context, vmapi.AdoptSourceRequest) (vmapi.Source, error) {
	return vmapi.Source{}, vmapi.ErrUnsupported
}
func (orchardBackend) RemoveEndpoint(context.Context, vmapi.RemoveEndpointRequest) (vmapi.RemovalResult, error) {
	return vmapi.RemovalResult{}, vmapi.ErrUnsupported
}
func (orchardBackend) EndpointDiagnostics(context.Context, vmapi.EndpointDiagnosticRequest) (vmapi.EndpointDiagnosticPage, error) {
	return vmapi.EndpointDiagnosticPage{}, vmapi.ErrUnsupported
}
func (orchardBackend) RetainSource(context.Context, vmapi.RetainSourceRequest) error {
	return vmapi.ErrUnsupported
}
func (orchardBackend) Recover(context.Context) (vmapi.RecoveryResult, error) {
	return vmapi.RecoveryResult{}, vmapi.ErrUnsupported
}

type vzBackend struct{}

var _ vmapi.Backend = vzBackend{}

func (vzBackend) Info(context.Context) (vmapi.BackendInfo, error) {
	return vmapi.BackendInfo{ID: "vzBackend", Capabilities: vmapi.Capabilities{RestartDelivery: true, LifecycleEndpointCleanup: true}}, nil
}
func (vzBackend) List(context.Context, vmapi.ListRequest) (vmapi.InstancePage, error) {
	return vmapi.InstancePage{Complete: true}, nil
}
func (vzBackend) Get(_ context.Context, r vmapi.GetRequest) (vmapi.Instance, error) {
	return vmapi.Instance{Ref: vmapi.InstanceRef{Name: r.Name, UID: "mac-uid"}, EffectiveIntent: journalFor().Expected}, nil
}
func (vzBackend) Create(_ context.Context, r vmapi.CreateRequest) (vmapi.CreateResult, error) {
	return vmapi.CreateResult{Dispatch: vmapi.DispatchUnknown}, &vmapi.Error{Code: vmapi.CodeUnknownOutcome, RequestID: r.RequestID}
}
func (vzBackend) Suspend(_ context.Context, r vmapi.InstanceMutation) (vmapi.MutationResult, error) {
	return vmapi.MutationResult{Instance: vmapi.Instance{Ref: r.Instance}, Dispatch: vmapi.DispatchAccepted}, nil
}
func (vzBackend) Delete(_ context.Context, r vmapi.InstanceMutation) (vmapi.MutationResult, error) {
	return vmapi.MutationResult{Instance: vmapi.Instance{Ref: r.Instance}, Dispatch: vmapi.DispatchAccepted, Complete: true, AbsenceConfirmed: true}, nil
}
func (vzBackend) Events(_ context.Context, r vmapi.EventRequest) (vmapi.EventPage, error) {
	return vmapi.EventPage{Items: []vmapi.BootstrapEvent{{Instance: r.Instance, Kind: vmapi.EventTemplateReady, Template: &vmapi.TemplateEvidence{MAC: "02:00:00:00:00:01", BootIdentity: "boot"}}, {Instance: r.Instance, Kind: vmapi.EventDeliverySucceeded}}, CaughtUp: true}, nil
}
func (vzBackend) Close() error { return nil }
func (vzBackend) Stop(context.Context, vmapi.StopRequest) (vmapi.MutationResult, error) {
	return vmapi.MutationResult{}, vmapi.ErrUnsupported
}
func (vzBackend) Exec(context.Context, vmapi.ExecRequest) (vmapi.ExecResult, error) {
	return vmapi.ExecResult{}, vmapi.ErrUnsupported
}
func (vzBackend) EnsureEndpoint(context.Context, vmapi.EnsureEndpointRequest) (vmapi.Endpoint, error) {
	return vmapi.Endpoint{}, vmapi.ErrUnsupported
}
func (vzBackend) GetEndpoint(context.Context, vmapi.EndpointRequest) (vmapi.Endpoint, error) {
	return vmapi.Endpoint{}, vmapi.ErrUnsupported
}
func (vzBackend) AdoptSource(context.Context, vmapi.AdoptSourceRequest) (vmapi.Source, error) {
	return vmapi.Source{}, vmapi.ErrUnsupported
}
func (vzBackend) RemoveEndpoint(context.Context, vmapi.RemoveEndpointRequest) (vmapi.RemovalResult, error) {
	return vmapi.RemovalResult{}, vmapi.ErrUnsupported
}
func (vzBackend) EndpointDiagnostics(context.Context, vmapi.EndpointDiagnosticRequest) (vmapi.EndpointDiagnosticPage, error) {
	return vmapi.EndpointDiagnosticPage{}, vmapi.ErrUnsupported
}
func (vzBackend) RetainSource(context.Context, vmapi.RetainSourceRequest) error {
	return vmapi.ErrUnsupported
}
func (vzBackend) Recover(context.Context) (vmapi.RecoveryResult, error) {
	return vmapi.RecoveryResult{}, vmapi.ErrUnsupported
}

func selectConfiguration(c vmapi.Capabilities, docker, restart, cleanup bool) error {
	if docker && !(c.Exec && c.Stop && c.TCPGuestEndpoints) {
		return vmapi.ErrUnsupported
	}
	if restart && !(c.Recover || c.RestartDelivery) {
		return vmapi.ErrUnsupported
	}
	if cleanup && !(c.RemoveEndpoint || c.LifecycleEndpointCleanup) {
		return vmapi.ErrUnsupported
	}
	return nil
}

type backendBinding struct {
	Backend  vmapi.Backend
	Revision string
}

type backendRegistry map[string]backendBinding

type launchJournal struct {
	BackendKey, Revision string
	Uncertain            bool
	UID                  vmapi.InstanceUID
	Endpoint             vmapi.EndpointRef

	Request  vmapi.CreateRequest
	Expected vmapi.EffectiveIntent
}

func journalFor() launchJournal {
	source := vmapi.SourceSelector{Locator: "base", Mode: vmapi.CloneCold}
	return launchJournal{
		BackendKey: "selected", Revision: "store-v1",
		Endpoint: vmapi.EndpointRef{UID: "endpoint", Allocation: "allocation"},
		Request:  vmapi.CreateRequest{Mutation: vmapi.Mutation{RequestID: "request"}, Name: "worker", Allocation: "allocation", Owner: "core", Source: source, Mode: vmapi.CloneCold, IntentFingerprint: "journal"},
		// Nil requested hardware inherits these source facts, recorded before dispatch.
		Expected: vmapi.EffectiveIntent{Allocation: "allocation", Owner: "core", Source: source, Resolved: vmapi.ResolvedSource{Identity: "source-uid", NativeLocator: "native-base", Parent: "parent", TrustIdentity: "trust"}, Mode: vmapi.CloneCold, Hardware: &vmapi.Hardware{CPUs: 4, MemoryBytes: 8 << 30, DiskBytes: 32 << 30}, Fingerprint: "journal"},
	}
}

// The fingerprint correlates the journal; it does not establish ownership.
// Expected contains independently recorded, normalized source/hardware facts.
func reconcile(ctx context.Context, registry backendRegistry, journal *launchJournal) (vmapi.Instance, error) {
	binding, ok := registry[journal.BackendKey]
	if !ok || journal.BackendKey == "" || journal.Revision == "" || binding.Revision != journal.Revision || binding.Backend == nil {
		return vmapi.Instance{}, vmapi.ErrConflict
	}
	b := binding.Backend
	r := journal.Request
	if r.RequestID == "" || r.IntentFingerprint == "" || journal.Expected.Fingerprint != r.IntentFingerprint {
		return vmapi.Instance{}, vmapi.ErrInvalidArgument
	}
	var instance vmapi.Instance
	observe := journal.Uncertain
	if !journal.Uncertain {
		// Persist this state before dispatch; only definite rejection can clear it.
		journal.Uncertain = true
		result, err := b.Create(ctx, r)
		if err == nil {
			instance = result.Instance
		} else if errors.Is(err, vmapi.ErrUnknownOutcome) {
			observe = true
		} else {
			journal.Uncertain = false
			return vmapi.Instance{}, err
		}
	}
	if observe {
		var err error
		instance, err = b.Get(ctx, vmapi.GetRequest{Name: r.Name, ExpectedUID: journal.UID})
		if err != nil {
			return vmapi.Instance{}, err
		}
	}
	if instance.Ref.Name != r.Name || instance.Ref.UID == "" || (journal.UID != "" && instance.Ref.UID != journal.UID) || !reflect.DeepEqual(instance.EffectiveIntent, journal.Expected) {
		return vmapi.Instance{}, vmapi.ErrConflict
	}
	journal.UID = instance.Ref.UID
	return instance, nil
}
func TestConsumer(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input vmapi.Backend
	}{{"tart", tartBackend{}}, {"orchard", orchardBackend{}}, {"vz", vzBackend{}}} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			info, err := tt.input.Info(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if err := selectConfiguration(info.Capabilities, false, true, true); err != nil {
				t.Fatal(err)
			}
			journal := journalFor()
			instance, err := reconcile(ctx, backendRegistry{"selected": {tt.input, "store-v1"}}, &journal)
			if err != nil {
				t.Fatal(err)
			}
			events, err := tt.input.Events(ctx, vmapi.EventRequest{Instance: instance.Ref, Limit: 2})
			if err != nil || len(events.Items) != 2 || events.Items[0].Template == nil {
				t.Fatalf("events: %v %v", events, err)
			}
			if _, err := tt.input.EndpointDiagnostics(ctx, vmapi.EndpointDiagnosticRequest{}); !errors.Is(err, vmapi.ErrUnsupported) {
				t.Fatal(err)
			}
			if err := tt.input.Close(); err != nil {
				t.Fatal(err)
			}
			if err := tt.input.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestConfiguration(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input vmapi.Capabilities
		want  bool
	}{
		{"missing restart", vmapi.Capabilities{}, false},
		{"equivalent restart and cleanup", vmapi.Capabilities{RestartDelivery: true, LifecycleEndpointCleanup: true}, true},
		{"missing cleanup", vmapi.Capabilities{Recover: true}, false},
		{"native hooks", vmapi.Capabilities{Recover: true, RemoveEndpoint: true}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectConfiguration(tt.input, false, true, true) == nil; got != tt.want {
				t.Fatalf("got %v want %v", got, tt.want)
			}
		})
	}
	if selectConfiguration(vmapi.Capabilities{}, true, false, false) == nil {
		t.Fatal("accepted missing Docker requirements")
	}
}

func TestSensitiveExecValues(t *testing.T) {
	input := vmapi.ExecRequest{Argv: []string{"sh", "-s"}, Stdin: []byte{0, 10, 255}, OutputLimitBytes: 32}
	exit := 7
	result := vmapi.ExecResult{Delivery: vmapi.DeliveryDelivered, ExitCode: &exit, Stdout: []byte{255}, StderrTruncated: true}
	if string(input.Stdin) != string([]byte{0, 10, 255}) || input.Argv[1] != "-s" || input.OutputLimitBytes != 32 {
		t.Fatal("lost command intent")
	}
	if result.ExitCode == nil || *result.ExitCode != 7 || !result.StderrTruncated {
		t.Fatal("lost exit/output evidence")
	}
	if (vmapi.ExecResult{}).ExitCode != nil {
		t.Fatal("unknown exit became success")
	}
}
