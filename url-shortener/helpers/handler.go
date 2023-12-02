package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for path, url := range pathsToUrls {
			fmt.Println(r.URL.Path)
			if r.URL.Path == path {
				http.Redirect(w, r, url, http.StatusFound)
				return
			}
		}

		fallback.ServeHTTP(w, r)
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type PathUrl struct {
	Path string `yaml:"path" json:"path,omitempty"`
	URL  string `yaml:"url" json:"url,omitempty"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var pathUrls []PathUrl

	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}

	pathsToUrl := pathArrayToMap(pathUrls)

	return MapHandler(pathsToUrl, fallback), nil
}

func JSONHandler(jsonByte []byte, fallback http.Handler) (http.HandlerFunc, error) {

	pathsUrls := []PathUrl{}

	dec := json.NewDecoder(strings.NewReader(string(jsonByte)))
	for {
		if err := dec.Decode(&pathsUrls); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	pathsToUrl := pathArrayToMap(pathsUrls)

	return MapHandler(pathsToUrl, fallback), nil

}

func pathArrayToMap(arr []PathUrl) map[string]string {
	paths := make(map[string]string)

	for _, v := range arr {
		paths[v.Path] = v.URL
	}

	return paths
}
