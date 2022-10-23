package lol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetYMLPath(t *testing.T) {
	championName := "some Champion Name"

	path := getYMLPath(championName)
	expectedPath := BaseChampionPath + "somechampionname" + fileExtension

	assert.Equal(t, expectedPath, path)
}
