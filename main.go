package main

import (
	"errors"
	"fmt"
	"strings"
)

const Init = "init"
const Blue = "blue"
const Green = "green"
const Staging = "staging"

var states = []string{
	Init,
	Blue,
	Green,
	Staging,
}

var fromStateCheckers = map[string]func(Command) error{
	Init:    checkFromInitState,
	Blue:    checkFromBlueState,
	Green:   checkFromGreenState,
	Staging: checkFromStagingState,
}

var cmds = []string{
	"i2b",
	"i2g",
	"b2s",
	"g2s",
	"s2b",
	"s2g",
}

type Command struct {
	cmd               string
	fromState         string
	toState           string
	currentBlueCount  int
	currentGreenCount int
	desiredBlueCount  int
	desiredGreenCount int
}

func parseFromAndToFromCmd(cmd string) (string, string, error) {
	if !validCmd(cmd) {
		return "", "", fmt.Errorf("unrecognized cmd %s, should in %s", cmd, strings.Join(cmds, ", "))
	}

	nodes := strings.Split(cmd, "2")
	from := firstString(states, nodes[0], func(state string, fromState string) bool {
		return strings.HasPrefix(state, fromState)
	})
	to := firstString(states, nodes[1], func(state string, toState string) bool {
		return strings.HasPrefix(state, toState)
	})
	return *from, *to, nil
}

func validCmd(cmd string) bool {
	return firstString(cmds, cmd, func(s string, s2 string) bool {
		return s == s2
	}) != nil
}

func firstString(slice []string, argument string, prediction func(string, string) bool) *string {
	for _, item := range slice {
		if prediction(item, argument) {
			return &item
		}
	}
	return nil
}

func checkFromInitState(cmd Command) error {
	if cmd.currentBlueCount != 0 || cmd.currentGreenCount != 0 {
		return fmt.Errorf("current blue and green counts are %d and %d so you cannot from init state", cmd.currentBlueCount, cmd.currentGreenCount)
	}
	return nil
}

func checkFromBlueState(cmd Command) error {
	if cmd.currentBlueCount == 0 && cmd.currentGreenCount == 0 {
		return errors.New("both current blue and green counts are zero, so you are from init state")
	}
	if cmd.currentBlueCount > 0 && cmd.currentGreenCount > 0 {
		return errors.New("both current blue and green counts are greater than zero, so you are from staging state")
	}
	if cmd.currentBlueCount == 0 && cmd.currentGreenCount > 0 {
		return fmt.Errorf("current blue count is 0 and current green count is %d, so you are from green state", cmd.currentGreenCount)
	}
	return nil
}

func checkFromGreenState(cmd Command) error {
	if cmd.currentBlueCount == 0 && cmd.currentGreenCount == 0 {
		return errors.New("both current blue and green counts are zero, so you are from init state")
	}
	if cmd.currentBlueCount > 0 && cmd.currentGreenCount > 0 {
		return errors.New("both current blue and green counts are greater than zero, so you are from staging state")
	}
	if cmd.currentBlueCount > 0 && cmd.currentGreenCount == 0 {
		return fmt.Errorf("current blue count is %d and current green count is 0, so you are from blue state", cmd.currentBlueCount)
	}
	return nil
}

func checkFromStagingState(cmd Command) error {
	if cmd.currentBlueCount == 0 && cmd.currentGreenCount == 0 {
		return errors.New("both current blue and green counts are zero, so you are from init state")
	}
	if cmd.currentBlueCount > 0 && cmd.currentGreenCount == 0 {
		return fmt.Errorf("current blue count is %d and current green count is 0, so you are from blue state", cmd.currentBlueCount)
	}
	if cmd.currentBlueCount == 0 && cmd.currentGreenCount > 0 {
		return fmt.Errorf("current blue count is 0 and current green count is %d, so you are from green state", cmd.currentGreenCount)
	}
	return nil
}

func checkToStagingState(cmd Command) error {
	if cmd.desiredBlueCount == 0 || cmd.desiredGreenCount == 0 {
		return errors.New("either blue or green desired count is zero, are you really want transit to staging?")
	}
	if cmd.fromState == Blue && cmd.currentBlueCount != cmd.desiredBlueCount {
		return fmt.Errorf("current blue count is %d and desired blue count is %d, transit to staging state cannot change both blue and green counts", cmd.currentBlueCount, cmd.desiredBlueCount)
	}
	if cmd.fromState == Green && cmd.currentGreenCount != cmd.desiredGreenCount {
		return fmt.Errorf("current green count is %d and desired green count is %d, transit to staging state cannot change both blue and green counts", cmd.currentGreenCount, cmd.desiredGreenCount)
	}
	return nil
}

func checkToBlueState(cmd Command) error {
	if cmd.desiredBlueCount == 0 {
		return errors.New("cannot transit to blue state with zero blue desired count")
	}
	return nil
}
