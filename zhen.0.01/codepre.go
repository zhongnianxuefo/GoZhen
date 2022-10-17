package zhen_0_01

import (
	"fmt"
	"strings"
)

type CodePre struct {
	FileCodeBlock *CodeBlock2

	NowCodeBlock *CodeBlock2

	//state         *ZhenState
}

func NewCodePre(block *CodeBlock2) (codePre CodePre) {
	codePre.FileCodeBlock = block
	//codePre.NowFunCodeBlock = codePre.FileCodeBlock
	codePre.LoadBaseCodePre()
	codePre.LoadBaseGoFun()
	codePre.LoadBaseOperatorPre()
	//codePre.state = state
	return
}

func (codePre *CodePre) Preprocess() (err error) {
	//codePre.FileCodeBlock = block
	for _, code := range codePre.FileCodeBlock.items {

		//codePre.NowLineCodeBlock = code

		err = codePre.CodePreprocess(code)
		if err != nil {
			return
		}
	}
	return
}
func (codePre *CodePre) AddKeyWord(keyWord KeyWord) {
	codePre.FileCodeBlock.keyWords.SetByName(keyWord.Name, ZValue(keyWord))

}
func (codePre *CodePre) GetKeyWord(identifier string) (isKeyWord bool, keyWord KeyWord) {
	k := codePre.FileCodeBlock.keyWords.GetByName(identifier)
	keyWord, isKeyWord = k.(KeyWord)

	return
}
func (codePre *CodePre) AddFun(FunName string, fun ZFun) {
	codePre.FileCodeBlock.functions.SetByName(FunName, ZValue(fun))

}
func (codePre *CodePre) AddOperator(opeName string, opeType OperatorType, opePriority int) {
	ope := Operator{Name: opeName, Type: opeType, Priority: opePriority}
	codePre.FileCodeBlock.operators.SetByName(opeName, ope)

}
func (codePre *CodePre) GetOperator(identifier string) (isOperator bool, ope Operator) {
	k := codePre.FileCodeBlock.operators.GetByName(identifier)
	ope, isOperator = k.(Operator)

	return
}

func (codePre *CodePre) GetFun(FunName string) (exist bool, fun ZFun) {
	k := codePre.FileCodeBlock.functions.GetByName(FunName)
	fun, exist = k.(ZFun)
	return
}
func (codePre *CodePre) LoadBaseOperatorPre() {

	codePre.AddOperator("=", OtAs, 10)
	codePre.AddOperator("+", OtAdd, 20)
	codePre.AddOperator("-", OtSub, 20)
	codePre.AddOperator("*", OtMul, 30)
	codePre.AddOperator("/", OtDiv, 30)
	codePre.AddOperator(".", OtPoint, 50)

}
func (codePre *CodePre) LoadBaseGoFun() {
	codePre.AddFun("显示", ShowVar)
}
func (codePre *CodePre) LoadBaseCodePre() {

	codePre.AddKeyWord(KeyWord{Name: "程序名", PreFun: DefineTextPreFun})
	codePre.AddKeyWord(KeyWord{Name: "版本号", PreFun: DefineTextPreFun})
	//
	codePre.AddKeyWord(KeyWord{Name: "定义", PreFun: DefineVarPreFun})

	codePre.AddKeyWord(KeyWord{Name: "常量", PreFun: DefineConstantPreFun})
	codePre.AddKeyWord(KeyWord{Name: "定义常量", PreFun: DefineConstantPreFun})

	codePre.AddKeyWord(KeyWord{Name: "变量", PreFun: DefineVarPreFun})
	codePre.AddKeyWord(KeyWord{Name: "定义变量", PreFun: DefineVarPreFun})

	codePre.AddKeyWord(KeyWord{Name: "全局变量", PreFun: DefineGlobalVarPreFun})
	codePre.AddKeyWord(KeyWord{Name: "定义全局变量", PreFun: DefineGlobalVarPreFun})

	codePre.AddKeyWord(KeyWord{Name: "局部变量", PreFun: DefineLocalVarPreFun})
	codePre.AddKeyWord(KeyWord{Name: "定义局部变量", PreFun: DefineLocalVarPreFun})

	//codePre.addKeyWord(KeyWord{Names: "如果", TokenType: KwtIf, PreFun: DefineVar})
	//codePre.addKeyWord(KeyWord{Names: "否则", TokenType: KwtElse, PreFun: DefineVar})
	//
	//codePre.addKeyWord(KeyWord{Names: "循环", TokenType: KwtWhile, PreFun: DefineVar})
	//codePre.addKeyWord(KeyWord{Names: "按条件循环", TokenType: KwtWhile, PreFun: DefineVar})
	//codePre.addKeyWord(KeyWord{Names: "按次数循环", TokenType: KwtFor, PreFun: DefineVar})
	//
	codePre.AddKeyWord(KeyWord{Name: "定义函数", PreFun: DefineFunPreFun})
	//
	codePre.AddKeyWord(KeyWord{Name: "参数", PreFun: DefineFunParaPreFun})
	codePre.AddKeyWord(KeyWord{Name: "返回", PreFun: DefineFunReturnPreFun})
	//
	codePre.AddKeyWord(KeyWord{Name: "运行", PreFun: DefineCallFunPreFun})
	//

	return
}
func (codePre *CodePre) CodePreprocess(codeBlock *CodeBlock2) (err error) {
	codePre.NowCodeBlock = codeBlock
	//fmt.Println(codeBlock.Pos, codeBlock.getChars(), codeBlock.WordType)
	//codePre.state.NowCodeBlock = codeBlock
	if codeBlock.WordType == CwtUnSet {
		switch codeBlock.BlockType {
		case CbtLetter:
			codePre.CheckLetter(codeBlock)
		case CbtNumber:
		case CbtString:
		case CbtOperator, CbtPoint:
			codePre.CheckOperator(codeBlock)
		}
	}

	for _, code := range codeBlock.items {
		err = codePre.CodePreprocess(code)
		if err != nil {
			return
		}
	}
	return
}

