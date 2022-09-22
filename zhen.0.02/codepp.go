package zhen_0_02

import (
	"sort"
	"strings"
)

type KeyWord struct {
	Name   string
	PreFun func(*CodePre) (err error)
}

type CodePre struct {
	fileCode *CodeFile

	nowCodeBlockNo   int
	mainCodeBlock    *CodeBlock
	nowLineCodeBlock *CodeBlock
	AllKeyWords      map[string]KeyWordType
	AllOperators     map[string]Operator

	tempVarNo   int
	tempVarUses map[*CodeBlock]interface{}

	//state         *ZhenState
}

func NewCodePP(code *CodeFile) (codePP CodePre) {
	codePP.fileCode = code
	codePP.AllKeyWords = make(map[string]KeyWordType)
	codePP.AllOperators = make(map[string]Operator)
	codePP.mainCodeBlock = &codePP.fileCode.AllCodeBlock[0]

	codePP.LoadBaseKeyWord()
	codePP.LoadBaseOperator()

	//codePP.newArea(0)
	//codePre.NowFunCodeBlock = codePre.FileCodeBlock
	//codePre.LoadBaseCodePre()
	//codePre.LoadBaseGoFun()
	//codePre.LoadBaseOperator()
	//codePre.state = state

	return
}
func (codePre *CodePre) newArea(codeNo int) {
	//varNames := NewVarNames()
	//varNames.AddVar(CodeVarKey{Name: "test", Type: CvtKeyWords})
	//codePre.fileCode.AllVarNames = append(codePre.fileCode.AllVarNames, varNames)
	//varNames := codePre.fileCode.AllVarNames[codeNo]
	//
	//codePre.fileCode.AllVarNames[codeNo] = varNames
	//codePre.NowFunCodeBlock = codePre.FileCodeBlock
	//codePre.LoadBaseCodePre()
	//codePre.LoadBaseGoFun()
	//codePre.LoadBaseOperator()
	//codePre.state = state
	return
}

func (codePre *CodePre) Preprocess() (err error) {

	err = codePre.CodePreprocess(0)
	if err != nil {
		return
	}
	//codePre.FileCodeBlock = block
	//for _, code := range codePre.FileCodeBlock.items {
	//
	//	//codePre.NowLineCodeBlock = code
	//
	//	err = codePre.CodePre(code)
	//	if err != nil {
	//		return
	//	}
	//}
	return
}
func (codePre *CodePre) AddKeyWord(keyWord string, keyWordType KeyWordType) {
	codePre.AllKeyWords[keyWord] = keyWordType

}

func (codePre *CodePre) GetKeyWord(keyWord string) (keyWordType KeyWordType, ok bool) {
	keyWordType, ok = codePre.AllKeyWords[keyWord]
	return
}
func (codePre *CodePre) AddOperator(opeName string, opeType OperatorType, opePriority int) {
	ope := Operator{Type: opeType, Priority: opePriority}
	codePre.AllOperators[opeName] = ope
}
func (codePre *CodePre) GetOperator(opeName string) (operator Operator, ok bool) {
	operator, ok = codePre.AllOperators[opeName]
	return
}

func (codePre *CodePre) LoadBaseKeyWord() {
	codePre.AddKeyWord("程序名", KwtDefineText)
	codePre.AddKeyWord("版本号", KwtDefineText)

	codePre.AddKeyWord("定义", KwtDefineVar)

	codePre.AddKeyWord("常量", KwtDefineConstant)
	codePre.AddKeyWord("定义常量", KwtDefineConstant)

	codePre.AddKeyWord("变量", KwtDefineVar)
	codePre.AddKeyWord("定义变量", KwtDefineVar)

	codePre.AddKeyWord("全局变量", KwtDefineGlobalVar)
	codePre.AddKeyWord("定义全局变量", KwtDefineGlobalVar)

	codePre.AddKeyWord("局部变量", KwtDefineLocalVar)
	codePre.AddKeyWord("定义局部变量", KwtDefineLocalVar)

	codePre.AddKeyWord("如果", KwtIf)
	codePre.AddKeyWord("否则", KwtElse)

	codePre.AddKeyWord("循环", KwtWhile)
	codePre.AddKeyWord("按条件循环", KwtWhile)
	codePre.AddKeyWord("按次数循环", KwtFor)

	codePre.AddKeyWord("定义函数", KwtDefineFun)
	codePre.AddKeyWord("参数", KwtDefineFunPara)
	codePre.AddKeyWord("返回", KwtDefineFunReturn)

	codePre.AddKeyWord("运行", KwtCallFun)

	return
}

