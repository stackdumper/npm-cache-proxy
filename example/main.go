package main

import (
	"net/http"
	"time"

	npmproxy "github.com/emeralt/npm-cache-proxy/proxy"
	"github.com/go-redis/redis"
)

func main() {
	proxy := npmproxy.Proxy{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			DB:       0,
			Password: "",
		}),
		HttpClient: &http.Client{},
		GetOptions: func() (npmproxy.Options, error) {
			return npmproxy.Options{
				RedisPrefix:        "ncp-",
				RedisExpireTimeout: 1 * time.Hour,
				UpstreamAddress:    "https://registry.npmjs.org",
			}, nil
		},
	}

	proxy.Server(npmproxy.ServerOptions{
		ListenAddress: "localhost:8080",
	}).ListenAndServe()
}
