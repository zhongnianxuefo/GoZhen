package zhen_0_02

import (
	"sort"
	"strconv"
	"strings"
)

type KeyWord struct {
	Name   string
	PreFun func(*CodePre) (err error)
}

type CodePre struct {
	fileCode *CodeFile

	//nowCodeBlockNo   int
	mainCodeBlock    *CodeBlock
	nowLineCodeBlock *CodeBlock
	AllKeyWords      map[string]KeyWordType
	AllOperators     map[string]Operator

	tempVarNo  map[*CodeBlock]int
	tempVarMap map[*CodeBlock]interface{}

	//state         *State
}

func NewCodePP(code *CodeFile) (codePP CodePre) {
	codePP.fileCode = code
	codePP.AllKeyWords = make(map[string]KeyWordType)
	codePP.AllOperators = make(map[string]Operator)
	codePP.mainCodeBlock = codePP.getCodeBlock(0)

	codePP.tempVarNo = make(map[*CodeBlock]int)
	codePP.tempVarMap = make(map[*CodeBlock]interface{})

	codePP.loadBaseKeyWord()
	codePP.loadBaseOperator()

	return
}
func (codePre *CodePre) getCodeBlock(n int) (codeBlock *CodeBlock) {
	codeBlock = &codePre.fileCode.AllCodeBlock[n]
	return
}
func (codePre *CodePre) getCodeBlockArea(codeBlock *CodeBlock) (codeBlockArea *CodeBlock) {
	if codeBlock.ParNo >= 0 {
		parCodeBlock := codePre.getCodeBlock(codeBlock.ParNo)
		switch parCodeBlock.BlockType {
		case CbtFile, CbtFun:
			codeBlockArea = parCodeBlock
		case CbtLine, CbtChildLine, CbtColon:
			codeBlockArea = parCodeBlock
		default:
			codeBlockArea = codePre.getCodeBlockArea(parCodeBlock)
		}
	}
	return
}

func (codePre *CodePre) getCodeBlockArea2(codeBlock *CodeBlock) (codeBlockArea *CodeBlock) {
	switch codeBlock.BlockType {
	case CbtFile, CbtFun:
		codeBlockArea = codeBlock
	default:
		parCodeBlock := codePre.getCodeBlock(codeBlock.ParNo)
		if parCodeBlock != nil {
			codeBlockArea = codePre.getCodeBlockArea(parCodeBlock)
		}
	}
	return
}

func (codePre *CodePre) addCodeStep(codeBlock *CodeBlock, step CodeStep) {
	parCodeBlock := codePre.getCodeBlockArea(codeBlock)
	if parCodeBlock != nil {
		parCodeBlock.Steps = append(parCodeBlock.Steps, step)
	}

	return
}

func (codePre *CodePre) getNextEnableCodeBlock(codeBlock *CodeBlock) (nextEnableCodeBlock *CodeBlock) {
	if codeBlock.NextNo >= 0 {
		nextCodeBlock := codePre.getCodeBlock(codeBlock.NextNo)
		switch nextCodeBlock.BlockType {
		case CbtLetter:
			nextEnableCodeBlock = nextCodeBlock
		case CbtNumber, CbtString:
			nextEnableCodeBlock = nextCodeBlock
		case CbtOperator, CbtPoint:
			nextEnableCodeBlock = nextCodeBlock
		case CbtLeftBracket, CbtLeftSquareBracket, CbtLeftBigBracket:
			nextEnableCodeBlock = nextCodeBlock
		default:
			nextEnableCodeBlock = codePre.getNextEnableCodeBlock(nextCodeBlock)
		}
	}

	return
}

func (codePre *CodePre) Preprocess() (err error) {
	err = codePre.CodePreprocess(codePre.mainCodeBlock)
	if err != nil {
		return
	}

	return
}

func (codePre *CodePre) addKeyWord(keyWord string, keyWordType KeyWordType) {
	codePre.AllKeyWords[keyWord] = keyWordType

}

func (codePre *CodePre) getKeyWord(keyWord string) (keyWordType KeyWordType, ok bool) {
	keyWordType, ok = codePre.AllKeyWords[keyWord]
	return
}

