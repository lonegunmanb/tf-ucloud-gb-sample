package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFromAndToStatesFromCmd(t *testing.T) {
	inputs := []struct {
		cmd        string
		expectFrom string
		expectTo   string
	}{
		{"i2b", Init, Blue},
		{"i2g", Init, Green},
		{"b2s", Blue, Staging},
		{"g2s", Green, Staging},
		{"s2b", Staging, Blue},
		{"s2g", Staging, Green},
	}

	for _, input := range inputs {
		from, to, err := parseFromAndToFromCmd(input.cmd)
		assert.Nil(t, err)
		assert.Equal(t, input.expectFrom, from)
		assert.Equal(t, input.expectTo, to)
	}
}

func TestParseIncorrectCmdShouldGetError(t *testing.T) {
	incorrectCmds := []string{
		"i2i",
		"b2g",
		"g2b",
		"s2s",
		"whatever",
		"",
	}

	for _, cmd := range incorrectCmds {
		_, _, err := parseFromAndToFromCmd(cmd)
		assert.NotNil(t, err)
	}
}

func TestFromInitStateWithNonZeroCurrentCountShouldReturnError(t *testing.T) {
	incorrectCmds := []Command{
		{currentBlueCount: 1, currentGreenCount: 0},
		{currentBlueCount: 0, currentGreenCount: 1},
	}
	for _, incorrectCmd := range incorrectCmds {
		err := checkFromInitState(incorrectCmd)
		assert.NotNil(t, err)
	}
}

func TestFromInitStateWithZeroBlueAndGreenCountShouldSuccess(t *testing.T) {
	cmd := Command{
		currentBlueCount:  0,
		currentGreenCount: 0,
	}
	err := checkFromInitState(cmd)
	assert.Nil(t, err)
}
