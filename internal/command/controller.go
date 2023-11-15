package command

import (
	"fmt"
	"strings"

	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/logger"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/lol"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/riot"
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

func setFilePath(champion1, champion2 lol.Champion) string {
	return fmt.Sprintf("fights/%s_vs_%s.loltactics", champion1.Name, champion2.Name)
}

func getRoundSpellsToString(spells []lol.Spell, hp, benchmark float64) string {
	var spellsToString string
	for _, s := range spells {
		spellsToString += fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", s.ID, s.Damage, hp, hp-s.Damage[s.MaxRank-1])
		hp = hp - s.Damage[s.MaxRank-1]
	}
	spellsToString += fmt.Sprintf("\nEnemy defeated in %.2fs\n", benchmark)
	return spellsToString
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
