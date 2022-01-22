package yml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

const baseChampionPath = "champions/lol/"
const fileExtension = ".yml"

// LoLChampion LoL champion data struct
type LoLChampion struct {
	Name   string   `yaml:"name"`
	Stats  stat     `yaml:"stats"`
	Spells []Spells `yaml:"spells"`
}

type stat struct {
	Hp        float32 `yaml:"hp"`
	AtkDamage float32 `yaml:"atk-damage"`
}

type Spells struct {
	ID     string  `yaml:"id"`
	Damage float32 `yaml:"damage"`
}

func LoadLoLChampion(championName string) (champion LoLChampion, err error) {
	var championYmlFilename = getYMLPath(championName)

	yamlFile, err := ioutil.ReadFile(championYmlFilename)
	if err != nil {
		return champion, fmt.Errorf("error reading %s: %w", championYmlFilename, err)
	}

	err = yaml.Unmarshal(yamlFile, &champion)
	if err != nil {
		return champion, fmt.Errorf("error unmarshalling: %w", err)
	}

	return champion, nil
}

func getYMLPath(championName string) string {
	return baseChampionPath + strings.ToLower(championName) + fileExtension
}
