package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type grammer interface {
	print() []string
	printBytecode() []string
}

var curClassName string

type classGrammer struct {
	className token
	cvds      []classVarDec
	sds       []subroutineDec
}

func (c classGrammer) print() []string {
	out := []string{}
	out = append(out, "<class>")
	out = append(out, "  "+keywordToken{class}.print())
	out = append(out, "  "+c.className.print())
	out = append(out, "  "+symbolToken{llb}.print())
	for _, cvd := range c.cvds {
		for _, str := range cvd.print() {
			out = append(out, "  "+str)
		}
	}
	for _, sd := range c.sds {
		for _, str := range sd.print() {
			out = append(out, "  "+str)
		}
	}
	out = append(out, "  "+symbolToken{rlb}.print())
	out = append(out, "</class>")
	return out
}

func (c classGrammer) printBytecode() []string {
	classLevelTable := newSymbolTable(nil)
	fieldCnt := 0
	curClassName = getIdentifierName(c.className)
	for _, cvd := range c.cvds {
		if isKeyword(cvd.staticOrField, static) {
			for _, name := range cvd.varNames {
				varTypeName := ""
				if _, ok := cvd.varType.(identifierToken); ok {
					varTypeName = cvd.varType.(identifierToken).val
				}
				classLevelTable.put(getIdentifierName(name), symbolTypeStatic, varTypeName)
			}
		} else if isKeyword(cvd.staticOrField, field) {
			for _, name := range cvd.varNames {
				varTypeName := ""
				if _, ok := cvd.varType.(identifierToken); ok {
					varTypeName = cvd.varType.(identifierToken).val
				}
				classLevelTable.put(getIdentifierName(name), symbolTypeField, varTypeName)
				fieldCnt++
			}
		} else {
			panic(fmt.Sprint("class varible is neither a field or a static: ", cvd.staticOrField))
		}
	}
	out := []string{}
	for _, subR := range c.sds {
		out = append(out, subR.printBytecode(classLevelTable, getIdentifierName(c.className), fieldCnt)...)
	}
	return out
}

func getIdentifierName(t token) string {
	identifier, ok := t.(identifierToken)
	if !ok {
		panic(fmt.Sprint("token ", t, " is not an identifierToken"))
	}
	return identifier.val
}

func compileClass(tokens []token) (classGrammer, []token) {
	c := classGrammer{}
	tokens = tokens[1:] // consume class
	c.className = tokens[0]
	//fmt.Println(c.className.print())
	tokens = tokens[1:]
	tokens = tokens[1:] // consume {
	cvds := []classVarDec{}
	for isKeyword(tokens[0], static) || isKeyword(tokens[0], field) {
		cvd, restTokens := compileClassVarDec(tokens)
		tokens = restTokens
		cvds = append(cvds, cvd)
	}
	c.cvds = cvds
	sds := []subroutineDec{}
	for isKeyword(tokens[0], constructor) || isKeyword(tokens[0], function) || isKeyword(tokens[0], method) {
		//fmt.Println(tokens[0].print())
		//fmt.Println("debug", len(tokens))
		sd, restTokens := compileSubroutineDec(tokens)
		tokens = restTokens
		//fmt.Println("debug2", len(restTokens))

		sds = append(sds, sd)
	}
	c.sds = sds
	tokens = tokens[1:]
	return c, tokens
}

type classVarDec struct {
	staticOrField token
	varType       token
	varNames      []token
}

func (dec classVarDec) print() []string {
	out := []string{}
	out = append(out, ("<classVarDec>"))
	out = append(out, "  "+dec.staticOrField.print())
	out = append(out, "  "+dec.varType.print())
	for i, varName := range dec.varNames {
		if i > 0 {
			out = append(out, ("  " + symbolToken{comma}.print()))
		}
		out = append(out, "  "+varName.print())
	}
	out = append(out, ("  " + symbolToken{semicolon}.print()))
	out = append(out, ("</classVarDec>"))
	return out
}

func compileClassVarDec(tokens []token) (classVarDec, []token) {
	// comsume a static/field
	staticOrFieldToken := tokens[0]
	tokens = tokens[1:]
	// consume a varType
	varType := tokens[0]
	tokens = tokens[1:]
	// consume varNames
	varNames := []token{}
	varNames = append(varNames, tokens[0])
	tokens = tokens[1:]
	for len(tokens) > 0 && isSymbolToken(tokens[0], comma) {
		tokens = tokens[1:]
		varNames = append(varNames, tokens[0])
		tokens = tokens[1:]
	}
	// consume the last ';'
	tokens = tokens[1:]
	return classVarDec{staticOrFieldToken, varType, varNames}, tokens
}

type subroutineDec struct {
	subroutineType token // method or function
	returnType     token
	subRoutineName token
	pl             parameterList
	sb             subroutineBody
}

func (sd subroutineDec) printBytecode(fatherTable *symbolTable, fatherClassName string, fatherFieldNum int) []string {
	if isKeyword(sd.subroutineType, constructor) {
		return sd.printBytecodeForConstructor(fatherTable, fatherClassName, fatherFieldNum)
	} else if isKeyword(sd.subroutineType, method) {
		return sd.printBytecodeForMethod(fatherTable, fatherClassName)
	} else if isKeyword(sd.subroutineType, function) {
		return sd.printBytecodeForFunction(fatherTable, fatherClassName)
	}
	panic(fmt.Sprint("unrecgonized subroutineType", sd.subroutineType))
}

