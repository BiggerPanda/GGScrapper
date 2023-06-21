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
	games := []Game{}
	c := colly.NewCollector(colly.CacheDir("./ggdeals_cache"))
	initCollectro(c, &games)
	c.Visit(baseURL)

}

func initCollectro(c *colly.Collector, games *[]Game) {

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
		i := 1
		e.ForEach("div.similar-deals-container.items-with-top-border-desktop", func(_ int, el *colly.HTMLElement) {
			game := Game{}
			game.name = el.ChildAttr("div.relative.hoverable-box.d-flex.flex-wrap.flex-align-center.game-item.cta-full.item.game-deals-item.game-list-item.keep-unmarked-container", "data-game-name")
			game.shop = el.ChildAttr("div.relative.hoverable-box.d-flex.flex-wrap.flex-align-center.game-item.cta-full.item.game-deals-item.game-list-item.keep-unmarked-container", "data-shop-name")
			game.price = el.ChildText(".price-inner.game-price-current")
			link := el.ChildAttr(".full-link", "href")
			game.link = el.Request.AbsoluteURL(link)
			*games = append(*games, game)
			fmt.Println("Game: ", i)
			i++
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})
}
