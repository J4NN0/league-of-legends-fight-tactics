package command

import (
	"errors"
	"testing"

	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/logger/loggertest"
	lolMocks "github.com/J4NN0/league-of-legends-fight-tactics/pkg/lol/mocks"
	riotMocks "github.com/J4NN0/league-of-legends-fight-tactics/pkg/riot/mocks"
	"github.com/KnutZuidema/golio/datadragon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchChampion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRiot := &riotMocks.Client{}
		mockRiot.On("GetLoLChampion", mock.AnythingOfType("string")).Once().Return(getMockDDChampion(), nil)
		mockLol := &lolMocks.Tactics{}
		mockLol.On("WriteChampion", mock.AnythingOfType("lol.Champion"), mock.AnythingOfType("string")).Once().Return(nil)

		ctrl := New(&loggertest.Logger{}, mockRiot, mockLol)

		err := ctrl.fetchChampion("mockName")

		assert.Nil(t, err)
	})

	t.Run("fail GetLoLChampion", func(t *testing.T) {
		mockRiot := &riotMocks.Client{}
		mockRiot.On("GetLoLChampion", mock.AnythingOfType("string")).Once().Return(getMockDDChampion(), errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, mockRiot, nil)

		err := ctrl.fetchChampion("mockName")

		assert.NotNil(t, err)
	})

	t.Run("fail Write", func(t *testing.T) {
		mockRiot := &riotMocks.Client{}
		mockRiot.On("GetLoLChampion", mock.AnythingOfType("string")).Once().Return(getMockDDChampion(), nil)
		mockLol := &lolMocks.Tactics{}
		mockLol.On("WriteChampion", mock.AnythingOfType("lol.Champion"), mock.AnythingOfType("string")).Once().Return(errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, mockRiot, mockLol)

		err := ctrl.fetchChampion("mockName")

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

		err := ctrl.fetchAllChampions()

		assert.Nil(t, err)
	})

	t.Run("fail GetAllLoLChampions", func(t *testing.T) {
		mockRiot := &riotMocks.Client{}
		mockRiot.On("GetAllLoLChampions").Once().Return(nil, errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, mockRiot, nil)

		err := ctrl.fetchAllChampions()

		assert.NotNil(t, err)
	})
}
