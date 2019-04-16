package proxy

import (
	"net/http"
	"strings"
	"time"

	ginzap "github.com/gin-contrib/zap"
	gin "github.com/gin-gonic/gin"
	zap "go.uber.org/zap"
)

type ServerOptions struct {
	ListenAddress string
}

func (proxy Proxy) Server(options ServerOptions) *http.Server {
	router := gin.New()

	logger, _ := zap.NewProduction()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	router.GET("/:scope/:name", proxy.GetPackageHandler)
	router.GET("/:scope", proxy.GetPackageHandler)
	router.NoRoute(proxy.NoRouteHandler)

	return &http.Server{
		Handler: router,
		Addr:    options.ListenAddress,
	}
}

func (proxy Proxy) GetPackageHandler(c *gin.Context) {
	key := c.Request.URL.Path

	pkg, err := proxy.GetMetadata(key, c.Request.Header)

	if err != nil {
		c.AbortWithError(500, err)
	} else {
		c.Data(200, "application/json", pkg)
	}
}

func (proxy Proxy) NoRouteHandler(c *gin.Context) {
	if strings.Contains(c.Request.URL.Path, ".tgz") {
		proxy.GetPackageHandler(c)
	} else if c.Request.URL.Path == "/" {
		_, err := proxy.RedisClient.Ping().Result()

		if err != nil {
			c.AbortWithStatusJSON(503, err)
		} else {
			c.AbortWithStatusJSON(200, gin.H{"ok": true})
		}
	} else {
		options, err := proxy.GetOptions()

		if err != nil {
			c.AbortWithStatusJSON(500, err)
		} else {
			c.Redirect(http.StatusTemporaryRedirect, options.UpstreamAddress+c.Request.URL.Path)
		}
	}
}
