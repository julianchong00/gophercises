package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents a link (<a href="...">) in an HTML document.
type Link struct {
	Href string
	Text string
}

// Parse will take in an HTML document and will return a slice
// of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}

// DFS function to extract text from <a> tags since they can also contain other nodes
// or even comments (which we want to ignore)
func text(n *html.Node) string {
	// Base case to return just the text from text nodes
	if n.Type == html.TextNode {
		return n.Data
	}
	// Check to return empty string in the case of any tags that are not ElementNodes.
	// Also removes comments because CommentNode is not being checked for
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}

// DFS function to get all link nodes from HTML and return a list
func linkNodes(n *html.Node) []*html.Node {
	// Base case - checks if node is an ElementNode and if it is an <a> tag
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// Another variadic parameter:
		// ... after linkNodes(c) basically expands the list that is returned
		// from deeper iterations
		ret = append(ret, linkNodes(c)...)
	}

	return ret
}
