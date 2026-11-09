# Reviewing the VM Backend operations

A caller first selects a configured Backend, records its stable key and immutable configuration revision with the intended VM in its journal, then asks Backend to provision, run/resume and deliver bootstrap. It observes identity, power, convergence and delivery independently; the provider decides readiness. This guide describes the proposed contract in the current Go declarations. The [support comparison](vmapi-support.md) shows which behavior Pulse already uses through Tart and Orchard and which details are new adapter obligations.

## Shared rules and values

Each configured Backend binds one consistent inventory, source, identity, capacity, endpoint-host and recovery domain. `InstanceRef{Name, UID}` includes a backend-local immutable nonempty native identity. A name locates an observation but never authorizes mutation. Existing-instance mutations, Exec and endpoint work must atomically validate that identity; a same-name replacement conflicts. Create assigns the new UID rather than taking an existing one. The audited adapters' check-then-command/HTTP flows do not establish this atomic guarantee.

Hardware uses logical CPUs and bytes. Nil hardware means unknown in observations and inheritance where Create allows omitted overrides; it is not zero-sized hardware. Cold cloning may change hardware and randomize MAC. Suspended cloning must preserve saved memory/MAC and reject incompatible hardware, saved-state or collision requirements without a cold fallback. Core reservations do not replace adapter collision checks against foreign native objects.

`EffectiveIntent` contains owner/allocation, source selector and resolved binding, clone mode, normalized hardware and caller fingerprint. It contains no bootstrap bytes or credentials. A selector can refer to a replaceable published source; resolved identity/native locator/parent/trust bind this particular clone. Replacing the published source must not change existing clone data or inherited trust.

Desired power records intent; observed power records evidence. A guest can shut down while Orchard's desired state still says running. Freshness and convergence must remain explicit: unknown/stale data cannot become stopped or absent. `Convergence.Known` distinguishes evidence from unspecified zero generations; generations are not general mutation tokens. Conditions are bounded sanitized diagnostics, not readiness policy.

All calls support concurrency and return detached snapshots. Context cancellation bounds a call; it does not prove that dispatched work was cancelled. Reads never create directories, repair resources or restart delivery. Backend configuration owns credentials, controller TLS/authentication, SSH/agent choice and inherited host-key verification. Shared requests do not select a guest transport or carry host authentication.

## Info: learn the backend's scope and supported work

**Call/result.** `Info(ctx)` returns `BackendInfo`: informational ID, explicit capabilities, positive retrieval/input bounds, event retention, observed capacity and observation time/freshness. Check capabilities before any mutation. A successful health probe alone is insufficient to populate these fields honestly.

**Concrete use.** Before preparing a Docker pair, require Exec, Stop and TCPGuestEndpoints. If they are missing, reject the configuration before Create. Require a supported guest OS and clone mode rather than discovering incompatibility after cloning. Capacity is advisory: core still counts unobserved durable intents and foreign consumers; incomplete accounting cannot justify a new slot.

**Current mapping.** Tart's `--version` and Orchard's nonempty controller version supply health evidence only. The expanded result is new common reporting, not an existing production Pulse Info payload.

## List: enumerate one backend domain without changing it

**Call/result.** `List(ctx, ListRequest{Page})` returns a stable-snapshot `InstancePage`. Limit must be positive and within Info's maximum. Cursor binds backend domain/query/snapshot; expiration returns ErrCursorExpired. Empty Next ends enumeration. Complete describes snapshot coverage; a single page is not an entire inventory when Next remains nonempty. Issues explain bounded, sanitized coverage gaps.

**Concrete use.** Enumerate all pages before concluding a name is absent or estimating capacity. Include managed workloads, stopped/prepared sources and foreign capacity consumers. If an observation fails or a cursor expires, restart a complete observation; do not treat the partial list as free capacity. Orchard worker-local source probing remains separate from listing foreign VM consumers.

**Current mapping.** Tart lists local VMs including foreign objects; Orchard lists controller VM records with worker/label filtering in the provider. Stable cursors, explicit completeness and bounds require adapter implementation beyond these current list calls.

## Get: observe a name, optionally requiring an exact UID

**Call/result.** `Get(ctx, GetRequest{Name, ExpectedUID})` returns one Instance. A supplied UID mismatch conflicts; omitting it permits discovery after an uncertain Create. Confirmed absence is not-found; unavailable/incomplete observation is unavailable. Neither response authorizes replay of a possibly dispatched Create.

**Concrete use.** A journal records backend key `worker-a` and immutable configuration revision `v1`, name `job-42` and expected normalized source/owner/hardware. After losing a Create response, Get by name on that journal-selected Backend. Accept ownership only with exact name and journaled backend association, a nonempty UID and matching full effective intent. A matching fingerprint with a different owner/source/hardware still conflicts. The test consumer in `vmapi/contract_test.go` and cases in `vmapi/reconcile_test.go` exercise these comparisons; expected hardware inheritance is recorded independently before dispatch.

