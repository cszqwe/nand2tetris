package main

import "fmt"

type symbolManager struct {
	symbolMap     map[string]int
	usedValue     map[int]bool
	nextAvailable int
}

func (m *symbolManager) findSymbol(symbol string) int {
	value, ok := m.symbolMap[symbol]
	if !ok {
		m.symbolMap[symbol] = m.nextAvailable
		m.usedValue[m.nextAvailable] = true
		value = m.nextAvailable
		for m.usedValue[m.nextAvailable] {
			m.nextAvailable++
		}
	}
	return value
}

// For label symbol, the value is determined by the line number in the program.
// So it could be duplicated with predefined symbol.
func (m *symbolManager) registerLabelSymbol(symbol string, value int) {
	m.symbolMap[symbol] = value
}

func getSymbolManager() *symbolManager {
	return &symbolManager{
		map[string]int{
			"R0":     0,
			"R1":     1,
			"R2":     2,
			"R3":     3,
			"R4":     4,
			"R5":     5,
			"R6":     6,
			"R7":     7,
			"R8":     8,
			"R9":     9,
			"R10":    10,
			"R11":    11,
			"R12":    12,
			"R13":    13,
			"R14":    14,
			"R15":    15,
			"SCREEN": 16384,
			"KBD":    24576,
			"SP":     0,
			"LCL":    1,
			"ARG":    2,
			"THIS":   3,
			"THAT":   4,
		},
		map[int]bool{
			0:     true,
			1:     true,
			2:     true,
			3:     true,
			4:     true,
			5:     true,
			6:     true,
			7:     true,
			8:     true,
			9:     true,
			10:    true,
			11:    true,
			12:    true,
			13:    true,
			14:    true,
			15:    true,
			16384: true,
			24576: true,
		},
		16,
	}
}

func processLabelSymbol(lines []string, symbolManager *symbolManager) []string {
	realCode := []string{}
	for _, line := range lines {
		if line[0] == '(' {
			symbol := line[1 : len(line)-1]
			symbolManager.registerLabelSymbol(symbol, len(realCode))
		} else {
			realCode = append(realCode, line)
		}
	}
	return realCode
}

func processSymbolsInACommand(lines []string, symbolManager *symbolManager) []string {
	noSymbolCode := []string{}
	for _, line := range lines {
		if line[0] == '@' {
			if !isNumber(line[1:]) {
				symbol := line[1:]
				//fmt.Println("findSymbol", symbol)
				noSymbolCode = append(noSymbolCode, fmt.Sprintf("@%d", symbolManager.findSymbol(symbol)))
				continue
			}
		}
		noSymbolCode = append(noSymbolCode, line)
	}
	return noSymbolCode
}

func isNumber(s string) bool {
	//fmt.Println("isNumber", s)
	for _, b := range s {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}
