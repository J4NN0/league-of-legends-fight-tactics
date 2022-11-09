package lol

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/logger/loggertest"
	"github.com/stretchr/testify/assert"

	"gopkg.in/yaml.v2"
)

const (
	mockChampionPathTmpl = "../../tests/champions/%s.yml"
	mockChampion1        = "mock01"
	mockChampion2        = "mock02"
	mockChampion3        = "mock03"
)

func readMockTestChampion(mockFileName string) (champion Champion, err error) {
	yamlFile, err := os.ReadFile(fmt.Sprintf(mockChampionPathTmpl, mockFileName))
	if err != nil {
		return Champion{}, err
	}

	err = yaml.Unmarshal(yamlFile, &champion)
	if err != nil {
		return Champion{}, fmt.Errorf("error unmarshalling: %w", err)
	}

	return champion, nil
}

func TestGetBestRoundOfSpells(t *testing.T) {
	fightTactics := FightTactics{&loggertest.Logger{}}

	t.Run("aa", func(t *testing.T) {
		var sol []Spell
		var bestSol = TacticsSol{Benchmark: math.MaxFloat64, RoundOfSpells: []Spell{}}
		var enemyHp = 200.0

		champion, err := readMockTestChampion(mockChampion1)
		assert.Nil(t, err)

		fightTactics.getBestRoundOfSpells(0, champion.Spells, sol, enemyHp, &bestSol)

		aa := champion.Spells[0]

		aaUsedTimes := math.Ceil(enemyHp / aa.Damage[0])
		aaCastTime := aa.Cast

		for _, s := range bestSol.RoundOfSpells {
			assert.Equal(t, aa, s)
		}
		assert.Equal(t, int(aaUsedTimes), len(bestSol.RoundOfSpells))
		assert.Equal(t, aaUsedTimes*aaCastTime, bestSol.Benchmark)
	})

	t.Run("only one spell", func(t *testing.T) {
		var sol []Spell
		var bestSol = TacticsSol{Benchmark: math.MaxFloat64, RoundOfSpells: []Spell{}}
		var enemyHp = 200.0

		champion, err := readMockTestChampion(mockChampion2)
		assert.Nil(t, err)

		fightTactics.getBestRoundOfSpells(0, champion.Spells, sol, enemyHp, &bestSol)

		qSpell := champion.Spells[0]
		maxRank := qSpell.MaxRank - 1

		spellUsedTimes := math.Ceil(enemyHp / qSpell.Damage[maxRank])
		spellCastTime := qSpell.Cast
		spellCooldownTime := qSpell.Cooldown[maxRank]

		assert.Equal(t, int(spellUsedTimes), len(bestSol.RoundOfSpells))
		assert.Equal(t, qSpell, bestSol.RoundOfSpells[0])
		assert.Equal(t, spellCastTime+(spellUsedTimes-1)*(spellCastTime+spellCooldownTime), bestSol.Benchmark)
	})

	t.Run("only spells (no re-usage)", func(t *testing.T) {
		var sol []Spell
		var bestSol = TacticsSol{Benchmark: math.MaxFloat64, RoundOfSpells: []Spell{}}
		var enemyHp = 450.0

		champion, err := readMockTestChampion(mockChampion3)
		assert.Nil(t, err)

		fightTactics.getBestRoundOfSpells(0, champion.Spells, sol, enemyHp, &bestSol)

		usedSpells := []Spell{champion.Spells[0], champion.Spells[1], champion.Spells[2], champion.Spells[3]}

		totCastTime := usedSpells[0].Cast + usedSpells[1].Cast + usedSpells[2].Cast + usedSpells[3].Cast

		assert.Equal(t, len(usedSpells), len(bestSol.RoundOfSpells))
		assert.Equal(t, usedSpells[0], bestSol.RoundOfSpells[0])
		assert.Equal(t, usedSpells[1], bestSol.RoundOfSpells[1])
		assert.Equal(t, usedSpells[2], bestSol.RoundOfSpells[2])
		assert.Equal(t, usedSpells[3], bestSol.RoundOfSpells[3])
		assert.Equal(t, totCastTime, bestSol.Benchmark)
	})

	t.Run("only spells (with re-usage)", func(t *testing.T) {
		var sol []Spell
		var bestSol = TacticsSol{Benchmark: math.MaxFloat64, RoundOfSpells: []Spell{}}
		var enemyHp = 600.0

		champion, err := readMockTestChampion(mockChampion3)
		assert.Nil(t, err)

		fightTactics.getBestRoundOfSpells(0, champion.Spells, sol, enemyHp, &bestSol)

		spells := []Spell{champion.Spells[0], champion.Spells[1], champion.Spells[2], champion.Spells[3]}
		usedSpells := []Spell{spells[0], spells[2], spells[0], spells[3], spells[0]}

		totCastTime := spells[0].Cast + spells[2].Cast + spells[0].Cast + spells[3].Cast + spells[0].Cast // 0 cooldown time

		assert.Equal(t, len(usedSpells), len(bestSol.RoundOfSpells))
		assert.Equal(t, usedSpells[0], bestSol.RoundOfSpells[0])
		assert.Equal(t, usedSpells[1], bestSol.RoundOfSpells[1])
		assert.Equal(t, usedSpells[2], bestSol.RoundOfSpells[2])
		assert.Equal(t, usedSpells[3], bestSol.RoundOfSpells[3])
		assert.Equal(t, usedSpells[4], bestSol.RoundOfSpells[4])
		assert.Equal(t, totCastTime, bestSol.Benchmark)
	})
}