func (codePre *CodePre) addOperator(opeName string, opeType OperatorType, opePriority int) {
	ope := Operator{Type: opeType, Priority: opePriority}
	codePre.AllOperators[opeName] = ope
}

func (codePre *CodePre) getOperator(opeName string) (operator Operator, ok bool) {
	operator, ok = codePre.AllOperators[opeName]
	return
}

func (codePre *CodePre) loadBaseKeyWord() {
	codePre.addKeyWord("程序名", KwtDefineText)
	codePre.addKeyWord("版本号", KwtDefineText)

	codePre.addKeyWord("定义", KwtDefineVar)

	codePre.addKeyWord("常量", KwtDefineConstant)
	codePre.addKeyWord("定义常量", KwtDefineConstant)

	codePre.addKeyWord("变量", KwtDefineVar)
	codePre.addKeyWord("定义变量", KwtDefineVar)

	codePre.addKeyWord("全局变量", KwtDefineGlobalVar)
	codePre.addKeyWord("定义全局变量", KwtDefineGlobalVar)

	codePre.addKeyWord("局部变量", KwtDefineLocalVar)
	codePre.addKeyWord("定义局部变量", KwtDefineLocalVar)

	codePre.addKeyWord("如果", KwtIf)
	codePre.addKeyWord("否则", KwtElse)

	codePre.addKeyWord("循环", KwtWhile)
	codePre.addKeyWord("按条件循环", KwtWhile)
	codePre.addKeyWord("按次数循环", KwtFor)

	codePre.addKeyWord("定义函数", KwtDefineFun)
	codePre.addKeyWord("参数", KwtDefineFunPara)
	codePre.addKeyWord("返回", KwtDefineFunReturn)

	codePre.addKeyWord("运行", KwtCallFun)
	codePre.addKeyWord("显示", KwtCallFun)

	return
}

func (codePre *CodePre) loadBaseOperator() {

	codePre.addOperator("=", OtAs, 10)
	codePre.addOperator("+", OtAdd, 20)
	codePre.addOperator("-", OtSub, 20)
	codePre.addOperator("*", OtMul, 30)
	codePre.addOperator("/", OtDiv, 30)
	codePre.addOperator(".", OtPoint, 50)

}
func (codePre *CodePre) getOperatorResultNo(codeBlock *CodeBlock) (no int, ok bool) {
	item, ok := codePre.tempVarMap[codeBlock]
	if ok {
		switch value := item.(type) {
		case int:
			no = value
		case *CodeBlock:
			c := value
			no, ok = codePre.getOperatorResultNo(c)
		}
	}
	return
}

func (codePre *CodePre) changeOperatorResultNo(codeBlock *CodeBlock, resultCodeBlock *CodeBlock) {
	item, ok := codePre.tempVarMap[codeBlock]
	if ok {
		switch value := item.(type) {
		case *CodeBlock:
			c := value
			codePre.changeOperatorResultNo(c, resultCodeBlock)
			return
		}
	}
	codePre.tempVarMap[codeBlock] = resultCodeBlock

	return
}

