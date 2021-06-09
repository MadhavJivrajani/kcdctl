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

package notifier

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

// feedback watches the state of the system
// and lists events generated on change in
// system state
type feedback struct {
	Events <-chan events.Message
	Errors <-chan error
}

// newFeedback is a constructor for the feedback type
func newFeedback(ctx context.Context, cli *client.Client) *feedback {
	events, errs := cli.Events(ctx, types.EventsOptions{})

	return &feedback{
		events,
		errs,
	}
}
