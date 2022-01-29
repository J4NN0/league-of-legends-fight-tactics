package main

import (
	"flag"
	"fmt"
	"league-of-legends-fight-tactics/internal/loltactics"
	"league-of-legends-fight-tactics/pkg/yml"
	"os"
	"path/filepath"
	"strings"
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
	fmt.Printf("[-] Reading %s champion data ...\n", c1Name)
	c1, err := yml.LoadLoLChampion(c1Name)
	if err != nil {
		fmt.Printf("Error loading champion: %v", err)
		return
	}
	// fmt.Printf("%+v\n", c1)

	fmt.Printf("[-] Reading %s champion data ...\n", c2Name)
	c2, err := yml.LoadLoLChampion(c2Name)
	if err != nil {
		fmt.Printf("Error loading champion: %v", err)
		return
	}
	// fmt.Printf("%+v\n", c2)

	loltactics.Fight(c1, c2)
}

func allChampionsFight() {
	var champions []string

	err := filepath.Walk(yml.BaseChampionPath, func(path string, info os.FileInfo, err error) error {
		champions = append(champions, strings.TrimSuffix(filepath.Base(path), filepath.Ext(filepath.Base(path))))
		return nil
	})
	if err != nil {
		panic(err)
	}

	champions = champions[1:]

	for _, c1 := range champions {
		for _, c2 := range champions {
			if c1 != c2 {
				fightChampion(c1, c2)
			}
		}
	}
}
