package lol

import (
	"fmt"
	"os"
	"strings"

	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/file"

	"gopkg.in/yaml.v2"
)

const (
	BaseChampionPath = "champions/lol/"
	fileExtension    = ".yml"
)

// Champion LoL champion data struct
type Champion struct {
	Version string  `json:"version"`
	Name    string  `yaml:"name"`
	Tags    string  `yaml:"tags"`
	Stats   Stat    `yaml:"stats"`
	Spells  []Spell `yaml:"spells"`
}

type Stat struct {
	Hp                   float64 `yaml:"hp"`
	HpPerLevel           float64 `yaml:"hp_per_level"`
	Armor                float64 `yaml:"armor"`
	ArmorPerLevel        float64 `yaml:"armor_per_level"`
	SpellBlock           float64 `yaml:"spell_block"`
	SpellBlockPerLevel   float64 `yaml:"spell_block_per_level"`
	AttackDamage         float64 `yaml:"attack_damage"`
	AttackDamagePerLevel float64 `yaml:"attack_damage_per_level"`
	AttackSpeed          float64 `yaml:"attack_speed"`
	AttackSpeedPerLevel  float64 `yaml:"attack_speed_per_level"`
}

type Spell struct {
	ID       string    `yaml:"id"`
	Name     string    `yaml:"name"`
	Damage   []float64 `yaml:"damage"`
	MaxRank  int       `yaml:"max_rank"`
	Cooldown []float64 `yaml:"cooldown"`
	Cast     float64   `yaml:"cast"`
}

func Read(championName string) (champion Champion, err error) {
	yamlFile, err := os.ReadFile(getYMLPath(championName))
	if err != nil {
		return Champion{}, err
	}

	err = yaml.Unmarshal(yamlFile, &champion)
	if err != nil {
		return Champion{}, fmt.Errorf("error unmarshalling: %w", err)
	}

	return champion, nil
}

func Write(champion Champion) error {
	fileName := getYMLPath(champion.Name)
	file.Create(fileName)

	data, err := yaml.Marshal(&champion)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, data, 0)
	if err != nil {
		return err
	}

	return nil
}

func getYMLPath(championName string) string {
	return BaseChampionPath + strings.ReplaceAll(strings.ToLower(championName), " ", "") + fileExtension
}
