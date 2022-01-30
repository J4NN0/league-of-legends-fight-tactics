package main

import (
	"flag"
	"fmt"
	"league-of-legends-fight-tactics/internal/lolchampion"
	"league-of-legends-fight-tactics/internal/loltactics"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	all := flag.Bool("all", false, "generate all fights tactics (default false)")
	c1Name := flag.String("c1", "", "first champion name")
	c2Name := flag.String("c2", "", "second champion name")
	flag.Parse()

	if *c1Name != "" && *c2Name != "" {
		fightChampion(*c1Name, *c2Name)
	} else if *all {
		allChampionsFight()
	} else {
		fmt.Printf("Usage: main.go -c1 champion1 -c2 champion2\n")
		flag.PrintDefaults()
	}
}

func fightChampion(c1Name, c2Name string) {
	fmt.Printf("[-] Loading %s vs %s champions data ...\n", c1Name, c2Name)

	c1, err := lolchampion.Load(c1Name)
	if err != nil {
		fmt.Printf("Error loading champion: %v", err)
		return
	}
	// fmt.Printf("%+v\n", c1)

	c2, err := lolchampion.Load(c2Name)
	if err != nil {
		fmt.Printf("Error loading champion: %v", err)
		return
	}
	// fmt.Printf("%+v\n", c2)

	loltactics.Fight(c1, c2)
}

func allChampionsFight() {
	var wg sync.WaitGroup
	var champions []string

	err := filepath.Walk(lolchampion.BaseChampionPath, func(path string, info os.FileInfo, err error) error {
		champions = append(champions, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
		return nil
	})
	if err != nil {
		panic(err)
	}

	champions = champions[1:]

	for _, c1 := range champions {
		for _, c2 := range champions {
			if c1 != c2 {
				wg.Add(1)
				c1 := c1
				c2 := c2
				go func() {
					defer wg.Done()
					fightChampion(c1, c2)
				}()
			}
		}
	}

	wg.Wait()
}