func (codePre *CodePre) addCodeStepOperator(leftVar *CodeBlock, operator *CodeBlock, rightVar *CodeBlock, resultVarNo int) {

	step := CodeStep{}
	switch operator.Operator.Type {
	case OtAs:
		step.CodeStepType = CstAs
	case OtAdd:
		step.CodeStepType = CstAdd
	case OtSub:
		step.CodeStepType = CstSub
	case OtMul:
		step.CodeStepType = CstMul
	case OtDiv:
		step.CodeStepType = CstDiv
	case OtPoint:
		step.CodeStepType = CstPoint
	case OtUnSet:
		step.CodeStepType = CstTryCall
	}

	if leftVar != nil {
		tempVar, ok := codePre.getOperatorResultNo(leftVar)

		if ok {
			step.TempVarNo1 = tempVar
		} else {
			switch leftVar.BlockType {
			case CbtLetter:
				step.VarName1 = leftVar.Chars
			case CbtString, CbtNumber:
				//step.TempVarNo1 = resultVarNo + 1
				step.ValueString1 = leftVar.Chars
				//codePre.setTempValue(step.TempVarNo1, leftVar)
			}

		}
	}
	if rightVar != nil {
		tempVar, ok := codePre.getOperatorResultNo(rightVar)

		if ok {
			step.TempVarNo2 = tempVar
		} else {
			switch rightVar.BlockType {
			case CbtLetter:
				step.VarName2 = rightVar.Chars
			case CbtString, CbtNumber:
				step.ValueString2 = rightVar.Chars
				//step.TempVarNo2 = resultVarNo + 2
				//codePre.setTempValue(step.TempVarNo2, rightVar)
			}
		}
	}
	step.ReturnVarNo = resultVarNo
	if leftVar != nil && leftVar != operator {
		codePre.changeOperatorResultNo(leftVar, operator)
	}
	if rightVar != nil && rightVar != operator {
		codePre.changeOperatorResultNo(rightVar, operator)
	}
	codePre.tempVarMap[operator] = resultVarNo
	codePre.addCodeStep(operator, step)
	//codePre.nowLineCodeBlock.Steps = append(codePre.nowLineCodeBlock.Steps, step)
}

func (codePre *CodePre) getOperatorPriorities(codeBlock *CodeBlock) (priorities []int) {
	priorityMap := make(map[int]bool)
	childNo := codeBlock.FirstChildNo
	for childNo >= 0 {
		childCodeBlock := codePre.getCodeBlock(childNo)

		p := childCodeBlock.Operator.Priority
		_, ok := priorityMap[p]
		if !ok {
			priorityMap[p] = true
			priorities = append(priorities, p)
		}

		childNo = childCodeBlock.NextNo
	}
	return
}
func (codePre *CodePre) checkOperatorCodeStep(codeBlock *CodeBlock, priority int) {
	if codeBlock == nil {
		return
	}
	var leftVar *CodeBlock
	var operator *CodeBlock
	var rightVar *CodeBlock
	var nextCheck *CodeBlock
	//leftUseTemp :=false
	//rightUseTemp :=true

	switch codeBlock.BlockType {
	case CbtLetter, CbtString, CbtNumber:
		leftVar = codeBlock
	case CbtLeftBracket, CbtLeftSquareBracket, CbtLeftBigBracket:
		leftVar = codeBlock
	case CbtOperator, CbtPoint:
		operator = codeBlock
	}
	if operator == nil {
		nextCodeBlock := codePre.getNextEnableCodeBlock(codeBlock)
		if nextCodeBlock != nil {
			switch nextCodeBlock.BlockType {
			case CbtOperator, CbtPoint:
				operator = nextCodeBlock
				codeBlock = nextCodeBlock
			}
		}
	}
	if operator != nil {
		nextCodeBlock := codePre.getNextEnableCodeBlock(codeBlock)
		if nextCodeBlock != nil {
			switch nextCodeBlock.BlockType {
			case CbtLetter, CbtString, CbtNumber:
				rightVar = nextCodeBlock
				//codeBlock = nextCodeBlock
			case CbtLeftBracket, CbtLeftSquareBracket, CbtLeftBigBracket:
				rightVar = nextCodeBlock
				//codeBlock = nextCodeBlock
			}
		}
		if operator.Operator.Priority == priority {
			codePre.tempVarNo[codePre.getCodeBlockArea(operator)] += 1
			returnVarNo := codePre.tempVarNo[codePre.getCodeBlockArea(operator)]
			codePre.addCodeStepOperator(leftVar, operator, rightVar, returnVarNo)
		}
		if rightVar != nil {
			nextCheck = rightVar
		} else {
			nextCheck = codeBlock
		}
	}

	if leftVar != nil && operator == nil {

		nextCodeBlock := codePre.getNextEnableCodeBlock(codeBlock)
		if nextCodeBlock != nil {
			switch nextCodeBlock.BlockType {
			case CbtLetter, CbtString, CbtNumber:
				rightVar = nextCodeBlock
				//codeBlock = nextCodeBlock
			case CbtLeftBracket, CbtLeftSquareBracket, CbtLeftBigBracket:
				rightVar = nextCodeBlock
				//codeBlock = nextCodeBlock
			}
		}
		if priority == 0 {
			codePre.tempVarNo[codePre.getCodeBlockArea(leftVar)] += 1
			returnVarNo := codePre.tempVarNo[codePre.getCodeBlockArea(leftVar)]
			codePre.addCodeStepOperator(leftVar, leftVar, rightVar, returnVarNo)
		}

		nextCheck = nextCodeBlock
	}

	if nextCheck != nil {

		codePre.checkOperatorCodeStep(nextCheck, priority)
	}

}

