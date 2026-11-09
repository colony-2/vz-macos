package vmapi

// SourceSelector is replaceable and contains no credentials. ProbeLocalSources
// selects caller-controlled disposable probes, never mutation during List/Get.
type SourceSelector struct {
	Locator string
	Mode    CloneMode
}

// ResolvedSource binds a clone to its native object at clone time. Published
// replacement must preserve existing clone data, identity and inherited trust.
// TrustIdentity references private adapter-managed host-key trust, not key material.
type ResolvedSource struct{ Identity, NativeLocator, Parent, TrustIdentity string }
type Source struct {
	Selector    SourceSelector
	OS          GuestOS
	Hardware    *Hardware
	Resolved    ResolvedSource
	Observation Observation
}

// Adoption registers a reusable existing source, not image conversion.
type AdoptSourceRequest struct {
	Mutation
	Selector SourceSelector
	Expected ResolvedSource
}

// Retention only marks an exact source reusable/delivered; it is not GC policy.
type RetainSourceRequest struct {
	Mutation
	Source              Source
	Reusable, Delivered bool
}
