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
