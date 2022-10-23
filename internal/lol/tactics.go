package lol

import (
	"fmt"
	"math"

	"github.com/J4NN0/league-of-legends-fight-tactics/internal/log"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/file"
)

// TODO: only considering spell's max rank atm (i.e. spell.MaxRank-1), but need to consider all (e.g. 'q' has 5 ranks, etc.)

type Tactics interface {
	Fight(champion1, champion2 Champion)
}

type FightTactics struct {
	log log.Logger
}

func New(log log.Logger) Tactics {
	return &FightTactics{log: log}
}

type bestSolution struct {
	Benchmark     float64 // time (in seconds) taken to slay the enemy
	RoundOfSpells []Spell
}

// Fight Champion1 vs Champion2 health point
func (f *FightTactics) Fight(champion1, champion2 Champion) {
	var sol []Spell
	var bestSol = bestSolution{Benchmark: math.MaxFloat64, RoundOfSpells: []Spell{}}

	getBestRoundOfSpells(0, champion1.Spells, sol, champion2.Stats.Hp, &bestSol)

	f.log.Printf("[%s vs %s] Best solution found: enemy slayed in %.2fs\n", champion1.Name, champion2.Name, bestSol.Benchmark)

	fileName := setFilePath(champion1, champion2)
	file.Create(fileName)
	file.Write(fileName, getRoundSpellsToString(bestSol.RoundOfSpells, champion2.Stats.Hp, bestSol.Benchmark))
}

func getBestRoundOfSpells(pos int, spells, sol []Spell, hp float64, bestSol *bestSolution) {
	if isHpZero(sol, hp) {
		setBenchmark(sol, bestSol)
		return
	}

	for i := 0; i < len(spells); i++ {
		if spells[i].Damage[0] != 0 {
			// TODO: excluding spells with zero damage atm, but need to take their passive into account
			sol = append(sol, spells[i])
			getBestRoundOfSpells(pos+1, spells, sol, hp, bestSol)
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

func setBenchmark(spells []Spell, bestSol *bestSolution) {
	tmpBench := getBenchmark(spells)
	if tmpBench < bestSol.Benchmark {
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

// getRoundSpellsToString Format spells into string
func getRoundSpellsToString(spells []Spell, hp, benchmark float64) (spellsToString string) {
	for _, s := range spells {
		spellsToString += fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", s.ID, s.Damage, hp, hp-s.Damage[s.MaxRank-1])
		hp = hp - s.Damage[s.MaxRank-1]
	}
	spellsToString += fmt.Sprintf("\nEnemy defeated in %.2fs\n", benchmark)

	return spellsToString
}

// setFilePath Set filename path
func setFilePath(champion1, champion2 Champion) string {
	return fmt.Sprintf("fights/%s_vs_%s.loltactics", champion1.Name, champion2.Name)
}
