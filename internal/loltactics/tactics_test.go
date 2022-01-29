package loltactics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"league-of-legends-fight-tactics/pkg/yml"
	"testing"
)

func TestSetFilePath(t *testing.T) {
	filename := setFilePath(yml.LoLChampion{Name: "Name1"}, yml.LoLChampion{Name: "Name2"})

	assert.Equal(t, "fights/Name1_vs_Name2.loltactics", filename)
}

func TestIsHpZero(t *testing.T) {
	spells := []yml.Spell{
		{
			ID:     "q",
			Damage: 10,
		},
		{
			ID:     "w",
			Damage: 10,
		},
		{
			ID:     "e",
			Damage: 10,
		},
		{
			ID:     "r",
			Damage: 10,
		},
	}

	isZero := isHpZero(spells, 30)
	assert.Equal(t, true, isZero)

	isZero = isHpZero(spells, 40)
	assert.Equal(t, true, isZero)

	isZero = isHpZero(spells, 50)
	assert.Equal(t, false, isZero)
}

func TestGetBenchmark(t *testing.T) {
	spells := []yml.Spell{
		{
			ID:       "q",
			Damage:   10,
			Cooldown: 1.0,
		},
		{
			ID:       "w",
			Damage:   20,
			Cooldown: 2.0,
		},
		{
			ID:       "e",
			Damage:   30,
			Cooldown: 3.0,
		},
		{
			ID:       "r",
			Damage:   40,
			Cooldown: 4.0,
		},
	}

	benchmark := getBenchmark(spells)
	assert.Equal(t, float32(10), benchmark)
}

func TestGetRoundSpellsToString(t *testing.T) {
	var benchmark, hp float32 = 3.0, 15.0
	spells := []yml.Spell{
		{
			ID:       "q",
			Damage:   10,
			Cooldown: 1.0,
		},
		{
			ID:       "w",
			Damage:   20,
			Cooldown: 2.0,
		},
	}

	spellsToString := getRoundSpellsToString(spells, hp, benchmark)
	expectedString := fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", spells[0].ID, spells[0].Damage, hp, hp-spells[0].Damage)
	expectedString += fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", spells[1].ID, spells[1].Damage, hp-spells[0].Damage, hp-spells[0].Damage-spells[1].Damage)
	expectedString += fmt.Sprintf("\nEnemy defeated in %.2fs\n", benchmark)

	assert.Equal(t, expectedString, spellsToString)
}