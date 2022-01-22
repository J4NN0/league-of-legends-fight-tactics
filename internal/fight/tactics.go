package fight

import (
	"fmt"
	"league-of-legends-fight-tactics/pkg/yml"
	"math"
	"os"
)

// Number of spell used to take hp to zero
var fileName string
var bestBenchmark = math.MaxInt

type Tactics struct {
}

func New() *Tactics {
	return &Tactics{}
}

// Fight Champion1 vs Champion2 health point
func (t *Tactics) Fight(champion1, champion2 yml.LoLChampion) {
	fileName = setFilePath(champion1, champion2)
	fi, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	sol := make([]yml.Spell, getWorstCase(champion1.Spells, champion2.Stats.Hp))

	getBestRoundOfSpells(0, champion1.Spells, sol, champion2.Stats.Hp)
}

// setFilePath Set filename path
func setFilePath(champion1, champion2 yml.LoLChampion) string {
	return fmt.Sprintf("fights/%s_vs_%s.fight", champion1.Name, champion2.Name)
}

// getWorstCase Compute the worst case for allocate 'sol' array (i.e. hp / spell with less damage)
func getWorstCase(spells []yml.Spell, hp float32) int {
	minDamage := math.MaxFloat32

	for _, spell := range spells {
		if float64(spell.Damage) < minDamage {
			minDamage = float64(spell.Damage)
		}
	}

	return int(float64(hp) / minDamage)
}

func getBestRoundOfSpells(pos int, spells, sol []yml.Spell, hp float32) {
	if isHpZero(sol, hp) {
		setBenchmark(sol, hp)
		return
	}

	for i := 0; i < len(spells); i++ {
		sol[pos] = spells[i]
		getBestRoundOfSpells(pos+1, spells, sol, hp)
	}
}

// isHpZero True if hp is zero, false otherwise
func isHpZero(sol []yml.Spell, hp float32) bool {
	for _, spell := range sol {
		hp = hp - spell.Damage
		if hp <= 0 {
			return true
		}
	}

	return false
}

func setBenchmark(spells []yml.Spell, hp float32) {
	tmpBench := getBenchmark(spells, hp)
	if tmpBench < bestBenchmark {
		bestBenchmark = tmpBench

		fo, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()

		fmt.Printf("New solution found: %v (%d)\n", spells[:bestBenchmark], bestBenchmark)
		_, err = fo.WriteString(getRoundSpellsToString(spells, bestBenchmark, hp))
		if err != nil {
			panic(err)
		}
	}
}

// getBenchmark Get the exact number of spells needed to take hp down to zero
func getBenchmark(spells []yml.Spell, hp float32) int {
	var benchmark int

	for i, spell := range spells {
		hp = hp - spell.Damage
		if hp <= 0 {
			benchmark = i + 1
			break
		}
	}

	return benchmark
}

// getRoundSpellsToString Format spells into string
func getRoundSpellsToString(spells []yml.Spell, benchmark int, hp float32) string {
	var spellsToString string

	for i := 0; i < benchmark; i++ {
		spellsToString += fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", spells[i].ID, spells[i].Damage, hp, hp-spells[i].Damage)
		hp = hp - spells[i].Damage
	}

	return spellsToString
}
