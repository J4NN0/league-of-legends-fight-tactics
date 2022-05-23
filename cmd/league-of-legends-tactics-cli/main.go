package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"league-of-legends-fight-tactics/internal/controller"
	"league-of-legends-fight-tactics/internal/log"
	"league-of-legends-fight-tactics/internal/lol"
	"league-of-legends-fight-tactics/internal/riot"
	"net/http"
	"os"
	"strings"
)

const appName string = "lol-tactics"

func main() {
	var all, fetchAll bool
	var c1Name, c2Name, fetch string

	logger := log.New(appName)
	riotClient := riot.NewClient(logger, &http.Client{})
	fightTactics := lol.New(logger)

	ctrl := controller.New(logger, riotClient, fightTactics)

	app := &cli.App{
		Name:    "loltactics",
		Usage:   "League of Legends Tactics",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "champion1",
				Aliases:     []string{"c1"},
				Value:       "",
				Usage:       "first league of legends champion name",
				Required:    false,
				Destination: &c1Name,
			},
			&cli.StringFlag{
				Name:        "champion2",
				Aliases:     []string{"c2"},
				Value:       "",
				Usage:       "second league of legends champion name",
				Required:    false,
				Destination: &c2Name,
			},
			&cli.BoolFlag{
				Name:        "all",
				Aliases:     []string{"a"},
				Value:       false,
				Usage:       "generate all fight tactics",
				Required:    false,
				Destination: &all,
			},
			&cli.StringFlag{
				Name:        "fetch",
				Aliases:     []string{"f"},
				Value:       "",
				Usage:       "fetch and update a specific league of legends champion (name must not to contain spaces)",
				Required:    false,
				Destination: &fetch,
			},
			&cli.BoolFlag{
				Name:        "fetchall",
				Aliases:     []string{"fa"},
				Value:       false,
				Usage:       "fetch and update all league of legends champions",
				Required:    false,
				Destination: &fetchAll,
			},
		},
		Action: func(c *cli.Context) error {
			if c1Name != "" && c2Name != "" {
				ctrl.ChampionsFight(strings.ToLower(c1Name), strings.ToLower(c2Name))
			} else if all {
				ctrl.AllChampionsFight()
			} else if fetch != "" {
				ctrl.FetchChampion(strings.ToLower(fetch))
			} else if fetchAll {
				ctrl.FetchAllChampions()
			} else {
				return fmt.Errorf("no flags provided")
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalf("%v", err)
	}
}
