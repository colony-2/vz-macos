# A Go VM contract bound to configured Backend instances

[Analysis](vmapi-analysis.md) separates current code from later native obligations; [coverage](vmapi-coverage.md) maps the audited consumer needs.

`github.com/tmc/vz-macos/vmapi` lets callers submit provisioning and bootstrap intent through one Backend and observe shared results across future Tart, Orchard and VZ adapters. Each configured Backend binds one resource domain. The package contains the contract and typed errors; native implementation follows separately. The [support comparison](vmapi-support.md) records the audited consumer requirements.

## Public entry point

Backend exposes these 17 methods. Optional methods remain present; capabilities and typed errors replace structural interface assertions. No public constructor, production implementation, network client or generic operation service is added.

```go
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
```

Baseline is Info, List, Get, Create, Suspend, Delete, Events and Close. Docker configurations additionally require Stop, Exec and endpoint ensure/observe. Backend configuration owns authentication, local processes, remote calls, guest transport selection, credentials and trusted host-key provisioning. There is no per-command SSH/agent selector or new Orchard service requirement. All methods support concurrent calls, return detached value snapshots, use context cancellation and return errors rather than panic. Close is idempotent, cancels and joins adapter workers as appropriate, and releases adapter resources without stopping independently running VMs. A newly configured adapter owns restart recovery; Close does not discard its durable delivery evidence.

## Changes from the previous candidate

The previous unreleased candidate exposed Placement. This revision removes it without a replacement request-level selector; request wrappers and other shared values remain.

| Declaration | Remove | Retain |
|---|---|---|
| InstanceRef | Placement | Name, UID |
| EndpointRef | Placement | UID, Allocation |
| SourceSelector | Placement | Locator, Mode |
| ListRequest | Placement | Page |
| GetRequest | Placement | Name, ExpectedUID |
| CreateRequest | Placement | Mutation, Name, Allocation, Owner, Source, Mode, Hardware, IntentFingerprint, Bootstrap |
| BackendInfo | Placements | ID, Capabilities, Limits, Capacity, Observation |

Nested instance/endpoint references, event values and effective intent inherit these removals. `BackendInfo.ID` is informational: it cannot authorize routing, prove host affinity or replace a consumer configuration binding.

## Backend selection and restart safety

One instance binds inventory, VM identities, source resolution, capacity, endpoint host, cursors and durable recovery to the same fixed domain. A future local adapter binds its native store/host; an Orchard adapter privately constrains its worker domain. Multiple unconstrained clients to one controller do not establish separate domains. Unsupported affinity must reject configuration rather than choose an arbitrary worker. No new remote service or live affinity claim belongs here.

The test consumer journals its own BackendKey and configuration binding before dispatch, alongside VM/source/endpoint reservations. This key is not a vmapi field. In the fixture, a registry entry couples Backend with a consumer-owned immutable configuration revision; the journal stores both key and expected revision. Missing entries or changed revisions reject before any backend call. Reopening an unchanged configuration may supply a new Go object with the same revision. A revision changes whenever the underlying store/host/worker domain changes; pointer identity or BackendInfo.ID alone cannot establish continuity. Production consumers must enforce equivalent configuration integrity later; this assignment adds no registry library or constructor.

Resolve the journal's key once, validate its revision, and use that Backend for Create and all later observations, mutations, source, event, recovery and endpoint work. Persist dispatch uncertainty. On restart, reconcile existing uncertain work by Get on the same backend, without issuing another Create. Never search other backends by name or move reservations automatically. Reconcile name, known UID and the entire normalized effective intent: allocation, owner, selector locator/mode, resolved identity/native locator/parent/trust, clone mode, all hardware dimensions and fingerprint. A matching fingerprint alone is insufficient. Keep accepted and uncertain paths plus not-found, unavailable, cancellation, missing correlation and mismatch coverage.

Names and UIDs may collide across backends. If the consumer routes a bare reference to the wrong backend containing the same Name/UID, the interface cannot always detect the mistake. Consumer routing must preserve this association. Within a selected domain, immutable UID validation still rejects replacement objects atomically; declaration tests cannot establish native race safety.

List remains complete for its domain, including stopped sources and foreign capacity consumers. Multiple handles may overlap physical capacity; admission must account for that overlap rather than sum each handle's available count. Core reservations include durable unobserved intent.

Source locators resolve locally, with no automatic source transfer. Exact endpoint host binding belongs to the selected domain; reservation and cleanup use the same journal association. Cursors bind originating domain, query, identity and snapshot/stream; foreign cursors must be rejected without establishing continuation or completeness. Expiry remains ErrCursorExpired; no new cursor error code is required. Durable recovery cannot scan or resume another domain, and closing/reopening preserves originating records.

## Shared values and observation

Retain useful opaque scalar IDs, Hardware (logical CPUs and bytes), OS, power, freshness, clone modes and stop modes from the candidate. Each configured Backend binds one stable resource domain. `InstanceRef{Name, UID}` requires a nonempty backend-local immutable UID for every mutation. Names locate observations, never authorize mutations. `GetRequest{Name, ExpectedUID}` allows name-only discovery after an ambiguous Create; a supplied UID mismatch returns conflict. `ListRequest{Page PageRequest}` includes all discoverable managed, stopped/source and foreign capacity consumers in that backend domain.

