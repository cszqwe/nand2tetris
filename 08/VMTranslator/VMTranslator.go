package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var curClassName string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	name := os.Args[1]
	fi, err := os.Stat(name)
	if err != nil {
		panic(err)
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		processDir(name)
	case mode.IsRegular():
		processFile(name)
	}

	return
}

func processDir(dirName string) {
	if dirName[len(dirName)-1] != '/' {
		dirName += "/"
	}
	names := strings.Split(dirName, "/")
	outputFileName := names[len(names)-2] + ".asm"
	//fmt.Println("debug", outputFileName)
	of, err := os.Create(dirName + outputFileName)
	check(err)
	out := bufio.NewWriter(of)
	defer out.Flush()
	fmt.Fprintf(out, "@256\n")
	fmt.Fprintf(out, "D=A\n")
	fmt.Fprintf(out, "@SP\n")
	fmt.Fprintf(out, "M=D\n")
	callSys := &callCommand{}
	callSys.nArgs = 0
	callSys.funcName = "Sys.init"
	for _, s := range callSys.write() {
		fmt.Fprintln(out, s)
	}
	filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".vm") {
			f, err := os.Open(path)
			check(err)
			in := bufio.NewReader(f)
			curClassName = getClassName(path)
			processSingleFile(in, out)
		}
		return nil
	})
}

func processFile(fileName string) {
	f, err := os.Open(fileName)
	check(err)
	outputFileName := strings.Replace(fileName, ".vm", ".asm", 1)
	of, err := os.Create(outputFileName)
	check(err)
	in := bufio.NewReader(f)
	out := bufio.NewWriter(of)
	curClassName = getClassName(fileName)
	processSingleFile(in, out)
}