func (sd subroutineDec) printBytecodeForConstructor(fatherTable *symbolTable, fatherClassName string, fatherFieldNum int) []string {
	out := []string{}
	numOfLocalVariable := 0
	tableHere := newSymbolTable(fatherTable)
	for i, name := range sd.pl.varNames {
		varTypeName := ""
		if _, ok := sd.pl.types[i].(identifierToken); ok {
			varTypeName = sd.pl.types[i].(identifierToken).val
		}
		tableHere.put(getIdentifierName(name), symbolTypeArg, varTypeName)
	}
	for _, varDec := range sd.sb.varDecs {
		numOfLocalVariable += len(varDec.varNames)
	}
	for _, varDec := range sd.sb.varDecs {
		for _, name := range varDec.varNames {
			varTypeName := ""
			if _, ok := varDec.varType.(identifierToken); ok {
				varTypeName = varDec.varType.(identifierToken).val
			}

			tableHere.put(getIdentifierName(name), symbolTypeLocal, varTypeName)
		}
	}
	out = append(out, fmt.Sprintf("function %s.%s %d\n", fatherClassName, getIdentifierName(sd.subRoutineName), numOfLocalVariable))
	out = append(out, fmt.Sprintf("push constant %d", fatherFieldNum))
	out = append(out, fmt.Sprintf("call Memory.alloc 1"))
	out = append(out, fmt.Sprintf("pop pointer 0"))
	out = append(out, sd.sb.sts.printBytecode(tableHere)...)
	return out
}

func (sd subroutineDec) printBytecodeForMethod(fatherTable *symbolTable, fatherClassName string) []string {
	out := []string{}
	numOfLocalVariable := 0
	for _, varDec := range sd.sb.varDecs {
		numOfLocalVariable += len(varDec.varNames)
	}
	out = append(out, fmt.Sprintf("function %s.%s %d\n", fatherClassName, getIdentifierName(sd.subRoutineName), numOfLocalVariable))
	out = append(out, "push argument 0")
	out = append(out, "pop pointer 0")

	subroutineLevelTable := newSymbolTable(fatherTable)
	subroutineLevelTable.put("this", symbolTypeArg, fatherClassName)
	for i, name := range sd.pl.varNames {
		varTypeName := ""
		if _, ok := sd.pl.types[i].(identifierToken); ok {
			varTypeName = sd.pl.types[i].(identifierToken).val
		}

		subroutineLevelTable.put(getIdentifierName(name), symbolTypeArg, varTypeName)
	}
	for _, varDec := range sd.sb.varDecs {
		for _, name := range varDec.varNames {
			varTypeName := ""
			if _, ok := varDec.varType.(identifierToken); ok {
				varTypeName = varDec.varType.(identifierToken).val
			}

			subroutineLevelTable.put(getIdentifierName(name), symbolTypeLocal, varTypeName)
		}
	}
	out = append(out, sd.sb.sts.printBytecode(subroutineLevelTable)...)
	return out
}

func (sd subroutineDec) printBytecodeForFunction(fatherTable *symbolTable, fatherClassName string) []string {
	//mt.Println("debug here")
	out := []string{}
	numOfLocalVariable := 0
	for _, varDec := range sd.sb.varDecs {
		numOfLocalVariable += len(varDec.varNames)
	}
	out = append(out, fmt.Sprintf("function %s.%s %d\n", fatherClassName, getIdentifierName(sd.subRoutineName), numOfLocalVariable))
	subroutineLevelTable := newSymbolTable(fatherTable)
	for i, name := range sd.pl.varNames {
		varTypeName := ""
		if _, ok := sd.pl.types[i].(identifierToken); ok {
			varTypeName = sd.pl.types[i].(identifierToken).val
		}

		subroutineLevelTable.put(getIdentifierName(name), symbolTypeArg, varTypeName)
	}
	for _, varDec := range sd.sb.varDecs {
		for _, name := range varDec.varNames {
			varTypeName := ""
			if _, ok := varDec.varType.(identifierToken); ok {
				varTypeName = varDec.varType.(identifierToken).val
			}
			subroutineLevelTable.put(getIdentifierName(name), symbolTypeLocal, varTypeName)
		}
	}
	out = append(out, sd.sb.sts.printBytecode(subroutineLevelTable)...)
	return out
}

func (sd subroutineDec) print() []string {
	out := []string{}
	out = append(out, "<subroutineDec>")
	out = append(out, "  "+sd.subroutineType.print())
	out = append(out, "  "+sd.returnType.print())
	out = append(out, "  "+sd.subRoutineName.print())
	out = append(out, "  "+symbolToken{lsb}.print())
	for _, str := range sd.pl.print() {
		out = append(out, "  "+str)
	}
	out = append(out, "  "+symbolToken{rsb}.print())
	for _, str := range sd.sb.print() {
		out = append(out, "  "+str)
	}
	out = append(out, "</subroutineDec>")
	return out
}

func compileSubroutineDec(tokens []token) (subroutineDec, []token) {
	sd := subroutineDec{}
	sd.subroutineType = tokens[0]
	tokens = tokens[1:]
	sd.returnType = tokens[0]
	tokens = tokens[1:]
	sd.subRoutineName = tokens[0]
	tokens = tokens[1:]
	tokens = tokens[1:] // consume (
	pl, restTokens := compileParameterList(tokens)
	sd.pl = pl
	tokens = restTokens
	tokens = tokens[1:] // consume )
	sb, restTokens := compileSubroutineBody(tokens)
	sd.sb = sb
	tokens = restTokens
	return sd, tokens
}

type parameterList struct {
	types    []token
	varNames []token
}

func (pl parameterList) print() []string {
	out := []string{}
	out = append(out, "<parameterList>")
	for i := range pl.types {
		if i != 0 {
			out = append(out, "  "+symbolToken{comma}.print())
		}
		out = append(out, "  "+pl.types[i].print())
		out = append(out, "  "+pl.varNames[i].print())
	}
	out = append(out, "</parameterList>")
	return out
}

func compileParameterList(tokens []token) (parameterList, []token) {
	pl := parameterList{}
	pl.types = []token{}
	pl.varNames = []token{}
	for !isSymbolToken(tokens[0], rsb) {
		if isSymbolToken(tokens[0], comma) {
			tokens = tokens[1:]
		}
		pl.types = append(pl.types, tokens[0])
		tokens = tokens[1:]
		pl.varNames = append(pl.varNames, tokens[0])
		tokens = tokens[1:]
	}
	return pl, tokens
}

