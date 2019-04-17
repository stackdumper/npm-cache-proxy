package cli

import (
	"fmt"
	"net/http"

	"github.com/emeralt/npm-cache-proxy/proxy"
	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

// start a server
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cached packages",
	Run: func(cmd *cobra.Command, args []string) {
		prx := proxy.Proxy{
			RedisClient: redis.NewClient(&redis.Options{}),
			HttpClient: &http.Client{
				Transport: http.DefaultTransport,
			},
			GetOptions: getOptions,
		}

		metadatas, err := prx.ListMetadata()
		if err != nil {
			panic(err)
		}

		for _, metadata := range metadatas {
			fmt.Println(metadata)
		}
	},
}
