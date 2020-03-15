package main

import "fmt"

type symbolType int

const (
	symbolTypeField symbolType = iota
	symbolTypeStatic
	symbolTypeLocal
	symbolTypeArg
)

type jackSymbol struct {
	name       string
	symbolType symbolType
	idx        int
	typeName   string
}

func (s jackSymbol) print() string {
	switch s.symbolType {
	case symbolTypeField:
		return fmt.Sprintf("this %d", s.idx)
	case symbolTypeStatic:
		return fmt.Sprintf("static %d", s.idx)
	case symbolTypeLocal:
		return fmt.Sprintf("local %d", s.idx)
	case symbolTypeArg:
		return fmt.Sprintf("argument %d", s.idx)
	default:
		panic(fmt.Sprint("invalid symbol:", s))
	}
}

type symbolTable struct {
	next       *symbolTable
	table      map[string]jackSymbol
	nextField  int
	nextStatic int
	nextLocal  int
	nextArg    int
}

func newSymbolTable(lastTable *symbolTable) *symbolTable {
	return &symbolTable{lastTable, make(map[string]jackSymbol), 0, 0, 0, 0}
}

func (s *symbolTable) put(name string, st symbolType, sn string) {
	switch st {
	case symbolTypeField:
		newSymbol := jackSymbol{name, st, s.nextField, sn}
		s.table[name] = newSymbol
		s.nextField++
	case symbolTypeArg:
		newSymbol := jackSymbol{name, st, s.nextArg, sn}
		s.table[name] = newSymbol
		s.nextArg++
	case symbolTypeLocal:
		newSymbol := jackSymbol{name, st, s.nextLocal, sn}
		s.table[name] = newSymbol
		s.nextLocal++
	case symbolTypeStatic:
		newSymbol := jackSymbol{name, st, s.nextStatic, sn}
		s.table[name] = newSymbol
		s.nextStatic++
	}
}

func (s *symbolTable) lookup(name string) jackSymbol {
	if s == nil {
		panic("cannot find symbol:" + name)
	}
	if found, ok := s.table[name]; ok {
		return found
	}
	return s.next.lookup(name)
}

func (s *symbolTable) contains(name string) bool {
	if s == nil {
		return false
	}
	if _, ok := s.table[name]; ok {
		return true
	}
	return s.next.contains(name)
}
