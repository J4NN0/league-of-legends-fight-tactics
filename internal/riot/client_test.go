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
	var format, version, championName = "standAloneComplex", "1.0.0", "TestName"
	var tags = []string{"Fighter", "Tank"}
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
}

func TestGetChampionUrl(t *testing.T) {
	championName := "someChampionName"
	expectedUrl := fmt.Sprintf("%s/%s.json", dDragonLolChampionBaseUrl, championName)

	url := getChampionUrl(championName)

	assert.Equal(t, expectedUrl, url)
}
