package riot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"league-of-legends-fight-tactics/internal/log"
	"net/http"
	"testing"
)

func mockResponse(obj interface{}, status int) *http.Response {
	jsonMarshal, _ := json.Marshal(obj)
	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(bytes.NewReader(jsonMarshal)),
		Header:     make(http.Header),
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestGetAllLoLChampionsSuccess(t *testing.T) {
	var format, version, championName = "standAloneComplex", "1.0.0", "TestName"
	var tags = []string{"Fighter", "Tank"}
	var hp, armor, atkDamage float32 = 1, 2, 3
	var expectedChampionsData = []DDragonChampionResponse{
		{
			Format:   format,
			Version:  version,
			DataName: championName,
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
	client := NewApiClient(log.New("testApp"), newTestClient(func(req *http.Request) *http.Response {
		callCount = callCount + 1

		// Get all champions
		if callCount == 1 {
			return mockResponse(dDragonLoLAllChampionsResponse{
				Format:  format,
				Version: version,
				Data:    map[string]interface{}{championName: "some data"},
			}, 200)
		}

		// Get champion (i.e. championName)
		if callCount == 2 {
			return mockResponse(DDragonChampionResponse{
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
	client := NewApiClient(log.New("testApp"), newTestClient(func(req *http.Request) *http.Response {
		return mockResponse(dDragonLoLAllChampionsResponse{}, 403)
	}))

	championsResponse, err := client.GetAllLoLChampions()

	assert.Equal(t, []DDragonChampionResponse{}, championsResponse)
	assert.NotNil(t, err)
}

func TestGetAllLoLChampionsFail_GetLoLChampion(t *testing.T) {
	callCount := 0
	client := NewApiClient(log.New("testApp"), newTestClient(func(req *http.Request) *http.Response {
		callCount = callCount + 1
		if callCount == 2 {
			return mockResponse(DDragonChampionResponse{}, 403)
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
	var expectedChampionResponse = DDragonChampionResponse{
		Format:   format,
		Version:  version,
		DataName: championName,
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
						Leveltip: leveltip{
							Label: labels,
						},
						EffectBurn: effects,
						Damage:     []float32{95, 130, 165, 200, 235},
					},
				},
			},
		},
	}

	client := NewApiClient(log.New("testApp"), newTestClient(func(req *http.Request) *http.Response {
		return mockResponse(DDragonChampionResponse{
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
							Leveltip: leveltip{
								Label: labels,
							},
							EffectBurn: effects,
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
	client := NewApiClient(log.New("testApp"), newTestClient(func(req *http.Request) *http.Response {
		return mockResponse(DDragonChampionResponse{}, 403)
	}))

	championResponse, err := client.GetLoLChampion("someChampionName")

	assert.Equal(t, DDragonChampionResponse{}, championResponse)
	assert.NotNil(t, err)
}

func TestSanitizeChampionName(t *testing.T) {
	championName := sanitizeChampionName("jHiN")
	assert.Equal(t, "Jhin", championName)

	championName = sanitizeChampionName("Aurelionsol")
	assert.Equal(t, "AurelionSol", championName)

	championName = sanitizeChampionName("dRMundo")
	assert.Equal(t, "DrMundo", championName)

	championName = sanitizeChampionName("jarvanIV")
	assert.Equal(t, "JarvanIV", championName)

	championName = sanitizeChampionName("Kogmaw")
	assert.Equal(t, "KogMaw", championName)

	championName = sanitizeChampionName("Leesin")
	assert.Equal(t, "LeeSin", championName)

	championName = sanitizeChampionName("Masteryi")
	assert.Equal(t, "MasterYi", championName)

	championName = sanitizeChampionName("Missfortune")
	assert.Equal(t, "MissFortune", championName)

	championName = sanitizeChampionName("Monkeyking")
	assert.Equal(t, "MonkeyKing", championName)

	championName = sanitizeChampionName("Reksai")
	assert.Equal(t, "RekSai", championName)

	championName = sanitizeChampionName("Tahmkench")
	assert.Equal(t, "TahmKench", championName)

	championName = sanitizeChampionName("Twistedfate")
	assert.Equal(t, "TwistedFate", championName)

	championName = sanitizeChampionName("Xinzhao")
	assert.Equal(t, "XinZhao", championName)
}

func TestGetChampionUrl(t *testing.T) {
	championName := "someChampionName"
	expectedUrl := fmt.Sprintf("%s/%s.json", dDragonLolChampionBaseUrl, championName)

	url := getChampionUrl(championName)

	assert.Equal(t, expectedUrl, url)
}

func TestGetSpellDamage(t *testing.T) {
	spellTest := spell{
		ID: "sampleSpell",
		Leveltip: leveltip{
			Label: []string{"Damage", "Attack Damage", "Cooldown"},
		},
		EffectBurn: []string{"", "95/130/165/200/235", "60/75/90/105/120"},
	}
	expectedDamage := []float32{95, 130, 165, 200, 235}

	riotClient := NewApiClient(log.New("test"), &http.Client{})
	damage := riotClient.getSpellDamage(spellTest)

	assert.Equal(t, expectedDamage, damage)
}

func TestGetSpellDamage_NoDamageLable(t *testing.T) {
	spellTest := spell{
		ID: "sampleSpell",
		Leveltip: leveltip{
			Label: []string{"Attack Damage", "Cooldown"},
		},
		EffectBurn: []string{"", "60/75/90/105/120"},
	}
	expectedDamage := []float32{0, 0, 0, 0, 0}

	riotClient := NewApiClient(log.New("test"), &http.Client{})
	damage := riotClient.getSpellDamage(spellTest)

	assert.Equal(t, expectedDamage, damage)
}