func (codePre *CodePre) CheckLineCodeBlock(codeBlock *CodeBlock) {
	if codeBlock != nil {
		if codeBlock.BlockType == CbtLine {
			codePre.nowLineCodeBlock = codeBlock

			//codePre.tempVarNo[codeBlock] = 0 //= make(map[*CodeBlock]int)
			//codePP.tempVarMap = make(map[*CodeBlock]interface{})

			//parCodeBlock := codePre.getCodeBlock(codeBlock.ParNo)
			//if parCodeBlock != nil {
			//	switch parCodeBlock.BlockType {
			//	case CbtFile, CbtFun:
			//		codePre.nowLineCodeBlock = codeBlock
			//		codePre.tempVarNo = 0
			//		codePre.tempVarMap = make(map[*CodeBlock]interface{})
			//	}
			//}
		}
	}
	return
}
func (codePre *CodePre) CodePreprocess(codeBlock *CodeBlock) (err error) {
	codePre.CheckLineCodeBlock(codeBlock)

	if codeBlock.WordType == CwtUnSet {
		switch codeBlock.BlockType {
		case CbtLetter:
			codePre.checkLetter(codeBlock)
		case CbtNumber:
			codePre.checkNumber(codeBlock)
		case CbtString:
			codePre.checkString(codeBlock)
		case CbtOperator, CbtPoint:
			codePre.checkOperator(codeBlock)
		case CbtLeftBracket, CbtLeftSquareBracket, CbtLeftBigBracket:
			codePre.checkBracket(codeBlock)
		}
	}
	childNo := codeBlock.FirstChildNo
	for childNo >= 0 {
		child := codePre.getCodeBlock(childNo)
		err = codePre.CodePreprocess(child)
		childNo = child.NextNo

	}
	codePre.CheckLineCodeBlock(codeBlock)
	operatorPriorities := codePre.getOperatorPriorities(codeBlock)
	if len(operatorPriorities) > 0 {
		sort.Sort(sort.Reverse(sort.IntSlice(operatorPriorities)))
		childCodeBlock := codePre.getCodeBlock(codeBlock.FirstChildNo)
		for _, p := range operatorPriorities {
			codePre.checkOperatorCodeStep(childCodeBlock, p)
		}
		codePre.changeOperatorResultNo(codeBlock, childCodeBlock)

	}
	if codeBlock.BlockType == CbtLine {

	}
	return
}
func (codePre *CodePre) setTempValue(tempValueNo int, codeBlock *CodeBlock) {
	step := CodeStep{}
	step.ValueString1 = codeBlock.Chars
	step.TempVarNo1 = tempValueNo
	step.CodeStepType = CstAs

	codePre.addCodeStep(codeBlock, step)
	//codePre.nowLineCodeBlock.Steps = append(codePre.nowLineCodeBlock.Steps, step)
	return
}

func (codePre *CodePre) checkLetter(codeBlock *CodeBlock) (err error) {

	check := false
	if !check {
		check = codePre.isKeyWord(codeBlock)
	}

	if !check {
		check = codePre.isConstant(codeBlock)

	}
	//if !check {
	//	check = codePre.CheckIsLocalVarName(codeBlock, codeBlock.getCodeArea())
	//
	//}
	//if !check {
	//	check = codePre.CheckIsGlobalVarName(codeBlock, codeBlock.getCodeArea())
	//
	//}

	//if !check {
	//	check = codePre.CheckIsFunName(codeBlock, codeBlock.getCodeArea())
	//
	//}

	//if err != nil {
	//	return
	//}
	//if isKeyWord {
	//	codeBlock.WordType = CwtKeyWord
	//	//k := ZhenValueToKeyWord(keyWord)
	//	//k.PreFun(codePre.state)
	//}
	//_ = keyWord
	return
}

