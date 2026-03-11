package main

import (
	"fmt"
	"log"

	"www.github.com/maxbrt/colibri/internal/rss"
)

func main() {
	sources, err := rss.ReadSources("./feeds")
	if err != nil {
		log.Fatalf("%s", err)
	}
	for _, s := range sources {
		fmt.Println(s.ID)
		fmt.Println(s.Name)
		fmt.Println(s.URL)
		fmt.Println(s.Category)
	}
}
