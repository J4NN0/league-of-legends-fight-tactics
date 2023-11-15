package main

import (
	"fmt"
	"github.com/J4NN0/league-of-legends-fight-tactics/internal/command"
	"net/http"
	"os"

	"github.com/J4NN0/league-of-legends-fight-tactics/internal/config"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/logger"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/lol"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/riot"

	"github.com/spf13/cobra"
)

const appName string = "lol-tactics"

func main() {
	log := logger.New(appName)

	appConfig, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("config reading failed: %v", err)
		return
	}

	riotClient := riot.NewClient(log, &http.Client{}, appConfig.RiotAPIKey, appConfig.LoLRegion)
	lolTactics := lol.NewTactics(log)

	ctrl := command.New(log, riotClient, lolTactics)

	rootCmd := &cobra.Command{
		Use:   "loltactics",
		Short: "league of legends fight tactics tool",
	}
	rootCmd.AddCommand(ctrl.FightCommand())
	rootCmd.AddCommand(ctrl.TacticsCommand())
	rootCmd.AddCommand(ctrl.DownloadCommand())
	rootCmd.AddCommand(ctrl.DownloadAllCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
