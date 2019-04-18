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
...

## Usage

### `ncp`
Start proxy server.

| Options                   | Env           | Default                 | Description                 |
| ------------------------- | ------------- | ----------------------- | --------------------------- |
| `-p, --port <port>`       | `HTTP_PORT`   | `8080`                  | Port to listen to           |
| `-l, --limit <count>`     | `CACHE_LIMIT` | -                       | Cached packages count limit |
| `-t, --ttl <timeout>`     | `CACHE_TTL`   | `3600`                  | Cache expiration timeout    |
| `-h, --host  <address>`   | `REDIS_HOST`  | `http://localhost:6379` | Redis address               |
| `-d, --db <database>`     | `REDIS_DB`    | `0`                     | Redis database              |
| `-a, --access <password>` | `REDIS_PASS`  | -                       | Redis password              |


### `ncp list`
List cached packages.

### `ncp purge`
Purge cached pacakges.

## Programmatic usage
...

## License
[MIT](./license)
