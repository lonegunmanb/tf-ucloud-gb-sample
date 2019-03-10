package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
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

func TestCheckFromInitStateWithNonZeroCurrentCountShouldReturnError(t *testing.T) {
	incorrectCmds := []Command{
		{currentBlueCount: 1, currentGreenCount: 0},
		{currentBlueCount: 0, currentGreenCount: 1},
	}
	for _, incorrectCmd := range incorrectCmds {
		err := checkFromInitState(incorrectCmd)
		assert.NotNil(t, err)
	}
}

func TestCheckFromInitStateWithZeroBlueAndGreenCountShouldSuccess(t *testing.T) {
	cmd := Command{
		currentBlueCount:  0,
		currentGreenCount: 0,
	}
	err := checkFromInitState(cmd)
	assert.Nil(t, err)
}

func TestCheckFromBlueStateWithIncorrectStatusShouldReturnError(t *testing.T) {
	incorrectCmds := []struct {
		cmd              Command
		expectedErrorMsg string
	}{
		{cmd: Command{currentBlueCount: 0, currentGreenCount: 0}, expectedErrorMsg: Init},
		{cmd: Command{currentBlueCount: 0, currentGreenCount: 1}, expectedErrorMsg: Green},
		{cmd: Command{currentBlueCount: 1, currentGreenCount: 1}, expectedErrorMsg: Staging},
	}
	for _, input := range incorrectCmds {
		err := checkFromBlueState(input.cmd)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), input.expectedErrorMsg))
	}
}

func TestCheckFromBlueStateWithCorrectStatusShouldReturnNil(t *testing.T) {
	cmd := Command{
		currentBlueCount:  1,
		currentGreenCount: 0,
	}
	err := checkFromBlueState(cmd)
	assert.Nil(t, err)
}

func TestCheckFromGreenStateWithIncorrectStatusShouldReturnError(t *testing.T) {
	incorrectCmds := []struct {
		cmd              Command
		expectedErrorMsg string
	}{
		{cmd: Command{currentBlueCount: 0, currentGreenCount: 0}, expectedErrorMsg: Init},
		{cmd: Command{currentBlueCount: 1, currentGreenCount: 0}, expectedErrorMsg: Blue},
		{cmd: Command{currentBlueCount: 1, currentGreenCount: 1}, expectedErrorMsg: Staging},
	}
	for _, input := range incorrectCmds {
		err := checkFromGreenState(input.cmd)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), input.expectedErrorMsg))
	}
}

func TestCheckFromGreenStateWithCorrectStatusShouldReturnNil(t *testing.T) {
	cmd := Command{
		currentBlueCount:  0,
		currentGreenCount: 1,
	}
	err := checkFromGreenState(cmd)
	assert.Nil(t, err)
}

func TestCheckFromStagingStateWithIncorrectStatusShouldReturnError(t *testing.T) {
	incorrectCmds := []struct {
		cmd              Command
		expectedErrorMsg string
	}{
		{cmd: Command{currentBlueCount: 0, currentGreenCount: 0}, expectedErrorMsg: Init},
		{cmd: Command{currentBlueCount: 1, currentGreenCount: 0}, expectedErrorMsg: Blue},
		{cmd: Command{currentBlueCount: 0, currentGreenCount: 1}, expectedErrorMsg: Green},
	}
	for _, input := range incorrectCmds {
		err := checkFromStagingState(input.cmd)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), input.expectedErrorMsg))
	}
}

func TestCheckFromStagingStateWithCorrectStatusShouldReturnNil(t *testing.T) {
	cmd := Command{
		currentBlueCount:  1,
		currentGreenCount: 1,
	}
	err := checkFromStagingState(cmd)
	assert.Nil(t, err)
}

func TestCheckToStagingStateWithIncorrectStatusShouldReturnError(t *testing.T) {
	incorrectCmds := []Command{
		{desiredBlueCount: 0, desiredGreenCount: 0},
		{desiredBlueCount: 0, desiredGreenCount: 1},
		{desiredBlueCount: 1, desiredGreenCount: 0},
		{desiredBlueCount: 1, desiredGreenCount: 1, currentBlueCount: 2, currentGreenCount: 0, fromState: Blue},
		{desiredBlueCount: 1, desiredGreenCount: 1, currentBlueCount: 0, currentGreenCount: 2, fromState: Green},
	}
	for _, cmd := range incorrectCmds {
		err := checkToStagingState(cmd)
		assert.NotNil(t, err)
	}
}

func TestCheckToStagingStateWithCorrectStatusShouldReturnNil(t *testing.T) {
	correctCmds := []Command{
		{desiredBlueCount: 1, desiredGreenCount: 1, currentBlueCount: 1, currentGreenCount: 0, fromState: Blue},
		{desiredBlueCount: 1, desiredGreenCount: 1, currentBlueCount: 0, currentGreenCount: 1, fromState: Green},
	}
	for _, cmd := range correctCmds {
		err := checkToStagingState(cmd)
		assert.Nil(t, err)
	}
}

func TestCheckToBlueWithIncorrectStatusShouldReturnError(t *testing.T) {
	incorrectCmd := Command{
		desiredBlueCount: 0,
	}
	err := checkToBlueState(incorrectCmd)
	assert.NotNil(t, err)
}

func TestCheckToBlueWithCorrectStatusShouldReturnNil(t *testing.T) {
	cmds := []Command{
		{desiredBlueCount: 1, currentGreenCount: 1, currentBlueCount: 1, fromState: Staging},
		{desiredBlueCount: 1, currentGreenCount: 0, currentBlueCount: 0, fromState: Init},
	}
	for _, cmd := range cmds {
		err := checkToBlueState(cmd)
		assert.Nil(t, err)
	}
}
