package cli

import (
	"net/http"

	npmproxy "github.com/emeralt/npm-cache-proxy/proxy"
	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

// start a server
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge all cached packages",
	Run: func(cmd *cobra.Command, args []string) {
		proxy := npmproxy.Proxy{
			RedisClient: redis.NewClient(&redis.Options{}),
			HttpClient: &http.Client{
				Transport: http.DefaultTransport,
			},
			GetOptions: getOptions,
		}

		err := proxy.PurgeMetadata()
		if err != nil {
			panic(err)
		}
	},
}
