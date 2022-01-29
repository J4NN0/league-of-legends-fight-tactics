package yml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

const BaseChampionPath = "champions/lol/"
const fileExtension = ".yml"

// LoLChampion LoL champion data struct
type LoLChampion struct {
	Name   string  `yaml:"name"`
	Stats  stat    `yaml:"stats"`
	Spells []Spell `yaml:"spells"`
}

type stat struct {
	Hp        float32 `yaml:"hp"`
	AtkDamage float32 `yaml:"atk-damage"`
}

type Spell struct {
	ID       string  `yaml:"id"`
	Damage   float32 `yaml:"damage"`
	Cooldown float32 `yaml:"cooldown"`
}

func LoadLoLChampion(championName string) (champion LoLChampion, err error) {
	yamlFile, err := ioutil.ReadFile(getYMLPath(championName))
	if err != nil {
		return LoLChampion{}, err
	}

	err = yaml.Unmarshal(yamlFile, &champion)
	if err != nil {
		return LoLChampion{}, fmt.Errorf("error unmarshalling: %w", err)
	}

	return champion, nil
}

func getYMLPath(championName string) string {
	return BaseChampionPath + strings.ToLower(championName) + fileExtension
}
