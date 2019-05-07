package cli

import (
	npmproxy "github.com/pkgems/npm-cache-proxy/proxy"
	"github.com/spf13/cobra"
)

// start a server
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge all cached paths",
	Run: func(cmd *cobra.Command, args []string) {
		proxy := getProxy()

		err := proxy.PurgeCachedPaths(npmproxy.Options{
			DatabasePrefix: persistentOptions.RedisPrefix,
		})
		if err != nil {
			panic(err)
		}
	},
}
