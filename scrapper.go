package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"web-scrapper/utility"

	"github.com/gocolly/colly"
)

type Game struct {
	Name     string
	ShopName string
	Price    []float64
	Link     string
}

func main() {
	checkLinks()
	checkHistory()
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
			price := el.ChildText(".price-inner.game-price-current")
			game.Price = utility.ParsePrice(price)
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

func checkLinks() {
	dat, err := os.Open("games.txt")
	utility.Check(err)
	defer dat.Close()

	scanner := bufio.NewScanner(dat)
	games := []Game{}
	c := colly.NewCollector(colly.CacheDir("./ggdeals_cache"))
	initCollectro(c, &games)

	for scanner.Scan() {
		baseURL := scanner.Text()
		fmt.Println("Scraping: ", baseURL)
		c.Visit(baseURL)
		gameJSON, err := json.Marshal(games)
		utility.Check(err)

		filename := "./data/" + strings.Trim(baseURL[22:], "/") + ".json"

		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			f, err := os.Create(filename)
			utility.Check(err)
			defer f.Close()
		} else if err == nil {
			e := os.Rename(filename, "./data/"+strings.Trim(baseURL[22:], "/")+"_old"+".json")
			utility.Check(e)
			f, err := os.Create(filename)
			utility.Check(err)
			defer f.Close()
		} else {
			utility.Check(err)
		}

		err = ioutil.WriteFile(filename, gameJSON, 0644)
		utility.Check(err)
	}
}

func checkHistory() {
	dat, err := os.Open("games.txt")
	utility.Check(err)
	defer dat.Close()
	scanner := bufio.NewScanner(dat)
	for scanner.Scan() {
		//filename := strings.Trim(scanner.Text()[22:], "/") + ".json"
		//content, err = ioutil.ReadFile("./data/" + filename)
		utility.Check(err)
	}

}
