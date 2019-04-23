package cli

import (
	"log"
	"time"

	npmproxy "github.com/emeralt/npm-cache-proxy/proxy"
	"github.com/spf13/cobra"
)

// start a server
var rootCmd = &cobra.Command{
	Use:   "ncp",
	Short: "ncp is a fast npm cache proxy that stores data in Redis",
	Run:   run,
}

var rootOptions struct {
	Silent          bool
	ListenAddress   string
	UpstreamAddress string
	CacheLimit      string
	CacheTTL        int
}

func init() {
	rootCmd.Flags().BoolVar(&rootOptions.Silent, "silent", getEnvBool("SILENT", "0"), "Disable logging")
	rootCmd.Flags().StringVar(&rootOptions.ListenAddress, "listen", getEnvString("LISTEN_ADDRESS", "localhost:8080"), "Address to listen")
	rootCmd.Flags().StringVar(&rootOptions.UpstreamAddress, "upstream", getEnvString("UPSTREAM_ADDRESS", "https://registry.npmjs.org"), "Upstream registry address")
	rootCmd.Flags().StringVar(&rootOptions.CacheLimit, "cache-limit", getEnvString("CACHE_LIMIT", "0"), "Cached packages count limit")
	rootCmd.Flags().IntVar(&rootOptions.CacheTTL, "cache-ttl", getEnvInt("CACHE_TTL", "3600"), "Cache expiration timeout in seconds")
}

func run(cmd *cobra.Command, args []string) {
	proxy := getProxy()

	log.Print("Listening on " + rootOptions.ListenAddress)

	err := proxy.Server(npmproxy.ServerOptions{
		ListenAddress: rootOptions.ListenAddress,
		Silent:        rootOptions.Silent,

		GetOptions: func() (npmproxy.Options, error) {
			return npmproxy.Options{
				DatabasePrefix:     persistentOptions.RedisPrefix,
				DatabaseExpiration: time.Duration(rootOptions.CacheTTL) * time.Second,
				UpstreamAddress:    rootOptions.UpstreamAddress,
			}, nil
		},
	}).ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