type subroutineBody struct {
	varDecs []varDec
	sts     statements
}

func (s subroutineBody) print() []string {
	out := []string{}
	out = append(out, "<subroutineBody>")
	out = append(out, "  "+symbolToken{llb}.print())
	for _, vc := range s.varDecs {
		for _, str := range vc.print() {
			out = append(out, "  "+str)
		}
	}
	for _, str := range s.sts.print() {
		out = append(out, "  "+str)
	}
	out = append(out, "  "+symbolToken{rlb}.print())
	out = append(out, "</subroutineBody>")
	return out
}

func compileSubroutineBody(tokens []token) (subroutineBody, []token) {
	//fmt.Println("debug body", tokens[0].print(), tokens[len(tokens)-2].print())
	sb := subroutineBody{}
	tokens = tokens[1:] // consume {
	vds := []varDec{}
	for isKeyword(tokens[0], _var) {
		vd, restTokens := compileVarDec(tokens)
		tokens = restTokens
		vds = append(vds, vd)
	}
	sts, restTokens := compileStatements(tokens)
	tokens = restTokens
	tokens = tokens[1:] // consume }
	sb.varDecs = vds
	sb.sts = sts
	return sb, tokens
}

func isKeyword(t token, keywordT keywordTokenType) bool {
	k, ok := t.(keywordToken)
	if !ok {
		return false
	}
	return k.keywordType == keywordT
}

type varDec struct {
	varType  token
	varNames []token
}

func (v varDec) print() []string {
	out := []string{}
	out = append(out, "<varDec>")
	out = append(out, "  "+keywordToken{_var}.print())
	out = append(out, "  "+v.varType.print())
	for i, name := range v.varNames {
		if i != 0 {
			out = append(out, "  "+symbolToken{comma}.print())
		}
		out = append(out, "  "+name.print())
	}
	out = append(out, "  "+symbolToken{semicolon}.print())
	out = append(out, "</varDec>")
	return out
}

func compileVarDec(tokens []token) (varDec, []token) {
	vD := varDec{}
	tokens = tokens[1:] // consume var
	vD.varType = tokens[0]
	tokens = tokens[1:]
	names := []token{}
	names = append(names, tokens[0])
	tokens = tokens[1:]
	for isSymbolToken(tokens[0], comma) {
		tokens = tokens[1:] //consume ,
		names = append(names, tokens[0])
		tokens = tokens[1:]
	}
	tokens = tokens[1:] // consume ;
	vD.varNames = names
	return vD, tokens
}

type statement interface {
	print() []string
	printBytecode(table *symbolTable) []string
}

func compileStatement(tokens []token) (statement, []token) {
	firstToken := tokens[0]
	keyword := firstToken.(keywordToken)
	switch keyword.keywordType {
	case _if:
		ifS := ifStatement{}
		tokens = tokens[2:] // consume if (
		exp, restTokens := compileExpression(tokens)
		tokens = restTokens
		ifS.condition = exp
		tokens = tokens[2:] // consume ) {
		statementsInIf, restTokens := compileStatements(tokens)
		ifS.ifStatements = statementsInIf
		tokens = restTokens
		tokens = tokens[1:] // consume }
		ifS.hasElse = false
		if isElseToken(tokens[0]) {
			ifS.hasElse = true
			tokens = tokens[2:] // consume else {
			statementsInElse, restTokens := compileStatements(tokens)
			tokens = restTokens
			tokens = tokens[1:] // consume }
			ifS.elseStatements = statementsInElse
		}
		return ifS, tokens
	case let:
		letS := letStatement{}
		tokens = tokens[1:] // consume let
		letS.varName = tokens[0]
		tokens = tokens[1:]
		if isSymbolToken(tokens[0], lmb) {
			tokens = tokens[1:] // consume [
			letS.hasArrayExp = true
			arrayExp, restTokens := compileExpression(tokens)
			letS.arrayExp = arrayExp
			tokens = restTokens
			tokens = tokens[1:] // consume ]
		}
		tokens = tokens[1:] // consume =
		exp, restTokens := compileExpression(tokens)
		tokens = restTokens
		letS.exp = exp
		tokens = tokens[1:] // consume ;
		return letS, tokens
	case while:
		whileS := whileStatement{}
		tokens = tokens[2:] // consume while (
		exp, restTokens := compileExpression(tokens)
		whileS.condition = exp
		tokens = restTokens
		tokens = tokens[2:] // consume ) {
		sts, restTokens := compileStatements(tokens)
		tokens = restTokens
		whileS.sts = sts
		tokens = tokens[1:] // consume }
		return whileS, tokens
	case _return:
		returnS := returnStatement{}
		returnS.hasExpression = false
		tokens = tokens[1:] // consume return
		if !isSymbolToken(tokens[0], semicolon) {
			returnS.hasExpression = true
			exp, restTokens := compileExpression(tokens)
			returnS.exp = exp
			tokens = restTokens
		}
		tokens = tokens[1:] // consume ;
		return returnS, tokens
	case do:
		doS := doStatement{}
		tokens = tokens[1:] // consume do
		t, restTokens := compileTerm(tokens)
		tokens = restTokens
		doS.subC = t.(subroutineCallTerm).sb
		tokens = tokens[1:] // consume ;
		return doS, tokens
	default:
		panic("invalid statement:" + tokens[0].print())
	}
}

func isSymbolToken(t token, st symbolTokenType) bool {
	symbolT, ok := t.(symbolToken)
	if !ok {
		return false
	}
	return symbolT.symbolType == st
}

func isElseToken(t token) bool {
	keyword, ok := t.(keywordToken)
	if ok {
		return keyword.keywordType == _else
	}
	return false
}

func compileStatements(tokens []token) (statements, []token) {
	sts := []statement{}
	for isStatement(tokens[0]) {
		st, restTokens := compileStatement(tokens)
		tokens = restTokens
		sts = append(sts, st)
	}
	return statements{sts}, tokens
}

