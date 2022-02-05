package lol

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetYMLPath(t *testing.T) {
	championName := "someChampionName"

	path := getYMLPath(championName)
	expectedPath := BaseChampionPath + strings.ToLower(championName) + fileExtension

	assert.Equal(t, expectedPath, path)
}
