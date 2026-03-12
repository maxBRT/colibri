package main

import (
	"fmt"
	"log"
	"sync"

	"www.github.com/maxbrt/colibri/internal/rss"
)

func main() {
	sources, err := rss.ReadSources("./sources")
	if err != nil {
		log.Fatalf("%s", err)
	}

	var wg sync.WaitGroup

	for _, s := range sources {
		wg.Add(1)

		go func(s rss.Source) {
			defer wg.Done()

			posts, err := rss.FetchAndParse(s)
			if err != nil {
				log.Fatalf("%s", err)
			}
			for _, p := range posts {
				fmt.Printf("%s, %s\n", p.Title, p.Description)
			}
		}(s)
	}
	wg.Wait()
}
