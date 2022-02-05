package lol

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"league-of-legends-fight-tactics/pkg/file"
	"strings"
)

const BaseChampionPath = "champions/lol/"
const fileExtension = ".yml"

// Champion LoL champion data struct
type Champion struct {
	Version string  `json:"version"`
	Name    string  `yaml:"name"`
	Tags    string  `yaml:"tags"`
	Stats   Stat    `yaml:"stats"`
	Spells  []Spell `yaml:"spells"`
}

type Stat struct {
	Hp                   float32 `yaml:"hp"`
	HpPerLevel           float32 `yaml:"hp_per_level"`
	Armor                float32 `yaml:"armor"`
	ArmorPerLevel        float32 `yaml:"armor_per_level"`
	SpellBlock           float32 `yaml:"spell_block"`
	SpellBlockPerLevel   float32 `yaml:"spell_block_per_level"`
	AttackDamage         float32 `yaml:"attack_damage"`
	AttackDamagePerLevel float32 `yaml:"attack_damage_per_level"`
	AttackSpeed          float32 `yaml:"attack_speed"`
	AttackSpeedPerLevel  float32 `yaml:"attack_speed_per_level"`
}

type Spell struct {
	ID       string    `yaml:"id"`
	Name     string    `yaml:"name"`
	Damage   float32   `yaml:"damage"`
	MaxRank  int       `yaml:"max_rank"`
	Cooldown []float32 `yaml:"cooldown"`
	Cast     float32   `yaml:"cast"`
}

func Read(championName string) (champion Champion, err error) {
	yamlFile, err := ioutil.ReadFile(getYMLPath(championName))
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

	err = ioutil.WriteFile(fileName, data, 0)
	if err != nil {
		return err
	}

	return nil
}

func getYMLPath(championName string) string {
	return BaseChampionPath + strings.ToLower(championName) + fileExtension
}
