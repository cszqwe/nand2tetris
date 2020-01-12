package main

import (
	"bufio"
	"fmt"
	"os"
)

// commandType of the hack assembly language.
type commandType int

const (
	a commandType = iota
	c
)

// assemblyCommand interface which would transfer a assemblyCommand to a hack machine language.
type assemblyCommand interface {
	toHackCommand() string
}

func toBinary(a int, length int) string {
	bytes := []byte{}
	for a > 0 {
		bytes = append(bytes, byte('0')+byte(a%2))
		a /= 2
	}
	zeroStr := []byte{}
	for i := len(bytes); i < length; i++ {
		zeroStr = append(zeroStr, '0')
	}
	return string(zeroStr) + reverse(string(bytes))
}

func reverse(s string) string {
	ans := []byte{}
	for i := 0; i < len(s); i++ {
		ans = append(ans, s[len(s)-i-1])
	}
	return string(ans)
}

type hackCommend struct {
	commandType
	aCommand
	cCommand
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	lines := []string{}
	for true {
		line, err := in.ReadBytes('\n')
		if err != nil {
			break
		}
		lineStr := trim(string(line))
		if len(lineStr) == 0 {
			continue
		}
		if len(lineStr) >= 2 && lineStr[:2] == "//" {
			continue
		}
		lines = append(lines, lineStr)
	}
	machineCode := process(lines)
	for _, code := range machineCode {
		fmt.Println(code)
	}
	return
}

func process(lines []string) []string {
	symbolManager := getSymbolManager()
	// process label symbol
	realCode := processLabelSymbol(lines, symbolManager)
	// process symbol in a command
	noSymbolCode := processSymbolsInACommand(realCode, symbolManager)
	// convert to command
	commands := processNoSymbolCode(noSymbolCode)
	machineCode := []string{}
	for _, command := range commands {
		machineCode = append(machineCode, command.toHackCommand())
	}
	return machineCode
}

func processNoSymbolCode(code []string) []assemblyCommand {
	commands := []assemblyCommand{}
	for _, c := range code {
		if c[0] == '@' {
			commands = append(commands, processACommand(c))
		} else {
			commands = append(commands, processCCommand(c))
		}
	}
	return commands
}

func trim(s string) string {
	ans := []byte{}
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\n' || s[i] == 13 {
			continue
		}
		if s[i] == '/' {
			break
		}
		ans = append(ans, s[i])
	}
	return string(ans)
}
