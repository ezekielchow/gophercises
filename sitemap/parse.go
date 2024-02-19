package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	Run()
}

func Run() {
	fileName := flag.String("file", "", "the name of the html file to parse")
	flag.Parse()

	links, err := ParseHtml(*fileName)
	if err != nil {
		log.Fatalf("Error parsing html file: %s", err)
	}

	for _, link := range links {
		log.Printf("Link: %v", link)
	}
}

type Link struct {
	Href string
	Text string
}

func parseLink(n *html.Node) Link {
	link := Link{}

	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			break
		}
	}

	link.Text = text(n)

	return link
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var ret []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}

	return ret
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = ret + text(c)
	}

	return strings.Join(strings.Fields(ret), " ")
}

func ParseHtml(fileName string) ([]Link, error) {

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("Error opening file: %s", err)
		return nil, err
	}

	doc, err := html.Parse(file)

	if err != nil {
		log.Fatalf("Error parsing html: %s", err)
		return nil, err
	}

	nodes := linkNodes(doc)

	links := []Link{}

	for _, node := range nodes {
		links = append(links, parseLink(node))
	}

	return links, nil
}
