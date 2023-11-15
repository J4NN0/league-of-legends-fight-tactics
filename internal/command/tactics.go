package command

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

func (c *Controller) TacticsCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "tactics",
		Aliases: []string{"t"},
		Short:   "generate all fight tactics",
		Args:    cobra.ExactArgs(0),
		Run:     c.allChampionsFight,
	}
}

func (c *Controller) allChampionsFight(cmd *cobra.Command, args []string) {
	var championsName []string
	err := filepath.Walk(baseChampionPath, func(path string, info os.FileInfo, err error) error {
		if path != baseChampionPath {
			championsName = append(championsName, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
		}
		return nil
	})
	if err != nil {
		cmd.PrintErrf("listing champions data files in path %s: %v", baseChampionPath, err)
		os.Exit(-1)
	}

	var wg sync.WaitGroup
	for _, c1 := range championsName {
		for _, c2 := range championsName {
			if c1 != c2 {
				wg.Add(1)
				c1 := c1
				c2 := c2
				go func() {
					defer wg.Done()
					err = c.championsFight(c1, c2)
					if err != nil {
						c.log.Warningf("Could not generate fight tactics between %s vs %s: %v", c1, c2, err)
					}
				}()
			}
		}
	}
	wg.Wait()
}