**Current mapping.** Tart combines journal/native binding evidence; Orchard returns desired state, conditions, generations and lineage. The common ExpectedUID and freshness semantics are explicit new contract obligations.

## Create: provision, run or resume, and begin bootstrap

**Call/result.** `CreateRequest` supplies RequestID, backend-local unique name, allocation/owner, source selector, explicit clone mode, optional hardware overrides, intent fingerprint and sensitive BootstrapIntent. Bootstrap contains recipe identity, exact stdin bytes and absolute deadline within advertised bounds. Success returns `CreateResult{Instance, Dispatch}` with stable UID and normalized nonsecret effective intent. It does not wait for application readiness.

**Concrete use.** To launch a prepared macOS slot, request suspended cloning only if advertised and compatible. Saved memory/MAC and exclusive collision checks are required. Unsupported suspended cloning rejects before dispatch; it must not silently perform a cold boot. For cold clones, omitted hardware can inherit independently established source facts. The caller compares the returned normalized intent to its journal, not raw pointer identity or the fingerprint alone.

**Failure/reconciliation.** Persist native intent/delivery evidence before irreversible work. Definite validation/unsupported rejection dispatches nothing. After possible dispatch, ErrUnknownOutcome and safe RequestID correlation preserve uncertainty, including cancellation; a context error can be joined. Observe by journaled name and verify full intent rather than call Create again. A missing journal/retained record or temporary not-found is not permission to clone twice. DispatchAccepted means responsibility was accepted, not running, bootstrap delivered or ready.

**Current mapping.** Tart already sequences clone/configure/detached run/asynchronous delivery; Orchard posts a VM with startup intent and execution follows remotely. The new package does not add a generic operation registry. Native implementation and safe recovery records remain later work.

## Suspend: request saved state and observe completion separately

**Call/result.** `Suspend(ctx, InstanceMutation)` carries RequestID and exact InstanceRef. MutationResult separates Dispatch from Complete and the returned observation.

**Concrete use.** After template-ready evidence with validated MAC/boot identity, request suspension of that exact VM. An accepted request is still pending until observed nonrunning, desired power suspended, known reached convergence, and no active suspending transition. A paused process or merely desired suspended state is insufficient evidence of saved-state completion.

**Current mapping.** Production core requires this operation for templates. Tart invokes suspend after UID checks. Orchard preserves the full object when changing desired power, rather than losing unmodeled fields in a partial PUT. Adapter qualification must still establish saved-memory behavior and replacement-race safety.

## Delete: remove an exact VM and its owned resources

**Call/result.** `Delete(ctx, InstanceMutation)` returns MutationResult. Only Complete plus AbsenceConfirmed establishes cleanup of the requested UID and owned native resources. DispatchAccepted alone is not deletion completion.

**Concrete use.** If `job-42` now refers to a different UID, return conflict and leave the replacement alone. If removal is still converging, keep the caller's reservation rather than free capacity. Endpoint lifecycle cleanup can be used as an explicit fallback only when equivalently guaranteed, including owned helpers.

**Current mapping.** Tart checks ownership/stopped state then removes the VM and checks directory absence. Orchard's name-based DELETE followed by absent controller record is not proof of physical cleanup or atomic UID preconditions. The stronger contract requires additional implementation evidence; it must not be inferred from the current method signature.

## Events: read bounded typed bootstrap evidence

**Call/result.** `Events(ctx, EventRequest{Instance, After, Limit})` performs a nonblocking ordered read for an exact instance. Limit is positive/bounded. Expired instance-stream cursors return ErrCursorExpired. Next continues retained evidence; CaughtUp means no further retained events now, not permanent completion or readiness.

**Concrete use.** Base-ready records base preparation. Template-ready carries MAC and boot identity; Docker-engine-ready carries engine identity. Delivery-succeeded corresponds to `PULSE_HANDOFF_OK`, distinct from those readiness stages. Delivery-failed is distinct from a power transition. Invalid template/Docker evidence cannot be accepted as success. The provider decides which typed evidence qualifies its recipe or paired workload.

**Current mapping.** Core currently parses exact marker lines in Message/Payload. Tart's bounded event file also contains startup stdout/stderr, and Orchard exposes string events. Future adapters must translate only permitted markers into typed values and preserve equivalent evidence across restart. Arbitrary output, scripts, credentials and leases must not enter this public stream; existing raw events do not already satisfy that guarantee.

## Close: dispose the adapter without deleting VMs

