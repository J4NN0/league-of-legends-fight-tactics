package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/J4NN0/league-of-legends-fight-tactics/internal/logger"
	"github.com/J4NN0/league-of-legends-fight-tactics/internal/lol"
	"github.com/J4NN0/league-of-legends-fight-tactics/internal/riot"
	"github.com/KnutZuidema/golio/datadragon"
)

const (
	baseChampionPath = "champions/lol"
	fileExtension    = "yml"
)

type Controller struct {
	log        logger.Logger
	riotClient riot.Client
	lolTactics lol.Tactics
}

func New(log logger.Logger, riotClient riot.Client, lolTactics lol.Tactics) *Controller {
	return &Controller{log: log, riotClient: riotClient, lolTactics: lolTactics}
}

func (c *Controller) ChampionsFight(championName1, championName2 string) error {
	c.log.Printf("Loading %s vs %s champions data ...\n", championName1, championName2)

	lolChampion1, err := c.lolTactics.ReadChampion(championName1)
	if err != nil {
		return fmt.Errorf("loading champion %s: %v", championName1, err)
	}

	lolChampion2, err := c.lolTactics.ReadChampion(championName2)
	if err != nil {
		return fmt.Errorf("loading champion %s: %v", championName2, err)
	}

	c.lolTactics.Fight(lolChampion1, lolChampion2)

	return nil
}

func (c *Controller) AllChampionsFight() error {
	var championsName []string
	err := filepath.Walk(baseChampionPath, func(path string, info os.FileInfo, err error) error {
		if path != baseChampionPath {
			championsName = append(championsName, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("listing champions data files in path %s: %v", baseChampionPath, err)
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
					err = c.ChampionsFight(c1, c2)
					if err != nil {
						c.log.Warningf("Could not generate fight tactics between %s vs %s: %v", c1, c2, err)
					}
				}()
			}
		}
	}
	wg.Wait()

	return nil
}

func (c *Controller) FetchChampion(championName string) error {
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

func (c *Controller) FetchAllChampions() error {
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

func (c *Controller) storeChampionToYMLFile(ddChampion datadragon.ChampionDataExtended) error {
	lolChampion := mapChampionResponseToLolChampionStruct(ddChampion)
	filePath := getYMLPath(lolChampion.ID)

	err := c.lolTactics.WriteChampion(lolChampion, filePath)
	if err != nil {
		return err
	}

	return nil
}

func mapChampionResponseToLolChampionStruct(ddChampion datadragon.ChampionDataExtended) lol.Champion {
	lolChampion := lol.Champion{
		ID:    ddChampion.ID,
		Name:  ddChampion.ChampionData.Name,
		Title: ddChampion.Title,
		Tags:  strings.Join(ddChampion.Tags, ", "),
		Passive: lol.Passive{
			Name:        ddChampion.Passive.Name,
			Description: ddChampion.Passive.Description,
		},
		Stats: lol.Stats{
			HealthPoints: ddChampion.Stats.HealthPoints,
			AttackDamage: ddChampion.Stats.AttackDamage,
			AttackSpeed:  ddChampion.Stats.AttackSpeedOffset,
		},
		Spells: []lol.Spell{
			{
				ID:       "aa",
				Name:     "Auto Attack",
				Damage:   []float64{ddChampion.Stats.AttackDamage},
				MaxRank:  1, // it has no upgrade
				Cooldown: []float64{ddChampion.Stats.AttackSpeedOffset},
				Cast:     0, // it cannot be retrieved from DataDragon APIs
			},
		},
	}

	// Add remaining spells
	for _, spell := range ddChampion.Spells {
		var spellDamages []float64
		if len(spell.Effect) >= 1 {
			spellDamages = spell.Effect[1] // effect and effectBurn arrays have a null value in the 0 index (aka they are arrays 1-based)
		}

		lolChampion.Spells = append(lolChampion.Spells, lol.Spell{
			ID:       spell.ID,
			Name:     spell.Name,
			Damage:   spellDamages,
			MaxRank:  spell.MaxRank,
			Cooldown: spell.Cooldown,
			Cast:     0.0, // it cannot be retrieved from DataDragon APIs
		})
	}

	return lolChampion
}

func getYMLPath(championName string) string {
	return fmt.Sprintf("%s/%s.%s", baseChampionPath, strings.ReplaceAll(strings.ToLower(championName), " ", ""), fileExtension)
}
