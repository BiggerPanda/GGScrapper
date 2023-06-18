package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("Hello World")
	c := colly.NewCollector()
	initCollectro(c)
	c.Visit("https://gg.deals/game/satisfactory/")

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

	c.OnHTML(".d-flex.header-game-prices-content.active", func(e *colly.HTMLElement) {
		fmt.Println("Found something: ", e.ChildText(".price"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})
}
