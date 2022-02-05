package main

import (
	"flag"
	"fmt"
	"league-of-legends-fight-tactics/internal/controller"
	"league-of-legends-fight-tactics/internal/log"
	"league-of-legends-fight-tactics/internal/lol"
	"league-of-legends-fight-tactics/internal/riot"
	"net/http"
	"strings"
)

const appName string = "lol-tactics"

func main() {
	all := flag.Bool("all", false, "generate all fights tactics (default false)")
	c1Name := flag.String("c1", "", "first champion name")
	c2Name := flag.String("c2", "", "second champion name")

	fetch := flag.String("fetch", "", "fetch and update league of legends champion")
	fetchAll := flag.Bool("fetchall", false, "fetch and update all league of legends champions")

	flag.Parse()

	logger := log.New(appName)
	riotClient := riot.NewApiClient(logger, &http.Client{})
	fightTactics := lol.New(logger)

	ctrl := controller.New(logger, riotClient, fightTactics)

	if *c1Name != "" && *c2Name != "" {
		ctrl.ChampionsFight(strings.ToLower(*c1Name), strings.ToLower(*c2Name))
	} else if *all {
		ctrl.AllChampionsFight()
	} else if *fetch != "" {
		ctrl.FetchChampion(strings.ToLower(*fetch))
	} else if *fetchAll {
		ctrl.FetchAllChampions()
	} else {
		fmt.Printf("Usage: main.go -c1 champion1 -c2 champion2\n")
		flag.PrintDefaults()
	}
}
