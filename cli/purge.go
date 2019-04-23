package cli

import (
	npmproxy "github.com/emeralt/npm-cache-proxy/proxy"
	"github.com/spf13/cobra"
)

// start a server
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge all cached paths",
	Run: func(cmd *cobra.Command, args []string) {
		proxy := getProxy(func() (npmproxy.Options, error) {
			return npmproxy.Options{
				DatabasePrefix: persistentOptions.RedisPrefix,
			}, nil
		})

		err := proxy.PurgeCachedPaths()
		if err != nil {
			panic(err)
		}
	},
}
