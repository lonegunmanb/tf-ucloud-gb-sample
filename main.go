package main

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"go-funk"
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
	if !funk.Contains(cmds, cmd) {
		return "", "", fmt.Errorf("unrecognized cmd %s, should in %s", cmd, strings.Join(cmds, ", "))
	}

	nodes := strings.Split(cmd, "2")
	from := linq.From(states).FirstWith(func(state interface{}) bool {
		return strings.HasPrefix(state.(string), nodes[0])
	}).(string)
	to := linq.From(states).FirstWith(func(state interface{}) bool {
		return strings.HasPrefix(state.(string), nodes[1])
	}).(string)
	return from, to, nil
}
