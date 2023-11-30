package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"urlshortener/helpers"
)

func main() {

	ymlFile := flag.String("yaml", "", "yaml file to map paths for redirection")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := helpers.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := readYaml(*ymlFile)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := helpers.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Starting the server on :8080")
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

func readYaml(fileName string) ([]byte, error) {
	if fileName == "" {
		return nil, nil
	}

	yaml, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return yaml, nil
}
