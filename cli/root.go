package cli

import (
	"net/http"
	"time"

	npmproxy "github.com/emeralt/npm-cache-proxy/proxy"
	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

// start a server
var rootCmd = &cobra.Command{
	Use:   "ncp",
	Short: "ncp is a fast npm cache proxy that stores data in Redis",
	Run: func(cmd *cobra.Command, args []string) {
		proxy := npmproxy.Proxy{
			RedisClient: redis.NewClient(&redis.Options{}),
			HttpClient: &http.Client{
				Transport: http.DefaultTransport,
			},
			GetOptions: getOptions,
		}

		proxy.Server(npmproxy.ServerOptions{
			ListenAddress: "localhost:8080",
		}).ListenAndServe()
	},
}

func getOptions() (npmproxy.Options, error) {
	return npmproxy.Options{
		RedisPrefix:        "ncp-",
		RedisExpireTimeout: 1 * time.Hour,

		UpstreamAddress:     "http://registry.npmjs.org",
		ReplaceAddress:      "https://registry.npmjs.org",
		StaticServerAddress: "http://localhost:8080",
	}, nil
}
