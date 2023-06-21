package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type Game struct {
	name  string
	shop  string
	price string
	link  string
}

func main() {
	baseURL := "https://gg.deals/game/satisfactory/"

	fmt.Println("Hello World")
	c := colly.NewCollector(colly.CacheDir("./ggdeals_cache"))
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

	c.OnHTML("div.offer-section.with-filters", func(e *colly.HTMLElement) {

		e.ForEach("div.similar-deals-container.items-with-top-border-desktop", func(_ int, el *colly.HTMLElement) {
			fmt.Println("Found Name: ", el.ChildAttr("div.relative.hoverable-box.d-flex.flex-wrap.flex-align-center.game-item.cta-full.item.game-deals-item.game-list-item.keep-unmarked-container", "data-game-name"))
			fmt.Println("Found Shop: ", el.ChildAttr("div.relative.hoverable-box.d-flex.flex-wrap.flex-align-center.game-item.cta-full.item.game-deals-item.game-list-item.keep-unmarked-container", "data-shop-name"))
			fmt.Println("Found price: ", el.ChildText(".price-inner.game-price-current"))
			link := el.ChildAttr(".full-link", "href")
			fmt.Println("Found link: ", el.Request.AbsoluteURL(link))
			fmt.Println(" ")

		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})
}
