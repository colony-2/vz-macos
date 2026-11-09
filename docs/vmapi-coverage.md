# Audited consumer coverage

For detailed review, see the [Pulse/Tart/Orchard support comparison](vmapi-support.md) and [operation behavior guide](vmapi-operations.md).

Every row of the pulse-tart owner's coverage evidence retained with the review maps below to revised declarations and source acceptance. That evidence checked pulse-tart revision ac235e01e7ca4b3bb65be8d39204d92985af4680. These are coverage commitments, not production implementation results.

| Audited dependency | Proposed contract and preserved behavior | Acceptance |
|---|---|---|
| Info, List, Get, Create, Delete | Backend methods; complete inventory, effective intent/UID result, accepted versus complete deletion | API-001–009, API-028 |
| Suspend, Events | Baseline; convergence criteria and typed persistent bootstrap evidence | API-016, API-028 |
| Exec, Stop, endpoint ensure/observe | Present methods; required capabilities for Docker configurations | API-013–019 |
| AdoptSource | Optional method; verified Get fallback, inherited trust and replaceable published selectors | API-005, API-015, API-019–020 |
| RetainSource | Optional reusable/delivered marker; no-op fallback where marking is unnecessary | API-005, API-019–020 |
| ProbeLocalSources | Info capability selects ordinary discovery or caller-controlled disposable probes | API-019–020, API-028 |
| Recover | Context-only durable adapter recovery; eligible bootstrap only, equivalent restart-delivery requirement | API-025–026 |
| RemoveEndpoint | Exact identity cleanup; fallback only with equivalent lifecycle cleanup | API-017–018, API-026 |
| EndpointDiagnostics | Optional bounded endpoint-specific qualification records | API-029 |
| Close | Idempotent resource/worker disposal independent from VM lifetime | API-027 |

Cross-cutting requirements: ambiguity and cancellation (API-008–009), inspectable safe errors (API-011–012), bounded retrieval and explicit expiry (API-010), at-least-once delivery without arbitrary Exec replay (API-014, API-025), source trust and output secrecy (API-005, API-015–016), complete foreign capacity accounting (API-004), ownership separation and later native qualification (API-024). API-021–023 and API-030 establish packaging, regression scope and candidate migration.

Configured-instance binding adds API-031–035: backend-local references, journal routing and configuration integrity, independent sources/cursors/recovery, overlapping capacity accounting, collision limits and restart reconciliation without Create replay. Every audited dependency above remains covered.
