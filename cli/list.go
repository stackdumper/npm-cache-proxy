package cli

import (
	"fmt"

	npmproxy "github.com/emeralt/npm-cache-proxy/proxy"
	"github.com/spf13/cobra"
)

// start a server
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cached packages",
	Run: func(cmd *cobra.Command, args []string) {
		proxy := getProxy(func() (npmproxy.Options, error) {
			return npmproxy.Options{
				RedisPrefix: persistentOptions.RedisPrefix,
			}, nil
		})

		metadatas, err := proxy.ListMetadata()
		if err != nil {
			panic(err)
		}

		for _, metadata := range metadatas {
			fmt.Println(metadata)
		}
	},
}
