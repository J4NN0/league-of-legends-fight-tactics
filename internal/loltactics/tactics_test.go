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

func TestGetWorstCase(t *testing.T) {
	var hp float32 = 500
	spells := []yml.Spell{
		{
			ID:     "q",
			Damage: 40,
		},
		{
			ID:     "w",
			Damage: 30,
		},
		{
			ID:     "e",
			Damage: 20,
		},
		{
			ID:     "r",
			Damage: 10,
		},
	}

	worstCase := getWorstCase(spells, hp)

	assert.Equal(t, int(hp/10), worstCase)
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

	benchmark := getBenchmark(spells, 10)
	assert.Equal(t, 1, benchmark)

	benchmark = getBenchmark(spells, 30)
	assert.Equal(t, 3, benchmark)

	benchmark = getBenchmark(spells, 35)
	assert.Equal(t, 4, benchmark)
}

func TestGetRoundSpellsToString(t *testing.T) {
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

	spellsToString := getRoundSpellsToString(spells, 1, 10)
	assert.Equal(t, fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", spells[0].ID, spells[0].Damage, 10.00, 10-spells[0].Damage), spellsToString)

	spellsToString = getRoundSpellsToString(spells, 2, 15)
	assert.Equal(t, fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n%s: %.2f (hp: %.2f -> %.2f)\n", spells[0].ID, spells[0].Damage, 15.00, 15.00-spells[0].Damage, spells[1].ID, spells[1].Damage, 5.00, 5-spells[0].Damage), spellsToString)
}
