package controller

import (
	"league-of-legends-fight-tactics/internal/lol"
	"league-of-legends-fight-tactics/internal/riot"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Controller struct {
	log          Logger
	riotClient   RiotClient
	fightTactics LolFightTactics
}

type Logger interface {
	Printf(fmt string, args ...interface{})
	Warningf(fmt string, args ...interface{})
	Fatalf(fmt string, args ...interface{})
}

type RiotClient interface {
	GetAllLoLChampions() (championsData []riot.DDragonChampionResponse, err error)
	GetLoLChampion(championName string) (championResponse riot.DDragonChampionResponse, err error)
}

type LolFightTactics interface {
	Fight(champion1, champion2 lol.Champion)
}

func New(log Logger, riotClient RiotClient, fightTactics *lol.FightTactics) *Controller {
	return &Controller{log: log, riotClient: riotClient, fightTactics: fightTactics}
}

func (c *Controller) ChampionsFight(c1Name, c2Name string) {
	c.log.Printf("Loading %s vs %s champions data ...\n", c1Name, c2Name)

	c1, err := lol.Read(c1Name)
	if err != nil {
		c.log.Fatalf("Error loading champion: %v", err)
		return
	}

	c2, err := lol.Read(c2Name)
	if err != nil {
		c.log.Fatalf("Error loading champion: %v", err)
		return
	}

	c.fightTactics.Fight(c1, c2)
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

	championsData, err := c.riotClient.GetAllLoLChampions()
	if err != nil {
		c.log.Fatalf("Error while fetching all league of legends champions: %v", err)
		return
	}

	for _, champion := range championsData {
		err = storeChampionToYMLFile(champion)
		if err != nil {
			c.log.Warningf("Could not store %s champion data: %v", champion.DataName, err)
		} else {
			c.log.Printf("%s successfully stored", champion.DataName)
		}
	}
}

func storeChampionToYMLFile(championData riot.DDragonChampionResponse) error {
	err := lol.Write(mapChampionResponseToLolChampionStruct(championData))
	if err != nil {
		return err
	}
	return nil
}

func mapChampionResponseToLolChampionStruct(championResponse riot.DDragonChampionResponse) lol.Champion {
	championData := championResponse.Data[championResponse.DataName]
	lolChampion := lol.Champion{
		Version: championResponse.Version,
		Name:    championResponse.DataName,
		Tags:    strings.Join(championData.Tags, ", "),
		Stats: lol.Stat{
			Hp:                   championData.Stats.Hp,
			HpPerLevel:           championData.Stats.HpPerLevel,
			Armor:                championData.Stats.Armor,
			ArmorPerLevel:        championData.Stats.ArmorPerLevel,
			SpellBlock:           championData.Stats.SpellBlock,
			SpellBlockPerLevel:   championData.Stats.SpellBlockPerLevel,
			AttackDamage:         championData.Stats.AttackDamage,
			AttackDamagePerLevel: championData.Stats.AttackDamagePerLevel,
			AttackSpeed:          championData.Stats.AttackSpeed,
			AttackSpeedPerLevel:  championData.Stats.AttackSpeedPerLevel,
		},
	}

	// Add auto attack to list of spells
	lolChampion.Spells = []lol.Spell{{
		ID:       "aa",
		Name:     "Auto Attack",
		Damage:   []float32{championData.Stats.AttackDamage},
		MaxRank:  1, // it has no rank
		Cooldown: []float32{championData.Stats.AttackSpeed},
		Cast:     0,
	}}

	// Add champion spells
	for _, spell := range championData.Spells {
		lolChampion.Spells = append(lolChampion.Spells, lol.Spell{
			ID:       spell.ID,
			Name:     spell.Name,
			Damage:   spell.Damage,
			MaxRank:  spell.MaxRank,
			Cooldown: spell.Cooldown,
			Cast:     0.0,
		})
	}

	return lolChampion
}
