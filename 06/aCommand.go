package main

import (
	"fmt"
	"strconv"
)

type aCommand struct {
	value int //@value
}

func (a *aCommand) toHackCommand() string {
	valueStr := toBinary(a.value, 15)
	return fmt.Sprintf("0%s", valueStr)
}

func processACommand(s string) *aCommand {
	num, _ := strconv.Atoi(s[1:])
	return &aCommand{num}
}
