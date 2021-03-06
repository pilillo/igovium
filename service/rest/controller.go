package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilillo/igovium/cache"
	"github.com/pilillo/igovium/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	keyParam = "key"
)

func Get(c *gin.Context) {
	key := c.Param(keyParam)
	value, response := cacheService.Get(key)
	if response != nil {
		c.JSON(response.Status, response)
	} else {
		c.JSON(http.StatusOK, value)
	}
}

func Put(c *gin.Context) {
	cacheEntry := cache.CacheEntry{}
	if err := c.BindJSON(&cacheEntry); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		response := cacheService.Put(&cacheEntry)
		c.JSON(response.Status, response)
	}
}

func Delete(c *gin.Context) {
	key := c.Param(keyParam)
	response := cacheService.Delete(key)
	c.JSON(response.Status, response)
}

var router = gin.Default()

// Ping ... replies to a ping message for healthcheck purposes
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

var cacheService cache.CacheService

// on module import, get singleton cache service
func init() {
	cacheService = cache.GetCacheService()
}

func StartEndpoint(config *utils.Config) {

	err := cacheService.Init(config)
	if err != nil {
		panic(err)
	}

	router.GET("/healthcheck", Ping)
	router.GET("/metrics", prometheusHandler())

	// cache service
	router.GET(fmt.Sprintf("/:%s", keyParam), Get)
	router.DELETE(fmt.Sprintf("/:%s", keyParam), Delete)
	router.PUT("/", Put)

	router.Run(fmt.Sprintf(":%d", config.RESTConfig.Port))
}
