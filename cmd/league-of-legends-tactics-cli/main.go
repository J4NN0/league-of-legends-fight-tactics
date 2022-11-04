package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/J4NN0/league-of-legends-fight-tactics/internal/controller"
	"github.com/J4NN0/league-of-legends-fight-tactics/internal/log"
	"github.com/J4NN0/league-of-legends-fight-tactics/internal/lol"
	"github.com/J4NN0/league-of-legends-fight-tactics/internal/riot"
	"github.com/urfave/cli/v2"
)

const appName string = "lol-tactics"

func main() {
	var tactics, downloadAll bool
	var download string
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
				Name:        "fight",
				Aliases:     []string{"f"},
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
				Name:        "tactics",
				Aliases:     []string{"t"},
				Value:       false,
				Usage:       "generate all fight tactics",
				Required:    false,
				Destination: &tactics,
				Action: func(context *cli.Context, b bool) error {
					if b {
						ctrl.AllChampionsFight()
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "download",
				Aliases:     []string{"d"},
				Value:       "",
				Usage:       "download and update a specific league of legends champion (name must not to contain spaces)",
				Required:    false,
				Destination: &download,
				Action: func(context *cli.Context, s string) error {
					if s != "" {
						ctrl.FetchChampion(strings.ToLower(download))
						return nil
					}
					return errors.New("champion name is empty")
				},
			},
			&cli.BoolFlag{
				Name:        "download_all",
				Aliases:     []string{"da", "a"},
				Value:       false,
				Usage:       "download and update all league of legends champions",
				Required:    false,
				Destination: &downloadAll,
				Action: func(context *cli.Context, b bool) error {
					if b {
						ctrl.FetchAllChampions()
					}
					return nil
				},
			},
		},
		Action: func(context *cli.Context) error {
			if !isInputProvided(tactics, downloadAll, download, championsName) {
				return errors.New("no input provided")
			}
			return nil
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

func isInputProvided(tactics, downloadAll bool, download string, championsName *cli.StringSlice) bool {
	if tactics == false && downloadAll == false && download == "" && championsName == nil {
		return false
	}
	return true
}
