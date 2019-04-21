# npm-cache-proxy

![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/emeralt/npm-cache-proxy.svg?style=for-the-badge)

- [npm-cache-proxy](#npm-cache-proxy)
	- [Download](#download)
	- [Usage](#usage)
		- [`ncp`](#ncp)
		- [`ncp list`](#ncp-list)
		- [`ncp purge`](#ncp-purge)
	- [Programmatic usage](#programmatic-usage)
	- [License](#license)

## Download
You can download binary for your platform on the [Releases](https://github.com/emeralt/npm-cache-proxy/releases) page. Alternatively, you can use [Docker Image](https://cloud.docker.com/u/emeralt/repository/docker/emeralt/npm-cache-proxy).

## Usage

### `ncp`

Start proxy server.

| Options                       | Env                | Default                      | Description                         |
| ----------------------------- | ------------------ | ---------------------------- | ----------------------------------- |
| `--listen <address>`          | `LISTEN_ADDRESS`   | `locahost:8080`              | Address to listen                   |
| `--upstream <address>`        | `UPSTREAM_ADDRESS` | `https://registry.npmjs.org` | Upstream registry address           |
| `--cache-limit <count>`       | `CACHE_LIMIT`      | -                            | Cached packages count limit         |
| `--cache-ttl <timeout>`       | `CACHE_TTL`        | `3600`                       | Cache expiration timeout in seconds |
| `--redis-address <address>`   | `REDIS_ADDRESS`    | `http://localhost:6379`      | Redis address                       |
| `--redis-database <database>` | `REDIS_DATABASE`   | `0`                          | Redis database                      |
| `--redis-password <password>` | `REDIS_PASS`       | -                            | Redis password                      |
| `--redis-prefix <prefix>`     | `REDIS_PREFIX`     | `ncp-`                       | Redis keys prefix                   |

### `ncp list`

List cached packages.

### `ncp purge`

Purge cached packages.

## Programmatic usage

```golang
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
```

## License

[MIT](./license)
