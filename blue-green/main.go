package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const Init = "init"
const Blue = "blue"
const Green = "green"
const Staging = "staging"
const Destroy = "destroy"

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

var toStateCheckers = map[string]func(Command) error{
	Blue:    checkToBlueState,
	Green:   checkToGreenState,
	Staging: checkToStagingState,
}

var cmds = []string{
	"i2b",
	"i2g",
	"b2s",
	"g2s",
	"s2b",
	"s2g",
	"b2b",
	"g2g",
	"destroy",
}

func main() {
	command, err := parseCommand()
	executedCmd, err := command.execute()
	if err != nil {
		exitError(err)
	}
	executedCmd.DesiredBlueCount = strconv.Itoa(executedCmd.desiredBlueCount)
	executedCmd.DesiredGreenCount = strconv.Itoa(executedCmd.desiredGreenCount)
	outputCommand(executedCmd, err)
}

func parseCommand() (Command, error) {
	cmd := flag.String("operation", "", "")
	currentBlueCount := flag.Int("currentBlue", 0, "")
	currentGreenCount := flag.Int("currentGreen", 0, "")
	desiredBlueCount := flag.Int("desiredBlue", 0, "")
	desiredGreenCount := flag.Int("desiredGreen", 0, "")
	flag.Parse()
	if *cmd == Destroy {
		//desired blue and green count equal to current blue and green count in case of we use "destroy" operation but use terraform apply by mistake
		outputCommand(Command{
			DesiredBlueCount:  strconv.Itoa(*currentBlueCount),
			DesiredGreenCount: strconv.Itoa(*currentGreenCount),
		}, nil)
		os.Exit(0)
	}
	from, to, err := parseFromAndToFromCmd(*cmd)
	if err != nil {
		exitError(err)
	}
	command := Command{
		cmd:               *cmd,
		fromState:         from,
		toState:           to,
		currentBlueCount:  *currentBlueCount,
		currentGreenCount: *currentGreenCount,
		desiredBlueCount:  *desiredBlueCount,
		desiredGreenCount: *desiredGreenCount,
	}
	return command, err
}

func outputCommand(executedCmd Command, err error) {
	bytes, err := json.Marshal(executedCmd)
	if err != nil {
		exitError(err)
	}
	fmt.Print(string(bytes))
}

func exitError(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(-1)
}

type Command struct {
	cmd                  string
	fromState            string
	toState              string
	currentBlueCount     int
	currentGreenCount    int
	desiredBlueCount     int
	desiredGreenCount    int
	DesiredBlueCount     string `json:"desiredBlueCount"`
	DesiredGreenCount    string `json:"desiredGreenCount"`
	LoadBalanceDirection string `json:"loadBalanceDirection"`
}

func (cmd Command) execute() (Command, error) {
	if err := fromStateCheckers[cmd.fromState](cmd); err != nil {
		return Command{}, err
	}
	if err := toStateCheckers[cmd.toState](cmd); err != nil {
		return Command{}, err
	}
	return cmd.setLoadBalanceDirection(), nil
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
	if cmd.desiredBlueCount <= 0 || cmd.desiredGreenCount <= 0 {
		return errors.New("either blue or green desired count is lesser or equal to zero, potentially misconfig")
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
	if cmd.desiredBlueCount <= 0 {
		return errors.New("cannot transit to blue state with zero or lesser blue desired count")
	}
	if cmd.desiredGreenCount != 0 {
		return fmt.Errorf("desired green count is %d not zero, potentially misconfig", cmd.desiredGreenCount)
	}
	return nil
}

func checkToGreenState(cmd Command) error {
	if cmd.desiredGreenCount <= 0 {
		return errors.New("cannot transit to green state with zero or lesser green desired count")
	}
	if cmd.desiredBlueCount != 0 {
		return fmt.Errorf("desired blue count is %d not zero, potentially misconfig", cmd.desiredBlueCount)
	}
	return nil
}

func (cmd Command) setLoadBalanceDirection() Command {
	if cmd.toState == Staging {
		cmd.LoadBalanceDirection = cmd.fromState
	} else {
		cmd.LoadBalanceDirection = cmd.toState
	}
	return cmd
}
