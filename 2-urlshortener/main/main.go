package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/julianchong00/urlshortener"
)

type Config struct {
	YamlFile string
	JsonFile string
}

const (
	DefaultYaml = "urlpath.yaml"
	DefaultJson = "urlpath.json"
)

func main() {
	// Set the default configuration
	config := Config{
		YamlFile: DefaultYaml,
		JsonFile: DefaultJson,
	}

	// Parse command line flags
	// yaml file flag
	flag.StringVar(
		&config.YamlFile,
		"yaml",
		DefaultYaml,
		"a yaml file to read url paths from",
	)

	// json file flag
	flag.StringVar(
		&config.JsonFile,
		"json",
		DefaultJson,
		"a json file to read url paths from",
	)
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortener.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshortener.YAMLHandler(config.YamlFile, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
