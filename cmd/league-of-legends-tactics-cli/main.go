package main

import (
	"context"
	"errors"
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
	var tacticsAll, fetchAll bool
	var fetch string
	var championsName *cli.StringSlice

	ctx := context.Background()

	logger := log.New(appName)
	riotClient := riot.NewClient(logger, &http.Client{})
	fightTactics := lol.New(logger)

	ctrl := controller.New(logger, riotClient, fightTactics)

	app := &cli.App{
		Name:    "loltactics",
		Usage:   "League of Legends Tactics",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "tactics",
				Aliases:     []string{"t"},
				Value:       nil,
				Usage:       "league of legends champions name",
				Required:    false,
				Destination: championsName,
				Action: func(context *cli.Context, i []string) error {
					c1Name := i[0]
					c2Name := i[1]
					if c1Name != "" && c2Name != "" {
						ctrl.ChampionsFight(strings.ToLower(c1Name), strings.ToLower(c2Name))
						return nil
					}
					return errors.New("champion name is empty")
				},
			},
			&cli.BoolFlag{
				Name:        "tacticsall",
				Aliases:     []string{"ta"},
				Value:       false,
				Usage:       "generate all fight tactics",
				Required:    false,
				Destination: &tacticsAll,
				Action: func(context *cli.Context, b bool) error {
					if b {
						ctrl.AllChampionsFight()
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "fetch",
				Aliases:     []string{"f"},
				Value:       "",
				Usage:       "fetch and update a specific league of legends champion (name must not to contain spaces)",
				Required:    false,
				Destination: &fetch,
				Action: func(context *cli.Context, s string) error {
					if s != "" {
						ctrl.FetchChampion(strings.ToLower(fetch))
						return nil
					}
					return errors.New("champion name is empty")
				},
			},
			&cli.BoolFlag{
				Name:        "fetchall",
				Aliases:     []string{"fa"},
				Value:       false,
				Usage:       "fetch and update all league of legends champions",
				Required:    false,
				Destination: &fetchAll,
				Action: func(context *cli.Context, b bool) error {
					if b {
						ctrl.FetchAllChampions()
					}
					return nil
				},
			},
		},
		Action: func(context *cli.Context) error {
			return errors.New("no input provided")
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		_ = cli.ShowAppHelp(&cli.Context{
			Context: ctx,
			App:     app,
		})
		logger.Fatalf("%v", err)
	}
}
