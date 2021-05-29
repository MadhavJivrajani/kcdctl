package core

// Container represents a docker contianer that would be spawned.
type Container struct {
	Name          string
	Image         string
	HostPort      string
	ContainerPort string
	Network       string
}

// DesiredState represents the desired state of the system,
// in terms of the number of containers required to be
// running and what the type of this container is that
// should be running.
type DesiredState struct {
	DesiredNum    int
	ContainerType Container
}

// CurrentState represents the current state of the system,
// in terms of the number of containers that are currently
// running in the system and what the type of this contianer
// is that is running.
type CurrentState struct {
	CurrentNum    int
	ContainerType Container
}

// Diff represents a drift of the Current state of the system
// from the Desired state of the system.
type Diff struct {
	Current *CurrentState
	Desired *DesiredState
}

// LoadBalancer represents the configuration of the loadbalancer
// that will be created at the time of system bootsrapping.
type LoadBalancer struct {
	// Image is the image from which
	// the load balancer container will
	// be spawned.
	Image string `yaml:"image"`
	// Name given to the lb container.
	Name string `yaml:"name"`
	// ExposedPort is the port that
	// is exposed and available to
	// users.
	ExposedPort string `yaml:"exposedPort"`
	// ContainerPort is the port that
	// the ExposedPort will be mapped
	// to.
	ContainerPort string `yaml:"containerPort"`
	// TargetPort is the port that
	// the load balancer will proxy
	// requests to.
	TargetPort string `yaml:"targetPort"`
}
