package controller

import (
	"errors"
	"fmt"
	"testing"

	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/logger/loggertest"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/lol"
	lolMocks "github.com/J4NN0/league-of-legends-fight-tactics/pkg/lol/mocks"
	riotMocks "github.com/J4NN0/league-of-legends-fight-tactics/pkg/riot/mocks"
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

func TestGetRoundSpellsToString(t *testing.T) {
	var benchmark, hp = 3.0, 15.0
	spells := []lol.Spell{
		{
			ID:       "q",
			Damage:   []float64{6, 7, 8, 9, 10},
			MaxRank:  5,
			Cooldown: []float64{4, 3, 2, 1, 1},
		},
		{
			ID:       "w",
			Damage:   []float64{6, 7, 8, 9, 20},
			MaxRank:  5,
			Cooldown: []float64{5, 4, 3, 2, 2},
		},
	}

	spellsToString := getRoundSpellsToString(spells, hp, benchmark)
	expectedString := fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", spells[0].ID, spells[0].Damage, hp, hp-spells[0].Damage[spells[0].MaxRank-1])
	expectedString += fmt.Sprintf("%s: %.2f (hp: %.2f -> %.2f)\n", spells[1].ID, spells[1].Damage, hp-spells[0].Damage[spells[0].MaxRank-1], hp-spells[0].Damage[spells[0].MaxRank-1]-spells[1].Damage[spells[0].MaxRank-1])
	expectedString += fmt.Sprintf("\nEnemy defeated in %.2fs\n", benchmark)

	assert.Equal(t, expectedString, spellsToString)
}

func TestSetFilePath(t *testing.T) {
	filename := setFilePath(lol.Champion{Name: "Name1"}, lol.Champion{Name: "Name2"})

	assert.Equal(t, "fights/Name1_vs_Name2.loltactics", filename)
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
