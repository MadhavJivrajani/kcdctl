package imperative

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/utils"
	"github.com/docker/docker/client"
)

const processorPort = "9000"

var (
	ctx          context.Context
	cli          *client.Client
	desiredState *core.DesiredState
)

func spawn(w http.ResponseWriter, req *http.Request) {
	log.Println("Spawning container...")
	err := utils.SpawnContainer(ctx, cli, desiredState.ContainerType)
	if err != nil {
		fmt.Fprintf(w, "error spawning container %v\n", err)
		return
	}
	fmt.Fprintf(w, "spawn successful!\n")
	log.Println("Spawn successful...")
}

func help(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome!\nAvailable endpoints are: /spawn\n")
}

// Processor takes action based on the command issued by the user.
func Processor(_ctx context.Context, _cli *client.Client, _desiredState *core.DesiredState, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx = _ctx
	cli = _cli
	desiredState = _desiredState
	http.HandleFunc("/", help)
	http.HandleFunc("/spawn", spawn)
	log.Println("Listening on port:", processorPort)
	log.Fatal(http.ListenAndServe(":"+processorPort, nil))
}
