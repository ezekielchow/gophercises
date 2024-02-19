package main

import (
	"testing"
)

func TestParsingHtml(t *testing.T) {

	testCases := []struct {
		fileName string
		expected []Link
	}{
		{
			fileName: "ex1.html",
			expected: []Link{
				{
					Href: "/other-page",
					Text: "A link to another page",
				},
			},
		},
		{
			fileName: "ex2.html",
			expected: []Link{
				{
					Href: "https://www.twitter.com/joncalhoun",
					Text: "Check me out on twitter",
				},
				{
					Href: "https://github.com/gophercises",
					Text: "Gophercises is on Github!",
				},
			},
		},
		{
			fileName: "ex3.html",
			expected: []Link{
				{
					Href: "#",
					Text: "Login",
				},
				{
					Href: "/lost",
					Text: "Lost? Need help?",
				},
				{
					Href: "https://twitter.com/marcusolsson",
					Text: "@marcusolsson",
				},
			},
		},
		{
			fileName: "ex4.html",
			expected: []Link{
				{
					Href: "/dog-cat",
					Text: "dog cat",
				},
			},
		},
	}

	for _, testCase := range testCases {
		links, err := ParseHtml(testCase.fileName)

		if err != nil || len(links) != len(testCase.expected) {
			t.Errorf("Error parsing html file: %s", err)
		}

		for i, link := range links {
			if testCase.expected[i] != link {
				t.Errorf("Expected %v but got %v", testCase.expected[i], link)
			}
		}
	}
}
