package proxy

import (
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

type Proxy struct {
	RedisClient *redis.Client
	HttpClient  *http.Client

	GetOptions func() (Options, error)
}

type Options struct {
	RedisPrefix        string
	RedisExpireTimeout time.Duration

	UpstreamAddress     string
	ReplaceAddress      string
	StaticServerAddress string
}
