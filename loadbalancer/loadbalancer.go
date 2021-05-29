package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var (
	exposedPort          string
	targetPort           string
	serviceSelector      string
	serviceSelectorValue string
)

const networkName = "kcd-bangalore-demo"

func discoverServices() ([]string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	listOpts := types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "label",
				Value: fmt.Sprintf("%s=%s", serviceSelector, serviceSelectorValue),
			},
		),
	}
	services, err := cli.ContainerList(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	var ipOfServices []string
	for _, service := range services {
		ipOfServices = append(ipOfServices, service.NetworkSettings.Networks[networkName].IPAddress)
	}

	return ipOfServices, nil
}

func proxyAndRespond(w http.ResponseWriter, req *http.Request, toSendTo string) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s", toSendTo, targetPort))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(body))
}

func serve(w http.ResponseWriter, req *http.Request) {
	services, err := discoverServices()
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())
	// pick a service randomly, can be done in a much
	// better manner.
	serviceToProxyTo := services[rand.Intn(len(services))]
	proxyAndRespond(w, req, serviceToProxyTo)
}

func startLoadBalancer() {
	http.HandleFunc("/", serve)
	log.Fatal(http.ListenAndServe(":"+exposedPort, nil))
}

func main() {
	exposedPort = os.Getenv("exposedPort")
	targetPort = os.Getenv("targetPort")
	serviceSelector = os.Getenv("serviceSelector")
	serviceSelectorValue = os.Getenv("serviceSelectorValue")

	startLoadBalancer()
}