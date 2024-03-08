package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func getValidUrl(baseUrl string, url string) string {

	if strings.HasPrefix(url, baseUrl) {
		return url
	}

	if strings.HasPrefix(url, "http") && !strings.HasPrefix(url, baseUrl) {
		return ""
	}

	if strings.HasPrefix(url, "/") {
		return baseUrl + url
	}

	if strings.HasPrefix(url, "#") {
		return ""
	}

	if strings.Contains(url, "mailto") {
		return ""
	}

	return ""
}

func getLinksFromUrl(baseUrl string, url string) ([]Link, error) {

	validUrl := getValidUrl(baseUrl, url)

	if validUrl == "" {
		return []Link{}, nil
	}

	r, err := http.Get(validUrl)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	s := string(b)

	docLinks, err := ParseHtml(strings.NewReader(s))

	if err != nil {
		return nil, err
	}

	links := []Link{}
	for _, link := range docLinks {
		link.Href = getValidUrl(baseUrl, link.Href)

		if link.Href == "" {
			continue
		}

		links = append(links, link)
	}

	return links, nil
}

func hasDuplicate(links []Link, needle Link) bool {

	for _, link := range links {
		if link.Href == needle.Href {
			return true
		}
	}

	return false
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Loc string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func printXml(links []Link) {

	toXml := urlset{
		Xmlns: xmlns,
	}

	for _, link := range links {
		toXml.Urls = append(toXml.Urls, loc{Loc: link.Href})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}

	fmt.Println(enc)
}

func BuildSitemap() {

	url := flag.String("url", "https://www.calhoun.io", "the url to build a sitemap for")
	depth := flag.Int("depth", 3, "the depth to traverse the site")
	flag.Parse()

	allLinks := []Link{}
	allLinks = append(allLinks, Link{Href: *url, Text: "Home Page"})

	// Without verification
	links, err := getLinksFromUrl(*url, *url)
	uniques := []Link{}
	for _, newLink := range links {
		if !hasDuplicate(allLinks, newLink) {
			uniques = append(uniques, newLink)
		}
	}

	allLinks = append(allLinks, uniques...)

	if err != nil {
		log.Fatalf("Error getting links from url: %s", err)
	}

	i := 1

	for len(links) > 0 {

		if i > *depth {
			break
		}

		newLinks := []Link{}

		for i := 0; i < len(links); i++ {
			link := links[i]

			isValidUrl := getValidUrl(*url, link.Href)

			if isValidUrl == "" {
				continue
			}

			l, err := getLinksFromUrl(*url, links[i].Href)

			if err != nil {
				log.Printf("Error getting links from %s: %s", links[i].Href, err)
				continue
			}

			newLinks = append(newLinks, l...)
		}

		uniques := []Link{}
		for _, newLink := range newLinks {
			if !hasDuplicate(allLinks, newLink) {
				uniques = append(uniques, newLink)
			}
		}

		allLinks = append(allLinks, uniques...)
		links = uniques
		i++

		continue
	}

	printXml(allLinks)
}