**Call/result.** `Close()` is idempotent. Cancel/join adapter workers and release local resources while preserving durable bootstrap evidence. It does not request stop/delete of independently running VMs.

**Concrete use.** Shutting down the provider waits for local delivery workers without interpreting shutdown as a VM teardown. A newly configured adapter can later recover eligible delivery from durable records. Deleting those records on Close would break that recovery.

**Current mapping.** Runtime uses Close for both paths. Tart cancels/waits for workers while detached VM execution continues; Orchard closes idle HTTP connections. No new production call is needed.

## Stop: request nonrunning power without deletion

**Call/result.** `Stop(ctx, StopRequest)` includes exact identity and explicit graceful/force mode. MutationResult distinguishes dispatch from completion. A stop does not imply deletion or reusable suspended memory.

**Concrete use.** Docker cleanup may need the guest nonrunning before further cleanup. Keep observing when only acceptance is known; do not treat a changed desired state as observed termination. An adapter must not claim a force/graceful distinction it cannot implement.

**Current mapping.** Tart uses stop; Orchard changes desired stopped power after UID observation. Existing callers supply no stop mode. Mode-specific behavior is an approved extension that requires later implementation/qualification rather than an already-used Pulse parameter.

## Exec: deliver one command without automatic replay

**Call/result.** `ExecRequest` preserves exact argv, stdin, optional environment/cwd, required absolute deadline and positive aggregate output limit. ExecResult contains exact identity, delivery phase, optional exit status, sensitive stdout/stderr and truncation flags. A nonzero guest exit is a result; transport failure is an error, potentially with uncertain delivery. Nil exit status is not exit zero.

**Concrete use.** A command exits 7 after returning bounded stderr: preserve exit 7 rather than report a transport failure. If transport drops after dispatch, expose uncertain delivery; do not run it again or let Recover replay it. Context bounds the call while the absolute deadline bounds guest work independently where supported. Only the explicit result exposes output; events/ordinary diagnostics must not repeat it.

**Current mapping.** Docker core currently sends script bytes over stdin via SSH, directly after Tart IP lookup or over Orchard's authenticated WebSocket tunnel. The new argv/env/cwd, separate outputs/exit status and deadline shape is broader than that audited production call. Transport/trust stay private adapter choices. Exact delivery and bounded result guarantees need later adapter evidence.

## EnsureEndpoint: establish one exact owned TCP listener

**Call/result.** `EnsureEndpointRequest` includes InstanceMutation, allocation, exact worker IP/port Bind and nonzero guest TCP target port. Endpoint returns immutable endpoint/allocation/instance identity, actual Bind, target, state and freshness. Identical owner/instance/binding/target converges; a substituted or occupied port conflicts.

**Concrete use.** The provider reserves worker TCP port 32001 for its Docker pair and requests forwarding to guest TCP 2222. A listener on 32002 is not an acceptable result. Restore forwarding after guest address changes or report uncertainty. EndpointListening proves a listener, not authenticated Docker readiness; core still verifies the engine and pairing.

**Current mapping.** Tart uses a launchd bridge following guest IP; Orchard manages a named endpoint and checks exact requested port against observed listening state. The shared allocation/endpoint UIDs and exact IP bind are more explicit than the existing combined bool-returning call.

## GetEndpoint: observe without ensuring or repairing

**Call/result.** `GetEndpoint(ctx, EndpointRequest)` requires exact endpoint/allocation and instance identity. Ownership mismatches conflict. The returned Endpoint is observation evidence only.

**Concrete use.** Poll the endpoint returned by EnsureEndpoint while it is pending. GetEndpoint must not start a helper, change a binding or silently recreate a missing endpoint. Choose a separate ensure/recovery action when repair is needed.

**Current mapping.** Neither audited adapter has a standalone GetEndpoint call; both combine observation with possible creation/reconfiguration in Endpoint. The split is new API structure for already-used observation, and a future adapter must implement genuinely read-only behavior.

## AdoptSource: register an existing reusable source

**Call/result.** `AdoptSourceRequest` carries a selector and expected nonsecret resolved lineage/trust references plus RequestID. The result describes Source, hardware/OS, observation and resolved binding. Secret host-key material is configured privately, not copied into the observation.

**Concrete use.** Register a configured stopped/suspended source with inherited trust. A published selector may intentionally be replaced later, while clones remain bound to their original native identity/data/trust. This operation does not download/convert an image or take ownership of arbitrary foreign data.

**Current mapping/fallback.** Tart has an adoption hook accepting configured trust. Orchard lacks it; production core falls back to Get. With the new API, unsupported means no dispatch and the caller explicitly verifies the observed source/trust facts before using that fallback.

## RetainSource: mark reuse or successful delivery

