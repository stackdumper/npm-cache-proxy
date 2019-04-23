package proxy

import (
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// GetCachedPath returns cached upstream response for a given url path.
func (proxy Proxy) GetCachedPath(options Options, path string, request *http.Request) ([]byte, error) {
	key := options.DatabasePrefix + path

	// get package from database
	pkg, err := proxy.Database.Get(key)

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

		req.Header = request.Header
		req.Header.Set("Accept-Encoding", "gzip")

		res, err := proxy.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if res.Header.Get("Content-Encoding") == "gzip" {
			zr, err := gzip.NewReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			res.Body = zr
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		pkg = string(body)

		// // TODO: avoid calling MustCompile every time
		// // find "dist": "https?://.*/ and replace to "dist": "{localurl}/
		// pkg = regexp.MustCompile(`(?U)"tarball":"https?://.*/`).ReplaceAllString(string(body), `"dist": "http://localhost:8080/`)

		// save to redis
		err = proxy.Database.Set(key, pkg, options.DatabaseExpiration)
		if err != nil {
			return nil, err
		}
	}

	return []byte(pkg), nil
}

// ListCachedPaths returns list of all cached url paths.
func (proxy Proxy) ListCachedPaths(options Options) ([]string, error) {
	metadata, err := proxy.Database.Keys(options.DatabasePrefix)
	if err != nil {
		return nil, err
	}

	deprefixedMetadata := make([]string, 0)
	for _, record := range metadata {
		deprefixedMetadata = append(deprefixedMetadata, strings.Replace(record, options.DatabasePrefix, "", 1))
	}

	return deprefixedMetadata, nil
}

// PurgeCachedPaths deletes all cached url paths.
func (proxy Proxy) PurgeCachedPaths(options Options) error {
	metadata, err := proxy.Database.Keys(options.DatabasePrefix)
	if err != nil {
		return err
	}

	for _, record := range metadata {
		err := proxy.Database.Delete(record)
		if err != nil {
			return err
		}
	}

	return nil
}
