package loltactics

import (
	"fmt"
	"league-of-legends-fight-tactics/pkg/file"
	"league-of-legends-fight-tactics/pkg/yml"
	"math"
)

var fileName string
var bestBenchmark float32 = math.MaxFloat32 // How much time (in seconds) to slay the enemy
var bestRoundOfSpells []yml.Spell

// Fight Champion1 vs Champion2 health point
func Fight(champion1, champion2 yml.LoLChampion) {
	var sol []yml.Spell

	getBestRoundOfSpells(0, champion1.Spells, sol, champion2.Stats.Hp)

	fmt.Printf("[+] Best solution found (%s vs %s): %v (slayed in %.2fs)\n", champion1.Name, champion2.Name, bestRoundOfSpells, bestBenchmark)

	fileName = setFilePath(champion1, champion2)
	file.Create(fileName)
	file.Write(fileName, getRoundSpellsToString(bestRoundOfSpells, champion2.Stats.Hp, bestBenchmark))
}

// setFilePath Set filename path
func setFilePath(champion1, champion2 yml.LoLChampion) string {
	return fmt.Sprintf("fights/%s_vs_%s.loltactics", champion1.Name, champion2.Name)
}

func getBestRoundOfSpells(pos int, spells, sol []yml.Spell, hp float32) {
	if isHpZero(sol, hp) {
		setBenchmark(sol)
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

func setBenchmark(spells []yml.Spell) {
	tmpBench := getBenchmark(spells)
	if tmpBench < bestBenchmark {
		bestBenchmark = tmpBench
		bestRoundOfSpells = spells
	}
}

// getBenchmark Given a set of spells (that take hp down to zero), return the related benchmark
func getBenchmark(spells []yml.Spell) float32 {
	var benchmark float32
	var usedSpells []yml.Spell

	for _, spell := range spells {
		benchmark += spell.Cast + getAdditionalTimeIfSpellIsInCooldown(spell, usedSpells)
		usedSpells = append(usedSpells, spell)
	}

	return benchmark
}

// getAdditionalTimeIfSpellIsInCooldown Get additional time to wait if spell has been used before and is therefore still in cooldown
func getAdditionalTimeIfSpellIsInCooldown(currentSpell yml.Spell, usedSpells []yml.Spell) float32 {
	var usedSpellPos = -1
	var timeToWait float32

	// Check if spell has been used already
	for i, usedSpell := range usedSpells {
		if currentSpell.ID == usedSpell.ID {
			usedSpellPos = i + 1 // (+1 to exclude this spell from the count below, i.e. cooldown starts after cast)
			break
		}
	}

	// If it has been used before, check if spell is still in cooldown or not
	if usedSpellPos != -1 {
		var timePassed float32

		// Compute how much time has passed between the first usage of the spell and its re-usage
		for i := usedSpellPos; i < len(usedSpells); i++ {
			timePassed += usedSpells[i].Cast
		}

		if timePassed >= currentSpell.Cooldown {
			// Spell is ready to be used
			timeToWait = 0
		} else {
			// Spell is still in cooldown
			timeToWait = currentSpell.Cooldown - timePassed
		}
	}

	return timeToWait
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
