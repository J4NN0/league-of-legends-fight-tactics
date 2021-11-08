package champion_reader

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const baseChampionPath = "./champions/"
const fileExtension = ".yml"

// Champion LoL champion data struct
type Champion struct {
	Name   string   `yaml:"name"`
	Stats  Stat     `yaml:"stats"`
	Spells []string `yaml:"spells"`
}

type Stat struct {
	Hp        float32 `yaml:"hp"`
	AtkDamage float32 `yaml:"atk-damage"`
}

// Reader YAML champion reader
type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) LoadChampion(championName string) (champion Champion, err error) {
	var championYmlFilename = baseChampionPath + championName + fileExtension

	fmt.Printf("Loading %s ...\n", championName)

	yamlFile, err := ioutil.ReadFile(championYmlFilename)
	if err != nil {
		return champion, fmt.Errorf("error while loading %s: %v", championYmlFilename, err)
	}

	err = yaml.Unmarshal(yamlFile, &champion)
	if err != nil {
		return champion, fmt.Errorf("error while unmarshalling: %v", err)
	}

	fmt.Printf("Champion loaded successfully\n")

	return champion, nil
}