func (codePre *CodePre) CheckOperator(codeBlock *CodeBlock2) (err error) {
	identifier := codeBlock.getChars()
	isOperator, operator := codePre.GetOperator(identifier)

	if isOperator {
		codeBlock.Operator = operator
	}
	return
}
func (codePre *CodePre) CheckLetter(codeBlock *CodeBlock2) (err error) {
	check := false

	if !check {
		check = codePre.CheckIsKeyWord(codeBlock)

	}
	if !check {
		check = codePre.CheckIsLocalVarName(codeBlock, codeBlock.getCodeArea())

	}
	if !check {
		check = codePre.CheckIsGlobalVarName(codeBlock, codeBlock.getCodeArea())

	}
	if !check {
		check = codePre.CheckIsConstantVarName(codeBlock, codeBlock.getCodeArea())

	}
	if !check {
		check = codePre.CheckIsFunName(codeBlock, codeBlock.getCodeArea())

	}

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

func (codePre *CodePre) CheckIsLocalVarName(codeBlock *CodeBlock2, codeArea *CodeBlock2) (check bool) {
	identifier := codeBlock.getChars()
	n := codeArea.localVars.FindByName(identifier)
	if n >= 0 {
		codeBlock.Word = identifier
		codeBlock.WordType = CwtLocalVar
		check = true
	}

	return
}

func (codePre *CodePre) CheckIsConstantVarName(codeBlock *CodeBlock2, codeArea *CodeBlock2) (check bool) {
	identifier := codeBlock.getChars()
	n := codeArea.constants.FindByName(identifier)
	if n >= 0 {
		codeBlock.Word = identifier
		codeBlock.WordType = CwtConstant
		check = true
	}
	if !check {
		if codeArea.BlockType != CbtFile {
			codeArea = codeArea.getCodeArea()
			return codePre.CheckIsConstantVarName(codeBlock, codeArea)
		}
	}

	return
}
func (codePre *CodePre) CheckIsGlobalVarName(codeBlock *CodeBlock2, codeArea *CodeBlock2) (check bool) {
	identifier := codeBlock.getChars()
	n := codeArea.globalVars.FindByName(identifier)
	if n >= 0 {
		codeBlock.Word = identifier
		codeBlock.WordType = CwtLocalVar
		check = true
	}
	if n < 0 && codeArea.BlockType != CbtFile {
		n := codeArea.getCodeFile().globalVars.FindByName(identifier)
		//todo 警告函数中使用全局变量需要声明
		if n >= 0 {
			codeBlock.Word = identifier
			codeBlock.WordType = CwtLocalVar
			check = true
		}
	}

	return
}
func (codePre *CodePre) CheckIsFunName(codeBlock *CodeBlock2, codeArea *CodeBlock2) (check bool) {
	identifier := codeBlock.getChars()
	v := codeArea.functions.GetByName(identifier)

	switch v.(type) {
	case ZFun:
		//fmt.Println("CheckIsFunName", codeBlock.Pos.LineNo, identifier, v, reflect.TypeOf(v))
		codeBlock.Word = identifier
		codeBlock.WordType = CwtFunName
		//keyWord.PreFun(codePre)
		//_ = keyWord
		check = true
	case *CodeBlock2:
		codeBlock.Word = identifier
		codeBlock.WordType = CwtFunName
		//keyWord.PreFun(codePre)
		//_ = keyWord
		check = true
	}
	if !check {

		if codeArea.BlockType != CbtFile {
			codeArea = codeArea.getCodeArea()
			return codePre.CheckIsFunName(codeBlock, codeArea)
		}
	}

	return
}
func (codePre *CodePre) CheckIsKeyWord(codeBlock *CodeBlock2) (check bool) {
	identifier := codeBlock.getChars()
	isKeyWord, keyWord := codePre.GetKeyWord(identifier)

	if isKeyWord {
		codeBlock.Word = identifier
		codeBlock.WordType = CwtKeyWord
		keyWord.PreFun(codePre)
		//_ = keyWord
		check = true
	}

	return
}

//func (zhen *ZhenState) SetVarFunction(valueName string, value ZhenValue) {
//	//todo 局部变量如何处理
//	zhen.globalValues[valueName] = value
//}

func DefineTextPreFun(codePre *CodePre) (err error) {
	nowCodeBlock := codePre.NowCodeBlock
	nowCodeBlock.Word = nowCodeBlock.getChars()
	var values []string
	isDef := false
	next, ok := nowCodeBlock.getNext()
	if ok {
		if next.BlockType == CbtColon {
			isDef = true
			for _, c := range next.items {
				c.WordType = CwtTxt
				c.Word = c.getChars()
				values = append(values, c.Word)
			}
		}
	}
	if isDef {
		nowCodeBlock.WordType = CwtKeyWord
		nowCodeBlock.globalVars.SetByName(nowCodeBlock.Word, strings.Join(values, " "))
	} else {
		nowCodeBlock.WordType = CwtConstant
	}

	return
}
func VarStatement(codeBlock *CodeBlock2, varType CodeWordType) (err error) {
	switch codeBlock.BlockType {
	case CbtLetter:
		name := codeBlock.getChars()
		//value := None
		//blnEnable := false
		//next, ok := codeBlock.getNext()
		codeBlock.WordType = varType
		codeBlock.Word = name

		switch varType {
		case CwtLocalVar:
			codeBlock.getCodeArea().localVars.SetByName(name, nil)
		case CwtGlobalVar:
			codeBlock.getCodeArea().globalVars.SetByName(name, nil)
			codeBlock.getCodeFile().globalVars.SetByName(name, nil)
		case CwtConstant:
			codeBlock.getCodeArea().constants.SetByName(name, nil)
		case CwtFunPara:
			codeBlock.getCodeArea().localVars.SetByName(name, nil)
		case CwtFunReturn:
			codeBlock.getCodeArea().localVars.SetByName(name, nil)
		}

		//if ok && next.BlockType == CbtOperator {
		//	if next.getChars() == "=" {
		//		next2, ok2 := next.getNext()
		//		if ok2 {
		//			switch next2.BlockType {
		//			case CbtString:
		//				value = StringToZhenValue(next2.getChars())
		//				blnEnable = true
		//			case CbtNumber:
		//				n, err := strconv.ParseFloat(next2.getChars(), 64)
		//				if err != nil {
		//					return err
		//				}
		//				value = NewZhenValueNumber(ZFloat(n))
		//				blnEnable = true
		//			}
		//		}
		//	}
		//}
		//if blnEnable {
		//	cs := ZhenCodeStep{codeStepType: ZCS_Var, valueName1: name, value: value}
		//	codeBlock.codeSteps = append(codeBlock.codeSteps, cs)
		//}

	case CbtLine, CbtChildLine:
		for _, c := range codeBlock.items {
			VarStatement(c, varType)
		}
	}
	return
}
func CodeBlockDefineVar(nowCodeBlock *CodeBlock2, varType CodeWordType) (err error) {
	nowCodeBlock.Word = nowCodeBlock.getChars()
	next, ok := nowCodeBlock.getNext()
	if ok {
		if next.BlockType == CbtColon {
			for _, c := range next.items {
				VarStatement(c, varType)
			}
		}
	}
	nowCodeBlock.WordType = CwtKeyWord

	return
}
func DefineVarPreFun(codePre *CodePre) (err error) {
	cwt := CwtLocalVar
	if codePre.NowCodeBlock.getCodeArea() == codePre.FileCodeBlock {
		cwt = CwtGlobalVar
	}
	return CodeBlockDefineVar(codePre.NowCodeBlock, cwt)
}
func DefineLocalVarPreFun(codePre *CodePre) (err error) {
	return CodeBlockDefineVar(codePre.NowCodeBlock, CwtGlobalVar)
}

func DefineGlobalVarPreFun(codePre *CodePre) (err error) {
	return CodeBlockDefineVar(codePre.NowCodeBlock, CwtGlobalVar)
}

func DefineConstantPreFun(codePre *CodePre) (err error) {
	return CodeBlockDefineVar(codePre.NowCodeBlock, CwtConstant)
}
func DefineFunPreFun(codePre *CodePre) (err error) {
	nowCodeBlock := codePre.NowCodeBlock
	nowCodeBlock.Word = nowCodeBlock.getChars()
	var funNameWords []string
	var returnCodeBlock *CodeBlock2
	next := nowCodeBlock
	ok := true
	for {
		next, ok = next.getNext()
		if ok == false {
			break
		}
		if next.BlockType == CbtColon {
			break
		}
		name := next.getChars()
		next.Word = name
		next.WordType = CwtFunName
		funNameWords = append(funNameWords, name)
		returnCodeBlock = next
	}

	nowCodeBlock.WordType = CwtKeyWord
	//nowCodeBlock.parCodeBlock.BlockType = CbtFun
	funCodeBlock := nowCodeBlock.parCodeBlock
	funCodeBlock.BlockType = CbtFun
	funName := strings.Join(funNameWords, "")

	funCodeBlock.getCodeArea().functions.SetByName(funName, funCodeBlock)
	//funCodeBlock.getCodeArea().localVars.SetByName(returnCodeBlock.getChars(), funCodeBlock)
	//func (codePre *CodePre) AddFun(FunName string, fun ZFun) {
	//	codePre.FileCodeBlock.functions.SetByName(FunName, ZValue(fun))
	//
	//}
	//func (codePre *CodePre) GetFun(FunName string) (exist bool, fun ZFun) {
	//	k := codePre.FileCodeBlock.functions.GetByName(FunName)
	//	fun, exist = k.(ZFun)
	//	return
	//}

	//todo 返回参数名添加到本地变量列表中
	_ = returnCodeBlock

	return
}

func DefineFunParaPreFun(codePre *CodePre) (err error) {
	return CodeBlockDefineVar(codePre.NowCodeBlock, CwtFunPara)
}
func DefineFunReturnPreFun(codePre *CodePre) (err error) {
	return CodeBlockDefineVar(codePre.NowCodeBlock, CwtFunReturn)
}
func DefineCallFunPreFun(codePre *CodePre) (err error) {
	nowCodeBlock := codePre.NowCodeBlock
	nowCodeBlock.Word = nowCodeBlock.getChars()
	next, ok := nowCodeBlock.getNext()
	if ok {
		if next.BlockType == CbtLetter {
			next.Word = next.getChars()
			next.WordType = CwtFunName
		}
	}
	nowCodeBlock.WordType = CwtKeyWord
	return
}

func DefineGoFunPreFun(codePre *CodePre) (err error) {
	nowCodeBlock := codePre.NowCodeBlock
	nowCodeBlock.Word = nowCodeBlock.getChars()

	nowCodeBlock.WordType = CwtFunName
	return
}

func ShowVar(state *ZhenState) (err error) {
	//todo 具体内容待开发
	fmt.Println("ShowVar")
	return
}
