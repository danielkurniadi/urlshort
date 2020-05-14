package urlshort

import (
	yaml "gopkg.in/yaml.v2"
	"net/http"
)

// YAMLPathURL represent YAML file which
// maps a URI Path to another valid URL
type YAMLPathURL struct {
	URL  string `yaml:"url"`
	Path string `yaml:"path"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
	return handler
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var handler func(http.ResponseWriter, *http.Request)

	yamlPathURLs, err := parseYAML(yml)
	if err != nil {
		return handler, err
	}

	pathsToUrls := makeYAMLToMap(yamlPathURLs)
	handler = MapHandler(pathsToUrls, fallback)
	return handler, err
}

// parseYAML take YAML file content in []byte (array) and map the content
// into a YAMLPathURL struct above
func parseYAML(yml []byte) ([]YAMLPathURL, error) {
	var yamlPathURLs []YAMLPathURL
	err := yaml.Unmarshal(yml, &yamlPathURLs)

	return yamlPathURLs, err
}

// makeYAMLToMap convert array of struct, []YAMLPathURL, into
// native map of Path to URLs
func makeYAMLToMap(yamlPathURLs []YAMLPathURL) map[string]string {
	var pathsToUrls = make(map[string]string)
	for _, yamlPathURL := range yamlPathURLs {
		pathsToUrls[yamlPathURL.Path] = yamlPathURL.URL
	}
	return pathsToUrls
}
