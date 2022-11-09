package lol

import (
	"math"

	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/logger"
)

//TODO: only considering spell's max rank atm (i.e. spell.MaxRank-1), but need to consider all (e.g. 'q' has 5 ranks, etc.)

type Tactics interface {
	ReadChampion(filePath string) (champion Champion, err error)
	WriteChampion(champion Champion, filePath string) error
	Fight(champion1, champion2 Champion) TacticsSol
}

type TacticsSol struct {
	Benchmark     float64 // time (in seconds) taken to slay the enemy
	RoundOfSpells []Spell
}

type FightTactics struct {
	log logger.Logger
}

func NewTactics(log logger.Logger) Tactics {
	return &FightTactics{log: log}
}

// Fight Champion1 vs Champion2 health point
func (f *FightTactics) Fight(champion1, champion2 Champion) TacticsSol {
	var sol []Spell
	var bestSol = TacticsSol{Benchmark: math.MaxFloat64, RoundOfSpells: []Spell{}}

	f.getBestRoundOfSpells(0, champion1.Spells, sol, champion2.Stats.HealthPoints, &bestSol)

	f.log.Printf("[%s vs %s] Best solution found: enemy slayed in %.2fs\n", champion1.Name, champion2.Name, bestSol.Benchmark)

	return bestSol
}

func (f *FightTactics) getBestRoundOfSpells(pos int, spells, sol []Spell, hp float64, bestSol *TacticsSol) {
	if isHpZero(sol, hp) {
		f.setBenchmark(sol, bestSol)
		return
	}

	for i := 0; i < len(spells); i++ {
		if spells[i].Damage[0] != 0 {
			// TODO: excluding spells with zero damage atm, but need to take their passive into account
			sol = append(sol, spells[i])
			f.getBestRoundOfSpells(pos+1, spells, sol, hp, bestSol)
			sol = sol[:len(sol)-1] // pop value
		}
	}
}

// isHpZero True if hp is zero, false otherwise
func isHpZero(sol []Spell, hp float64) bool {
	for _, spell := range sol {
		hp = hp - spell.Damage[spell.MaxRank-1]
		if hp <= 0 {
			return true
		}
	}
	return false
}

func (f *FightTactics) setBenchmark(spells []Spell, bestSol *TacticsSol) {
	tmpBench := getBenchmark(spells)
	if tmpBench < bestSol.Benchmark {
		f.log.Printf("Found new best round of spells. Enemy slayed in %.2f seconds", tmpBench)

		bestSol.Benchmark = tmpBench

		bestSol.RoundOfSpells = make([]Spell, len(spells))
		copy(bestSol.RoundOfSpells, spells)
	}
}

// getBenchmark Given a set of spells (which brings the hp to zero), return the relevant benchmark (i.e. time needed to kill the enemy)
func getBenchmark(spells []Spell) (benchmark float64) {
	usedSpells := make([]Spell, len(spells))
	for i, s := range spells {
		benchmark += s.Cast + getAdditionalTimeIfSpellIsInCooldown(s, usedSpells)
		usedSpells[i] = s
	}
	return benchmark
}

// getAdditionalTimeIfSpellIsInCooldown Get additional waiting time if the spell has been used previously and is therefore still in cooldown.
func getAdditionalTimeIfSpellIsInCooldown(currentSpell Spell, usedSpells []Spell) (timeToWait float64) {
	// Check last time spell was used
	usedSpellPos := -1
	for i := len(usedSpells) - 1; i >= 0; i-- {
		if currentSpell.ID == usedSpells[i].ID {
			usedSpellPos = i + 1 // +1 to exclude this spell from the count below, i.e. cooldown starts after cast
			break
		}
	}

	// If it has been used before, check if spell is still in cooldown or not
	if usedSpellPos != -1 {
		var timePassed float64

		// Compute how much time has passed between the first usage of the spell and its re-usage
		for i := usedSpellPos; i < len(usedSpells); i++ {
			timePassed += usedSpells[i].Cast
		}

		if timePassed >= currentSpell.Cooldown[currentSpell.MaxRank-1] {
			// Spell is ready to be used
			timeToWait = 0
		} else {
			// Spell is still in cooldown
			timeToWait = currentSpell.Cooldown[currentSpell.MaxRank-1] - timePassed
		}
	}

	return timeToWait
}
