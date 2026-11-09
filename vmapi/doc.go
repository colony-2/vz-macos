// Package vmapi defines one VM Backend and shared values without an operational
// constructor or implementation. See docs/vmapi-design.md for the contract.
//
// Pulse-tart owns recipes, readiness, launch deduplication, journals, durable
// reservations and Docker pairing. Adapters own native identity, lifecycle,
// authentication, transport, inherited host-key verification and private durable
// delivery records. Pulse owns placement/leases; c2j owns execution/renewal and
// application-secret retrieval. Bootstrap inputs and credentials never appear in
// public observations, events, errors or ordinary diagnostics. Conditions and
// issues must be bounded and sanitized by adapters.
//
// Each configured Backend binds one stable inventory, identity, source, capacity,
// endpoint-host, cursor and recovery domain. BackendInfo.ID is informational,
// not routing or affinity evidence. Consumers journal a stable backend key and
// immutable configuration binding before dispatch, including source and endpoint
// reservations. Missing or rebound configuration rejects before backend calls.
// Reopening unchanged configuration may use a new Go object. Never search another
// backend by name or replay an uncertain Create; observe the journaled domain and
// validate known UID and full effective intent. Names and UIDs can collide across
// domains, so a bare reference cannot always detect wrong-backend routing.
// Multiple configured handles can overlap physical capacity; admission must
// account for that overlap rather than summing their reported availability.
// Source locators resolve locally, without automatic transfer. Cursors bind the
// originating domain, query, identity and snapshot/stream; reject foreign cursors
// without continuation or completeness. Recover reads only originating durable
// records. Unsupported worker affinity must reject adapter configuration.
//
// Consult Info before mutations. Unsupported means no dispatch. Without adoption,
// observe the configured source with Get and verify lineage/trust. Without retention,
// explicitly select a no-op only where marking is unnecessary. ProbeLocalSources
// selects core-managed disposable probes, never observational mutation. Recovery
// may be omitted only with equivalent adapter/service restart delivery; otherwise
// reject configurations requiring it. Endpoint removal may be omitted only with
// advertised equivalent VM lifecycle cleanup. Unsupported diagnostics are skipped,
// never health success. Reject Docker configurations missing Exec/Stop/endpoints.
//
// Native UID races, memory-preserving cloning, collision checks, delivery durability
// and secrecy, endpoint connectivity and restart recovery require later adapter
// qualification. Test fixtures establish source usability, not native safety.
// No CLI or wire behavior is changed. This replaces an unreleased candidate; local
// tags/imports and available publication records are recorded in validation docs.
package vmapi
