package vmapi

import "time"

type (
	InstanceUID  string
	RequestID    string
	EndpointUID  string
	AllocationID string
	Cursor       string
	SnapshotID   string
)

// InstanceRef requires a nonempty immutable UID; names alone never authorize mutation.
// UIDs are never reused for replacements in a backend domain.
type InstanceRef struct {
	Name string
	UID  InstanceUID
}
type EndpointRef struct {
	UID        EndpointUID
	Allocation AllocationID
}

// Hardware uses logical CPUs and bytes. Nil hardware means unknown.
type Hardware struct {
	CPUs                   uint32
	MemoryBytes, DiskBytes uint64
}
type GuestOS string

const (
	GuestOSUnknown GuestOS = ""
	GuestOSMacOS   GuestOS = "macos"
	GuestOSLinux   GuestOS = "linux"
)

type PowerState string

const (
	PowerUnknown   PowerState = ""
	PowerStopped   PowerState = "stopped"
	PowerRunning   PowerState = "running"
	PowerSuspended PowerState = "suspended"
)

type Freshness string

const (
	FreshnessUnknown Freshness = ""
	FreshnessFresh   Freshness = "fresh"
	FreshnessStale   Freshness = "stale"
)

type Management string

const (
	ManagementUnknown Management = ""
	ManagementManaged Management = "managed"
	ManagementForeign Management = "foreign"
)

type Observation struct {
	At        time.Time
	Freshness Freshness
}
type Transition string

const (
	TransitionUnknown    Transition = ""
	TransitionNone       Transition = "none"
	TransitionStarting   Transition = "starting"
	TransitionSuspending Transition = "suspending"
	TransitionStopping   Transition = "stopping"
)

// Generations are observation evidence, not mutation tokens; zero is unspecified
// when Known is false.
type Convergence struct {
	DesiredGeneration, ObservedGeneration uint64
	Known                                 bool
	Transition                            Transition
}
type Condition struct {
	Code    ErrorCode
	Message string
}
type Instance struct {
	Ref                         InstanceRef
	OS                          GuestOS
	Hardware                    *Hardware
	Management                  Management
	Source                      ResolvedSource
	EffectiveIntent             EffectiveIntent
	DesiredPower, ObservedPower PowerState
	Observation                 Observation
	Convergence                 Convergence
	Delivery                    DeliveryPhase
	Conditions                  []Condition
}
type Mutation struct{ RequestID RequestID }
type InstanceMutation struct {
	Mutation
	Instance InstanceRef
}
