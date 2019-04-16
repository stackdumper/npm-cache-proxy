package main

import (
	"net/http"
	"time"

	npmproxy "github.com/emeralt/npm-cache-proxy/proxy"
	"github.com/go-redis/redis"
)

func getOptions() (npmproxy.Options, error) {
	return npmproxy.Options{
		RedisPrefix:        "",
		RedisExpireTimeout: 1 * time.Hour,

		UpstreamAddress:     "http://registry.npmjs.org",
		ReplaceAddress:      "https://registry.npmjs.org",
		StaticServerAddress: "http://localhost:8080",
	}, nil
}

func main() {
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
}
