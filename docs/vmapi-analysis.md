# vz-macos API: current contract and remaining native work

The API is already importable at `github.com/tmc/vz-macos/vmapi`. I inspected HEAD `2939ddfa06319e61576dc00d4324085116f75f88`; the consolidated contract has no production implementation or constructor. The revised contract removes public Placement and binds resources through configured Backend instances.

| Area and source | Exists now | Needs adding or qualification later |
|---|---|---|
| vmapi/backend.go | One concurrent Backend with Info, List, Get, Create, Suspend, Delete, Events, Close, Stop, Exec, EnsureEndpoint, GetEndpoint, AdoptSource, RetainSource, Recover, RemoveEndpoint, EndpointDiagnostics | Configured-instance binding is documented; implement native behavior later |
| vmapi/types.go, inventory.go, lifecycle.go, sources.go | Shared identity, power, observation, effective intent, source and capability values, with backend-local references | Domain invariants are documented; fixtures exercise backend configuration and caller journals |
| vmapi/errors.go and errors_test.go | Executable typed error Is/As behavior | Native adapters must map rejection, conflict, unavailable observation and uncertain dispatch accurately |
| vmapi/contract_test.go, reconcile_test.go, example_test.go | Three independent test-only Backends and full effective-intent comparison/no-replay cases | Fixtures cover collisions, missing/rebound handles, durable journal reopening and no fallback dispatch |
| vm_registry.go:257 | ListVMs creates its directory and silently skips invalid entries | Observational complete inventory including foreign capacity consumers; explicit partial/unavailable evidence |
| vm_ops.go:14 | DeleteVM checks a name/path and removes it | Atomic native UID preconditions and replacement-race qualification |
| clone.go:23,109 | Disk/config cloning explicitly removes suspend.vmstate | Saved-memory clone compatibility, collision exclusion and source lineage/trust preservation; never claim cold clone equivalence |
| runtime_lifecycle.go, control_http.go, agent_client.go, control_client.go | Lifecycle/control/guest primitives in package main | Importable native adapter with transport/trust hidden behind conceptual intent; existing primitives alone do not implement Backend |
| control_operations_store.go | File-backed operation records | Durable, secret-safe eligible bootstrap delivery and restart recovery; never replay ambiguous clone/run or arbitrary Exec |
| port_forward.go:39 | localhost TCP listener forwarding to guest vsock through ControlServer | Exact allocation/VM-owned TCP guest endpoints, address-change recovery and identity-safe cleanup |
| go.mod | Portable vmapi package within a native application module; machine-local Apple replacement | Keep vmapi dependency isolation. Native root gates require the Apple source and compatible Mac/signing environment |

The recovered audit identified the same core gaps: complete inventory, identity-safe mutation, suspended cloning, delivery evidence and owned TCP endpoints. The supplied pulse-tart owner consultation confirmed the consolidated consumer coverage at audited revision `ac235e01e7ca4b3bb65be8d39204d92985af4680`. It does not prove native VZ behavior or private Orchard worker binding.

I own this contract and later VZ implementation in vz-macos. Pulse-tart owns its consumer journals, admission, recipes/readiness and eventual provider migration. Its existing consultation is retained unchanged; no new migration assignment or remote behavior change is proposed. Native qualification remains a separate implementation task.