func (codePre *CodePre) LoadBaseOperator() {

	codePre.AddOperator("=", OtAs, 10)
	codePre.AddOperator("+", OtAdd, 20)
	codePre.AddOperator("-", OtSub, 20)
	codePre.AddOperator("*", OtMul, 30)
	codePre.AddOperator("/", OtDiv, 30)
	codePre.AddOperator(".", OtPoint, 50)

}
func (codePre *CodePre) getOperatorResultNo(codeBlock *CodeBlock) (no int, ok bool) {
	item, ok := codePre.tempVarUses[codeBlock]
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
	item, ok := codePre.tempVarUses[codeBlock]
	if ok {
		switch value := item.(type) {
		case *CodeBlock:
			c := value
			codePre.changeOperatorResultNo(c, resultCodeBlock)
			return
		}
	}
	codePre.tempVarUses[codeBlock] = resultCodeBlock

	return
}

func (codePre *CodePre) addOperatorCodeStep(leftVar *CodeBlock, operator *CodeBlock, rightVar *CodeBlock, returnTempNo int) {

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
	}
	codePre.tempVarUses[operator] = returnTempNo
	if leftVar != nil {
		tempVar, ok := codePre.getOperatorResultNo(leftVar)
		codePre.changeOperatorResultNo(leftVar, operator)

		if ok {
			step.TempVarNo1 = tempVar
		} else {
			switch leftVar.BlockType {
			case CbtLetter:
				step.VarName1 = leftVar.Chars
			case CbtString, CbtNumber:
				step.TempVarNo1 = returnTempNo + 1
				codePre.setTempValue(step.TempVarNo1, leftVar)
			}

		}
	}
	if rightVar != nil {
		tempVar, ok := codePre.getOperatorResultNo(rightVar)
		codePre.changeOperatorResultNo(rightVar, operator)
		if ok {
			step.TempVarNo2 = tempVar
		} else {
			switch rightVar.BlockType {
			case CbtLetter:
				step.VarName2 = rightVar.Chars
			case CbtString, CbtNumber:
				step.TempVarNo2 = returnTempNo + 2
				codePre.setTempValue(step.TempVarNo2, rightVar)
			}
		}
	}
	step.ReturnVarNo = returnTempNo

	//codePre.tempVarNo = step.ReturnVarNo
	//fmt.Println(leftUseTempVar, rightUseTempVar, step)
	codePre.nowLineCodeBlock.Steps = append(codePre.nowLineCodeBlock.Steps, step)

}

