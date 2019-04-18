package proxy

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// GetMetadata returns cached NPM response for a given package path
func (proxy Proxy) GetMetadata(name string, originalPath string, header http.Header) ([]byte, error) {
	options, err := proxy.GetOptions()
	if err != nil {
		return nil, err
	}

	// get package from redis
	pkg, err := proxy.RedisClient.Get(options.RedisPrefix + name).Result()

	// either package doesn't exist or there's some other problem
	if err != nil {

		// check if error is caused by nonexistend package
		// if no, return error
		if err.Error() != "redis: nil" {
			return nil, err
		}

		// error is caused by nonexistent package
		// fetch package
		req, err := http.NewRequest("GET", options.UpstreamAddress+originalPath, nil)

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
			options.RedisPrefix+name,
			pkg,
			options.RedisExpireTimeout,
		).Result()
		if err != nil {
			return nil, err
		}
	}

	// replace tarball urls
	// FIXME: unmarshall and replace only necessary fields
	// convertedPkg := strings.ReplaceAll(string(pkg), options.ReplaceAddress, options.StaticServerAddress)

	return []byte(pkg), nil
}

// ListMetadata returns list of all cached packages
func (proxy Proxy) ListMetadata() ([]string, error) {
	options, err := proxy.GetOptions()
	if err != nil {
		return nil, err
	}

	metadata, err := proxy.RedisClient.Keys(options.RedisPrefix + "*").Result()
	if err != nil {
		return nil, err
	}

	deprefixedMetadata := make([]string, 0)
	for _, record := range metadata {
		deprefixedMetadata = append(deprefixedMetadata, strings.Replace(record, options.RedisPrefix, "", 1))
	}

	return deprefixedMetadata, nil
}

// PurgeMetadata deletes all cached packages
func (proxy Proxy) PurgeMetadata() error {
	options, err := proxy.GetOptions()
	if err != nil {
		return err
	}

	metadata, err := proxy.RedisClient.Keys(options.RedisPrefix + "*").Result()
	if err != nil {
		return err
	}

	for _, record := range metadata {
		_, err := proxy.RedisClient.Del(record).Result()
		if err != nil {
			return err
		}
	}

	return nil
}
