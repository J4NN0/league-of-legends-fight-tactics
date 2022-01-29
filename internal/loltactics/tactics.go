package loltactics

import (
	"fmt"
	"league-of-legends-fight-tactics/pkg/file"
	"league-of-legends-fight-tactics/pkg/yml"
	"math"
)

var fileName string
var bestBenchmark float32 = math.MaxFloat32 // How many seconds to take down the enemy

// Fight Champion1 vs Champion2 health point
func Fight(champion1, champion2 yml.LoLChampion) {
	var sol []yml.Spell

	fileName = setFilePath(champion1, champion2)
	file.Create(fileName)

	getBestRoundOfSpells(0, champion1.Spells, sol, champion2.Stats.Hp)
}

// setFilePath Set filename path
func setFilePath(champion1, champion2 yml.LoLChampion) string {
	return fmt.Sprintf("fights/%s_vs_%s.loltactics", champion1.Name, champion2.Name)
}

func getBestRoundOfSpells(pos int, spells, sol []yml.Spell, hp float32) {
	if isHpZero(sol, hp) {
		setBenchmark(sol, hp)
		return
	}

	for i := 0; i < len(spells); i++ {
		sol = append(sol, spells[i])
		getBestRoundOfSpells(pos+1, spells, sol, hp)
		sol = sol[:len(sol)-1] // pop value
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
	tmpBench := getBenchmark(spells)
	if tmpBench < bestBenchmark {
		bestBenchmark = tmpBench
		fmt.Printf("New solution found: %v (%.2fs)\n", spells, bestBenchmark)
		file.Write(fileName, getRoundSpellsToString(spells, hp, bestBenchmark))
	}
}

// getBenchmark Given a set of spells (that take hp down to zero), return the related benchmark
func getBenchmark(spells []yml.Spell) float32 {
	var benchmark float32

	for _, spell := range spells {
		benchmark += spell.Cooldown
	}

	return benchmark
}

// getRoundSpellsToString Format spells into string
func getRoundSpellsToString(spells []yml.Spell, hp, benchmark float32) string {
	var spellsToString string

	for _, s := range spells {
		spellsToString += fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", s.ID, s.Damage, hp, hp-s.Damage)
		hp = hp - s.Damage
	}
	spellsToString += fmt.Sprintf("\nEnemy defeated in %.2fs\n", benchmark)

	return spellsToString
}