func isStatement(t token) bool {
	keyword, ok := t.(keywordToken)
	if !ok {
		return false
	}
	if keyword.keywordType == do || keyword.keywordType == let || keyword.keywordType == while || keyword.keywordType == _if ||
		keyword.keywordType == _return {
		return true
	}
	return false
}

type statements struct {
	sts []statement
}

func (s statements) print() []string {
	out := []string{}
	out = append(out, "<statements>")
	for _, statement := range s.sts {
		for _, str := range statement.print() {
			out = append(out, "  "+str)
		}
	}
	out = append(out, "</statements>")
	return out
}

func (s statements) printBytecode(table *symbolTable) []string {
	out := []string{}
	for _, st := range s.sts {
		out = append(out, st.printBytecode(table)...)
	}
	return out
}

type ifStatement struct {
	condition      expression
	ifStatements   statements
	hasElse        bool
	elseStatements statements
}

func (s ifStatement) printBytecode(table *symbolTable) []string {
	l1 := getNewLabel()
	l2 := getNewLabel()
	out := s.condition.printBytecode(table)
	out = append(out, "not")
	out = append(out, fmt.Sprintf("if-goto %s", l1))
	out = append(out, s.ifStatements.printBytecode(table)...)
	out = append(out, fmt.Sprintf("goto %s", l2))
	out = append(out, fmt.Sprintf("label %s", l1))
	if s.hasElse {
		out = append(out, s.elseStatements.printBytecode(table)...)
	}
	out = append(out, fmt.Sprintf("label %s", l2))
	return out
}

func (s ifStatement) print() []string {
	out := []string{}
	out = append(out, "<ifStatement>")
	out = append(out, "  "+keywordToken{_if}.print())
	out = append(out, "  "+symbolToken{lsb}.print())
	for _, str := range s.condition.print() {
		out = append(out, "  "+str)
	}
	out = append(out, "  "+symbolToken{rsb}.print())
	out = append(out, "  "+symbolToken{llb}.print())
	for _, str := range s.ifStatements.print() {
		out = append(out, "  "+str)
	}
	out = append(out, "  "+symbolToken{rlb}.print())
	if s.hasElse {
		out = append(out, "  "+keywordToken{_else}.print())
		out = append(out, "  "+symbolToken{llb}.print())
		for _, str := range s.elseStatements.print() {
			out = append(out, "  "+str)
		}
		out = append(out, "  "+symbolToken{rlb}.print())
	}
	out = append(out, "</ifStatement>")
	return out
}

type letStatement struct {
	varName     token
	hasArrayExp bool
	arrayExp    expression
	exp         expression
}

func (s letStatement) printBytecode(table *symbolTable) []string {
	out := s.exp.printBytecode(table)
	if !s.hasArrayExp {
		out = append(out, fmt.Sprintf("pop %s", table.lookup(getIdentifierName(s.varName)).print()))
	} else {
		out = append(out, fmt.Sprintf("push %s", table.lookup(getIdentifierName(s.varName)).print()))
		out = append(out, s.arrayExp.printBytecode(table)...)
		out = append(out, "add")
		out = append(out, "pop pointer 1")
		out = append(out, "pop that 0")
	}
	return out
}

func (s letStatement) print() []string {
	out := []string{}
	out = append(out, "<letStatement>")
	out = append(out, "  "+keywordToken{let}.print())
	out = append(out, "  "+s.varName.print())
	if s.hasArrayExp {
		out = append(out, "  "+symbolToken{lmb}.print())
		for _, str := range s.arrayExp.print() {
			out = append(out, "  "+str)
		}
		out = append(out, "  "+symbolToken{rmb}.print())
	}
	out = append(out, "  "+symbolToken{equal}.print())
	for _, str := range s.exp.print() {
		out = append(out, "  "+str)
	}
	out = append(out, "  "+symbolToken{semicolon}.print())
	out = append(out, "</letStatement>")
	return out
}

type whileStatement struct {
	condition expression
	sts       statements
}

func (s whileStatement) printBytecode(table *symbolTable) []string {
	l1, l2 := getNewLabel(), getNewLabel()
	out := []string{}
	out = append(out, fmt.Sprintf("label %s", l1))
	out = append(out, s.condition.printBytecode(table)...)
	out = append(out, fmt.Sprintf("not"))
	out = append(out, fmt.Sprintf("if-goto %s", l2))
	out = append(out, s.sts.printBytecode(table)...)
	out = append(out, fmt.Sprintf("goto %s", l1))
	out = append(out, fmt.Sprintf("label %s", l2))
	return out
}

func (s whileStatement) print() []string {
	out := []string{}
	out = append(out, "<whileStatement>")
	out = append(out, "  "+keywordToken{while}.print())
	out = append(out, "  "+symbolToken{lsb}.print())
	for _, str := range s.condition.print() {
		out = append(out, "  "+str)
	}
	out = append(out, "  "+symbolToken{rsb}.print())
	out = append(out, "  "+symbolToken{llb}.print())
	for _, str := range s.sts.print() {
		out = append(out, "  "+str)
	}
	out = append(out, "  "+symbolToken{rlb}.print())
	out = append(out, "</whileStatement>")
	return out
}

type doStatement struct {
	subC subroutineCall
}

func (s doStatement) printBytecode(table *symbolTable) []string {
	out := s.subC.printBytecode(table)
	out = append(out, "pop temp 0")
	return out
}

func (s doStatement) print() []string {
	out := []string{}
	out = append(out, "<doStatement>")
	out = append(out, "  "+keywordToken{do}.print())
	for _, str := range s.subC.print() {
		out = append(out, "  "+str)
	}
	out = append(out, "  "+symbolToken{semicolon}.print())
	out = append(out, "</doStatement>")
	return out
}

type returnStatement struct {
	hasExpression bool
	exp           expression
}

