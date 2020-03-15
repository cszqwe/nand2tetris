package main

import (
	"fmt"
	"strconv"
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