**Call/result.** `RetainSourceRequest` identifies an exact observed Source and reusable/delivered flags; success returns no separate operation handle. Validate that native binding rather than marking a same-name replacement.

**Concrete use.** After preparing a source, mark it reusable/delivered so local delivery recovery and published-source handling retain the intended state. This does not request “keep three versions,” delete old images or schedule garbage collection.

**Current mapping/fallback.** Tart persists these flags. Orchard lacks the hook and current core selects a successful no-op. The proposed backend still reports unsupported; core explicitly permits the no-op only when that backend has no equivalent marking requirement.

## Recover: resume eligible bootstrap from durable adapter evidence

**Call/result.** `Recover(ctx)` takes no caller-supplied name or replacement record. Adapter-owned durable evidence drives bounded per-instance recovered/skipped/uncertain records. Overflow makes Complete false; no cursor is implied.

**Concrete use.** After process restart, resume only known started, undelivered, eligible bootstrap with exact identity, intent, native phase and trust evidence and no existing worker. Creating/prepared/starting ambiguity cannot authorize another clone/run. Never replay arbitrary Exec. A crash after guest delivery but before recording success can repeat the idempotent recipe: this is at-least-once delivery, not exactly once.

**Current mapping/fallback.** Runtime calls Tart recovery before core reconciliation. Tart persists spec/UID/phase/process/native binding/trust and delivery/reuse flags; core journals remain necessary too. Orchard has no Recover hook. If restart delivery is required, accept an absent hook only with separately established equivalent remote-service delivery; otherwise reject that configuration. This audit did not establish controller restart guarantees.

## RemoveEndpoint: confirm owned listener cleanup

**Call/result.** `RemoveEndpointRequest` carries RequestID, exact instance and endpoint/allocation identity. RemovalResult distinguishes dispatch from Complete and AbsenceConfirmed for those exact owned resources.

**Concrete use.** Stop/remove the owned listener/helper, then confirm absence before releasing its port reservation. Do not remove a different allocation's listener now using that port. An accepted asynchronous cleanup is not confirmed absence.

**Current mapping/fallback.** Production Docker cleanup uses Tart's explicit hook. Orchard has no hook and relies on VM lifecycle cleanup. Allow that fallback only when equivalent endpoint cleanup is explicitly advertised; detached Tart helpers cannot be left behind under a generic no-op fallback.

## EndpointDiagnostics: inspect bounded sanitized helper evidence

**Call/result.** `EndpointDiagnosticRequest` combines exact endpoint/instance ownership, query cursor and positive byte bound. EndpointDiagnosticPage contains sanitized typed records, Next and truncation. Cursors bind query/endpoint/snapshot; expiry is explicit. Reads cannot repair or restart delivery.

**Concrete use.** When qualifying a pending listener, retrieve bounded helper diagnostics and continue if Next is present. Truncated evidence is not complete diagnosis. Never expose raw helper configuration containing secrets, bootstrap input or arbitrary Exec output. Unsupported means this optional evidence is unavailable; it is not proof of endpoint health.

**Current mapping.** Tart has helper diagnostic maps used by live tests, not production core; Orchard has no hook. The shared bounded, scoped, sanitized shape is an intentional refinement for qualification and does not imply a new production dependency.

## Errors and review decisions

Errors support wrapped `errors.Is` for invalid argument, not-found, conflict, unsupported, unavailable, unknown outcome and cursor expiry, plus `errors.As` for safe RequestID correlation. Public Error carries no raw backend cause/request/credential. Adapters must bound and sanitize messages. Unknown codes are not known sentinels; a nil typed error must not spuriously match them.

The decisions exposed by this API are whether a backend can supply the requested clone mode, Docker features, restart delivery and cleanup guarantees. Missing optional behavior selects only the documented equivalent fallback; it does not manufacture support. The present source tests establish representability, error inspection and conservative consumer decisions. Actual native lifecycle, suspended compatibility, secrecy, durable delivery and endpoint connectivity remain later implementation evidence.

## Configuration continuity

Before any dispatch, resolve the journal's backend key and verify its immutable configuration revision. Missing or changed configuration fails before backend calls. A new Go object may reopen the unchanged domain. BackendInfo.ID is informational. Persist uncertainty before Create; after restart observe only the same backend and validate known UID plus full effective intent. Never search other backends or replay uncertain Create. Source and endpoint reservations retain this association.

Names and UIDs can collide across domains; a bare reference cannot always detect wrong-backend routing. Source locators resolve locally. Foreign cursors fail without continuation or completeness. Recover reads only originating durable records. Handles may overlap physical capacity: do not sum their available counts without accounting for overlap. Private worker affinity requires later adapter qualification.
