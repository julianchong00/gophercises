package urlshortener

import (
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// Return function which implements the http.HandlerFunc method signature
	// i.e. any function with a ResponseWriter and Request argument
	return func(rw http.ResponseWriter, r *http.Request) {
		// If we can match path, redirect to it
		path := r.URL.Path
		if destination, ok := pathsToUrls[path]; ok {
			http.Redirect(rw, r, destination, http.StatusFound)
			return
		}
		// Otherwise call fallback handler
		fallback.ServeHTTP(rw, r)
	}
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
func YAMLHandler(yml string, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

// Capitalise first letter of fields to make them visible to entire program
type PathURL struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// Unmarshal YAML bytes into specified struct
func parseYAML(yml string) ([]PathURL, error) {
    // Read yaml file
	yamlFile, err := os.ReadFile(yml)
	if err != nil {
		return nil, err
	}

	var pathsToUrls []PathURL
    // Unmarshal yaml byte array into a list of PathURLs
	err = yaml.Unmarshal(yamlFile, &pathsToUrls)
	if err != nil {
		return nil, err
	}

	return pathsToUrls, nil
}

func buildMap(yml []PathURL) map[string]string {
	pathMap := make(map[string]string)
	for _, pu := range yml {
		pathMap[pu.Path] = pu.Url
	}

	return pathMap
}