func TestIsHpZero(t *testing.T) {
	spells := []Spell{
		{
			ID:      "q",
			Damage:  []float64{10},
			MaxRank: 1,
		},
		{
			ID:      "w",
			Damage:  []float64{10},
			MaxRank: 1,
		},
		{
			ID:      "e",
			Damage:  []float64{10},
			MaxRank: 1,
		},
		{
			ID:      "r",
			Damage:  []float64{10},
			MaxRank: 1,
		},
	}

	t.Run("hp below zero", func(t *testing.T) {
		isZero := isHpZero(spells, 30)
		assert.Equal(t, true, isZero)
	})

	t.Run("hp zero", func(t *testing.T) {
		isZero := isHpZero(spells, 40)
		assert.Equal(t, true, isZero)
	})

	t.Run("hp not zero", func(t *testing.T) {
		isZero := isHpZero(spells, 50)
		assert.Equal(t, false, isZero)
	})
}

func TestGetBenchmark(t *testing.T) {
	t.Run("no re-usage spells", func(t *testing.T) {
		spells := []Spell{
			{
				ID:       "aa",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     0.5,
			},
			{
				ID:       "q",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     1.0,
			},
			{
				ID:       "w",
				Damage:   []float64{6, 7, 8, 9, 20},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     2.0,
			},
			{
				ID:       "e",
				Damage:   []float64{6, 7, 8, 9, 30},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     3.0,
			},
			{
				ID:       "r",
				Damage:   []float64{6, 7, 8, 9, 40},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     4.0,
			},
		}

		benchmark := getBenchmark(spells)
		expectedBenchmark := 10.5 // sum of all spells cast time

		assert.Equal(t, expectedBenchmark, benchmark)
	})

	t.Run("re-usage spells in a row", func(t *testing.T) {
		spells := []Spell{
			{
				ID:       "aa",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     0.5,
			},
			{
				ID:       "q",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 1},
				Cast:     1.0,
			},
			{
				ID:       "w",
				Damage:   []float64{6, 7, 8, 9, 20},
				MaxRank:  5,
				Cooldown: []float64{5, 4, 3, 2, 2},
				Cast:     2.0,
			},
			{
				ID:       "w",
				Damage:   []float64{6, 7, 8, 9, 20},
				MaxRank:  5,
				Cooldown: []float64{5, 4, 3, 2, 2},
				Cast:     2.0,
			},
		}

		benchmark := getBenchmark(spells)
		expectedBenchmark := 7.5 // sum of all spells cast time + W cooldown

		assert.Equal(t, expectedBenchmark, benchmark)
	})

	t.Run("re-usage spells split by another one", func(t *testing.T) {
		spells := []Spell{
			{
				ID:       "aa",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     0.5,
			},
			{
				ID:       "q",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 1},
				Cast:     1.0,
			},
			{
				ID:       "w",
				Damage:   []float64{6, 7, 8, 9, 20},
				MaxRank:  5,
				Cooldown: []float64{5, 4, 3, 2, 2},
				Cast:     2.0,
			},
			{
				ID:       "aa",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     0.5,
			},
			{
				ID:       "w",
				Damage:   []float64{6, 7, 8, 9, 20},
				MaxRank:  5,
				Cooldown: []float64{5, 4, 3, 2, 2},
				Cast:     2.0,
			},
		}

		benchmark := getBenchmark(spells)
		expectedBenchmark := 7.5 // sum of all spells cast time + (W cooldown - A cast time)

		assert.Equal(t, expectedBenchmark, benchmark)
	})
}

