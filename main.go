package main

import (
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
