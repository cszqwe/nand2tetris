package main

import (
	"strconv"
	"strings"
)

func parse(strs []string) []command {
	commands := []command{}
	for _, s := range strs {
		commands = append(commands, parseSingle(s))
	}
	return commands
}

func parseSingle(cmd string) command {
	strs := strings.Split(cmd, " ")
	switch strs[0] {
	case "push":
		return parsePush(strs)
	case "pop":
		return parsePop(strs)
	case "add":
		return &addCommand{}
	case "sub":
		return &subCommand{}
	case "and":
		return &andCommand{}
	case "or":
		return &orCommand{}
	case "not":
		return &notCommand{}
	case "neg":
		return &negCommand{}
	case "lt":
		return &ltCommand{}
	case "gt":
		return &gtCommand{}
	case "eq":
		return &eqCommand{}
	default:
		panic("invalid command: " + strs[0])
	}
}

func parsePush(cmd []string) command {
	if len(cmd) != 3 {
		panic(cmd)
	}
	parsedCmd := &pushCommand{}
	parsedCmd.symbol = parseSymbol(cmd[1])
	var err error
	parsedCmd.offset, err = strconv.Atoi(cmd[2])
	if err != nil {
		panic(err)
	}
	return parsedCmd
}

func parsePop(cmd []string) command {
	if len(cmd) != 3 {
		panic(cmd)
	}
	parsedCmd := &popCommand{}
	parsedCmd.symbol = parseSymbol(cmd[1])
	var err error
	parsedCmd.offset, err = strconv.Atoi(cmd[2])
	if err != nil {
		panic(err)
	}
	return parsedCmd
}

func parseSymbol(s string) symbolType {
	switch s {
	case "static":
		return static
	case "local":
		return local
	case "this":
		return this
	case "that":
		return that
	case "argument":
		return argument
	case "temp":
		return temp
	case "pointer":
		return pointer
	case "constant":
		return constant
	default:
		panic("invalid symbol: " + s)
	}
}