func (s returnStatement) printBytecode(table *symbolTable) []string {
	if s.hasExpression {
		out := s.exp.printBytecode(table)
		out = append(out, "return")
		return out
	}
	out := []string{}
	out = append(out, "push constant 0")
	out = append(out, "return")
	return out
}

func (s returnStatement) print() []string {
	out := []string{}
	out = append(out, "<returnStatement>")
	out = append(out, "  "+keywordToken{_return}.print())
	if s.hasExpression {
		for _, str := range s.exp.print() {
			out = append(out, "  "+str)
		}
	}
	out = append(out, "  "+symbolToken{semicolon}.print())
	out = append(out, "</returnStatement>")
	return out
}

type term interface {
	print() []string
	printBytecode(table *symbolTable) []string
}

type singleTerm struct {
	id token
}

func (st singleTerm) printBytecode(table *symbolTable) []string {
	out := []string{}
	if st.id.getType() == identifier {
		out = append(out, fmt.Sprintf("push %s", table.lookup(getIdentifierName(st.id)).print()))
	} else if st.id.getType() == keyword {
		keywordHere := st.id.(keywordToken)
		if keywordHere.keywordType == _true {
			out = append(out, fmt.Sprintf("push constant %d", 1))
			out = append(out, "neg")
		} else if keywordHere.keywordType == null || keywordHere.keywordType == _false {
			out = append(out, fmt.Sprintf("push constant %d", 0))
		} else if keywordHere.keywordType == this {
			out = append(out, fmt.Sprintf("push pointer 0"))
		} else {
			panic(fmt.Sprint("invalid contant keyword", keywordHere))
		}
	} else if st.id.getType() == intergerConstant {
		out = append(out, fmt.Sprintf("push constant %d", st.id.(integerConstantToken).val))
	} else if st.id.getType() == stringConstant {
		str := st.id.(stringConstantToken)
		out = append(out, fmt.Sprintf("push constant %d", len(str.val)))
		out = append(out, fmt.Sprintf("call String.new 1"))
		for _, ch := range str.val {
			out = append(out, fmt.Sprintf("push constant %d", ch))
			out = append(out, fmt.Sprintf("call String.appendChar 2"))
		}
	} else {
		panic(fmt.Sprint("invalid singleTerm:", st))
	}
	return out
}

func (st singleTerm) print() []string {
	out := []string{}
	out = append(out, "<term>")
	out = append(out, "  "+st.id.print())
	out = append(out, "</term>")
	return out
}

type arrayTerm struct {
	varName token
	exp     expression
}

func (t arrayTerm) printBytecode(table *symbolTable) []string {
	out := []string{}
	out = append(out, fmt.Sprintf("push %s", table.lookup(getIdentifierName(t.varName)).print()))
	out = append(out, t.exp.printBytecode(table)...)
	out = append(out, "add")
	out = append(out, "pop pointer 1")
	out = append(out, "push that 0")
	return out
}

func (t arrayTerm) print() []string {
	out := []string{}
	out = append(out, "<term>")
	out = append(out, "  "+t.varName.print())
	out = append(out, "  "+symbolToken{lmb}.print())
	nextE := t.exp.print()
	for _, str := range nextE {
		out = append(out, "  "+str)
	}
	out = append(out, "  "+symbolToken{rmb}.print())
	out = append(out, "</term>")
	return out
}

type subroutineCallTerm struct {
	sb subroutineCall
}

func (st subroutineCallTerm) print() []string {
	out := []string{}
	out = append(out, "<term>")
	for _, str := range st.sb.print() {
		out = append(out, "  "+str)
	}
	out = append(out, "</term>")
	return out
}

func (st subroutineCallTerm) printBytecode(table *symbolTable) []string {
	return st.sb.printBytecode(table)
}

type subroutineCall struct {
	hasClassOrVarName bool
	classOrVarName    token
	subrountineName   token
	expList           expressionList
}

func (st subroutineCall) printBytecode(table *symbolTable) []string {
	out := []string{}
	classOrVarName := ""
	if st.hasClassOrVarName {
		classOrVarName = getIdentifierName(st.classOrVarName)
	} else {
		classOrVarName = curClassName
	}
	className := ""
	argNum := 0
	if table.contains(classOrVarName) {
		varName := table.lookup(classOrVarName)
		className = varName.typeName
		out = append(out, fmt.Sprintf("push %s", varName.print()))
		argNum++
	} else if !st.hasClassOrVarName {
		out = append(out, fmt.Sprintf("push pointer 0"))
		argNum++
		className = curClassName
	} else {
		className = classOrVarName
	}
	argNum += len(st.expList.expressions)
	out = append(out, st.expList.printBytecode(table)...)
	out = append(out, fmt.Sprintf("call %s.%s %d", className, getIdentifierName(st.subrountineName), argNum))
	return out
}

func (st subroutineCall) print() []string {
	out := []string{}
	if st.hasClassOrVarName {
		out = append(out, st.classOrVarName.print())
		// print .
		out = append(out, symbolToken{dot}.print())
	}
	out = append(out, st.subrountineName.print())
	out = append(out, symbolToken{lsb}.print())

	for _, str := range st.expList.print() {
		out = append(out, str)
	}
	out = append(out, symbolToken{rsb}.print())

	return out
}

type bracketTerm struct {
	exp expression
}

func (t bracketTerm) printBytecode(table *symbolTable) []string {
	return t.exp.printBytecode(table)
}

func (t bracketTerm) print() []string {
	out := []string{}
	out = append(out, "<term>")
	out = append(out, "  "+symbolToken{lsb}.print())
	nextE := t.exp.print()
	for _, str := range nextE {
		out = append(out, "  "+str)
	}
	out = append(out, "  "+symbolToken{rsb}.print())
	out = append(out, "</term>")
	return out
}

type unaryOpTerm struct {
	op token
	t  term
}

