package lol

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetFilePath(t *testing.T) {
	filename := setFilePath(Champion{Name: "Name1"}, Champion{Name: "Name2"})

	assert.Equal(t, "fights/Name1_vs_Name2.loltactics", filename)
}

func TestIsHpZero(t *testing.T) {
	spells := []Spell{
		{
			ID:      "q",
			Damage:  []float32{10},
			MaxRank: 1,
		},
		{
			ID:      "w",
			Damage:  []float32{10},
			MaxRank: 1,
		},
		{
			ID:      "e",
			Damage:  []float32{10},
			MaxRank: 1,
		},
		{
			ID:      "r",
			Damage:  []float32{10},
			MaxRank: 1,
		},
	}

	isZero := isHpZero(spells, 30)
	assert.Equal(t, true, isZero)

	isZero = isHpZero(spells, 40)
	assert.Equal(t, true, isZero)

	isZero = isHpZero(spells, 50)
	assert.Equal(t, false, isZero)
}

func TestGetBenchmark_WithNoReUsageSpells(t *testing.T) {
	spells := []Spell{
		{
			ID:       "aa",
			Damage:   []float32{6, 7, 8, 9, 10},
			MaxRank:  5,
			Cooldown: []float32{4, 3, 2, 1, 0},
			Cast:     0.5,
		},
		{
			ID:       "q",
			Damage:   []float32{6, 7, 8, 9, 10},
			MaxRank:  5,
			Cooldown: []float32{4, 3, 2, 1, 0},
			Cast:     1.0,
		},
		{
			ID:       "w",
			Damage:   []float32{6, 7, 8, 9, 20},
			MaxRank:  5,
			Cooldown: []float32{4, 3, 2, 1, 0},
			Cast:     2.0,
		},
		{
			ID:       "e",
			Damage:   []float32{6, 7, 8, 9, 30},
			MaxRank:  5,
			Cooldown: []float32{4, 3, 2, 1, 0},
			Cast:     3.0,
		},
		{
			ID:       "r",
			Damage:   []float32{6, 7, 8, 9, 40},
			MaxRank:  5,
			Cooldown: []float32{4, 3, 2, 1, 0},
			Cast:     4.0,
		},
	}

	benchmark := getBenchmark(spells)
	assert.Equal(t, float32(10.5), benchmark)
}

func TestGetBenchmark_WithReUsageSpells(t *testing.T) {
	spells := []Spell{
		{
			ID:       "aa",
			Damage:   []float32{6, 7, 8, 9, 10},
			MaxRank:  5,
			Cooldown: []float32{4, 3, 2, 1, 0},
			Cast:     0.5,
		},
		{
			ID:       "q",
			Damage:   []float32{6, 7, 8, 9, 10},
			MaxRank:  5,
			Cooldown: []float32{4, 3, 2, 1, 1},
			Cast:     1.0,
		},
		{
			ID:       "w",
			Damage:   []float32{6, 7, 8, 9, 20},
			MaxRank:  5,
			Cooldown: []float32{5, 4, 3, 2, 2},
			Cast:     2.0,
		},
		{
			ID:       "w",
			Damage:   []float32{6, 7, 8, 9, 20},
			MaxRank:  5,
			Cooldown: []float32{5, 4, 3, 2, 2},
			Cast:     2.0,
		},
	}

	benchmark := getBenchmark(spells)
	assert.Equal(t, float32(7.5), benchmark)
}

func TestGetAdditionalTimeIfSpellIsInCooldown(t *testing.T) {
	usedSpells := []Spell{
		{
			ID:       "aa",
			Damage:   []float32{6, 7, 8, 9, 10},
			MaxRank:  5,
			Cooldown: []float32{4, 3, 2, 1, 0},
			Cast:     0.5,
		},
		{
			ID:       "q",
			Damage:   []float32{6, 7, 8, 9, 50},
			MaxRank:  5,
			Cooldown: []float32{5, 5, 5, 5, 5},
			Cast:     2.0,
		},
		{
			ID:       "w",
			Damage:   []float32{6, 7, 8, 9, 20},
			MaxRank:  5,
			Cooldown: []float32{8, 8, 8, 8, 8},
			Cast:     1.0,
		},
	}
	spellToBeReused := Spell{ID: "w", Damage: []float32{6, 7, 8, 9, 20}, MaxRank: 5, Cooldown: []float32{8, 8, 8, 8, 8}, Cast: 1.0}

	timeToWait := getAdditionalTimeIfSpellIsInCooldown(spellToBeReused, usedSpells)
	assert.Equal(t, float32(8), timeToWait)
}

func TestGetRoundSpellsToString(t *testing.T) {
	var benchmark, hp float32 = 3.0, 15.0
	spells := []Spell{
		{
			ID:       "q",
			Damage:   []float32{6, 7, 8, 9, 10},
			MaxRank:  5,
			Cooldown: []float32{4, 3, 2, 1, 1},
		},
		{
			ID:       "w",
			Damage:   []float32{6, 7, 8, 9, 20},
			MaxRank:  5,
			Cooldown: []float32{5, 4, 3, 2, 2},
		},
	}

	spellsToString := getRoundSpellsToString(spells, hp, benchmark)
	expectedString := fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", spells[0].ID, spells[0].Damage, hp, hp-spells[0].Damage[spells[0].MaxRank-1])
	expectedString += fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", spells[1].ID, spells[1].Damage, hp-spells[0].Damage[spells[0].MaxRank-1], hp-spells[0].Damage[spells[0].MaxRank-1]-spells[1].Damage[spells[0].MaxRank-1])
	expectedString += fmt.Sprintf("\nEnemy defeated in %.2fs\n", benchmark)

	assert.Equal(t, expectedString, spellsToString)
}
