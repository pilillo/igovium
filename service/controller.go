package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilillo/igovium/cache"
	"github.com/pilillo/igovium/commons"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	keyParam = "key"
)

func Get(c *gin.Context) {
	key := c.Param(keyParam)
	value, err := restService.Get(key)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	} else {
		c.JSON(http.StatusOK, value)
	}
}

func Put(c *gin.Context) {
	cacheEntry := cache.CacheEntry{}
	if err := c.BindJSON(&cacheEntry); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		err := restService.Put(&cacheEntry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, cacheEntry)
		}
	}
}

func Delete(c *gin.Context) {
	key := c.Param(keyParam)
	err := restService.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, nil)
	}
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

func StartEndpoint(config *commons.Config) {

	err := restService.Init(config)
	if err != nil {
		panic(err)
	}

	router.GET("/healthcheck", Ping)
	router.GET("/metrics", prometheusHandler())

	// cache service
	router.GET(fmt.Sprintf("/:%s", keyParam), Get)
	router.DELETE(fmt.Sprintf("/:%s", keyParam), Delete)
	router.PUT("/", Put)

	router.Run(fmt.Sprintf(":%d", config.Port))
}
