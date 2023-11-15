package command

import (
	"errors"
	"testing"

	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/logger/loggertest"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/lol"
	lolMocks "github.com/J4NN0/league-of-legends-fight-tactics/pkg/lol/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChampionsFight(t *testing.T) {
	t.Run("fail first Read", func(t *testing.T) {
		mockLol := &lolMocks.Tactics{}
		mockLol.On("ReadChampion", mock.AnythingOfType("string")).Return(lol.Champion{}, errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, nil, mockLol)

		err := ctrl.championsFight("mockName1", "mockName2")

		assert.NotNil(t, err)
	})

	t.Run("fail second Read", func(t *testing.T) {
		mockLol := &lolMocks.Tactics{}
		mockLol.On("ReadChampion", mock.AnythingOfType("string")).Once().Return(getMockLoLChampion(), nil)
		mockLol.On("ReadChampion", mock.AnythingOfType("string")).Once().Return(lol.Champion{}, errors.New("some error"))

		ctrl := New(&loggertest.Logger{}, nil, mockLol)

		err := ctrl.championsFight("mockName1", "mockName2")

		assert.NotNil(t, err)
	})
}