`Instance` contains Ref, OS, Hardware, Management, Source lineage, EffectiveIntent, DesiredPower, ObservedPower, Observation, Convergence, Delivery and bounded Conditions. `Observation{At, Freshness}` cannot turn unavailable observation into stopped/absent. `Convergence{DesiredGeneration, ObservedGeneration, Known, Transition}` uses zero generations as unspecified when Known is false; Transition includes starting/suspending/stopping/none/unknown. Generation is observation evidence, not a generic mutation token. Suspend is complete only when observed nonrunning, desired suspended, convergence known and reached, with no active suspending condition. Delete completion confirms absence of the exact UID and owned resources, never deletion of a same-name replacement.

`BackendInfo` reports informational backend identity, explicit Capabilities, positive retrieval/input limits, retention bounds and observed capacity with a Complete flag. Unknown support is false. Capabilities cover supported OS/clone modes; Docker Exec/Stop/endpoints; source probing; adoption/retention; restart delivery/recovery; endpoint removal/diagnostics and lifecycle-coupled endpoint cleanup. Listing need not discover worker-local source names when `ProbeLocalSources` is true; this does not excuse missing foreign capacity consumers. Capacity observations remain advisory: core admission includes durable unobserved intents.

Use the existing bounded stable-snapshot pagination model: positive limits, opaque query/backend-domain/snapshot-scoped cursors, explicit completeness and safe bounded Issues, empty Next for enumeration end. Incomplete pages and expired cursors cannot prove absence or free capacity. Get/List/Events/diagnostics never create directories, repair resources or start delivery.

## Provisioning, sources and effective intent

`CreateRequest` contains caller RequestID, backend-local unique Name, allocation/owner identity, SourceSelector, explicit CloneMode, optional Hardware overrides, caller IntentFingerprint, and sensitive BootstrapIntent. It requests provisioning plus launch/resume plus bootstrap, without a readiness wait. `BootstrapIntent` carries recipe identity, exact stdin/script bytes and an absolute deadline, subject to advertised size bounds. It is privately persisted where delivery recovery requires it; secrets are absent from argv, ordinary diagnostics and public observations. The adapter selects SSH/agent and verifies inherited host-key trust.

`EffectiveIntent` returns normalized nonsecret fields: owner/allocation, source selector and resolved lineage, clone mode, hardware, and caller fingerprint. It excludes BootstrapIntent bytes, credentials, leases and raw backend objects. Backend normalization must be deterministic; the fingerprint is correlation evidence, not authentication or permission to mutate. The caller compares effective intent with its journal before accepting ownership.

`SourceSelector{Locator, Mode}` is a replaceable published name or opaque native locator, without credentials. `Source` describes the selector, OS, hardware, observation and `ResolvedSource{Identity, NativeLocator, Parent, TrustIdentity}`. Resolved identity is the adapter's binding evidence at clone time; it need not imply a globally versioned immutable published source object. TrustIdentity is a nonsecret reference to adapter-managed inherited host-key trust. Existing clones preserve their bound identity, data and trust when a published name changes. Disposable clone identity is strictly bound to its native object.

AdoptSource accepts a selector and expected nonsecret lineage/trust references; secret trust material is configured privately. It registers a reusable existing source, rather than introducing image conversion. RetainSource marks an exact observed source reusable/delivered; it is not a retention or garbage collection policy. Cold cloning permits hardware changes and randomized MAC. Suspended cloning preserves saved memory/MAC and requires compatibility and exclusive-slot/collision checks, including foreign objects. Unsupported or incompatible suspended requests fail without cold fallback. Core reservations and backend native collision checks remain separate responsibilities.

`CreateResult{Instance, Dispatch}` requires a stable UID and effective intent on successful acceptance. `MutationResult{Instance, Dispatch, Complete, AbsenceConfirmed}` separates accepted work from completion. Dispatch is accepted/rejected/unknown, not an operation ID. Suspend/Stop may return accepted while observations converge. Delete requires AbsenceConfirmed and Complete before declaring success of cleanup; acceptance alone is not that evidence. `RemovalResult` similarly identifies the exact endpoint and distinguishes accepted cleanup from confirmed absence.

Persist native intent and delivery evidence before the corresponding irreversible work. A request ID and matching journal support reconciliation, but the public API promises neither idempotent Create replay nor a generalized operation registry. Definitely rejected validation/unsupported requests dispatch nothing. After possible dispatch, return ErrUnknownOutcome with safe request correlation; context cancellation after dispatch must preserve this distinction, optionally joining the context error. The caller observes by unique name and validates UID/effective intent; absence of a retained record cannot authorize blind replay.

## Bootstrap events, Exec and recovery

Events are keyed by exact InstanceRef, not generic operation. `EventRequest{Instance, After, Limit}` and `EventPage{Items, Next, CaughtUp}` provide bounded nonblocking reads, cursor expiry, per-instance order and explicit retention. CaughtUp means no more retained events at observation time, not readiness or permanent stream completion.

