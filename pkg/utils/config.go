/*
Copyright © 2021 Madhav Jivrajani madhav.jiv@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"os"

	"github.com/MadhavJivrajani/kcdctl/pkg/core"
	"gopkg.in/yaml.v2"
)

// Template represents information nescessary to
// start a container that is a part of the desired
// state of the system that the user defines.
type Template struct {
	// Image is the image using which the container
	// will be created.
	Image string
	// Name is NOT the name of the container but an
	// identifier shared by the containers running
	// as part of this configuration.
	Name string
}

// Spec represents the desired state of the system.
type Spec struct {
	// Replicas is the number of instances of the
	// container that should be running.
	Replicas int
	// Template holds information about how to create
	// and start containers.
	Template Template
}

// Config represents the configuration that the user
// passes to the tool, it contains the configuration
// for the load balancer to be user and the desired
// state of the system.
type Config struct {
	// LoadBalancer contains the configuration for the
	// load balancer of the system, which is the point
	// of entry to the application running.
	LoadBalancer core.LoadBalancer `yaml:"loadBalancer"`
	// Spec is the desired state of the system.
	Spec Spec
}

// ReadConfig reads the provided config into the Config
// type and returns an error if any.
func ReadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
