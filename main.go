package main

import (
	"github.com/pilillo/igovium/service/grpc_server"
	"github.com/pilillo/igovium/service/rest"
	"github.com/pilillo/igovium/utils"
)

func main() {
	done := make(chan bool, 1)
	config := utils.LoadCfg()
	// start rest endpoint
	go rest.StartEndpoint(config)
	// start grpc endpoint
	go grpc_server.StartEndpoint(config)
	// wait for signal
	<-done
}
