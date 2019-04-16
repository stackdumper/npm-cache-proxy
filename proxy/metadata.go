package proxy

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// GetMetadata returns NPM response for a given package path
func (proxy Proxy) GetMetadata(path string, header http.Header) ([]byte, error) {
	options, err := proxy.GetOptions()
	if err != nil {
		return nil, err
	}

	// get package from redis
	pkg, err := proxy.RedisClient.Get(options.RedisPrefix + path).Result()

	// either package doesn't exist or there's some other problem
	if err != nil {

		// check if error is caused by nonexistend package
		// if no, return error
		if err.Error() != "redis: nil" {
			return nil, err
		}

		// error is caused by nonexistent package
		// fetch package
		req, err := http.NewRequest("GET", options.UpstreamAddress+path, nil)

		// inherit headers from request
		req.Header = header
		if err != nil {
			return nil, err
		}

		res, err := proxy.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// convert body to string
		pkg = string(body)

		// save to redis
		_, err = proxy.RedisClient.Set(
			options.RedisPrefix+path,
			pkg,
			options.RedisExpireTimeout,
		).Result()
		if err != nil {
			return nil, err
		}
	}

	// replace tarball urls
	// FIXME: unmarshall and replace only necessary fields
	convertedPkg := strings.ReplaceAll(string(pkg), options.ReplaceAddress, options.StaticServerAddress)

	return []byte(convertedPkg), nil
}
