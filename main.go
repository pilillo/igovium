package main

import (
	"github.com/pilillo/igovium/commons"
	"github.com/pilillo/igovium/service"
)

func main() {
	config := commons.LoadCfg()

	service.StartEndpoint(config)
}
