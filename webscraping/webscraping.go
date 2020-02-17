package webscraping

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

/*
 * Result is a struct of result
 */
type Result struct {
	title       string
	description string
	price       string
}

/*
 * FormatResutl is a format return
 */
func (r Result) FormatResult() string {
	return r.title
}

func getFirstElementByClass(htmlParsed *html.Node, elm, className string) *html.Node {
	for m := htmlParsed.FirstChild; m != nil; m = m.NextSibling {
		if m.Data == elm && hasClass(m.Attr, className) {
			return m
		}
		r := getFirstElementByClass(m, elm, className)
		if r != nil {
			return r
		}
	}
	return nil
}

func hasClass(attrs []html.Attribute, className string) bool {
	for _, attr := range attrs {
		if attr.Key == "class" && strings.Contains(attr.Val, className) {
			return true
		}
	}
	return false
}

func getFirstTextNode(htmlParsed *html.Node) *html.Node {
	if htmlParsed == nil {
		return nil
	}

	for m := htmlParsed.FirstChild; m != nil; m = m.NextSibling {
		if m.Type == html.TextNode {
			return m
		}
		r := getFirstTextNode(m)
		if r != nil {
			return r
		}
	}
	return nil
}

/*
 * Scrap is a parser of data
 */
func Scrap(url string, ch chan Result) {
	var result Result

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error:  %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error:  %v", err)
		}
	}()

	body := resp.Body
	htmlParser, err := html.Parse(body)
	if err != nil {
		log.Printf("Error:  %v", err)
	}

	title := getFirstTextNode(getFirstElementByClass(htmlParser, "h1", "title__title"))
	if title != nil {
		result.title = title.Data
	} else {
		log.Printf("Scrap error: Can't find title. title: %v", title)
	}

	description := getFirstTextNode(getFirstElementByClass(htmlParser, "p", "description__text"))
	if description != nil {
		result.description = description.Data
	} else {
		log.Printf("Scrap error: Can't find description. description: %v", description)
	}

	price := getFirstTextNode(getFirstElementByClass(htmlParser, "h3", "price__price-info"))
	if price != nil {
		result.price = price.Data
	} else {
		log.Printf("Scrap error: Can't find price. Price: %v", price)
	}

	ch <- result
}

func ScrapListURL(URLSet []string, resultCh chan Result) {
	ch := make(chan Result)
	var res Result
	for i, url := range URLSet {
		go Scrap(url, ch)
		res = <-ch
		resultCh <- res
		fmt.Println("enviei", i)
	}

	defer close(resultCh)
}
