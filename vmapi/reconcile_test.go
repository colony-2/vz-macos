package vmapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/netip"
	"testing"

	"github.com/tmc/vz-macos/vmapi"
)

type reconciliationBackend struct {
	orchardBackend
	instance          vmapi.Instance
	createErr, getErr error
	creates, gets     int
}

func (b *reconciliationBackend) Create(context.Context, vmapi.CreateRequest) (vmapi.CreateResult, error) {
	b.creates++
	return vmapi.CreateResult{Instance: b.instance, Dispatch: vmapi.DispatchAccepted}, b.createErr
}
func (b *reconciliationBackend) Get(context.Context, vmapi.GetRequest) (vmapi.Instance, error) {
	b.gets++
	return b.instance, b.getErr
}

func TestReconcileIntent(t *testing.T) {
	for _, path := range []struct {
		name     string
		input    error
		wantGets int
	}{
		{"accepted", nil, 0}, {"uncertain", vmapi.ErrUnknownOutcome, 1},
	} {
		t.Run(path.name, func(t *testing.T) {
			for _, tt := range []struct {
				name  string
				input func(*vmapi.Instance)
				want  error
			}{
				{"normalized source defaults", func(*vmapi.Instance) {}, nil},
				{"name", func(i *vmapi.Instance) { i.Ref.Name = "other" }, vmapi.ErrConflict},
				{"missing uid", func(i *vmapi.Instance) { i.Ref.UID = "" }, vmapi.ErrConflict},
				{"allocation", func(i *vmapi.Instance) { i.EffectiveIntent.Allocation = "other" }, vmapi.ErrConflict},
				{"owner", func(i *vmapi.Instance) { i.EffectiveIntent.Owner = "other" }, vmapi.ErrConflict},
				{"source locator", func(i *vmapi.Instance) { i.EffectiveIntent.Source.Locator = "other" }, vmapi.ErrConflict},
				{"source mode", func(i *vmapi.Instance) { i.EffectiveIntent.Source.Mode = vmapi.CloneSuspended }, vmapi.ErrConflict},
				{"resolved identity", func(i *vmapi.Instance) { i.EffectiveIntent.Resolved.Identity = "other" }, vmapi.ErrConflict},
				{"native locator", func(i *vmapi.Instance) { i.EffectiveIntent.Resolved.NativeLocator = "other" }, vmapi.ErrConflict},
				{"parent", func(i *vmapi.Instance) { i.EffectiveIntent.Resolved.Parent = "other" }, vmapi.ErrConflict},
				{"trust", func(i *vmapi.Instance) { i.EffectiveIntent.Resolved.TrustIdentity = "other" }, vmapi.ErrConflict},
				{"clone mode", func(i *vmapi.Instance) { i.EffectiveIntent.Mode = vmapi.CloneSuspended }, vmapi.ErrConflict},
				{"missing hardware", func(i *vmapi.Instance) { i.EffectiveIntent.Hardware = nil }, vmapi.ErrConflict},
				{"cpus", func(i *vmapi.Instance) { i.EffectiveIntent.Hardware.CPUs++ }, vmapi.ErrConflict},
				{"memory", func(i *vmapi.Instance) { i.EffectiveIntent.Hardware.MemoryBytes++ }, vmapi.ErrConflict},
				{"disk", func(i *vmapi.Instance) { i.EffectiveIntent.Hardware.DiskBytes++ }, vmapi.ErrConflict},
				{"fingerprint", func(i *vmapi.Instance) { i.EffectiveIntent.Fingerprint = "other" }, vmapi.ErrConflict},
			} {
				t.Run(tt.name, func(t *testing.T) {
					journal := journalFor()
					expected := journalFor().Expected
					instance := vmapi.Instance{Ref: vmapi.InstanceRef{Name: journal.Request.Name, UID: "uid"}, EffectiveIntent: expected}
					tt.input(&instance)
					backend := &reconciliationBackend{instance: instance, createErr: path.input}
					got, err := reconcile(context.Background(), backendRegistry{"selected": {backend, "store-v1"}}, &journal)
					if !errors.Is(err, tt.want) {
						t.Fatalf("got %v want %v", err, tt.want)
					}
					if err == nil && got.Ref.UID != "uid" {
						t.Fatal("lost identity")
					}
					if backend.creates != 1 || backend.gets != path.wantGets {
						t.Fatalf("creates %d gets %d", backend.creates, backend.gets)
					}
				})
			}
		})
	}
}

