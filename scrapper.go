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
	checkForInputFile()
	checkLinks()
	checkHistory()
}

func initCollector(c *colly.Collector, games *[]Game) {

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
	initCollector(c, &games)

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
			e := os.Rename(filename, "./data/"+strings.Trim(baseURL[22:], "/")+"_old.json")
			utility.Check(e)
			f, err := os.Create(filename)
			utility.Check(err)
			defer f.Close()
		} else {
			utility.Check(err)
		}

		err = ioutil.WriteFile(filename, gameJSON, 0644)
		utility.Check(err)
		games = []Game{}
	}
}

func checkHistory() {
	dat, err := os.Open("games.txt")
	utility.Check(err)
	defer dat.Close()
	scanner := bufio.NewScanner(dat)
	shouldNotify := false
	notifyGamesName := []string{}

	var gamesCurrent []Game
	var gamesOld []Game

	for scanner.Scan() {
		filenameCurrent := "./data/" + strings.Trim(scanner.Text()[22:], "/") + ".json"
		filenameOld := "./data/" + strings.Trim(scanner.Text()[22:], "/") + "_old.json"

		if _, err := os.Stat(filenameOld); errors.Is(err, os.ErrNotExist) {
			continue
		}

		contentCurrent, err := ioutil.ReadFile(filenameCurrent)
		utility.Check(err)
		contentOld, err := ioutil.ReadFile(filenameOld)
		utility.Check(err)

		err = json.Unmarshal(contentCurrent, &gamesCurrent)
		utility.Check(err)
		err = json.Unmarshal(contentOld, &gamesOld)
		utility.Check(err)

		for i := 0; i < len(gamesCurrent); i++ {
			if gamesCurrent[i].Name == gamesOld[i].Name {
				if gamesCurrent[i].Price[0] < gamesOld[i].Price[0] {
					fmt.Println("Game: ", gamesOld[i].Name, " is cheaper than before!")
					shouldNotify = true
					notifyGamesName = append(notifyGamesName, gamesOld[i].Name+" "+gamesOld[i].ShopName+" "+fmt.Sprintf("%.2f", gamesOld[i].Price[0])+" -> "+fmt.Sprintf("%.2f", gamesCurrent[i].Price[0]))

				} else if gamesCurrent[i].Price[0] >= gamesOld[i].Price[0] {
					fmt.Println("Game: ", gamesCurrent[i].Name, " is more expensive than before!")
					if err := os.Remove(filenameOld); err != nil {
						log.Fatal(err)
					}
					break
				}
			}
		}

		gamesCurrent = []Game{}
		gamesOld = []Game{}

		if shouldNotify {
			utility.Notify(notifyGamesName)
			shouldNotify = false
			notifyGamesName = []string{}
		}
	}
}

func checkForInputFile() {
	if _, err := os.Stat("games.txt"); errors.Is(err, os.ErrNotExist) {
		exapleLinks := []string{"https://gg.deals/game/satisfactory/",
			"https://gg.deals/game/diablo-iv/",
			"https://gg.deals/game/starfield/",
			"https://gg.deals/game/elden-ring/"}
		fmt.Println("File games.txt does not exist!")
		fmt.Println("Creating file games.txt")
		f, err := os.Create("games.txt")
		utility.Check(err)
		defer f.Close()
		fmt.Println("Please add links to games in games.txt")
		fmt.Println("Example: https://gg.deals/game/monster-hunter-world-iceborne/")
		fmt.Println("Adding example link to games.txt")
		for _, link := range exapleLinks {
			_, err := f.WriteString(link + "\n")
			utility.Check(err)
		}
		utility.Check(err)

	}
}
