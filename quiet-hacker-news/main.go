package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"quiet-hacker-news/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func getStory(id int, index int, ch chan<- item) {
	var client hn.Client

	i, err := client.GetItem(id)
	if err != nil {
		log.Default().Println("Error getting item")
		return
	}
	parsed := parseHNItem(i)
	if isStoryLink(parsed) {
		parsed.index = index
		ch <- parsed
		return
	}
	log.Default().Println("Not a story")
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		var stories []item

		rounds := 0
		storyCh := make(chan item)

		for len(stories) < numStories {

			itemsToQuery := int(math.Ceil(1.25 * float64(numStories)))

			for i := rounds * itemsToQuery; i < (rounds+1)*itemsToQuery; i++ {
				if i > len(ids)-1 {
					break
				}

				go func() {
					getStory(ids[i], i, storyCh)
				}()

			}

			for i := 0; i < itemsToQuery; i++ {
				stories = append(stories, <-storyCh)
			}

			rounds++
		}

		sort.Slice(stories, func(i, j int) bool {
			return stories[i].index < stories[j].index
		})

		data := templateData{
			Stories: stories[0:numStories],
			Time:    time.Since(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host  string
	index int
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
