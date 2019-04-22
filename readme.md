# npm-cache-proxy

![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/emeralt/npm-cache-proxy.svg?style=for-the-badge)

- [npm-cache-proxy](#npm-cache-proxy)
	- [Quick start](#quick-start)
	- [Deployment](#deployment)
	- [Usage](#usage)
		- [`ncp`](#ncp)
		- [`ncp list`](#ncp-list)
		- [`ncp purge`](#ncp-purge)
	- [Programmatic usage](#programmatic-usage)
		- [Example](#example)
	- [Benchmark](#benchmark)
	- [License](#license)

## Quick start
Release binaries for different platforms can be downloaded on the [Releases](https://github.com/emeralt/npm-cache-proxy/releases) page. Also, [Docker Image](https://cloud.docker.com/u/emeralt/repository/docker/emeralt/npm-cache-proxy) is provided.

```bash
docker run -e REDIS_ADDRESS=host.docker.internal:6379 -p 8080:8080 -it emeralt/npm-cache-proxy
```

## Deployment
NCP can be deployed using Kubernetes, Docker Compose or any other container orchestration platform. NCP supports running indefinite amount of instances simultaneously. 

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
| `--redis-password <password>` | `REDIS_PASSWORD`   | -                            | Redis password                      |
| `--redis-prefix <prefix>`     | `REDIS_PREFIX`     | `ncp-`                       | Redis keys prefix                   |

### `ncp list`

List cached packages.

### `ncp purge`

Purge cached packages.

## Programmatic usage
Along with the CLI, go package is provided. Documentation is available on [godoc.org](https://godoc.org/github.com/emeralt/npm-cache-proxy/proxy).

```bash
go get -u github.com/emeralt/npm-cache-proxy/proxy
```

### Example
```golang
package main

import (
	"net/http"
	"time"

	npmproxy "github.com/emeralt/npm-cache-proxy/proxy"
	redis "github.com/go-redis/redis"
)

func main() {
	proxy := npmproxy.Proxy{
		Database: npmproxy.DatabaseRedis{
			Client: redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
			}),
		},
		HttpClient: &http.Client{},
		GetOptions: func() (npmproxy.Options, error) {
			return npmproxy.Options{
				DatabasePrefix:     "ncp-",
				DatabaseExpiration: 1 * time.Hour,
				UpstreamAddress:    "https://registry.npmjs.org",
			}, nil
		},
	}

	proxy.Server(npmproxy.ServerOptions{
		ListenAddress: "localhost:8080",
	}).ListenAndServe()
}
```

## Benchmark

```bash
# GOMAXPROCS=1 GIN_MODE=release ncp --listen localhost:8080

$ go-wrk -c 100 -d 5 http://localhost:8080/ascii
Running 5s test @ http://localhost:8080/ascii
  100 goroutine(s) running concurrently
33321 requests in 5.00537759s, 212.69MB read
Requests/sec:		6657.04
Transfer/sec:		42.49MB
Avg Req Time:		15.02169ms
Fastest Request:	230.514µs
Slowest Request:	571.420487ms
Number of Errors:	0
```

## License

[MIT](./license)
