package main

import (
	"fmt"

	"github.com/silvergama/go_web_scraping/webscraping"
)

// URLToProcess armazena todas a urls para serem processadas
var URLToProcess = []string{
	"urls",
}

func main() {
	// ini := time.Now()
	r := make(chan webscraping.Result, 3)
	go webscraping.ScrapListURL(URLToProcess, r)
	// time.Sleep(5 * time.Second)
	for url := range r {
		fmt.Println(url)
	}
	// fmt.Printf("(Took %g secs)", time.Since(ini).Seconds())
}
