package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
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
			game.Price = parsePrice(price)
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

		filename := strings.Trim(baseURL[22:], "/") + ".json"
		f, err := os.Create(filename)
		utility.Check(err)
		defer f.Close()
		err = ioutil.WriteFile(filename, gameJSON, 0644)
		utility.Check(err)
	}
}

// parse price into list of floats
func parsePrice(price string) []float64 {
	r, _ := regexp.Compile("[+-]?([0-9]*[,])?[0-9]+")

	var prices []float64

	fmt.Println("Price: ", r.FindAllString(price, -1))
	for _, p := range r.FindAllString(price, -1) {
		p = strings.Replace(p, ",", ".", -1)
		if s, err := strconv.ParseFloat(p, 32); err == nil {
			prices = append(prices, s)
		}
	}

	return prices

}
