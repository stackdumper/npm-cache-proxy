<div align="center">
  <img width="450" src="./logo.png"> 

  <h1>npm-cache-proxy</h1>

  <a href="https://hub.docker.com/r/emeralt/npm-cache-proxy/tags">
    <img src="https://img.shields.io/github/release/emeralt/npm-cache-proxy.svg" alt="Current Release" />
  </a>
  <a href="https://hub.docker.com/r/emeralt/npm-cache-proxy/builds">
    <img src="https://img.shields.io/docker/cloud/build/emeralt/npm-cache-proxy.svg" alt="CI Build">
  </a>
  <a href="https://github.com/emeralt/npm-cache-proxy/blob/master/license">
    <img src="https://img.shields.io/github/license/emeralt/npm-cache-proxy.svg" alt="Licence">
  </a>
</div>

<br />
<br />

## Introduction

#### ⚡️ Performance
NCP is a tiny but very fast caching proxy written in Go. It uses Redis for data storage, which in combination with the speed of Go makes it incredibly fast. NCP is well-optimized and can be run on almost any platform, so if you have a Raspberry Pi, you can install NCP as your local cache there.

#### ✨ Modularity
NCP is modular. Now it has only one database adapter which is Redis. If you need support for any other database, feel free to open an issue or implement it [on your own](https://github.com/emeralt/npm-cache-proxy/blob/7c8b90ff6ba0656f60e3de915b9fb4eaabfb467b/proxy/proxy.go#L29) and then open a pull request (_bonus points_).

#### 💡 Simplicity
NCP is very simple. It just proxies requests to an upstream registry, caches response and returns cached response for next requests to the same package. Cached data are stored in Redis with an original request URL as a key.


<br />


## Installation
NCP binaries for different paltforms can be downloaded can be downloaded on the [Releases](https://github.com/emeralt/npm-cache-proxy/releases) page. Also, Docker image is provided on [Docker Hub](https://cloud.docker.com/u/emeralt/repository/docker/emeralt/npm-cache-proxy).

#### 💫 Quick Start
The quickies way to get started with NCP is to use Docker.

```bash
# run proxy inside of docker container in background
docker run -e REDIS_ADDRESS=host.docker.internal:6379 -p 8080:8080 -it -d emeralt/npm-cache-proxy

# configure npm to use caching proxy as registry
npm config set registry http://localhost:8080
```

<br />

## CLI
NCP provides command line interface for interaction with a cached data.

<details>
<summary>Options</summary>

| Options                       | Env                | Default                      | Description                         |
| ----------------------------- | ------------------ | ---------------------------- | ----------------------------------- |
| `--listen <address>`          | `LISTEN_ADDRESS`   | `locahost:8080`              | Address to listen                   |
| `--upstream <address>`        | `UPSTREAM_ADDRESS` | `https://registry.npmjs.org` | Upstream registry address           |
| `--silent <address>`          | `SILENT`           | `0`                          | Disable logs                        |
| `--cache-limit <count>`       | `CACHE_LIMIT`      | -                            | Cached packages count limit         |
| `--cache-ttl <timeout>`       | `CACHE_TTL`        | `3600`                       | Cache expiration timeout in seconds |
| `--redis-address <address>`   | `REDIS_ADDRESS`    | `http://localhost:6379`      | Redis address                       |
| `--redis-database <database>` | `REDIS_DATABASE`   | `0`                          | Redis database                      |
| `--redis-password <password>` | `REDIS_PASSWORD`   | -                            | Redis password                      |
| `--redis-prefix <prefix>`     | `REDIS_PREFIX`     | `ncp-`                       | Redis keys prefix                   |

</details>

#### `ncp`
Start NCP server.

#### `ncp list`
List cached url paths.

#### `ncp purge`
Purge cached url paths.


<br />


## Benchmark
Benchmark is run on Macbook Pro 15″ 2017, Intel Core i7-7700HQ.

#### 1️⃣ 1 process

```bash
# GOMAXPROCS=1 ncp --silent

$ go-wrk -c 100 -d 6 http://localhost:8080/tiny-tarball
Running 6s test @ http://localhost:8080/tiny-tarball
  100 goroutine(s) running concurrently

70755 requests in 5.998378587s, 91.16MB read

Requests/sec:		11795.69
Transfer/sec:		15.20MB
Avg Req Time:		8.477674ms
Fastest Request:	947.743µs
Slowest Request:	815.787409ms
Number of Errors:	0
```

#### ♾ unlimited processes

```bash
# ncp --silent

$ go-wrk -c 100 -d 6 http://localhost:8080/tiny-tarball
Running 6s test @ http://localhost:8080/tiny-tarball
  100 goroutine(s) running concurrently

115674 requests in 5.98485984s, 149.04MB read

Requests/sec:		19327.77
Transfer/sec:		24.90MB
Avg Req Time:		5.173902ms
Fastest Request:	273.015µs
Slowest Request:	34.777963ms
Number of Errors:	0
```


<br />


## Programmatic Usage
NCP provides `proxy` go package that can be used programmatically. Docs are available on [godoc.org](https://godoc.org/github.com/emeralt/npm-cache-proxy/proxy).

#### 🤖 Example
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


<br />


## License
[MIT](./license)