func (codePre *CodePre) isConstant(codeBlock *CodeBlock) (check bool) {
	//identifier := codeBlock.Chars
	//keyWordType, ok := codePre.getKeyWord(identifier)
	return
}
func (codePre *CodePre) isKeyWord(codeBlock *CodeBlock) (check bool) {
	identifier := codeBlock.Chars
	keyWordType, ok := codePre.getKeyWord(identifier)
	if ok {
		codeBlock.WordType = CwtKeyWord
		switch keyWordType {
		case KwtDefineText:
			codePre.DefineText(codeBlock)
		case KwtDefineVar:
			codePre.DefineVar(codeBlock, CwtVar)
		case KwtDefineGlobalVar:
			codePre.DefineVar(codeBlock, CwtGlobalVar)
		case KwtDefineLocalVar:
			codePre.DefineVar(codeBlock, CwtLocalVar)
		case KwtDefineConstant:
			codePre.DefineVar(codeBlock, CwtConstant)
		case KwtDefineFun:
			codePre.DefineFun(codeBlock)
		case KwtDefineFunPara:
			codePre.DefineVar(codeBlock, CwtFunPara)
		case KwtDefineFunReturn:
			codePre.DefineVar(codeBlock, CwtFunReturn)
		case KwtCallFun:
			codePre.DefineCallFun(codeBlock)
		}

		check = true
	}

	return
}

func (codePre *CodePre) DefineCallFun(codeBlock *CodeBlock) {

	//
	//nextNo := codeBlock.NextNo
	//for nextNo >= 0 {
	//	nextCodeBlock := codePre.getCodeBlock(nextNo)
	//	switch nextCodeBlock {
	//
	//
	//	}
	//	if nextCodeBlock.BlockType == CbtColon {
	//		step := CodeStep{}
	//		step.CodeStepType = CstDefineText
	//		step.VarName1 = codeBlock.Chars
	//		step.ValueString1 = strings.Join(values, " ")
	//		codePre.addCodeStep(codeBlock, step)
	//	}
	//	nextNo = codeBlock.NextNo
	//}
	//
	//codeBlock.WordType = CwtKeyWord

}
func (codePre *CodePre) DefineText(codeBlock *CodeBlock) {
	var values []string
	isDef := false
	nextNo := codeBlock.NextNo
	if nextNo >= 0 {
		nextCodeBlock := codePre.getCodeBlock(nextNo)
		if nextCodeBlock.BlockType == CbtColon {
			isDef = true
			childNo := nextCodeBlock.FirstChildNo
			for childNo > 0 {
				childCodeBlock := codePre.getCodeBlock(childNo)
				childCodeBlock.WordType = CwtTxt
				//childCodeBlock.Word = c.getChars()
				values = append(values, childCodeBlock.Chars)
				childNo = childCodeBlock.NextNo
			}
		}
	}
	if isDef {
		codeBlock.WordType = CwtKeyWord
		step := CodeStep{}
		step.CodeStepType = CstDefineText
		step.VarName1 = codeBlock.Chars
		step.ValueString1 = strings.Join(values, " ")
		codePre.addCodeStep(codeBlock, step)
		//codePre.nowLineCodeBlock.Steps = append(codePre.nowLineCodeBlock.Steps, step)
		//nowCodeBlock.globalVars.SetByName(nowCodeBlock.Word, strings.Join(values, " "))
	} else {
		codeBlock.WordType = CwtConstant
	}

	return
}
func (codePre *CodePre) DefineVar(codeBlock *CodeBlock, varType CodeWordType) {
	nextNo := codeBlock.NextNo
	if nextNo >= 0 {
		nextCodeBlock := codePre.getCodeBlock(nextNo)
		if nextCodeBlock.BlockType == CbtColon {

			childNo := nextCodeBlock.FirstChildNo
			for childNo >= 0 {
				child := codePre.getCodeBlock(childNo)
				codePre.VarStatement(child, varType)

				childNo = child.NextNo
			}
		}
	}

	codeBlock.WordType = CwtKeyWord

}
func (codePre *CodePre) VarStatement(codeBlock *CodeBlock, varType CodeWordType) (err error) {
	switch codeBlock.BlockType {
	case CbtLetter:
		step := CodeStep{}
		step.VarName1 = codeBlock.Chars
		switch varType {
		case CwtVar:
			step.CodeStepType = CstDefineVar
		case CwtLocalVar:
			step.CodeStepType = CstDefineLocalVar
		case CwtGlobalVar:
			step.CodeStepType = CstDefineGlobalVar
		case CwtConstant:
			step.CodeStepType = CstDefineConstant
		case CwtFunPara:
			step.CodeStepType = CstDefineFunPara
		case CwtFunReturn:
			step.CodeStepType = CstDefineFunReturn
		}
		if step.CodeStepType != CstNone {
			codePre.addCodeStep(codeBlock, step)
			//codePre.nowLineCodeBlock.Steps = append(codePre.nowLineCodeBlock.Steps, step)
		}

	case CbtLine, CbtChildLine:
		childNo := codeBlock.FirstChildNo
		for childNo >= 0 {
			child := codePre.getCodeBlock(childNo)
			codePre.VarStatement(child, varType)
			childNo = child.NextNo
		}
	}
	return
}

