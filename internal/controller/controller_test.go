package controller

import (
	"errors"
	"testing"

	"github.com/J4NN0/league-of-legends-fight-tactics/internal/logger/loggertest"
	"github.com/J4NN0/league-of-legends-fight-tactics/internal/lol"
	lolMocks "github.com/J4NN0/league-of-legends-fight-tactics/internal/lol/mocks"
	riotMocks "github.com/J4NN0/league-of-legends-fight-tactics/internal/riot/mocks"
	"github.com/KnutZuidema/golio/datadragon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getMockLoLChampion() lol.Champion {
	return lol.Champion{
		ID:    "mockID",
		Name:  "mockName",
		Title: "mockTitle",
		Tags:  "Some, tags, here",
		Passive: lol.Passive{
			Name:        "passiveName",
			Description: "passiveDescription",
		},
		Stats: lol.Stats{
			HealthPoints: 50,
			AttackDamage: 10,
			AttackSpeed:  2,
		},
		Spells: []lol.Spell{
			{
				ID:       "aa",
				Name:     "Auto Attack",
				MaxRank:  1,
				Damage:   []float64{10},
				Cooldown: []float64{2},
				Cast:     0,
			},
			{
				ID:       "q",
				Name:     "QName",
				MaxRank:  5,
				Cooldown: []float64{10, 8, 6, 4, 2},
				Damage:   []float64{8, 10, 12, 14, 16},
				Cast:     0,
			},
		},
	}
}

func getMockDDChampion() datadragon.ChampionDataExtended {
	return datadragon.ChampionDataExtended{
		ChampionData: datadragon.ChampionData{
			ID:    "mockID",
			Name:  "mockName",
			Title: "mockTitle",
			Tags:  []string{"Some", "tags", "here"},
			Stats: datadragon.ChampionDataStats{
				HealthPoints:      50,
				AttackDamage:      10,
				AttackSpeedOffset: 2,
			},
		},
		Passive: datadragon.PassiveData{
			Name:        "passiveName",
			Description: "passiveDescription",
		},
		Spells: []datadragon.SpellData{
			{
				ID:       "q",
				Name:     "QName",
				MaxRank:  5,
				Cooldown: []float64{10, 8, 6, 4, 2},
				Effect:   [][]float64{nil, {8, 10, 12, 14, 16}},
			},
		},
	}
}

func TestChampionsFight(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLol := &lolMocks.Tactics{}
		mockLol.On("ReadChampion", mock.AnythingOfType("string")).Return(getMockLoLChampion(), nil)
		mockLol.On("Fight", mock.AnythingOfType("lol.Champion"), mock.AnythingOfType("lol.Champion")).Once().Return()

		ctrl := New(&loggertest.Logger{}, nil, mockLol)

		err := ctrl.ChampionsFight("mockName1", "mockName2")

		assert.Nil(t, err)
	})

	t.Run("fail first Read", func(t *testing.T) {
		mockLol := &lolMocks.Tactics{}
		mockLol.On("ReadChampion", mock.AnythingOfType("string")).Return(lol.Champion{}, errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, nil, mockLol)

		err := ctrl.ChampionsFight("mockName1", "mockName2")

		assert.NotNil(t, err)
	})

	t.Run("fail second Read", func(t *testing.T) {
		mockLol := &lolMocks.Tactics{}
		mockLol.On("ReadChampion", mock.AnythingOfType("string")).Once().Return(getMockLoLChampion(), nil)
		mockLol.On("ReadChampion", mock.AnythingOfType("string")).Once().Return(lol.Champion{}, errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, nil, mockLol)

		err := ctrl.ChampionsFight("mockName1", "mockName2")

		assert.NotNil(t, err)
	})
}

func TestAllChampionsFight(t *testing.T) {
	mockLol := &lolMocks.Tactics{}
	mockLol.On("ReadChampion", mock.AnythingOfType("string")).Return(getMockLoLChampion(), nil)
	mockLol.On("Fight", mock.AnythingOfType("lol.Champion"), mock.AnythingOfType("lol.Champion")).Once().Return()

	ctrl := New(&loggertest.Logger{}, nil, mockLol)

	err := ctrl.AllChampionsFight()

	assert.Nil(t, err)
}

