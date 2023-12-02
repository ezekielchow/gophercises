package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"urlshortener/helpers"
)

func main() {

	filePath := flag.String("filePath", "", "file path to map paths for redirection")
	dsn := flag.String("dsn", "", "dsn to sqlite file")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := helpers.MapHandler(pathsToUrls, mux)

	fileBytes, err := readFile(*filePath)
	if err != nil {
		panic(err)
	}

	var handler http.HandlerFunc

	fileExtension := path.Ext(*filePath)
	switch fileExtension {
	case ".yaml", ".yml":
		handler, err = helpers.YAMLHandler([]byte(fileBytes), mapHandler)
		if err != nil {
			panic(err)
		}
	case ".json":
		handler, err = helpers.JSONHandler([]byte(fileBytes), mapHandler)
		if err != nil {
			panic(err)
		}
	default:
		handler = mapHandler
	}

	if *dsn != "" {
		handler, err = helpers.DBHandler(*dsn, mapHandler)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func readFile(fileName string) ([]byte, error) {
	if fileName == "" {
		return nil, nil
	}

	fileByte, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return fileByte, nil
}