func (codePre *CodePre) DefineFun(codeBlock *CodeBlock) {
	var funNameWords []string
	var returnCodeBlock *CodeBlock
	next := codeBlock
	for next.NextNo >= 0 {
		next = codePre.getCodeBlock(next.NextNo)
		if next.BlockType == CbtColon || next.BlockType == CbtLine {
			break
		}
		name := next.Chars
		returnCodeBlock = next
		next.WordType = CwtFunName
		funNameWords = append(funNameWords, name)
	}

	codeBlock.WordType = CwtKeyWord
	funName := strings.Join(funNameWords, "")

	if funName != "" {
		funCodeBlock := codePre.getCodeBlock(codeBlock.ParNo)
		funCodeBlock.BlockType = CbtFun
		step := CodeStep{}
		step.CodeStepType = CstDefineFun
		step.VarName1 = funName
		funCodeBlock.Steps = append(funCodeBlock.Steps, step)
		stepReturn := CodeStep{}
		stepReturn.CodeStepType = CstDefineFunReturn
		stepReturn.VarName1 = returnCodeBlock.Chars
		if len(funNameWords) > 1 {
			returnCodeBlock.WordType = CwtFunReturn
		}
		funCodeBlock.Steps = append(funCodeBlock.Steps, stepReturn)

	}

	return
}

func (codePre *CodePre) checkString(codeBlock *CodeBlock) (err error) {
	codeBlock.WordType = CwtString
	return
}
func (codePre *CodePre) checkBracket(codeBlock *CodeBlock) (err error) {
	codeBlock.WordType = CWtBracket
	return
}

func (codePre *CodePre) checkNumber(codeBlock *CodeBlock) (err error) {
	identifier := codeBlock.Chars
	if strings.Index(identifier, ".") >= 0 {
		_, e := strconv.ParseFloat(identifier, 64)
		if e == nil {
			codeBlock.WordType = CwtFloat
		}

	} else {
		_, e := strconv.Atoi(identifier)
		if e == nil {
			codeBlock.WordType = CwtInt
		}
	}
	return
}
func (codePre *CodePre) checkOperator(codeBlock *CodeBlock) (err error) {
	identifier := codeBlock.Chars
	operator, ok := codePre.getOperator(identifier)
	if ok {
		codeBlock.Operator = operator
		switch codeBlock.Operator.Type {

		case OtAs:
			codeBlock.WordType = CwtAs
		case OtAdd:
			codeBlock.WordType = CwtAdd
		case OtSub:
			codeBlock.WordType = CwtSub
		case OtMul:
			codeBlock.WordType = CwtMul
		case OtDiv:
			codeBlock.WordType = CwtDiv
		case OtPoint:
			codeBlock.WordType = CwtPoint

		}
	}

	return
}