func (codePre *CodePre) getCodeBlockOperatorPriorities(codeBlock *CodeBlock) (priorities []int) {
	priorityMap := make(map[int]bool)
	childNo := codeBlock.FirstChildNo
	for childNo >= 0 {
		childCodeBlock := &codePre.fileCode.AllCodeBlock[childNo]
		if childCodeBlock.Operator.Type != OtUnSet {
			p := childCodeBlock.Operator.Priority
			_, ok := priorityMap[p]
			if !ok {
				priorityMap[p] = true
				priorities = append(priorities, p)
			}
		}
		childNo = childCodeBlock.NextNo
	}
	return
}
func (codePre *CodePre) checkCodeBlockAddOperatorCodeStep(codeBlock *CodeBlock, priority int) {
	if codeBlock == nil {
		return
	}
	var leftVar *CodeBlock
	var operator *CodeBlock
	var rightVar *CodeBlock

	switch codeBlock.BlockType {
	case CbtLetter, CbtString, CbtNumber:
		leftVar = codeBlock
	case CbtLeftBracket, CbtLeftSquareBracket, CbtLeftBigBracket:
		leftVar = codeBlock
	case CbtOperator, CbtPoint:
		operator = codeBlock
	}
	if operator == nil {
		if codeBlock.NextNo >= 0 {
			codeBlock = &codePre.fileCode.AllCodeBlock[codeBlock.NextNo]
			switch codeBlock.BlockType {
			case CbtOperator, CbtPoint:
				operator = codeBlock
			}
		}
	}
	if operator != nil {
		if codeBlock.NextNo >= 0 {
			codeBlock = &codePre.fileCode.AllCodeBlock[codeBlock.NextNo]
			switch codeBlock.BlockType {
			case CbtLetter, CbtString, CbtNumber:
				rightVar = codeBlock
			case CbtLeftBracket, CbtLeftSquareBracket, CbtLeftBigBracket:
				rightVar = codeBlock
			}
		}
	}

	if operator != nil && operator.Operator.Priority == priority {
		codePre.tempVarNo += 1
		returnVarNo := codePre.tempVarNo
		codePre.addOperatorCodeStep(leftVar, operator, rightVar, returnVarNo)

	}
	if rightVar != nil {
		codePre.checkCodeBlockAddOperatorCodeStep(rightVar, priority)
	}

}
func (codePre *CodePre) CodePreprocess(codeBlockNo int) (err error) {

	codePre.nowCodeBlockNo = codeBlockNo
	codeBlock := &codePre.fileCode.AllCodeBlock[codePre.nowCodeBlockNo]
	if codeBlock.BlockType == CbtLine {
		codePre.nowLineCodeBlock = codeBlock
		codePre.tempVarNo = 0
		codePre.tempVarUses = make(map[*CodeBlock]interface{})
	}
	//fmt.Println(codeBlock.Pos, codeBlock.getChars(), codeBlock.WordType)
	//codePre.state.NowCodeBlock = codeBlock
	if codeBlock.WordType == CwtUnSet {
		switch codeBlock.BlockType {
		case CbtLetter:
			codePre.CheckLetter(codeBlockNo)
		case CbtNumber:
		case CbtString:
		case CbtOperator, CbtPoint:
			codePre.CheckOperator(codeBlock)
		}
	}
	childNo := codeBlock.FirstChildNo
	for childNo >= 0 {
		err = codePre.CodePreprocess(childNo)
		childNo = codePre.fileCode.AllCodeBlock[childNo].NextNo

	}

	//tt := 1
	//lastPriority := 0
	//leftUseTempVar := 0
	//rightUseTempVar := 0
	operatorPriorities := codePre.getCodeBlockOperatorPriorities(codeBlock)
	if len(operatorPriorities) > 0 {
		sort.Sort(sort.Reverse(sort.IntSlice(operatorPriorities)))
		for _, p := range operatorPriorities {
			childCodeBlock := &codePre.fileCode.AllCodeBlock[codeBlock.FirstChildNo]
			codePre.checkCodeBlockAddOperatorCodeStep(childCodeBlock, p)

		}
		childCodeBlock := &codePre.fileCode.AllCodeBlock[codeBlock.FirstChildNo]
		codePre.changeOperatorResultNo(codeBlock, childCodeBlock)

	}
	if codeBlock.BlockType == CbtLine {

	}
	return
}
func (codePre *CodePre) setTempValue(tempValueNo int, codeBlock *CodeBlock) {
	step := CodeStep{}
	step.ValueString = codeBlock.Chars
	step.TempVarNo1 = tempValueNo
	step.CodeStepType = CstAs
	codePre.nowLineCodeBlock.Steps = append(codePre.nowLineCodeBlock.Steps, step)
	return
}

func (codePre *CodePre) CheckLetter(codeBlockNo int) (err error) {
	codeBlock := &codePre.fileCode.AllCodeBlock[codeBlockNo]

	check := false
	if !check {
		check = codePre.CheckIsKeyWord(codeBlock)
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
	//	check = codePre.CheckIsConstantVarName(codeBlock, codeBlock.getCodeArea())
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
func (codePre *CodePre) CheckIsKeyWord(codeBlock *CodeBlock) (check bool) {
	identifier := codeBlock.Chars
	keyWordType, ok := codePre.GetKeyWord(identifier)

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
		}

		//keyWord.PreFun(codePre)
		//_ = keyWord
		check = true
	}

	return
}

