# kcdctl
<p align="center">
    <img src="https://github.com/ashleymcnamara/gophers/blob/master/GOPHER_DAD.png?raw=true" width="300" height="300">
</p>

This is the source code for the demo done as part of the talk "Imperative, Declarative and Kubernetes" at the [Kubernetes Community Days, Bengaluru 2021](https://community.cncf.io/events/details/cncf-kcd-bengaluru-presents-kubernetes-community-days-bengaluru) conferernce.

## Talk and Slides
- The talk recording can be found [here](https://www.youtube.com/watch?v=hB3H4_YRnFc).
- The slides used for the talk can be found [here](https://speakerdeck.com/madhavjivrajani/imperative-declarative-and-kubernetes).

## Pre-requisites
To run the demo, you will need the following installed:
- Docker
- Go

## Installation
### Install using `go get`
```
go get -u github.com/MadhavJivrajani/kcdctl
```
Note: This will install the `kcdctl` executable which can be used anywhere in the termnial provided $GOPATH/bin is in your PATH.

### Install from source
```
git clone https://github.com/MadhavJivrajani/kcdctl.git && cd kcdctl
go build -o kcdctl main.go
mv kcdctl /usr/local/bin # to make the executable available system wide
```
## Usage
```
██╗  ██╗ ██████╗██████╗  ██████╗████████╗██╗     
██║ ██╔╝██╔════╝██╔══██╗██╔════╝╚══██╔══╝██║     
█████╔╝ ██║     ██║  ██║██║        ██║   ██║     
██╔═██╗ ██║     ██║  ██║██║        ██║   ██║     
██║  ██╗╚██████╗██████╔╝╚██████╗   ██║   ███████╗
╚═╝  ╚═╝ ╚═════╝╚═════╝  ╚═════╝   ╚═╝   ╚══════╝

A tool to help demo imperative and declarative systems as part of the KCD Bangalore Conference

Usage:
  kcdctl [command]

Available Commands:
  apply       applies a configuarion to the system
  cleanup     cleanup removes all containers that were created as part of a given configuration file
  help        Help about any command
  spawn       spawn creates containers

Flags:
      --config string   config file (default is $HOME/.kcdctl.yaml)
  -h, --help            help for kcdctl

Use "kcdctl [command] --help" for more information about a command.
```

### `apply`
```
applies a configuarion to the system

Usage:
  kcdctl apply [flags]

Flags:
  -f, --file string   file is the yaml configuarion file that the user wants to apply on the system
  -h, --help          help for apply
  -m, --mode string   mode is the mode the system should be run in: imperative/declarative (default "declarative")

Global Flags:
      --config string   config file (default is $HOME/.kcdctl.yaml)
```
An example of a configuration file can be found [here](./examples/cluster.yaml).
#### To run in declarative mode:
```
# don't really need to mention mode for declarative.
kcdctl apply -f examples/cluster.yaml --mode=declarative
```
#### To run in imperative mode:
```
kcdctl apply -f examples/cluster.yaml --mode=imperative
```

The first step in either of the `mode`s is `BootstrapHost`, in which:
- A docker network called `kcd-bangalore-demo` is created and the load balancer container is attached to this network.
- A load balancer is created in accordance with the config provided by the user.

#### Note:
If an excess number of containers than the desired state are created, nothing is done, excess container deletion isn't taken care of, but ideally should happen.

### `spawn`
```
spawn creates containers

Usage:
  kcdctl spawn [flags]

Flags:
  -d, --diff int      diff is the number of containers to spawn (right now, assumed to be +ve) (default 1)
  -f, --file string   file is the yaml configuarion file that the user wants to apply on the system
  -h, --help          help for spawn

Global Flags:
      --config string   config file (default is $HOME/.kcdctl.yaml)
```
`spawn` is assumed to be used when in an `imperative` mode. In a `declarative` mode, spawning of containers is handled automatically.

### `cleanup`
```
cleanup removes all containers that were created as part of a given configuration file

Usage:
  kcdctl cleanup [flags]

Flags:
  -f, --file string   file is the yaml configuarion file using which containers will be cleaned up
  -h, --help          help for cleanup

Global Flags:
      --config string   config file (default is $HOME/.kcdctl.yaml)
```
`cleanup` performs the following steps:
- First stop containers that are running the application.
- Next stop the load balancer container.
- Finally delete the docker network that was created during `BootstrapHost`

## Images used:
- The image used for the application running can be found [here](https://hub.docker.com/repository/docker/maddyoii/kcd-blr-example).
- The image used for the loadbalancer can be found [here](https://hub.docker.com/repository/docker/maddyoii/kcd-loadbalancer).