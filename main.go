package main

import (
	"fmt"
	"sync"
	"github.com/duolok/blue-jay/engine"
)


func main() {
	scraperNames := engine.LoadScrapers()
	for _, name := range scraperNames {
		fmt.Println(name)
	}

	game, err := engine.LoadLastSearch("games.csv")
	if err != nil {
		fmt.Println("Error: ", err, game)
		return
	}
	fmt.Println(game)

	var wg sync.WaitGroup
	wg.Add(1)

	go engine.Search(scraperNames, "hollow knight", &wg)
	wg.Wait()

	fmt.Println("Scraping completed for all games.")
}

