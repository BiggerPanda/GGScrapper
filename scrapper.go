package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Game struct {
	Name     string
	ShopName string
	Price    string
	Link     string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// read the base url from file
	dat, err := os.Open("games.txt")
	check(err)
	defer dat.Close()

	scanner := bufio.NewScanner(dat)

	for scanner.Scan() {
		baseURL := scanner.Text()
		fmt.Println("Scraping: ", baseURL)
		games := []Game{}
		c := colly.NewCollector(colly.CacheDir("./ggdeals_cache"))
		initCollectro(c, &games)
		c.Visit(baseURL)
		// err2 := beeep.Notify("Title", "Message body", "assets/information.png")
		// check(err2)
		gameJSON, err := json.Marshal(games)
		check(err)
		fmt.Println(string(gameJSON))
		gamename := baseURL[22:]
		filename := strings.Trim(gamename, "/") + ".json"
		f, err := os.Create(filename)
		check(err)
		defer f.Close()
		err3 := ioutil.WriteFile(filename, gameJSON, 0644)
		check(err3)
	}
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
			game.Name = el.ChildAttr("div.relative.hoverable-box.d-flex.flex-wrap.flex-align-center.game-item.cta-full.item.game-deals-item.game-list-item.keep-unmarked-container", "data-game-name")
			game.ShopName = el.ChildAttr("div.relative.hoverable-box.d-flex.flex-wrap.flex-align-center.game-item.cta-full.item.game-deals-item.game-list-item.keep-unmarked-container", "data-shop-name")
			game.Price = el.ChildText(".price-inner.game-price-current")
			link := el.ChildAttr(".full-link", "href")
			game.Link = el.Request.AbsoluteURL(link)
			*games = append(*games, game)
			fmt.Println("Game: ", i)
			i++
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})
}
