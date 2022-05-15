package lol

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetYMLPath(t *testing.T) {
	championName := "some Champion Name"

	path := getYMLPath(championName)
	expectedPath := BaseChampionPath + "somechampionname" + fileExtension

	assert.Equal(t, expectedPath, path)
}
