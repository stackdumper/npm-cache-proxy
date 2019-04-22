package proxy

import (
	"net/http"
	"time"
)

type Proxy struct {
	Database   Database
	HttpClient *http.Client

	GetOptions func() (Options, error)
}

type Options struct {
	DatabasePrefix     string
	DatabaseExpiration time.Duration
	UpstreamAddress    string
}

type Database interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
	Delete(key string) error
	Keys(prefix string) ([]string, error)
	Health() error
}
