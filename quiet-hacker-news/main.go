package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
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

func getStory(client hn.Client, id int, index int) (item, error) {
	i, err := client.GetItem(id)
	if err != nil {
		return item{}, err
	}
	parsed := parseHNItem(i)
	if isStoryLink(parsed) {
		parsed.index = index
		fmt.Println("received", index)
		return parsed, nil
	}
	return item{}, errors.New("not a story")
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
		for len(stories) < numStories {

			itemsToQuery := int(math.Ceil(1.25 * float64(numStories)))

			wg := sync.WaitGroup{}

			for i := rounds * itemsToQuery; i < (rounds+1)*itemsToQuery; i++ {
				if i > len(ids)-1 {
					break
				}

				wg.Add(1)
				go func() {
					defer wg.Done()
					story, err := getStory(client, ids[i], i)
					if err != nil {
						fmt.Println("error", err)
						return
					}
					stories = append(stories, story)
				}()

			}
			wg.Wait()

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