func (t unaryOpTerm) printBytecode(table *symbolTable) []string {
	out := []string{}
	out = append(out, t.t.printBytecode(table)...)
	switch t.op.(symbolToken).symbolType {
	case not:
		out = append(out, "not")
	case minus:
		out = append(out, "neg")
	default:
		panic("invalid unaryOp")
	}
	return out
}

func (t unaryOpTerm) print() []string {
	out := []string{}
	out = append(out, "<term>")
	out = append(out, "  "+t.op.print())
	nextT := t.t.print()
	for _, str := range nextT {
		out = append(out, "  "+str)
	}
	out = append(out, "</term>")
	return out
}

type expression struct {
	firstTerm term
	ops       []token
	restTerms []term
}

func (e expression) printBytecode(table *symbolTable) []string {
	out := []string{}
	out = append(out, e.firstTerm.printBytecode(table)...)
	for i := range e.ops {
		out = append(out, e.restTerms[i].printBytecode(table)...)
		switch e.ops[i].(symbolToken).symbolType {
		case plus:
			out = append(out, "add")
		case minus:
			out = append(out, "sub")
		case mul:
			out = append(out, "call Math.multiply 2")
		case div:
			out = append(out, "call Math.divide 2")
		case lesser:
			out = append(out, "lt")
		case bigger:
			out = append(out, "gt")
		case equal:
			out = append(out, "eq")
		case and:
			out = append(out, "and")
		case or:
			out = append(out, "or")
		}
	}
	return out
}

func (e expression) print() []string {
	out := []string{}
	out = append(out, "<expression>")
	for _, str := range e.firstTerm.print() {
		out = append(out, "  "+str)
	}
	for i := range e.ops {
		out = append(out, "  "+e.ops[i].print())
		for _, str := range e.restTerms[i].print() {
			out = append(out, "  "+str)
		}
	}
	out = append(out, "</expression>")
	return out
}

func compileExpression(tokens []token) (expression, []token) {
	exp := expression{}
	firstTerm, restTokens := compileTerm(tokens)
	tokens = restTokens
	exp.firstTerm = firstTerm
	ops := []token{}
	restTerms := []term{}
	for len(tokens) > 0 && isOpToken(tokens[0]) {
		ops = append(ops, tokens[0])
		tokens = tokens[1:]
		theTerm, restTokens := compileTerm(tokens)
		restTerms = append(restTerms, theTerm)
		tokens = restTokens
	}
	exp.ops = ops
	exp.restTerms = restTerms
	return exp, tokens
}

func isOpToken(t token) bool {
	symbol, ok := t.(symbolToken)
	if !ok {
		return false
	}
	if symbol.symbolType == plus || symbol.symbolType == minus || symbol.symbolType == mul || symbol.symbolType == div ||
		symbol.symbolType == and || symbol.symbolType == or || symbol.symbolType == lesser || symbol.symbolType == bigger || symbol.symbolType == equal {
		return true
	}
	return false
}

func compileExpressionList(tokens []token) (expressionList, []token) {
	expressions := []expression{}
	for len(tokens) > 0 {
		firstToken := tokens[0]
		if isSymbolToken(firstToken, rsb) {
			return expressionList{expressions}, tokens
		}
		if isSymbolToken(firstToken, comma) {
			tokens = tokens[1:] // consume a comma
		}

		curExpression, newTokens := compileExpression(tokens)
		tokens = newTokens
		expressions = append(expressions, curExpression)
	}
	// would never reach this step
	return expressionList{}, nil
}

func (el expressionList) printBytecode(table *symbolTable) []string {
	out := []string{}
	for _, e := range el.expressions {
		out = append(out, e.printBytecode(table)...)
	}
	return out
}

func (el expressionList) print() []string {
	out := []string{}
	out = append(out, "<expressionList>")
	for i, exp := range el.expressions {
		if i > 0 {
			out = append(out, "  "+symbolToken{comma}.print())
		}
		for _, str := range exp.print() {
			out = append(out, "  "+str)
		}
	}
	out = append(out, "</expressionList>")
	return out
}

func compileTerm(tokens []token) (term, []token) {
	firstToken := tokens[0]
	toOp, ok := firstToken.(symbolToken)
	if ok {
		// check whether this is unary op term
		if toOp.symbolType == minus || toOp.symbolType == not {
			t := unaryOpTerm{}
			t.op = firstToken
			tokens = tokens[1:]
			followedTerm, restTokens := compileTerm(tokens)
			t.t = followedTerm
			return t, restTokens
		}
		// check whether this is bracket
		if toOp.symbolType == lsb {
			t := bracketTerm{}
			// consume the (
			tokens = tokens[1:]
			followedExp, restTokens := compileExpression(tokens)
			// consume the )
			restTokens = restTokens[1:]
			t.exp = followedExp
			return t, restTokens
		}
		panic("invalid term:" + firstToken.print() + strconv.Itoa(len(tokens)))
	}
	secondToken := tokens[1]
	toOp, ok = secondToken.(symbolToken)
	if ok {
		// varName[expression]
		if toOp.symbolType == lmb {
			aT := arrayTerm{}
			aT.varName = firstToken
			tokens = tokens[2:] // consume varName + [
			followedExp, restTokens := compileExpression(tokens)
			// consume the ]
			restTokens = restTokens[1:]
			aT.exp = followedExp
			return aT, restTokens
		}
		// subRoutineName(ExpressionList)
		if toOp.symbolType == lsb {
			t := subroutineCall{}
			t.hasClassOrVarName = false
			t.subrountineName = firstToken
			tokens = tokens[2:]
			followedExpList, restTokens := compileExpressionList(tokens)
			t.expList = followedExpList
			restTokens = restTokens[1:] // consume the )
			return subroutineCallTerm{t}, restTokens
		}
		// varName.subRoutineName(ExpressionList)
		if toOp.symbolType == dot {
			t := subroutineCall{}
			t.hasClassOrVarName = true
			t.classOrVarName = firstToken
			t.subrountineName = tokens[2]
			tokens = tokens[4:]
			followedExpList, restTokens := compileExpressionList(tokens)
			t.expList = followedExpList
			restTokens = restTokens[1:] // consume the )
			return subroutineCallTerm{t}, restTokens
		}
	}
	st := singleTerm{}
	st.id = firstToken
	tokens = tokens[1:]
	return st, tokens
}

