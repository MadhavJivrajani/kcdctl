package core

// Container represents a docker contianer
// that would be spawned.
type Container struct {
	Name          string
	Image         string
	HostPort      uint32
	ContainerPort uint32
}

// DesiredState represents the desired state of the system,
// in terms of the number of containers required to be
// running and what the type of this container is that
// should be running.
type DesiredState struct {
	DesiredNum    uint32
	ContainerType *Container
}

// CurrentState represents the current state of the system,
// in terms of the number of containers that are currently
// running in the system and what the type of this contianer
// is that is running.
type CurrentState struct {
	CurrentNum    uint32
	ContainerType *Container
}
