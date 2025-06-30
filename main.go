package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <tracking code>")

		return
	}

	trackingCode := os.Args[1]
	fmt.Println("Checking package:", trackingCode)

	// URL to make the HTTP request to
	url := fmt.Sprintf("https://foxpost.hu/csomagkovetes/?code=%s", trackingCode)

	// Make the GET request
	resp, _ := http.Get(url)
	defer func() { _ = resp.Body.Close() }()

	// Parse the HTML response
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var results []struct {
		Title string
		Date  string
		Desc  string
	}

	// Extract the updates from the HTML document
	updates := getNodeWithClass(doc, "ul", "parcel-status-items__list") // Find the <ul> with class "parcel-status-items__list"
	if updates != nil {
		for item := updates.FirstChild; item != nil; item = item.NextSibling {
			if item.Type == html.ElementNode && item.Data == "li" {
				titleNode := getSubItemContentWithClass(item, "parcel-status-items__list-item-title")
				dateNode := getSubItemContentWithClass(item, "parcel-status-items__list-item-date")
				descNode := getSubItemContentWithClass(item, "parcel-status-items__list-item-description")

				if titleNode != "" && descNode != "" {
					results = append(results, struct {
						Title string
						Date  string
						Desc  string
					}{
						Title: titleNode,
						Date:  dateNode,
						Desc:  descNode,
					})
				}
			}
		}
	} else {
		fmt.Println("No updates found for the provided tracking code.")
		return
	}

	// Summarize the results
	if len(results) > 0 {
		fmt.Println("Found", len(results), "updates, results:")
		fmt.Println()
		for _, r := range results {
			fmt.Printf("Status: %s\n", r.Title)
			fmt.Printf("Date: %s\n", r.Date)
			fmt.Printf("Description: %s\n", r.Desc)
			fmt.Println("-----------------------------")
		}
	} else {
		fmt.Println("No tracking information found for the provided code.")
	}
}

func getNodeWithClass(n *html.Node, elementType string, class string) *html.Node {
	if n.Type == html.ElementNode && n.Data == elementType {
		for _, attr := range n.Attr {
			if attr.Key == "class" && hasClass(attr.Val, class) {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := getNodeWithClass(c, elementType, class); result != nil {
			return result
		}
	}

	return nil
}

func hasClass(classList, class string) bool {
	for _, p := range strings.Split(classList, " ") {
		if p == class {
			return true
		}
	}

	return false
}

func getSubItemContentWithClass(n *html.Node, class string) string {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" && hasClass(attr.Val, class) {
				return getTextContent(n)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := getSubItemContentWithClass(c, class); result != "" {
			return result
		}
	}

	return ""
}

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if text := getTextContent(c); text != "" {
			return text
		}
	}

	return ""
}
