package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fileName := os.Args[1]
	outputFileName := strings.Replace(fileName, ".vm", ".asm", 1)
	f, err := os.Open(fileName)
	check(err)
	of, err := os.Create(outputFileName)
	check(err)
	in := bufio.NewReader(f)
	out := bufio.NewWriter(of)
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
	commands := parse(lines)
	for _, command := range commands {
		tmp := command.write()
		//fmt.Fprintln(out, command)
		for _, s := range tmp {
			fmt.Fprintln(out, s)
		}
	}
	return
}

func process(lines []string) []string {

	return nil
}

func trim(s string) string {
	ans := []byte{}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == 13 {
			continue
		}
		if s[i] == '/' {
			break
		}
		ans = append(ans, s[i])
	}
	return string(ans)
}

// commandType of the hack assembly language.
type commandType int

var seq int

const (
	push commandType = iota
	pop
	add
	sub
	and
	or
	lt
	gt
	eq
	neg
	not
)

type command interface {
	write() []string
}

type symbolType int

const (
	argument symbolType = iota
	local
	this
	that
	static
	pointer
	constant
	temp
)

type pushCommand struct {
	symbol symbolType
	offset int
}

func (cmd *pushCommand) write() []string {
	switch cmd.symbol {
	case argument:
		return pushType1("ARG", cmd.offset)
	case local:
		return pushType1("LCL", cmd.offset)
	case this:
		return pushType1("THIS", cmd.offset)
	case that:
		return pushType1("THAT", cmd.offset)
	case temp:
		return pushType2(fmt.Sprintf("%d", cmd.offset+5))
	case static:
		return pushType2(fmt.Sprintf("%d", cmd.offset+16))
	case pointer:
		if cmd.offset == 0 {
			return pushType2("THIS")
		}
		return pushType2("THAT")
	default:
		ans := []string{}
		ans = append(ans, fmt.Sprintf("@%d", cmd.offset))
		ans = append(ans, "D=A")
		ans = append(ans, "@SP")
		ans = append(ans, "A=M")
		ans = append(ans, "M=D")
		ans = append(ans, "@SP")
		ans = append(ans, "M=M+1")
		return ans
	}
}

// For this, that, arg, local
func pushType1(symbol string, offset int) []string {
	ans := []string{}
	ans = append(ans, fmt.Sprintf("@%s", symbol))
	ans = append(ans, "D=M")
	ans = append(ans, fmt.Sprintf("@%d", offset))
	ans = append(ans, "A=D+A")
	ans = append(ans, "D=M")
	ans = append(ans, "@SP")
	ans = append(ans, "A=M")
	ans = append(ans, "M=D")
	ans = append(ans, "@SP")
	ans = append(ans, "M=M+1")
	return ans
}

// For this, that, arg, local
func pushType2(symbol string) []string {
	ans := []string{}
	ans = append(ans, fmt.Sprintf("@%s", symbol))
	ans = append(ans, "D=M")
	ans = append(ans, "@SP")
	ans = append(ans, "A=M")
	ans = append(ans, "M=D")
	ans = append(ans, "@SP")
	ans = append(ans, "M=M+1")
	return ans
}

type popCommand struct {
	symbol symbolType
	offset int
}

func (cmd *popCommand) write() []string {
	switch cmd.symbol {
	case argument:
		return popType1("ARG", cmd.offset)
	case local:
		return popType1("LCL", cmd.offset)
	case this:
		return popType1("THIS", cmd.offset)
	case that:
		return popType1("THAT", cmd.offset)
	case temp:
		return popType2(fmt.Sprintf("%d", cmd.offset+5))
	case static:
		return popType2(fmt.Sprintf("%d", cmd.offset+16))
	case pointer:
		if cmd.offset == 0 {
			return popType2("THIS")
		}
		return popType2("THAT")
	default:
		panic("cannot pop to constnat")
	}
}

// For this, that, arg, local
func popType1(symbol string, offset int) []string {
	ans := []string{}
	ans = append(ans, "@SP")
	ans = append(ans, "M=M-1")
	ans = append(ans, fmt.Sprintf("@%s", symbol))
	ans = append(ans, "D=M")
	ans = append(ans, fmt.Sprintf("@%d", offset))
	ans = append(ans, "D=D+A")
	ans = append(ans, "@Addr")
	ans = append(ans, "M=D")
	ans = append(ans, "@SP")
	ans = append(ans, "A=M")
	ans = append(ans, "D=M")
	ans = append(ans, "@Addr")
	ans = append(ans, "A=M")
	ans = append(ans, "M=D")
	return ans
}

// For this, that, arg, local
func popType2(symbol string) []string {
	ans := []string{}
	ans = append(ans, "@SP")
	ans = append(ans, "M=M-1")
	ans = append(ans, "A=M")
	ans = append(ans, "D=M")
	ans = append(ans, fmt.Sprintf("@%s", symbol))
	ans = append(ans, "M=D")
	return ans
}

type addCommand struct {
}

func (cmd *addCommand) write() []string {
	return arthiWithOneOP("M=M+D")
}

func arthiWithOneOP(op string) []string {
	ans := []string{}
	ans = append(ans, "@SP")
	ans = append(ans, "M=M-1")
	ans = append(ans, "A=M")
	ans = append(ans, "D=M")
	ans = append(ans, "A=A-1")
	ans = append(ans, op)
	return ans
}

type subCommand struct {
}

func (cmd *subCommand) write() []string {
	return arthiWithOneOP("M=M-D")
}

type andCommand struct {
}

func (cmd *andCommand) write() []string {
	return arthiWithOneOP("M=M&D")
}

type orCommand struct {
}

func (cmd *orCommand) write() []string {
	return arthiWithOneOP("M=M|D")
}

type notCommand struct {
}

func (cmd *notCommand) write() []string {
	ans := []string{}
	ans = append(ans, "@SP")
	ans = append(ans, "A=M-1")
	ans = append(ans, "M=!M")
	return ans
}

type negCommand struct {
}

func (cmd *negCommand) write() []string {
	ans := []string{}
	ans = append(ans, "@SP")
	ans = append(ans, "A=M-1")
	ans = append(ans, "M=-M")
	return ans
}

type ltCommand struct {
}

func (cmd *ltCommand) write() []string {
	return jumpCommand("D;JLT")
}

type gtCommand struct {
}

func (cmd *gtCommand) write() []string {
	return jumpCommand("D;JGT")
}

type eqCommand struct {
}

func (cmd *eqCommand) write() []string {
	return jumpCommand("D;JEQ")
}

func next() int {
	seq++
	return seq
}

func jumpCommand(condition string) []string {
	ans := []string{}
	// By default is true.
	ans = append(ans, "@tmp")
	ans = append(ans, "M=-1")
	ans = append(ans, "@SP")
	ans = append(ans, "M=M-1")
	ans = append(ans, "A=M")
	ans = append(ans, "D=M")
	ans = append(ans, "A=A-1")
	ans = append(ans, "D=M-D")
	id := next()
	ans = append(ans, fmt.Sprintf("@CONTINUE%d", id))
	ans = append(ans, condition)
	ans = append(ans, "@tmp")
	ans = append(ans, "M=0")
	ans = append(ans, fmt.Sprintf("(CONTINUE%d)", id))
	ans = append(ans, "@tmp")
	ans = append(ans, "D=M")
	ans = append(ans, "@SP")
	ans = append(ans, "A=M-1")
	ans = append(ans, "M=D")
	return ans
}

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
