package main

import (
	"github.com/pilillo/igovium/cache"
	"github.com/pilillo/igovium/utils"
)

func main() {
	config := utils.LoadCfg()
	cache.HistoricizeDBCache(config.DBCacheConfig)
}