func TestFetchChampion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRiot := &riotMocks.Client{}
		mockRiot.On("GetLoLChampion", mock.AnythingOfType("string")).Once().Return(getMockDDChampion(), nil)
		mockLol := &lolMocks.Tactics{}
		mockLol.On("WriteChampion", mock.AnythingOfType("lol.Champion"), mock.AnythingOfType("string")).Once().Return(nil)

		ctrl := New(&loggertest.Logger{}, mockRiot, mockLol)

		err := ctrl.FetchChampion("mockName")

		assert.Nil(t, err)
	})

	t.Run("fail GetLoLChampion", func(t *testing.T) {
		mockRiot := &riotMocks.Client{}
		mockRiot.On("GetLoLChampion", mock.AnythingOfType("string")).Once().Return(getMockDDChampion(), errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, mockRiot, nil)

		err := ctrl.FetchChampion("mockName")

		assert.NotNil(t, err)
	})

	t.Run("fail Write", func(t *testing.T) {
		mockRiot := &riotMocks.Client{}
		mockRiot.On("GetLoLChampion", mock.AnythingOfType("string")).Once().Return(getMockDDChampion(), nil)
		mockLol := &lolMocks.Tactics{}
		mockLol.On("WriteChampion", mock.AnythingOfType("lol.Champion"), mock.AnythingOfType("string")).Once().Return(errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, mockRiot, mockLol)

		err := ctrl.FetchChampion("mockName")

		assert.NotNil(t, err)
	})
}

func TestFetchAllChampions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRiot := &riotMocks.Client{}
		mockRiot.On("GetAllLoLChampions").Once().Return([]datadragon.ChampionDataExtended{getMockDDChampion()}, nil)
		mockLol := &lolMocks.Tactics{}
		mockLol.On("WriteChampion", mock.AnythingOfType("lol.Champion"), mock.AnythingOfType("string")).Return(nil)

		ctrl := New(&loggertest.Logger{}, mockRiot, mockLol)

		err := ctrl.FetchAllChampions()

		assert.Nil(t, err)
	})

	t.Run("fail GetAllLoLChampions", func(t *testing.T) {
		mockRiot := &riotMocks.Client{}
		mockRiot.On("GetAllLoLChampions").Once().Return(nil, errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, mockRiot, nil)

		err := ctrl.FetchAllChampions()

		assert.NotNil(t, err)
	})
}

func TestStoreChampionToYMLFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLol := &lolMocks.Tactics{}
		mockLol.On("WriteChampion", mock.AnythingOfType("lol.Champion"), mock.AnythingOfType("string")).Return(nil)

		ctrl := New(&loggertest.Logger{}, nil, mockLol)

		err := ctrl.storeChampionToYMLFile(getMockDDChampion())

		assert.Nil(t, err)
	})

	t.Run("fail Write", func(t *testing.T) {
		mockLol := &lolMocks.Tactics{}
		mockLol.On("WriteChampion", mock.AnythingOfType("lol.Champion"), mock.AnythingOfType("string")).Return(errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, nil, mockLol)

		err := ctrl.storeChampionToYMLFile(getMockDDChampion())

		assert.NotNil(t, err)
	})
}

func TestMapChampionResponseToLolChampionStruct(t *testing.T) {
	lolChampion := mapChampionResponseToLolChampionStruct(getMockDDChampion())
	assert.Equal(t, getMockLoLChampion(), lolChampion)
}

func TestGetYMLPath(t *testing.T) {
	t.Run("lowercase name", func(t *testing.T) {
		path := getYMLPath("name")
		expectedPath := baseChampionPath + "/name." + fileExtension

		assert.Equal(t, expectedPath, path)
	})

	t.Run("uppercase name", func(t *testing.T) {
		path := getYMLPath("NAME")
		expectedPath := baseChampionPath + "/name." + fileExtension

		assert.Equal(t, expectedPath, path)
	})

	t.Run("multi case plus spaces", func(t *testing.T) {
		path := getYMLPath("SomE nAme")
		expectedPath := baseChampionPath + "/somename." + fileExtension

		assert.Equal(t, expectedPath, path)
	})
}
