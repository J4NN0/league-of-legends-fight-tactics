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
	ID      string  `yaml:"id"`
	Name    string  `yaml:"name"`
	Title   string  `yaml:"title"`
	Tags    string  `yaml:"tags"`
	Passive Passive `yaml:"passive"`
	Stats   Stats   `yaml:"stats"`
	Spells  []Spell `yaml:"spells"`
}

type Passive struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Stats struct {
	HealthPoints float64 `yaml:"health_points"`
	AttackDamage float64 `yaml:"attack_damage"`
	AttackSpeed  float64 `yaml:"attack_speed"`
}

type Spell struct {
	ID       string    `yaml:"id"`
	Name     string    `yaml:"name"`
	MaxRank  int       `yaml:"max_rank"`
	Damage   []float64 `yaml:"damage"`
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
	fileName := getYMLPath(champion.ID)
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