`BootstrapEvent` includes cursor, time, instance, intent fingerprint and a typed Kind. Kinds are base-ready, template-ready, docker-engine-ready, delivery-succeeded and delivery-failed. Template-ready carries MAC and boot identity; docker-engine-ready carries engine identity. Delivery-succeeded is equivalent to PULSE_HANDOFF_OK, separate from base/template/application readiness. Invalid or unsupported evidence is not normalized into success. Delivery state accepted/dispatched/delivered/failed/unknown is separate from power and core readiness. No free-form script/stdout payload is a public event field. Adapters later translate exact permitted legacy markers, retain equivalent bounded evidence across restart, and exclude unrelated output. Public errors and Conditions contain only bounded sanitized messages.

`ExecRequest` contains InstanceMutation, Argv, exact Stdin bytes, optional environment/working directory, required absolute deadline and positive aggregate output bound. `ExecResult` contains exact identity, delivery phase, optional exit code, bounded stdout/stderr byte slices and truncation flags. It may be returned after terminal completion or with uncertainty and error; absent exit status is not exit zero. Nonzero guest exit is a result; transport failure is an error and may have unknown delivery. Caller context bounds the call; the requested deadline bounds guest work independently where supported. Arbitrary Exec is never retried automatically and is not resumed by Recover. Output is sensitive and accessible only in the explicit result, never ordinary events/diagnostics. No separate output service or readiness interpretation is needed.

Recover takes only context and uses adapter-owned durable records, not an untrusted name/request replacing those records. `RecoveryResult` reports bounded per-identity recovered/skipped/uncertain records and completeness, with cursor-free overflow reported as incomplete. Only known eligible, undelivered bootstrap with identity, intent, native phase and trust evidence may resume. Never repeat ambiguous clone/run or arbitrary Exec. A crash after delivery before recording success can retry the idempotent bootstrap recipe: delivery is at least once, never exactly once. Observation does not trigger Recover. Core journals still own launch deduplication and unobserved reservations.

## Endpoints

Reuse net/netip for exact TCP addresses. EnsureEndpointRequest carries InstanceMutation, AllocationID, exact worker Bind address/port and guest TCP target port. Endpoint returns immutable EndpointRef, exact instance/owner, actual worker Bind, guest target, observed state and freshness. GetEndpoint and diagnostics require exact endpoint/allocation/instance binding; mismatches conflict. Ensure is convergent only for identical owner, instance, binding and target; occupied or substituted ports conflict. Adapter restoration follows guest address changes or reports uncertainty. Listening is separate from authenticated Docker readiness.

RemoveEndpoint requires exact identities and returns confirmed absence only after native owned resources are absent. EndpointDiagnosticRequest includes exact ownership, cursor and byte bound; its page contains bounded sanitized helper records, continuation, expiry and truncation, without secrets, arbitrary Exec output or raw configuration. Diagnostics remains optional qualification support.

## Optional behavior and fallback

Info is consulted before mutations. Unsupported means no dispatch; partial feature support must not pass configuration validation silently.

| Feature | Explicit caller behavior when unsupported |
|---|---|
| AdoptSource | Fall back to observational Get of the configured source and verify source/trust facts; do not manufacture adoption. |
| RetainSource | Successful no-op is permitted for adapters with no reusable/delivered marking need, as on Orchard. Backend method still reports unsupported. |
| ProbeLocalSources | False selects ordinary discovery. True selects core-managed disposable clone probes; no automatic probe mutation in List/Get. |
| Recover | Accept only if restart delivery is supplied equivalently by the adapter/service, otherwise reject configurations requiring it. |
| RemoveEndpoint | Defer to VM lifecycle cleanup only when explicitly advertised as equivalent; otherwise fail configurations requiring cleanup. |
| EndpointDiagnostics | Skip optional endpoint qualification detail; never substitute success for health. |
| Docker operations | Reject Docker configuration before Create when required support is absent. |

Fallbacks belong to documented consumer selection, not a production helper in this package. Test-only consumer code exercises the decisions. No backend capability is qualified by an interface assertion.

## Errors and migration

Keep the existing inspectable codes/sentinels: invalid argument, not found, conflict, unsupported, unavailable, unknown outcome and cursor expired. Keep errors.Is/errors.As support through wrapping and safe RequestID correlation; do not reintroduce a generic OperationID. Do not retain raw request, credentials or backend causes in public errors. Unknown codes and unrelated targets must not match known sentinels. Identity mismatch is conflict; incomplete observation is unavailable, never not-found. All messages are lowercase, bounded and sanitized by adapters.

The [validation record](vmapi-validation.md) records the bounded publication and consumer search. No released dependency was found in those records. A later discovery requires a deprecation and migration decision before breaking released consumers.

## Evidence boundary

Independent test-only Tart-shaped, Orchard-shaped and VZ-shaped Backends and an isolated consumer module exercise source usability. Native identity races, suspended cloning, durable secret-safe delivery, affinity and endpoint connectivity remain unqualified. See [analysis](vmapi-analysis.md), [coverage](vmapi-coverage.md) and [validation](vmapi-validation.md).
