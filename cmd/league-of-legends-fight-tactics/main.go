package main

import (
	"fmt"
	"league-of-legends-fight-tactics/internal/loltactics"
	"league-of-legends-fight-tactics/pkg/yml"
)

func main() {
	championJhin, err := yml.LoadLoLChampion("Jhin")
	if err != nil {
		fmt.Printf("Error loading champion: %v", err)
		return
	}
	fmt.Printf("%+v\n", championJhin)

	championLucian, err := yml.LoadLoLChampion("Lucian")
	if err != nil {
		fmt.Printf("Error loading champion: %v", err)
		return
	}
	fmt.Printf("%+v\n", championLucian)

	loltactics.Fight(championJhin, championLucian)
}
