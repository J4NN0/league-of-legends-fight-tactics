package lolchampion

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

const BaseChampionPath = "champions/lol/"
const fileExtension = ".yml"

// Champion LoL champion data struct
type Champion struct {
	Name   string  `yaml:"name"`
	Stats  stat    `yaml:"stats"`
	Spells []Spell `yaml:"spells"`
}

type stat struct {
	Hp float32 `yaml:"hp"`
}

type Spell struct {
	ID       string  `yaml:"id"`
	Damage   float32 `yaml:"damage"`
	Cooldown float32 `yaml:"cooldown"`
	Cast     float32 `yaml:"cast"`
}

func Load(championName string) (champion Champion, err error) {
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

func getYMLPath(championName string) string {
	return BaseChampionPath + strings.ToLower(championName) + fileExtension
}
