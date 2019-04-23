package proxy

import (
	"net/http"
	"strings"
	"time"

	ginzap "github.com/gin-contrib/zap"
	gin "github.com/gin-gonic/gin"
	zap "go.uber.org/zap"
)

// ServerOptions provides configuration for Server method
type ServerOptions struct {
	ListenAddress string
	Silent        bool

	GetOptions func() (Options, error)
}

// Server creates http proxy server
func (proxy Proxy) Server(options ServerOptions) *http.Server {
	gin.SetMode("release")

	router := gin.New()

	if options.Silent {
		router.Use(gin.Recovery())
	} else {
		logger, _ := zap.NewProduction()
		router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		router.Use(ginzap.RecoveryWithZap(logger, true))
	}

	router.GET("/:scope/:name", proxy.getPackageHandler(options))
	router.GET("/:scope", proxy.getPackageHandler(options))
	router.NoRoute(proxy.noRouteHandler(options))

	return &http.Server{
		Handler: router,
		Addr:    options.ListenAddress,
	}
}

func (proxy Proxy) getPackageHandler(options ServerOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		options, err := options.GetOptions()

		if err != nil {
			c.AbortWithError(500, err)
		} else {
			pkg, err := proxy.GetCachedPath(options, c.Request.URL.Path, c.Request)

			if err != nil {
				c.AbortWithError(500, err)
			} else {
				c.Header("Cache-Control", "public, max-age="+string(int(options.DatabaseExpiration.Seconds())))
				c.Data(200, "application/json", pkg)
			}
		}
	}
}

func (proxy Proxy) noRouteHandler(options ServerOptions) gin.HandlerFunc {
	tarballHandler := proxy.getPackageHandler(options)

	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, ".tgz") {
			// get tarball
			tarballHandler(c)
		} else if c.Request.URL.Path == "/" {
			// get health
			err := proxy.Database.Health()

			if err != nil {
				c.AbortWithStatusJSON(503, err)
			} else {
				c.AbortWithStatusJSON(200, gin.H{"ok": true})
			}
		} else {
			// redirect
			options, err := options.GetOptions()

			if err != nil {
				c.AbortWithStatusJSON(500, err)
			} else {
				c.Redirect(http.StatusTemporaryRedirect, options.UpstreamAddress+c.Request.URL.Path)
			}
		}
	}
}
