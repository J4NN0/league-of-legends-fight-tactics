package fight

import (
	"fmt"
	"league-of-legends-fight-tactics/pkg/yml"
	"math"
)

var maxBenchmark = math.MaxInt

type Tactics struct {
}

func New() *Tactics {
	return &Tactics{}
}

// Fight Champion1 vs Champion2 HP
func (t *Tactics) Fight(champion1, champion2 yml.LoLChampion) {
	worstCase := getWorstCase(champion1.Spells, champion2.Stats.Hp)
	sol := make([]yml.Spells, worstCase)

	getBestRoundOfSpells(0, champion1.Spells, sol, champion2.Stats.Hp)
}

func getWorstCase(spells []yml.Spells, hp float32) int {
	minDamage := math.MaxFloat32

	for _, spell := range spells {
		if float64(spell.Damage) < minDamage {
			minDamage = float64(spell.Damage)
		}
	}

	return int(float64(hp) / minDamage)
}

func getBestRoundOfSpells(pos int, spells, sol []yml.Spells, hp float32) {
	if hpIsZero(sol, hp) {
		setBenchmark(sol, hp)
		return
	}

	for i := 0; i < len(spells); i++ {
		sol[pos] = spells[i]
		getBestRoundOfSpells(pos+1, spells, sol, hp)
	}
}

func hpIsZero(sol []yml.Spells, hp float32) bool {
	for _, spell := range sol {
		hp = hp - spell.Damage
	}

	if hp <= 0 {
		return true
	} else {
		return false
	}
}

func setBenchmark(solSpells []yml.Spells, hp float32) {
	tmpBench := getBenchmarkSol(solSpells, hp)
	if tmpBench < maxBenchmark {
		maxBenchmark = tmpBench
		fmt.Printf("New solution found: %v (%d)\n", solSpells, maxBenchmark)
	}
}

func getBenchmarkSol(solSpells []yml.Spells, hp float32) int {
	for i, spell := range solSpells {
		hp = hp - spell.Damage
		if hp <= 0 {
			return i
		}
	}
	return 0
}
