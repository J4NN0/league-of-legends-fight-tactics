package controller

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/J4NN0/league-of-legends-fight-tactics/internal/log"
	"github.com/J4NN0/league-of-legends-fight-tactics/internal/lol"
	"github.com/J4NN0/league-of-legends-fight-tactics/internal/riot"
	"github.com/KnutZuidema/golio/datadragon"
)

type Controller struct {
	log          log.Logger
	riotClient   riot.Client
	fightTactics lol.Tactics
}

func New(log log.Logger, riotClient riot.Client, fightTactics lol.Tactics) *Controller {
	return &Controller{log: log, riotClient: riotClient, fightTactics: fightTactics}
}

func (c *Controller) ChampionsFight(championName1, championName2 string) {
	c.log.Printf("Loading %s vs %s champions data ...\n", championName1, championName2)

	lolChampion1, err := lol.Read(championName1)
	if err != nil {
		c.log.Fatalf("Error loading champion %s: %v", championName1, err)
		return
	}

	lolChampion2, err := lol.Read(championName2)
	if err != nil {
		c.log.Fatalf("Error loading champion %s: %v", championName2, err)
		return
	}

	c.fightTactics.Fight(lolChampion1, lolChampion2)
}

func (c *Controller) AllChampionsFight() {
	var wg sync.WaitGroup
	var champions []string

	err := filepath.Walk(lol.BaseChampionPath, func(path string, info os.FileInfo, err error) error {
		champions = append(champions, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
		return nil
	})
	if err != nil {
		c.log.Fatalf("Error listing champions data files in path %s: %v", lol.BaseChampionPath, err)
		return
	}

	champions = champions[1:]

	for _, c1 := range champions {
		for _, c2 := range champions {
			if c1 != c2 {
				wg.Add(1)
				c1 := c1
				c2 := c2
				go func() {
					defer wg.Done()
					c.ChampionsFight(c1, c2)
				}()
			}
		}
	}

	wg.Wait()
}

func (c *Controller) FetchChampion(championName string) {
	c.log.Printf("Fetching %s ...", championName)

	championData, err := c.riotClient.GetLoLChampion(championName)
	if err != nil {
		c.log.Fatalf("Error while fetching league of legends champions: %v", err)
		return
	}

	err = storeChampionToYMLFile(championData)
	if err != nil {
		c.log.Fatalf("Could not store %s champion data: %v", championName, err)
		return
	}

	c.log.Printf("%s successfully stored", championName)
}

func (c *Controller) FetchAllChampions() {
	c.log.Printf("Fetching all league of legends champions ...\n")

	ddChampions, err := c.riotClient.GetAllLoLChampions()
	if err != nil {
		c.log.Fatalf("Error while fetching all league of legends champions: %v", err)
		return
	}

	for _, champion := range ddChampions {
		err = storeChampionToYMLFile(champion)
		if err != nil {
			c.log.Warningf("Could not store %s champion data: %v", champion.ChampionData.Name, err)
		} else {
			c.log.Printf("%s successfully stored", champion.ChampionData.Name)
		}
	}
}

func storeChampionToYMLFile(ddChampion datadragon.ChampionDataExtended) error {
	lolChampion := mapChampionResponseToLolChampionStruct(ddChampion)
	err := lol.Write(lolChampion)
	if err != nil {
		return err
	}
	return nil
}

func mapChampionResponseToLolChampionStruct(ddChampion datadragon.ChampionDataExtended) lol.Champion {
	lolChampion := lol.Champion{
		DataDragonID: ddChampion.ID,
		Name:         ddChampion.ChampionData.Name,
		Title:        ddChampion.Title,
		Tags:         strings.Join(ddChampion.Tags, ", "),
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
