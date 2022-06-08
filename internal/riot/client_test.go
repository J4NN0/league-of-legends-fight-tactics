package riot

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"league-of-legends-fight-tactics/tests/mock"
	"net/http"
	"testing"
)

func TestGetAllLoLChampionsSuccess(t *testing.T) {
	var format, version, championName = "standAloneComplex", "1.0.0", "TestName"
	var tags = []string{"Fighter", "Tank"}
	var hp, armor, atkDamage float32 = 1, 2, 3
	var expectedChampionsData = []DDragonChampionResponse{
		{
			Format:   format,
			Version:  version,
			DataName: sanitizeChampionName(championName),
			Data: map[string]championData{
				championName: {
					Name: championName,
					Tags: tags,
					Stats: stats{
						Hp:           hp,
						Armor:        armor,
						AttackDamage: atkDamage,
					},
				},
			},
		},
	}

	callCount := 0
	client := NewClient(&mock.Logger{}, mock.NewTestClient(func(req *http.Request) *http.Response {
		callCount = callCount + 1

		// Get all champions
		if callCount == 1 {
			return mock.Response(dDragonLoLAllChampionsResponse{
				Format:  format,
				Version: version,
				Data:    map[string]interface{}{championName: "some data"},
			}, 200)
		}

		// Get champion (i.e. championName)
		if callCount == 2 {
			return mock.Response(DDragonChampionResponse{
				Format:  format,
				Version: version,
				Data: map[string]championData{
					championName: {
						Name: championName,
						Tags: tags,
						Stats: stats{
							Hp:           hp,
							Armor:        armor,
							AttackDamage: atkDamage,
						},
					},
				},
			}, 200)
		}

		return nil
	}))

	championsResponse, _ := client.GetAllLoLChampions()

	assert.Equal(t, expectedChampionsData, championsResponse)
}

func TestGetAllLoLChampionsFail_GetAllChampions(t *testing.T) {
	client := NewClient(&mock.Logger{}, mock.NewTestClient(func(req *http.Request) *http.Response {
		return mock.Response(dDragonLoLAllChampionsResponse{}, 403)
	}))

	championsResponse, err := client.GetAllLoLChampions()

	assert.Equal(t, []DDragonChampionResponse{}, championsResponse)
	assert.NotNil(t, err)
}

func TestGetAllLoLChampionsFail_GetLoLChampion(t *testing.T) {
	callCount := 0
	client := NewClient(&mock.Logger{}, mock.NewTestClient(func(req *http.Request) *http.Response {
		callCount = callCount + 1
		if callCount == 2 {
			return mock.Response(DDragonChampionResponse{}, 403)
		}
		return nil
	}))

	championsResponse, err := client.GetAllLoLChampions()

	assert.Equal(t, []DDragonChampionResponse{}, championsResponse)
	assert.NotNil(t, err)
}

func TestGetLoLChampionSuccess(t *testing.T) {
	var format, version, championName, spellID = "standAloneComplex", "1.0.0", "TestName", "spellID"
	var tags, labels, effects = []string{"Fighter", "Tank"}, []string{"Damage", "Attack Damage", "Cooldown"}, []string{"", "95/130/165/200/235", "60/75/90/105/120"}
	var hp, armor, atkDamage float32 = 1, 2, 3
	var damages = []float32{95, 130, 165, 200, 235}
	var expectedChampionResponse = DDragonChampionResponse{
		Format:   format,
		Version:  version,
		DataName: sanitizeChampionName(championName),
		Data: map[string]championData{
			championName: {
				Name: championName,
				Tags: tags,
				Stats: stats{
					Hp:           hp,
					Armor:        armor,
					AttackDamage: atkDamage,
				},
				Spells: []spell{
					{
						ID: spellID,
						LevelTip: levelTip{
							Label: labels,
						},
						EffectBurn: effects,
						Damage:     damages,
					},
				},
			},
		},
	}

	client := NewClient(&mock.Logger{}, mock.NewTestClient(func(req *http.Request) *http.Response {
		return mock.Response(DDragonChampionResponse{
			Format:  format,
			Version: version,
			Data: map[string]championData{
				championName: {
					Name: championName,
					Tags: tags,
					Stats: stats{
						Hp:           hp,
						Armor:        armor,
						AttackDamage: atkDamage,
					},
					Spells: []spell{
						{
							ID: spellID,
							LevelTip: levelTip{
								Label: labels,
							},
							EffectBurn: effects,
							Damage:     damages,
						},
					},
				},
			},
		}, 200)
	}))

	var championResponse, _ = client.GetLoLChampion(championName)

	assert.Equal(t, expectedChampionResponse, championResponse)
}

func TestGetLoLChampionFail(t *testing.T) {
	client := NewClient(&mock.Logger{}, mock.NewTestClient(func(req *http.Request) *http.Response {
		return mock.Response(DDragonChampionResponse{}, 403)
	}))

	championResponse, err := client.GetLoLChampion("someChampionName")

	assert.Equal(t, DDragonChampionResponse{}, championResponse)
	assert.NotNil(t, err)
}

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

func TestGetChampionURL(t *testing.T) {
	championName := "someChampionName"
	expectedURL := fmt.Sprintf("%s/%s.json", dDragonLolChampionBaseURL, championName)

	url := getChampionURL(championName)

	assert.Equal(t, expectedURL, url)
}

func TestGetSpellDamage(t *testing.T) {
	spellTest := spell{
		ID: "sampleSpell",
		LevelTip: levelTip{
			Label: []string{"Damage", "Attack Damage", "Cooldown"},
		},
		EffectBurn: []string{"", "95/130/165/200/235", "60/75/90/105/120"},
	}
	expectedDamage := []float32{95, 130, 165, 200, 235}

	riotClient := NewClient(&mock.Logger{}, &http.Client{})
	damage := riotClient.getSpellDamage(spellTest)

	assert.Equal(t, expectedDamage, damage)
}

func TestGetSpellDamage_NoDamageLabel(t *testing.T) {
	spellTest := spell{
		ID: "sampleSpell",
		LevelTip: levelTip{
			Label: []string{"Attack Damage", "Cooldown"},
		},
		EffectBurn: []string{"", "60/75/90/105/120"},
	}
	expectedDamage := []float32{0, 0, 0, 0, 0}

	riotClient := NewClient(&mock.Logger{}, &http.Client{})
	damage := riotClient.getSpellDamage(spellTest)

	assert.Equal(t, expectedDamage, damage)
}

func TestGetSpellDamageFail(t *testing.T) {
	spellTest := spell{
		ID: "sampleSpell",
		LevelTip: levelTip{
			Label: []string{"Damage", "Attack Damage", "Cooldown"},
		},
		EffectBurn: []string{"", "not/numbers/in/here/lol", "60/75/90/105/120"},
	}
	expectedDamage := []float32{0, 0, 0, 0, 0}

	riotClient := NewClient(&mock.Logger{}, &http.Client{})
	damage := riotClient.getSpellDamage(spellTest)

	assert.Equal(t, expectedDamage, damage)
}