func TestReconcileFailures(t *testing.T) {
	for _, tt := range []struct {
		name          string
		input         func(*launchJournal, *reconciliationBackend)
		want          error
		creates, gets int
	}{
		{"missing backend key", func(j *launchJournal, b *reconciliationBackend) { j.BackendKey = "" }, vmapi.ErrConflict, 0, 0},
		{"missing revision", func(j *launchJournal, b *reconciliationBackend) { j.Revision = "" }, vmapi.ErrConflict, 0, 0},
		{"missing request correlation", func(j *launchJournal, b *reconciliationBackend) { j.Request.RequestID = "" }, vmapi.ErrInvalidArgument, 0, 0},
		{"empty matching fingerprints", func(j *launchJournal, b *reconciliationBackend) {
			j.Request.IntentFingerprint = ""
			j.Expected.Fingerprint = ""
		}, vmapi.ErrInvalidArgument, 0, 0},
		{"journal correlation mismatch", func(j *launchJournal, b *reconciliationBackend) { j.Expected.Fingerprint = "other" }, vmapi.ErrInvalidArgument, 0, 0},
		{"definite rejection", func(j *launchJournal, b *reconciliationBackend) { b.createErr = vmapi.ErrUnsupported }, vmapi.ErrUnsupported, 1, 0},
		{"not found after dispatch", func(j *launchJournal, b *reconciliationBackend) { b.getErr = vmapi.ErrNotFound }, vmapi.ErrNotFound, 1, 1},
		{"unavailable after dispatch", func(j *launchJournal, b *reconciliationBackend) { b.getErr = vmapi.ErrUnavailable }, vmapi.ErrUnavailable, 1, 1},
		{"cancelled after possible dispatch", func(j *launchJournal, b *reconciliationBackend) {
			b.createErr = errors.Join(vmapi.ErrUnknownOutcome, context.Canceled)
			b.getErr = context.Canceled
		}, context.Canceled, 1, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			journal := journalFor()
			backend := &reconciliationBackend{createErr: vmapi.ErrUnknownOutcome}
			tt.input(&journal, backend)
			_, err := reconcile(context.Background(), backendRegistry{"selected": {backend, "store-v1"}}, &journal)
			if !errors.Is(err, tt.want) {
				t.Fatalf("got %v want %v", err, tt.want)
			}
			if backend.creates != tt.creates || backend.gets != tt.gets {
				t.Fatalf("creates %d gets %d", backend.creates, backend.gets)
			}
		})
	}
}

func TestRestartBinding(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
		want  error
	}{
		{"unchanged", "store-v1", nil}, {"rebound", "store-v2", vmapi.ErrConflict}, {"missing", "", vmapi.ErrConflict},
	} {
		t.Run(tt.name, func(t *testing.T) {
			j := journalFor()
			instance := vmapi.Instance{Ref: vmapi.InstanceRef{Name: j.Request.Name, UID: "colliding"}, EffectiveIntent: j.Expected}
			original := &reconciliationBackend{instance: instance, createErr: vmapi.ErrUnknownOutcome, getErr: vmapi.ErrUnavailable}
			_, err := reconcile(context.Background(), backendRegistry{"selected": {original, "store-v1"}}, &j)
			if !errors.Is(err, vmapi.ErrUnavailable) || !j.Uncertain {
				t.Fatalf("dispatch: %v %+v", err, j)
			}
			// Serialize the journal as a durable consumer would, then reopen a new object.
			data, err := json.Marshal(j)
			if err != nil {
				t.Fatal(err)
			}
			var reopened launchJournal
			if err := json.Unmarshal(data, &reopened); err != nil {
				t.Fatal(err)
			}
			selected := &reconciliationBackend{instance: instance}
			other := &reconciliationBackend{instance: instance}
			registry := backendRegistry{"other": {other, "store-v1"}}
			if tt.input != "" {
				registry["selected"] = backendBinding{selected, tt.input}
			}
			got, err := reconcile(context.Background(), registry, &reopened)
			if !errors.Is(err, tt.want) {
				t.Fatalf("got %v want %v", err, tt.want)
			}
			wantGets := 0
			if tt.want == nil {
				wantGets = 1
				if got.Ref.UID != "colliding" {
					t.Fatal("lost uid")
				}
			}
			if selected.creates != 0 || selected.gets != wantGets || other.creates != 0 || other.gets != 0 {
				t.Fatal("wrong dispatch or replay")
			}
			if reopened.Endpoint != j.Endpoint || reopened.Expected.Source != j.Expected.Source {
				t.Fatal("lost reservations")
			}
			if original.creates != 1 || original.gets != 1 {
				t.Fatal("initial dispatch changed")
			}
		})
	}
}

func TestKnownUIDConflict(t *testing.T) {
	j := journalFor()
	j.Uncertain = true
	j.UID = "original"
	b := &reconciliationBackend{instance: vmapi.Instance{Ref: vmapi.InstanceRef{Name: j.Request.Name, UID: "replacement"}, EffectiveIntent: j.Expected}}
	_, err := reconcile(context.Background(), backendRegistry{"selected": {b, "store-v1"}}, &j)
	if !errors.Is(err, vmapi.ErrConflict) || b.creates != 0 || b.gets != 1 {
		t.Fatalf("replacement accepted: %v", err)
	}
}

type domainStore struct {
	source  vmapi.ResolvedSource
	pending bool
}

type domainBackend struct {
	orchardBackend
	domain                                string
	store                                 *domainStore
	lists, sources, recoveries, endpoints int
}

