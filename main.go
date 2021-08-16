package main

import (
	"github.com/pilillo/igovium/service/grpc_server"
	"github.com/pilillo/igovium/service/rest"
	"github.com/pilillo/igovium/utils"
)

func main() {
	done := make(chan bool, 1)
	config := utils.LoadCfg()
	// start rest endpoint (if conf defined)
	if config.RESTConfig != nil {
		go rest.StartEndpoint(config)
	}
	// start grpc endpoint (if conf defined)
	if config.GRPCConfig != nil {
		go grpc_server.StartEndpoint(config)
	}
	// wait for signal
	<-done
}
