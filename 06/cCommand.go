package main

import (
	"fmt"
	"strings"
)

type cDest int

const (
	destNULL cDest = 0
	destM    cDest = 1
	destD    cDest = 2
	destMD   cDest = 3
	destA    cDest = 4
	destAM   cDest = 5
	destAD   cDest = 6
	destAMD  cDest = 7
)

type cJump int

const (
	jumpNULL cJump = 0
	jumpJGT  cJump = 1
	jumpJEQ  cJump = 2
	jumpJGE  cJump = 3
	jumpJLT  cJump = 4
	jumpJNE  cJump = 5
	jumpJLE  cJump = 6
	jumpJMP  cJump = 7
)

type cComp string

const (
	comp0     = "0101010"
	comp1     = "0111111"
	compm1    = "0111010" //-1
	compD     = "0001100"
	compA     = "0110000"
	compND    = "0001101"
	compNA    = "0110001"
	compmD    = "0001111"
	compmA    = "0110011"
	compDp1   = "0011111"
	compAp1   = "0110111"
	compDm1   = "0001110"
	compAm1   = "0110010"
	compDpA   = "0000010"
	compDmA   = "0010011"
	compAmD   = "0000111"
	compDandA = "0000000"
	compDorA  = "0010101"
	compM     = "1110000"
	compNM    = "1110001"
	compmM    = "1110011"
	compMp1   = "1110111"
	compMm1   = "1110010"
	compDpM   = "1000010"
	compDmM   = "1010011"
	compMmD   = "1000111"
	compDandM = "1000000"
	compDorM  = "1010101"
)

type cCommand struct {
	cDest
	cJump
	cComp
}

func (c *cCommand) toHackCommand() string {
	return fmt.Sprintf("111%s%s%s", c.cComp, toBinary(int(c.cDest), 3), toBinary(int(c.cJump), 3))
}

func processCCommand(s string) *cCommand {
	strs1 := strings.Split(s, "=")
	destC := destNULL
	jumpC := jumpNULL
	rest := ""
	if len(strs1) == 2 {
		destStr := strs1[0]
		rest = strs1[1]
		destC = getDestC(destStr)
	} else {
		rest = strs1[0]
	}
	strs2 := strings.Split(rest, ";")
	compStr := strs2[0]
	if len(strs2) == 2 {
		jumpC = getJumpC(strs2[1])
	}
	compC := getCompC(compStr)
	return &cCommand{destC, jumpC, compC}
}

func getCompC(compStr string) cComp {
	switch compStr {
	case "0":
		return comp0
	case "1":
		return comp1
	case "-1":
		return compm1
	case "D":
		return compD
	case "A":
		return compA
	case "!D":
		return compND
	case "!A":
		return compNA
	case "-D":
		return compmD
	case "-A":
		return compmA
	case "D+1":
		return compDp1
	case "A+1":
		return compAp1
	case "D-1":
		return compDm1
	case "A-1":
		return compAm1
	case "D+A":
		return compDpA
	case "D-A":
		return compDmA
	case "A-D":
		return compAmD
	case "D&A":
		return compDandA
	case "D|A":
		return compDorA
	case "M":
		return compM
	case "!M":
		return compNM
	case "-M":
		return compmM
	case "M+1":
		return compMp1
	case "M-1":
		return compMm1
	case "D+M":
		return compDpM
	case "D-M":
		return compDmM
	case "M-D":
		return compMmD
	case "D&M":
		return compDandM
	case "D|M":
		return compDorM
	default:
		panic("invalid command " + compStr)
	}
}

func getDestC(destStr string) cDest {
	switch destStr {
	case "M":
		return destM
	case "D":
		return destD
	case "A":
		return destA
	case "MD":
		return destMD
	case "AD":
		return destAD
	case "AM":
		return destAM
	case "AMD":
		return destAMD
	default:
		panic("invalid command " + destStr)
	}
}

func getJumpC(jumpStr string) cJump {
	switch jumpStr {
	case "JGT":
		return jumpJGT
	case "JEQ":
		return jumpJEQ
	case "JGE":
		return jumpJGE
	case "JLE":
		return jumpJLE
	case "JLT":
		return jumpJLT
	case "JNE":
		return jumpJNE
	case "JMP":
		return jumpJMP
	default:
		panic("invalid command " + jumpStr)
	}
}