type expressionList struct {
	expressions []expression
}

var labelNo = 0

func getNewLabel() string {
	defer func() {
		labelNo++
	}()
	return fmt.Sprintf("Label%d", labelNo)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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
	filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".jack") {
			f, err := os.Open(path)
			check(err)
			in := bufio.NewReader(f)
			outFileName := path[:len(path)-5] + ".xml"
			outFile, err := os.Create(outFileName)
			check(err)
			out := bufio.NewWriter(outFile)
			outTokenFileName := path[:len(path)-5] + "T.xml"
			outTokenFile, err := os.Create(outTokenFileName)
			check(err)
			outToken := bufio.NewWriter(outTokenFile)
			outTokenFileNameVM := path[:len(path)-5] + ".vm"
			outTokenFileVM, err := os.Create(outTokenFileNameVM)
			check(err)
			outTokenVM := bufio.NewWriter(outTokenFileVM)

			processSingleFile(in, outToken, out, outTokenVM)
		}
		return nil
	})
}

func processFile(fileName string) {
	f, err := os.Open(fileName)
	check(err)
	outputFileName := strings.Replace(fileName, ".jack", ".xml", 1)
	outputTokenFileName := strings.Replace(fileName, ".jack", "T.xml", 1)
	outputVMName := strings.Replace(fileName, ".jack", ".vm", 1)

	of, err := os.Create(outputFileName)
	check(err)
	otf, err := os.Create(outputTokenFileName)
	check(err)
	ovm, err := os.Create(outputVMName)
	check(err)
	in := bufio.NewReader(f)
	out := bufio.NewWriter(of)
	outToken := bufio.NewWriter(otf)
	processSingleFile(in, outToken, out, bufio.NewWriter(ovm))
}

func processSingleFile(in *bufio.Reader, outToken *bufio.Writer, out *bufio.Writer, outByteCode *bufio.Writer) {
	defer out.Flush()
	defer outToken.Flush()
	defer outByteCode.Flush()
	lines := []string{}
	for true {
		line, err := in.ReadBytes('\n')
		//fmt.Println("check ", string(line))

		if err != nil {
			if len(line) > 0 {
				lines = append(lines, string(line))
			}
			break
		}
		//fmt.Println("checkAfter ", string(line))

		lineStr := trim(string(line))
		if len(lineStr) == 0 {
			continue
		}
		if len(lineStr) >= 2 && (lineStr[:2] == "//" || lineStr[:2] == " *") {
			continue
		}
		firstChar := byte(0)
		for _, ch := range lineStr {
			if ch != ' ' {
				firstChar = byte(ch)
				break
			}
		}

		if firstChar == '*' {
			continue
		}
		lines = append(lines, lineStr)
	}
	tokens := []token{}
	for _, line := range lines {
		tmp := getTokens(line)
		tokens = append(tokens, tmp...)
	}
	fmt.Fprintln(outToken, "<tokens>")
	for _, token := range tokens {
		fmt.Fprintln(outToken, "	"+token.print())
	}
	fmt.Fprintln(outToken, "</tokens>")
	c, _ := compileClass(tokens)
	for _, str := range c.print() {
		fmt.Fprintln(out, str)
	}
	for _, str := range c.printBytecode() {
		fmt.Fprintln(outByteCode, str)
	}
}

func trim(s string) string {
	ans := []byte{}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == 13 {
			continue
		}
		if s[i] == '/' && i < len(s)-1 && (s[i+1] == '/' || s[i+1] == '*') {
			break
		}
		ans = append(ans, s[i])
	}
	for len(ans) > 0 && (ans[len(ans)-1] == ' ' || ans[len(ans)-1] == '\t') {
		ans = ans[:len(ans)-1]
	}
	return string(ans)
}

type jackTokenType int

const (
	keyword jackTokenType = iota
	symbol
	intergerConstant
	stringConstant
	identifier
)

type keywordTokenType string

const (
	notAKeyWord keywordTokenType = "this is not a keyword"
	class       keywordTokenType = "class"
	constructor keywordTokenType = "constructor"
	function    keywordTokenType = "function"
	method      keywordTokenType = "method"
	field       keywordTokenType = "field"
	static      keywordTokenType = "static"
	_var        keywordTokenType = "var"
	_int        keywordTokenType = "int"
	char        keywordTokenType = "char"
	boolean     keywordTokenType = "boolean"
	void        keywordTokenType = "void"
	_true       keywordTokenType = "true"
	_false      keywordTokenType = "false"
	null        keywordTokenType = "null"
	this        keywordTokenType = "this"
	let         keywordTokenType = "let"
	do          keywordTokenType = "do"
	_if         keywordTokenType = "if"
	_else       keywordTokenType = "else"
	while       keywordTokenType = "while"
	_return     keywordTokenType = "return"
)

type symbolTokenType byte

const (
	lsb       symbolTokenType = '('
	rsb       symbolTokenType = ')'
	lmb       symbolTokenType = '['
	rmb       symbolTokenType = ']'
	llb       symbolTokenType = '{'
	rlb       symbolTokenType = '}'
	dot       symbolTokenType = '.'
	comma     symbolTokenType = ','
	semicolon symbolTokenType = ';'
	plus      symbolTokenType = '+'
	minus     symbolTokenType = '-'
	mul       symbolTokenType = '*'
	div       symbolTokenType = '/'
	and       symbolTokenType = '&'
	or        symbolTokenType = '|'
	bigger    symbolTokenType = '>'
	lesser    symbolTokenType = '<'
	not       symbolTokenType = '~'
	equal     symbolTokenType = '='
)

type token interface {
	getType() jackTokenType
	print() string
}

