package cli

import (
	"fmt"
	"net/http"
	"os"

	npmproxy "github.com/emeralt/npm-cache-proxy/proxy"
	"github.com/go-redis/redis"
)

// global options
var persistentOptions struct {
	RedisAddress  string
	RedisDatabase int
	RedisPassword string
	RedisPrefix   string
}

// initialize global options
func init() {
	rootCmd.PersistentFlags().StringVar(&persistentOptions.RedisAddress, "redis-address", getEnvString("REDIS_ADDRESS", "localhost:6379"), "Redis address")
	rootCmd.PersistentFlags().IntVar(&persistentOptions.RedisDatabase, "redis-database", getEnvInt("REDIS_DATABASE", "0"), "Redis database")
	rootCmd.PersistentFlags().StringVar(&persistentOptions.RedisPassword, "redis-password", getEnvString("REDIS_PASSWORD", ""), "Redis password")
	rootCmd.PersistentFlags().StringVar(&persistentOptions.RedisPrefix, "redis-prefix", getEnvString("REDIS_PREFIX", "ncp-"), "Redis prefix")
}

func getProxy() *npmproxy.Proxy {
	return &npmproxy.Proxy{
		Database: npmproxy.DatabaseRedis{
			Client: redis.NewClient(&redis.Options{
				Addr:     persistentOptions.RedisAddress,
				DB:       persistentOptions.RedisDatabase,
				Password: persistentOptions.RedisPassword,
			}),
		},
		HttpClient: &http.Client{
			Transport: http.DefaultTransport,
		},
	}
}

// Run starts the CLI
func Run() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(purgeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
