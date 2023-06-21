package main

import (
	"fmt"
	"log"
	"net/url"
	"path"

	"github.com/gocolly/colly"
)

func main() {
	baseURL := "https://gg.deals/game/satisfactory/"

	fmt.Println("Hello World")
	c := colly.NewCollector()
	initCollectro(c)
	c.Visit(baseURL)

}

func initCollectro(c *colly.Collector) {

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	c.OnHTML(".similar-deals-container.items-with-top-border-desktop", func(e *colly.HTMLElement) {
		u, _ := url.Parse("https://gg.deals")
		u.Path = path.Join(u.Path, e.ChildAttr("a", "href"))
		c.Visit(u.String())
		fmt.Println("Found something: ", u)
		fmt.Println("Found something: ", e.ChildText(".price"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})
}
