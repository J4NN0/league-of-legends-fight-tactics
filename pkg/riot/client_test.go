package riot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeChampionName(t *testing.T) {
	t.Run("sanitizeChampionName Jhin", func(t *testing.T) {
		championName := sanitizeChampionName("jHiN")
		assert.Equal(t, "Jhin", championName)
	})

	t.Run("sanitizeChampionName AurelionSol", func(t *testing.T) {
		championName := sanitizeChampionName("Aurelionsol")
		assert.Equal(t, "AurelionSol", championName)
	})

	t.Run("sanitizeChampionName DrMundo", func(t *testing.T) {
		championName := sanitizeChampionName("dRMundo")
		assert.Equal(t, "DrMundo", championName)
	})

	t.Run("sanitizeChampionName JarvanIV", func(t *testing.T) {
		championName := sanitizeChampionName("jarvanIV")
		assert.Equal(t, "JarvanIV", championName)
	})

	t.Run("sanitizeChampionName KogMaw", func(t *testing.T) {
		championName := sanitizeChampionName("Kogmaw")
		assert.Equal(t, "KogMaw", championName)
	})

	t.Run("sanitizeChampionName LeeSin", func(t *testing.T) {
		championName := sanitizeChampionName("Leesin")
		assert.Equal(t, "LeeSin", championName)
	})

	t.Run("sanitizeChampionName MasterYi", func(t *testing.T) {
		championName := sanitizeChampionName("Masteryi")
		assert.Equal(t, "MasterYi", championName)
	})

	t.Run("sanitizeChampionName MissFortune", func(t *testing.T) {
		championName := sanitizeChampionName("Missfortune")
		assert.Equal(t, "MissFortune", championName)
	})

	t.Run("sanitizeChampionName MonkeyKing", func(t *testing.T) {
		championName := sanitizeChampionName("Monkeyking")
		assert.Equal(t, "MonkeyKing", championName)
	})

	t.Run("sanitizeChampionName Wukong", func(t *testing.T) {
		championName := sanitizeChampionName("Monkeyking")
		assert.Equal(t, "MonkeyKing", championName)
	})

	t.Run("sanitizeChampionName RekSai", func(t *testing.T) {
		championName := sanitizeChampionName("Reksai")
		assert.Equal(t, "RekSai", championName)
	})

	t.Run("sanitizeChampionName TahmKench", func(t *testing.T) {
		championName := sanitizeChampionName("Tahmkench")
		assert.Equal(t, "TahmKench", championName)
	})

	t.Run("sanitizeChampionName TwistedFate", func(t *testing.T) {
		championName := sanitizeChampionName("Twistedfate")
		assert.Equal(t, "TwistedFate", championName)
	})

	t.Run("sanitizeChampionName XinZhao", func(t *testing.T) {
		championName := sanitizeChampionName("Xinzhao")
		assert.Equal(t, "XinZhao", championName)
	})
}
