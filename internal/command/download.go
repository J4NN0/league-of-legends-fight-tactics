package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func (c *Controller) DownloadCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "download",
		Aliases: []string{"d"},
		Short:   "download and update a specific league of legends champion (name must not to contain spaces)",
		Args:    cobra.ExactArgs(1),
		Run:     c.download,
	}
}

func (c *Controller) DownloadAllCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "download_all",
		Aliases: []string{"a", "-da"},
		Short:   "download and update all league of legends champion",
		Args:    cobra.ExactArgs(0),
		Run:     c.downloadAll,
	}
}

func (c *Controller) download(cmd *cobra.Command, args []string) {
	championName := strings.ToLower(args[0])

	err := c.fetchChampion(championName)
	if err != nil {
		cmd.PrintErr(err)
		os.Exit(-1)
	}
}

func (c *Controller) fetchChampion(championName string) error {
	c.log.Printf("Fetching %s ...", championName)

	championData, err := c.riotClient.GetLoLChampion(championName)
	if err != nil {
		return fmt.Errorf("fetching league of legends champions: %v", err)
	}

	err = c.storeChampionToYMLFile(championData)
	if err != nil {
		return fmt.Errorf("could not store %s champion data: %v", championName, err)
	}

	c.log.Printf("%s successfully stored", championName)

	return nil
}

func (c *Controller) downloadAll(cmd *cobra.Command, args []string) {
	err := c.fetchAllChampions()
	if err != nil {
		cmd.PrintErr(err)
		os.Exit(-1)
	}
}

func (c *Controller) fetchAllChampions() error {
	c.log.Printf("Fetching all league of legends champions ...\n")

	ddChampions, err := c.riotClient.GetAllLoLChampions()
	if err != nil {
		return fmt.Errorf("fetching all league of legends champions: %v", err)
	}

	for _, champion := range ddChampions {
		err = c.storeChampionToYMLFile(champion)
		if err != nil {
			c.log.Warningf("Could not store %s champion data: %v", champion.ChampionData.Name, err)
		} else {
			c.log.Printf("%s successfully stored", champion.ChampionData.Name)
		}
	}

	return nil
}
