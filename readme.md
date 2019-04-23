<h1 align="center">
<img width="400" src="./logo.png"> 

npm-cache-proxy
</h1>

<p align="center">
  <a href="https://hub.docker.com/r/emeralt/npm-cache-proxy/tags">
    <img src="https://img.shields.io/github/release/emeralt/npm-cache-proxy.svg" alt="Current Release" />
  </a>
  <a href="https://hub.docker.com/r/emeralt/npm-cache-proxy/builds">
    <img src="https://img.shields.io/docker/cloud/build/emeralt/npm-cache-proxy.svg" alt="CI Build">
  </a>
  <a href="https://github.com/emeralt/npm-cache-proxy/blob/master/liscense">
    <img src="https://img.shields.io/github/license/emeralt/npm-cache-proxy.svg" alt="Licence">
  </a>
</p>

> `npm-cache-proxy` is a lightweight npm caching proxy written in Go that achieves warp speeds by using Redis for data storage. B-b-blazing fast!


<details>
<summary>Table of Contents</summary>
<p>

- [Getting started](#getting-started)
- [CLI](#cli)
	- [`ncp`](#ncp)
	- [`ncp list`](#ncp-list)
	- [`ncp purge`](#ncp-purge)
- [Programmatic usage](#programmatic-usage)
	- [Example](#example)
- [Deployment](#deployment)
- [Benchmark](#benchmark)
- [License](#license)

</p>
</details>

## Getting started
Release binaries for different platforms can be downloaded on the [Releases](https://github.com/emeralt/npm-cache-proxy/releases) page. Also, [Docker image](https://cloud.docker.com/u/emeralt/repository/docker/emeralt/npm-cache-proxy) is provided.

The fastest way to get started with NCP is to use Docker image:
```bash
# run proxy inside of docker container in background
docker run -e REDIS_ADDRESS=host.docker.internal:6379 -p 8080:8080 -it -d emeralt/npm-cache-proxy

# configure npm to use caching proxy as registry
npm config set registry http://localhost:8080
```

## CLI
Additionally, NCP provides a command line utility for proxy configuration and data management.

---

| Global Options                | Env              | Default                 | Description       |
| ----------------------------- | ---------------- | ----------------------- | ----------------- |
| `--redis-address <address>`   | `REDIS_ADDRESS`  | `http://localhost:6379` | Redis address     |
| `--redis-database <database>` | `REDIS_DATABASE` | `0`                     | Redis database    |
| `--redis-password <password>` | `REDIS_PASSWORD` | -                       | Redis password    |
| `--redis-prefix <prefix>`     | `REDIS_PREFIX`   | `ncp-`                  | Redis keys prefix |

---

### `ncp`

Start proxy server.

```bash
ncp --listen "localhost:1234" # listen on port 1234
```

| Options                 | Env                | Default                      | Description                         |
| ----------------------- | ------------------ | ---------------------------- | ----------------------------------- |
| `--listen <address>`    | `LISTEN_ADDRESS`   | `locahost:8080`              | Address to listen                   |
| `--upstream <address>`  | `UPSTREAM_ADDRESS` | `https://registry.npmjs.org` | Upstream registry address           |
| `--cache-limit <count>` | `CACHE_LIMIT`      | -                            | Cached packages count limit         |
| `--cache-ttl <timeout>` | `CACHE_TTL`        | `3600`                       | Cache expiration timeout in seconds |
| `--silent <address>`    | `SILENT`           | `0`                          | Disable logs                        |

---

### `ncp list`
List all cached packages.

---

### `ncp purge`
Purge all cached packages.

---

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
	"github.com/go-redis/redis"
)

func main() {
	// create proxy
	proxy := npmproxy.Proxy{
		// use redis as database
		Database: npmproxy.DatabaseRedis{
			// see github.com/go-redis/redis
			Client: redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
			}),
		},

		// reuse connections
		HttpClient: &http.Client{},
	}

	// create and start server
	proxy.Server(npmproxy.ServerOptions{
		ListenAddress: "localhost:8080",

		// allow fetching options dynamically on each request
		GetOptions: func() (npmproxy.Options, error) {
			return npmproxy.Options{
				DatabasePrefix:     "ncp-",
				DatabaseExpiration: 1 * time.Hour,
				UpstreamAddress:    "https://registry.npmjs.org",
			}, nil
		},
	}).ListenAndServe()
}
```

## Deployment
NCP can be deployed using Kubernetes, Docker Compose or any other container orchestration platform. NCP supports running indefinite amount of instances simultaneously. 

## Benchmark
Macbook Pro 15″ 2017, Intel Core i7-7700HQ. Note `GOMAXPROCS=1`. 

```bash
# SILENT=1 GIN_MODE=release GOMAXPROCS=1 ncp

$ go-wrk -c 100 -d 10 http://localhost:8080/ascii
Running 10s test @ http://localhost:8080/ascii
  100 goroutine(s) running concurrently
120227 requests in 9.996788553s, 136.21MB read
Requests/sec:		12026.56
Transfer/sec:		13.63MB
Avg Req Time:		8.314928ms
Fastest Request:	1.157793ms
Slowest Request:	782.917402ms
Number of Errors:	0
```

## License

[MIT](./license)
