package lol

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
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

func (f *FightTactics) ReadChampion(filePath string) (champion Champion, err error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return Champion{}, err
	}

	err = yaml.Unmarshal(yamlFile, &champion)
	if err != nil {
		return Champion{}, fmt.Errorf("error unmarshalling: %w", err)
	}

	return champion, nil
}

func (f *FightTactics) WriteChampion(champion Champion, filePath string) error {
	data, err := yaml.Marshal(&champion)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0700)
	if err != nil {
		return err
	}

	return nil
}
