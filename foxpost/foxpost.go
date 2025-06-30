package foxpost

import (
	"errors"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

// PackageUpdate .
type PackageUpdate struct {
	Title string
	Date  string
	Desc  string
}

const foxpostBaseURL = "https://foxpost.hu/csomagkovetes/?code="

// Logic .
type Logic struct {
	httpClient *http.Client
}

// New .
func New() *Logic {
	return &Logic{
		httpClient: &http.Client{},
	}
}

func (l *Logic) TrackPackage(trackingCode string) ([]PackageUpdate, error) {
	resp, err := l.httpClient.Get(foxpostBaseURL + trackingCode)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Parse the HTML response
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []PackageUpdate

	updates := getNodeWithClass(doc, "ul", "parcel-status-items__list") // Find the <ul> with class "parcel-status-items__list"
	if updates != nil {
		for item := updates.FirstChild; item != nil; item = item.NextSibling {
			if item.Type == html.ElementNode && item.Data == "li" {
				titleNode := getSubItemContentWithClass(item, "parcel-status-items__list-item-title")
				dateNode := getSubItemContentWithClass(item, "parcel-status-items__list-item-date")
				descNode := getSubItemContentWithClass(item, "parcel-status-items__list-item-description")

				if titleNode != "" && descNode != "" {
					results = append(results, PackageUpdate{
						Title: titleNode,
						Date:  dateNode,
						Desc:  descNode,
					})
				}
			}
		}
	} else {
		return nil, errors.New("no updates found for the provided tracking code")
	}

	return results, nil
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