func processSingleFile(in *bufio.Reader, out *bufio.Writer) {
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
	for i, command := range commands {
		tmp := command.write()
		fmt.Fprintln(out, fmt.Sprintf("// %s", lines[i]))
		for _, s := range tmp {
			fmt.Fprintln(out, s)
		}
	}
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
	for len(ans) > 0 && (ans[len(ans)-1] == ' ' || ans[len(ans)-1] == '\t') {
		ans = ans[:len(ans)-1]
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
		return pushType2(fmt.Sprintf("%s.%d", curClassName, cmd.offset))
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
		return popType2(fmt.Sprintf("%s.%d", curClassName, cmd.offset))
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

type labelCommand struct {
	labelName string
}

func (cmd *labelCommand) write() []string {
	return []string{fmt.Sprintf("(%s)", cmd.labelName)}
}

type gotoCommand struct {
	labelName string
}

func (cmd *gotoCommand) write() []string {
	return []string{fmt.Sprintf("@%s", cmd.labelName), "0;JMP"}
}

type ifGotoCommand struct {
	labelName string
}

func (cmd *ifGotoCommand) write() []string {
	ans := []string{}
	ans = append(ans, "@SP")
	ans = append(ans, "M=M-1")
	ans = append(ans, "A=M")
	ans = append(ans, "D=M")
	ans = append(ans, "@"+cmd.labelName)
	ans = append(ans, "D;JNE")
	return ans
}

type callCommand struct {
	funcName string
	nArgs    int
}

func (cmd *callCommand) write() []string {
	ans := []string{}
	id := next()
	returnAddress := fmt.Sprintf("RETURN_ADDRESS_%d", id)
	// push ReturnAddress
	ans = append(ans, fmt.Sprintf("@%s", returnAddress))
	ans = append(ans, fmt.Sprintf("D=A"))
	ans = append(ans, "@SP")
	ans = append(ans, "A=M")
	ans = append(ans, "M=D")
	ans = append(ans, "@SP")
	ans = append(ans, "M=M+1")
	// push LCL
	ans = append(ans, fmt.Sprintf("@LCL"))
	ans = append(ans, fmt.Sprintf("D=M"))
	ans = append(ans, "@SP")
	ans = append(ans, "A=M")
	ans = append(ans, "M=D")
	ans = append(ans, "@SP")
	ans = append(ans, "M=M+1")
	// push ARG
	ans = append(ans, fmt.Sprintf("@ARG"))
	ans = append(ans, fmt.Sprintf("D=M"))
	ans = append(ans, "@SP")
	ans = append(ans, "A=M")
	ans = append(ans, "M=D")
	ans = append(ans, "@SP")
	ans = append(ans, "M=M+1")
	// push THIS
	ans = append(ans, fmt.Sprintf("@THIS"))
	ans = append(ans, fmt.Sprintf("D=M"))
	ans = append(ans, "@SP")
	ans = append(ans, "A=M")
	ans = append(ans, "M=D")
	ans = append(ans, "@SP")
	ans = append(ans, "M=M+1")
	// push THAT
	ans = append(ans, fmt.Sprintf("@THAT"))
	ans = append(ans, fmt.Sprintf("D=M"))
	ans = append(ans, "@SP")
	ans = append(ans, "A=M")
	ans = append(ans, "M=D")
	ans = append(ans, "@SP")
	ans = append(ans, "M=M+1")
	// ARG = SP - 5 - nArgs
	ans = append(ans, "@SP")
	ans = append(ans, "D=M")
	ans = append(ans, "@5")
	ans = append(ans, "D=D-A")
	ans = append(ans, fmt.Sprintf("@%d", cmd.nArgs))
	ans = append(ans, "D=D-A")
	ans = append(ans, "@ARG")
	ans = append(ans, "M=D")
	// LCL = SP
	ans = append(ans, "@SP")
	ans = append(ans, "D=M")
	ans = append(ans, "@LCL")
	ans = append(ans, "M=D")
	// goto FunctionName
	ans = append(ans, "@"+cmd.funcName)
	ans = append(ans, "0;JMP")
	//(returnAddress)
	ans = append(ans, fmt.Sprintf("(%s)", returnAddress))
	return ans
}

type funcCommand struct {
	funcName string
	nVars    int
}

func (cmd *funcCommand) write() []string {
	initLabel := fmt.Sprintf("%s_init", cmd.funcName)
	endInitLabel := fmt.Sprintf("%s_end_init", cmd.funcName)
	ans := []string{}
	// (functionName)
	ans = append(ans, fmt.Sprintf("(%s)", cmd.funcName))
	// repeat N Vars Time: push 0
	ans = append(ans, fmt.Sprintf("@%d", cmd.nVars))
	ans = append(ans, "D=A")
	ans = append(ans, "@tmp")
	ans = append(ans, "M=D")
	ans = append(ans, fmt.Sprintf("(%s)", initLabel))
	ans = append(ans, "@tmp")
	ans = append(ans, "D=M")
	ans = append(ans, fmt.Sprintf("@%s", endInitLabel))
	ans = append(ans, "D;JEQ") // if tmp == 0 break
	ans = append(ans, "@tmp")
	ans = append(ans, "M=M-1") // tmp--
	ans = append(ans, "@SP")   // push 0
	ans = append(ans, "A=M")
	ans = append(ans, "M=0")
	ans = append(ans, "@SP")
	ans = append(ans, "M=M+1")
	ans = append(ans, fmt.Sprintf("@%s", initLabel))
	ans = append(ans, "0;JMP")
	ans = append(ans, fmt.Sprintf("(%s)", endInitLabel)) // (endInit)
	return ans
}

type returnCommand struct {
}

func (cmd *returnCommand) write() []string {
	ans := []string{}
	// endFrame = LCL
	ans = append(ans, "@LCL")
	ans = append(ans, "D=M")
	ans = append(ans, "@endFrame")
	ans = append(ans, "M=D")
	// retAddr = *(endFrame - 5)
	ans = append(ans, "@5")
	ans = append(ans, "A=D-A")
	ans = append(ans, "D=M")
	ans = append(ans, "@retAddr")
	ans = append(ans, "M=D")
	// *ARG = pop()
	ans = append(ans, "@SP")
	ans = append(ans, "M=M-1")
	ans = append(ans, "A=M")
	ans = append(ans, "D=M")
	ans = append(ans, "@ARG")
	ans = append(ans, "A=M")
	ans = append(ans, "M=D")
	// SP = ARG+1
	ans = append(ans, "@ARG")
	ans = append(ans, "D=M")
	ans = append(ans, "D=D+1")
	ans = append(ans, "@SP")
	ans = append(ans, "M=D")
	// THAT = *(endFrame - 1)
	ans = append(ans, "@endFrame")
	ans = append(ans, "D=M")
	ans = append(ans, "A=D-1")
	ans = append(ans, "D=M")
	ans = append(ans, "@THAT")
	ans = append(ans, "M=D")
	// THIS = *(endFrame - 2)
	ans = append(ans, "@endFrame")
	ans = append(ans, "D=M")
	ans = append(ans, "@2")
	ans = append(ans, "A=D-A")
	ans = append(ans, "D=M")
	ans = append(ans, "@THIS")
	ans = append(ans, "M=D")
	// ARG = *(endFrame - 3)
	ans = append(ans, "@endFrame")
	ans = append(ans, "D=M")
	ans = append(ans, "@3")
	ans = append(ans, "A=D-A")
	ans = append(ans, "D=M")
	ans = append(ans, "@ARG")
	ans = append(ans, "M=D")
	// LCL = *(endFrame - 4)
	ans = append(ans, "@endFrame")
	ans = append(ans, "D=M")
	ans = append(ans, "@4")
	ans = append(ans, "A=D-A")
	ans = append(ans, "D=M")
	ans = append(ans, "@LCL")
	ans = append(ans, "M=D")
	// goto retAddr
	ans = append(ans, "@retAddr")
	ans = append(ans, "A=M;JMP")
	return ans
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
	//fmt.Println("debug", cmd)
	strs := strings.Split(cmd, " ")
	//fmt.Println("debug", len(strs), strs)

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
	case "label":
		return parseLabel(strs)
	case "goto":
		return parseGoto(strs)
	case "if-goto":
		return parseIfGoto(strs)
	case "function":
		return parseFunction(strs)
	case "call":
		return parseCall(strs)
	case "return":
		return &returnCommand{}
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

func parseLabel(cmd []string) command {
	if len(cmd) != 2 {
		panic(cmd)
	}
	parsedCmd := &labelCommand{}
	parsedCmd.labelName = cmd[1]
	return parsedCmd
}

func parseGoto(cmd []string) command {
	if len(cmd) != 2 {
		panic(cmd)
	}
	parsedCmd := &gotoCommand{}
	parsedCmd.labelName = cmd[1]
	return parsedCmd
}

func parseIfGoto(cmd []string) command {
	if len(cmd) != 2 {
		panic(cmd)
	}
	parsedCmd := &ifGotoCommand{}
	parsedCmd.labelName = cmd[1]
	return parsedCmd
}

func parseFunction(cmd []string) command {
	if len(cmd) != 3 {
		panic(cmd)
	}
	parsedCmd := &funcCommand{}
	parsedCmd.funcName = cmd[1]
	var err error
	parsedCmd.nVars, err = strconv.Atoi(cmd[2])
	check(err)
	return parsedCmd
}

func parseCall(cmd []string) command {
	if len(cmd) != 3 {
		panic(cmd)
	}
	parsedCmd := &callCommand{}
	parsedCmd.funcName = cmd[1]
	var err error
	parsedCmd.nArgs, err = strconv.Atoi(cmd[2])
	check(err)
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

func getClassName(s string) string {
	ss := strings.Split(s, "/")
	return ss[len(ss)-1][:len(ss[len(ss)-1])-3]
}
