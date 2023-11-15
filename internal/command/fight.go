package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/J4NN0/league-of-legends-fight-tactics/internal/file"
	"github.com/spf13/cobra"
)

func (c *Controller) FightCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "fight",
		Aliases: []string{"f"},
		Short:   "league of legends champions name",
		Args:    cobra.ExactArgs(2),
		Run:     c.fight,
	}
}

func (c *Controller) fight(cmd *cobra.Command, args []string) {
	championName1 := strings.ToLower(args[0])
	championName2 := strings.ToLower(args[1])

	err := c.championsFight(championName1, championName2)
	if err != nil {
		cmd.PrintErr(err)
		os.Exit(-1)
	}
}

func (c *Controller) championsFight(championName1, championName2 string) error {
	c.log.Printf("Loading %s champion data ...\n", championName1)
	lolChampion1, err := c.lolTactics.ReadChampion(getYMLPath(championName1))
	if err != nil {
		return fmt.Errorf("loading champion %s: %v", championName1, err)
	}

	c.log.Printf("Loading %s champion data ...\n", championName2)
	lolChampion2, err := c.lolTactics.ReadChampion(getYMLPath(championName2))
	if err != nil {
		return fmt.Errorf("loading champion %s: %v", championName2, err)
	}

	c.log.Printf("Finding fight tactics (%s vs %s) ...\n", championName1, championName2)
	tacticsSol := c.lolTactics.Fight(lolChampion1, lolChampion2)

	fileName := setFilePath(lolChampion1, lolChampion2)
	file.Create(fileName)
	file.Write(fileName, getRoundSpellsToString(tacticsSol.RoundOfSpells, lolChampion2.Stats.HealthPoints, tacticsSol.Benchmark))

	return nil
}
