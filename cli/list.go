package cli

import (
	"fmt"

	npmproxy "github.com/pkgems/npm-cache-proxy/proxy"
	"github.com/spf13/cobra"
)

// start a server
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cached paths",
	Run: func(cmd *cobra.Command, args []string) {
		proxy := getProxy()

		metadatas, err := proxy.ListCachedPaths(npmproxy.Options{
			DatabasePrefix: persistentOptions.RedisPrefix,
		})
		if err != nil {
			panic(err)
		}

		for _, metadata := range metadatas {
			fmt.Println(metadata)
		}
	},
}
