package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/J4NN0/league-of-legends-fight-tactics/internal/config"
	"github.com/J4NN0/league-of-legends-fight-tactics/internal/controller"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/logger"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/lol"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/riot"
	"github.com/urfave/cli/v2"
)

var noInputErr = errors.New("no input provided")

const appName string = "lol-tactics"

func main() {
	ctx := context.Background()

	// Logger
	log := logger.New(appName)

	// Config
	appConfig, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("config reading failed: %v", err)
		return
	}

	// Riot Client
	riotClient := riot.NewClient(log, &http.Client{}, appConfig.RiotAPIKey, appConfig.LoLRegion)

	// LoL Tactics
	lolTactics := lol.NewTactics(log)

	// Controller
	ctrl := controller.New(log, riotClient, lolTactics)

	// CLI App
	var tactics, downloadAll bool
	var download string
	var championsName *cli.StringSlice
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
					if c1Name == "" || c2Name == "" {
						return errors.New("champion name is empty")
					}

					err = ctrl.ChampionsFight(strings.ToLower(c1Name), strings.ToLower(c2Name))
					if err != nil {
						return fmt.Errorf("champion fight failed: %w", err)
					}
					return nil
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
					if !b {
						return nil
					}

					err = ctrl.AllChampionsFight()
					if err != nil {
						return fmt.Errorf("all champions fight failed: %w", err)
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
					if s == "" {
						return errors.New("champion name is empty")
					}

					err = ctrl.FetchChampion(strings.ToLower(download))
					if err != nil {
						return fmt.Errorf("fetch champion failed: %w", err)
					}
					return nil
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
					if !b {
						return nil
					}

					err = ctrl.FetchAllChampions()
					if err != nil {
						return fmt.Errorf("fetch all champions failed: %w", err)
					}
					return nil
				},
			},
		},
		Action: func(context *cli.Context) error {
			if !isInputProvided(tactics, downloadAll, download, championsName) {
				return noInputErr
			}
			return nil
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		if errors.Is(err, noInputErr) {
			_ = cli.ShowAppHelp(&cli.Context{
				Context: ctx,
				App:     app,
			})
		}
		log.Fatalf("%v", err)
	}
}

func isInputProvided(tactics, downloadAll bool, download string, championsName *cli.StringSlice) bool {
	if !tactics && !downloadAll && download == "" && championsName == nil {
		return false
	}
	return true
}