func TestGetAdditionalTimeIfSpellIsInCooldown(t *testing.T) {
	t.Run("re-usage spell one time", func(t *testing.T) {
		usedSpells := []Spell{
			{
				ID:       "aa",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     0.5,
			},
			{
				ID:       "q",
				Damage:   []float64{6, 7, 8, 9, 50},
				MaxRank:  5,
				Cooldown: []float64{5, 5, 5, 5, 5},
				Cast:     2.0,
			},
			{
				ID:       "w",
				Damage:   []float64{6, 7, 8, 9, 20},
				MaxRank:  5,
				Cooldown: []float64{8, 8, 8, 8, 8},
				Cast:     1.0,
			},
		}
		currentSpell := usedSpells[2]

		timeToWait := getAdditionalTimeIfSpellIsInCooldown(currentSpell, usedSpells)
		expectedTimeToWait := 8.0 // W cooldown

		assert.Equal(t, expectedTimeToWait, timeToWait)
	})

	t.Run("re-usage spell in a row", func(t *testing.T) {
		usedSpells := []Spell{
			{
				ID:       "aa",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     0.5,
			},
			{
				ID:       "q",
				Damage:   []float64{6, 7, 8, 9, 50},
				MaxRank:  5,
				Cooldown: []float64{5, 5, 5, 5, 5},
				Cast:     2.0,
			},
			{
				ID:       "w",
				Damage:   []float64{6, 7, 8, 9, 20},
				MaxRank:  5,
				Cooldown: []float64{8, 8, 8, 8, 8},
				Cast:     1.0,
			},
			{
				ID:       "w",
				Damage:   []float64{6, 7, 8, 9, 20},
				MaxRank:  5,
				Cooldown: []float64{8, 8, 8, 8, 8},
				Cast:     1.0,
			},
		}
		currentSpell := usedSpells[3]

		timeToWait := getAdditionalTimeIfSpellIsInCooldown(currentSpell, usedSpells)
		expectedTimeToWait := 8.0 // W cooldown

		assert.Equal(t, expectedTimeToWait, timeToWait)
	})

	t.Run("re-usage spell split by another one", func(t *testing.T) {
		usedSpells := []Spell{
			{
				ID:       "aa",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     0.5,
			},
			{
				ID:       "q",
				Damage:   []float64{6, 7, 8, 9, 50},
				MaxRank:  5,
				Cooldown: []float64{5, 5, 5, 5, 5},
				Cast:     2.0,
			},
			{
				ID:       "w",
				Damage:   []float64{6, 7, 8, 9, 20},
				MaxRank:  5,
				Cooldown: []float64{8, 8, 8, 8, 8},
				Cast:     1.0,
			},
			{
				ID:       "aa",
				Damage:   []float64{6, 7, 8, 9, 10},
				MaxRank:  5,
				Cooldown: []float64{4, 3, 2, 1, 0},
				Cast:     0.5,
			},
		}
		currentSpell := usedSpells[2]

		timeToWait := getAdditionalTimeIfSpellIsInCooldown(currentSpell, usedSpells)
		expectedTimeToWait := 7.5 // W cooldown - AA cast time

		assert.Equal(t, expectedTimeToWait, timeToWait)
	})
}