func (codePre *CodePre) DefineText(codeBlock *CodeBlock) {
	//nowCodeBlock := codePre.NowCodeBlock
	//nowCodeBlock.Word = nowCodeBlock.getChars()
	var values []string
	isDef := false
	nextNo := codeBlock.NextNo
	if nextNo >= 0 {
		nextCodeBlock := codePre.fileCode.AllCodeBlock[nextNo]
		if nextCodeBlock.BlockType == CbtColon {
			isDef = true
			childNo := nextCodeBlock.FirstChildNo
			for childNo > 0 {
				childCodeBlock := &codePre.fileCode.AllCodeBlock[childNo]
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
		step.ValueString = strings.Join(values, " ")
		codePre.nowLineCodeBlock.Steps = append(codePre.nowLineCodeBlock.Steps, step)
		//nowCodeBlock.globalVars.SetByName(nowCodeBlock.Word, strings.Join(values, " "))
	} else {
		codeBlock.WordType = CwtConstant
	}

	return
}
func (codePre *CodePre) DefineVar(codeBlock *CodeBlock, varType CodeWordType) {

	nextNo := codeBlock.NextNo
	if nextNo >= 0 {
		nextCodeBlock := codePre.fileCode.AllCodeBlock[nextNo]
		if nextCodeBlock.BlockType == CbtColon {

			childNo := nextCodeBlock.FirstChildNo
			for childNo >= 0 {
				child := &codePre.fileCode.AllCodeBlock[childNo]
				codePre.VarStatement(child, varType)

				childNo = child.NextNo
			}
		}
	}

	codeBlock.WordType = CwtKeyWord

	//cwt := CwtLocalVar
	//if codePre.NowCodeBlock.getCodeArea() == codePre.FileCodeBlock {
	//	cwt = CwtGlobalVar
	//}
	//return CodeBlockDefineVar(codePre.NowCodeBlock, cwt)
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
			codePre.nowLineCodeBlock.Steps = append(codePre.nowLineCodeBlock.Steps, step)
		}

	case CbtLine, CbtChildLine:
		childNo := codeBlock.FirstChildNo
		for childNo >= 0 {
			child := &codePre.fileCode.AllCodeBlock[childNo]
			codePre.VarStatement(child, varType)
			childNo = child.NextNo
		}
	}
	return
}

func (codePre *CodePre) DefineFun(codeBlock *CodeBlock) {
	//nowCodeBlock := codePre.NowCodeBlock
	//nowCodeBlock.Word = nowCodeBlock.getChars()
	var funNameWords []string
	var returnCodeBlock *CodeBlock
	next := codeBlock
	for next.NextNo >= 0 {
		next = &codePre.fileCode.AllCodeBlock[next.NextNo]
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
		funCodeBlock := &codePre.fileCode.AllCodeBlock[codeBlock.ParNo]
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

	//funCodeBlock.getCodeArea().functions.SetByName(funName, funCodeBlock)
	////funCodeBlock.getCodeArea().localVars.SetByName(returnCodeBlock.getChars(), funCodeBlock)
	////func (codePre *CodePre) AddFun(FunName string, fun ZFun) {
	////	codePre.FileCodeBlock.functions.SetByName(FunName, ZValue(fun))
	////
	////}
	////func (codePre *CodePre) GetFun(FunName string) (exist bool, fun ZFun) {
	////	k := codePre.FileCodeBlock.functions.GetByName(FunName)
	////	fun, exist = k.(ZFun)
	////	return
	////}
	//
	////todo 返回参数名添加到本地变量列表中
	//_ = returnCodeBlock

	return
}

func (codePre *CodePre) CheckOperator(codeBlock *CodeBlock) (err error) {
	identifier := codeBlock.Chars
	operator, ok := codePre.GetOperator(identifier)
	if ok {
		codeBlock.Operator = operator
	}
	return
}