type keywordToken struct {
	keywordType keywordTokenType
}

func (token keywordToken) getType() jackTokenType {
	return keyword
}

func (token keywordToken) print() string {
	return fmt.Sprintf("<keyword> %s </keyword>", string(token.keywordType))
}

type symbolToken struct {
	symbolType symbolTokenType
}

func (token symbolToken) getType() jackTokenType {
	return symbol
}

func (token symbolToken) print() string {
	if token.symbolType == bigger {
		return fmt.Sprintf("<symbol> &gt; </symbol>")
	}
	if token.symbolType == lesser {
		return fmt.Sprintf("<symbol> &lt; </symbol>")
	}
	if token.symbolType == and {
		return fmt.Sprintf("<symbol> &amp; </symbol>")
	}
	return fmt.Sprintf("<symbol> %s </symbol>", string(token.symbolType))
}

type integerConstantToken struct {
	val int
}

func (token integerConstantToken) getType() jackTokenType {
	return intergerConstant
}

func (token integerConstantToken) print() string {
	return fmt.Sprintf("<integerConstant> %d </integerConstant>", token.val)
}

type stringConstantToken struct {
	val string
}

func (token stringConstantToken) getType() jackTokenType {
	return stringConstant
}

func (token stringConstantToken) print() string {
	return fmt.Sprintf("<stringConstant> %s </stringConstant>", token.val)
}

type identifierToken struct {
	val string
}

func (token identifierToken) getType() jackTokenType {
	return identifier
}

func (token identifierToken) print() string {
	return fmt.Sprintf("<identifier> %s </identifier>", token.val)
}

func getTokens(s string) []token {
	tokens := []token{}
	var t token
	for len(s) != 0 {
		if s[0] == '(' || s[0] == ')' ||
			s[0] == '[' || s[0] == ']' ||
			s[0] == '{' || s[0] == '}' ||
			s[0] == '.' || s[0] == ',' || s[0] == ';' ||
			s[0] == '+' || s[0] == '-' ||
			s[0] == '*' || s[0] == '/' ||
			s[0] == '&' || s[0] == '|' ||
			s[0] == '~' || s[0] == '=' ||
			s[0] == '>' || s[0] == '<' {
			t, s = tokenizeSymbol(s)
			tokens = append(tokens, t)
		} else if s[0] >= '0' && s[0] <= '9' {
			t, s = tokenizeIntegerConstant(s)
			tokens = append(tokens, t)
		} else if s[0] == '"' {
			t, s = tokenizeStringConstant(s)
			tokens = append(tokens, t)
		} else if s[0] == '_' || s[0] >= 'a' && s[0] <= 'z' || s[0] >= 'A' && s[0] <= 'Z' {
			t, s = tokenizeIdentifierOrKeyword(s)
			tokens = append(tokens, t)
		} else if s[0] == ' ' || s[0] == 9 {
			s = s[1:]
		} else {
			panic(fmt.Sprintf("cannot tokenize:%d", s[0]))
		}
	}
	return tokens
}

func tokenizeSymbol(s string) (token, string) {
	ch := s[0]
	s = s[1:]
	switch ch {
	case '(':
		return symbolToken{lsb}, s
	case ')':
		return symbolToken{rsb}, s
	case '[':
		return symbolToken{lmb}, s
	case ']':
		return symbolToken{rmb}, s
	case '{':
		return symbolToken{llb}, s
	case '}':
		return symbolToken{rlb}, s
	case '.':
		return symbolToken{dot}, s
	case ',':
		return symbolToken{comma}, s
	case ';':
		return symbolToken{semicolon}, s
	case '+':
		return symbolToken{plus}, s
	case '-':
		return symbolToken{minus}, s
	case '*':
		return symbolToken{mul}, s
	case '/':
		return symbolToken{div}, s
	case '&':
		return symbolToken{and}, s
	case '|':
		return symbolToken{or}, s
	case '>':
		return symbolToken{bigger}, s
	case '<':
		return symbolToken{lesser}, s
	case '~':
		return symbolToken{not}, s
	case '=':
		return symbolToken{equal}, s
	default:
		panic(fmt.Sprintf("invalid symbol: %s", string(s[0])))
	}
}

func tokenizeIntegerConstant(s string) (token, string) {
	number := 0
	for len(s) > 0 && s[0] >= '0' && s[0] <= '9' {
		number = number*10 + int(s[0]-'0')
		s = s[1:]
	}
	return integerConstantToken{number}, s
}

func tokenizeStringConstant(s string) (token, string) {
	str := []byte{}
	s = s[1:]
	for len(s) > 0 && s[0] != '"' {
		str = append(str, s[0])
		s = s[1:]
	}
	s = s[1:]
	return stringConstantToken{string(str)}, s
}

func tokenizeIdentifierOrKeyword(s string) (token, string) {
	str := []byte{}
	for len(s) > 0 && (isNumber(s[0]) || isAlphaBeta(s[0]) || s[0] == '_') {
		str = append(str, s[0])
		s = s[1:]
	}
	curKeywordTokenType := getKeyword(string(str))
	if curKeywordTokenType == notAKeyWord {
		return identifierToken{string(str)}, s
	}
	return keywordToken{curKeywordTokenType}, s
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlphaBeta(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func getKeyword(s string) keywordTokenType {
	switch s {
	case "class":
		return class
	case "constructor":
		return constructor
	case "function":
		return function
	case "method":
		return method
	case "field":
		return field
	case "static":
		return static
	case "var":
		return _var
	case "int":
		return _int
	case "char":
		return char
	case "boolean":
		return boolean
	case "void":
		return void
	case "true":
		return _true
	case "false":
		return _false
	case "null":
		return null
	case "this":
		return this
	case "let":
		return let
	case "do":
		return do
	case "if":
		return _if
	case "else":
		return _else
	case "while":
		return while
	case "return":
		return _return
	default:
		return notAKeyWord
	}
}