func (b *domainBackend) List(_ context.Context, r vmapi.ListRequest) (vmapi.InstancePage, error) {
	b.lists++
	if r.Page.Cursor != "" && r.Page.Cursor != vmapi.Cursor(b.domain+"/snapshot") {
		return vmapi.InstancePage{}, vmapi.ErrCursorExpired
	}
	return vmapi.InstancePage{Items: []vmapi.Instance{{Ref: vmapi.InstanceRef{Name: "worker", UID: "colliding"}, Management: vmapi.ManagementForeign, ObservedPower: vmapi.PowerStopped}}, Snapshot: vmapi.SnapshotID(b.domain), Complete: true}, nil
}
func (b *domainBackend) AdoptSource(_ context.Context, r vmapi.AdoptSourceRequest) (vmapi.Source, error) {
	b.sources++
	return vmapi.Source{Selector: r.Selector, Resolved: b.store.source}, nil
}
func (b *domainBackend) Recover(context.Context) (vmapi.RecoveryResult, error) {
	b.recoveries++
	result := vmapi.RecoveryResult{Complete: true}
	if b.store.pending {
		result.Records = []vmapi.RecoveryRecord{{Instance: vmapi.InstanceRef{Name: "worker", UID: "colliding"}, IntentFingerprint: "journal", State: vmapi.RecoveryRecovered}}
		b.store.pending = false
	}
	return result, nil
}
func (b *domainBackend) GetEndpoint(_ context.Context, r vmapi.EndpointRequest) (vmapi.Endpoint, error) {
	b.endpoints++
	port := uint16(2201)
	if b.domain == "b" {
		port = 2202
	}
	return vmapi.Endpoint{Ref: r.Endpoint, Instance: r.Instance, Bind: netip.AddrPortFrom(netip.MustParseAddr("127.0.0.1"), port), GuestPort: 2222}, nil
}

func TestDomainIsolation(t *testing.T) {
	ctx := context.Background()
	aStore := &domainStore{source: vmapi.ResolvedSource{Identity: "source-a", TrustIdentity: "trust-a"}, pending: true}
	bStore := &domainStore{source: vmapi.ResolvedSource{Identity: "source-b", TrustIdentity: "trust-b"}, pending: true}
	a := &domainBackend{domain: "a", store: aStore}
	b := &domainBackend{domain: "b", store: bStore}
	registry := backendRegistry{"selected": {a, "store-v1"}, "other": {b, "store-v1"}}
	j := journalFor()
	selected := registry[j.BackendKey].Backend
	source, err := selected.AdoptSource(ctx, vmapi.AdoptSourceRequest{Selector: j.Request.Source})
	if err != nil || source.Resolved != aStore.source || b.sources != 0 {
		t.Fatal("source crossed domains")
	}
	otherSource, err := b.AdoptSource(ctx, vmapi.AdoptSourceRequest{Selector: j.Request.Source})
	if err != nil || otherSource.Resolved == source.Resolved {
		t.Fatal("source locator treated as global")
	}
	page, err := selected.List(ctx, vmapi.ListRequest{Page: vmapi.PageRequest{Limit: 1}})
	if err != nil || !page.Complete || len(page.Items) != 1 || page.Items[0].Management != vmapi.ManagementForeign {
		t.Fatal("inventory lost foreign stopped object")
	}
	foreign, err := b.List(ctx, vmapi.ListRequest{Page: vmapi.PageRequest{Limit: 1, Cursor: "a/snapshot"}})
	if !errors.Is(err, vmapi.ErrCursorExpired) || foreign.Next != "" || foreign.Complete || len(foreign.Items) != 0 {
		t.Fatal("foreign cursor continued")
	}
	local, err := a.List(ctx, vmapi.ListRequest{Page: vmapi.PageRequest{Limit: 1, Cursor: "a/snapshot"}})
	if err != nil || !local.Complete {
		t.Fatal("local cursor rejected")
	}
	endpoint, err := selected.GetEndpoint(ctx, vmapi.EndpointRequest{Instance: page.Items[0].Ref, Endpoint: j.Endpoint})
	if err != nil || endpoint.Ref != j.Endpoint || endpoint.Bind.Port() != 2201 || b.endpoints != 0 {
		t.Fatal("endpoint reservation crossed domains")
	}
	if err := a.Close(); err != nil {
		t.Fatal(err)
	}
	reopened := &domainBackend{domain: "a", store: aStore}
	registry[j.BackendKey] = backendBinding{reopened, "store-v1"}
	recovered, err := registry[j.BackendKey].Backend.Recover(ctx)
	if err != nil || len(recovered.Records) != 1 || aStore.pending || !bStore.pending || b.recoveries != 0 {
		t.Fatal("recovery crossed durable domains")
	}
	again, err := reopened.Recover(ctx)
	if err != nil || len(again.Records) != 0 {
		t.Fatal("completed record resumed")
	}
}
