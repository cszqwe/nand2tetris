package main

import "fmt"

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
