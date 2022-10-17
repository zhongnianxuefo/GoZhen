package zhen_0_01

import (
	"fmt"
	"github.com/beevik/etree"
	"io/ioutil"
	"os"
)

const ZhenTempValuesLen = 10

type ZhenState struct {
	debug bool

	globalValues *VarGroup
	localValues  *VarGroup

	tempValues [ZhenTempValuesLen]ZValue

	//txtCode  TxtCode
	allCodes []ZhenCodeOld

	allCodeStepPointers []*ZhenCodeStep
	runCodeStep         int

	NowCodeBlock *CodeBlock2

	MainCodeBlock *CodeBlock2
	KeyWords      *VarGroupWithDef
	Operators     *VarGroupWithDef
	Constants     *VarGroupWithDef
	Functions     *VarGroupWithDef
	GlobalVars    *VarGroupWithDef
	LocalVars     *VarGroupWithDef
	RunCodeBlock  []*CodeBlock2
}

func NewZhenState() (zhen ZhenState) {
	zhen.globalValues = NewVarGroup()
	zhen.localValues = NewVarGroup()

	return

}

//func (zhen *ZhenState) RunCodeBlock() (info string) {
//
//}
func (zhen *ZhenState) getRunInfo() (info string) {
	var codeCountStep = len(zhen.allCodeStepPointers)
	if zhen.runCodeStep >= codeCountStep {
		info = "程序正常退出"
	} else {
		//s := zhen.allCodeStepPointers[zhen.runCodeStep]
		info = fmt.Sprintf("程序运行到%d步", zhen.runCodeStep)
	}
	return
}
func (zhen *ZhenState) LoadTxtFile2(fileName string) (body string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	body = string(content)
	return
}
func (zhen *ZhenState) LoadTxtFile(fileName string) (err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	return zhen.LoadTxt(string(content))
}
func (zhen *ZhenState) LoadTxt(codes string) (err error) {

	Analyze := NewTxtCodeAnalyze(codes)

	err = Analyze.AnalyzeCode()
	if err != nil {
		return nil
	}

	zhen.MainCodeBlock = Analyze.MainCode
	if err != nil {
		return
	}

	return
}

func (zhen *ZhenState) LoadString(code string) (err error) {
	doc := etree.NewDocument()
	err = doc.ReadFromString(code)
	if err != nil {
		return
	}

	program := doc.SelectElement("程序")

	for _, line := range program.ChildElements() {
		var code ZhenCodeOld
		if line.Tag == "代码行" {
			code.needRun = true
			for _, word := range line.ChildElements() {
				if word.Tag != "代码" {
					var w ZhenCodeWord
					w.tag = word.Tag
					w.content = word.Text()
					code.codeWords = append(code.codeWords, w)
				} else {
					s, err := getCodeStepFromElement(word)
					if err != nil {
						return err
					}
					code.codeSteps = append(code.codeSteps, s)
				}
			}
			zhen.AddCode(code)
		}
	}
	return

}

func (zhen *ZhenState) AddCode(code ZhenCodeOld) {
	zhen.allCodes = append(zhen.allCodes, code)
	if code.needRun {
		for i := range code.codeSteps {
			//fmt.Println(i, s, &code.codeSteps[i])
			zhen.allCodeStepPointers = append(zhen.allCodeStepPointers, &code.codeSteps[i])
		}
	}

}

//func (zhen *ZhenState) SetVarValue(valueName string, value ZhenValue) {
//	//todo 局部变量如何处理
//	zhen.localValues.SetByName(valueName, value)
//}
//
//func (zhen *ZhenState) GetVarValue(valueName string) (value ZhenValue) {
//	//todo 局部变量如何处理
//	value = zhen.localValues.GetByName(valueName)
//	return
//}
//
//func (zhen *ZhenState) SetGlobalVarValue(valueName string, value ZhenValue) {
//
//	zhen.globalValues.SetByName(valueName, value)
//}
//
//func (zhen *ZhenState) GetGlobalVarValue(valueName string) (value ZhenValue) {
//
//	value = zhen.globalValues.GetByName(valueName)
//	return
//}
//
//func (zhen *ZhenState) SetTempVarValue(tempValueNo int, value ZhenValue) {
//
//	zhen.tempValues[tempValueNo] = value
//}
//
//func (zhen *ZhenState) GetTempVarValue(tempValueNo int) (value ZhenValue) {
//
//	value = zhen.tempValues[tempValueNo]
//	return
//}

//
//func (zhen *ZhenState) Var(s *ZhenCodeStep) {
//	//todo 局部变量如何处理
//	zhen.globalValues[s.valueName1] = s.value
//	if zhen.debug {
//		fmt.Printf("定义变量：%s，初始数值为：%s\n", s.valueName1, ZhenValueToString(s.value))
//	}
//
//}
//
//func (zhen *ZhenState) Assign(s *ZhenCodeStep) {
//	//todo 局部变量如何处理
//	zhen.globalValues[s.valueName] = s.value
//	if zhen.debug {
//		fmt.Printf("变量：%s，赋值为：%s\n", s.valueName, ZhenValueToString(s.value))
//	}
//}
//func (zhen *ZhenState) PrintVar(s *ZhenCodeStep) {
//	//todo 局部变量如何处理
//	v := zhen.globalValues[s.valueName]
//
//
//}
//
//func (zhen *ZhenState) codeStepRun(s *ZhenCodeStep) (err error) {
//	st := s.codeStepType
//
//	switch st {
//	case ZCS_None:
//
//	case ZCS_Var:
//		zhen.SetGlobalVarValue(s.valueName1, s.value)
//	case ZCS_As:
//		v := zhen.GetGlobalVarValue(s.valueName1)
//		zhen.SetGlobalVarValue(s.valueName2, v)
//
//	case ZCS_Add, ZCS_Sub, ZCS_Mul, ZCS_Div, ZCS_Eq, ZCS_Ne, ZCS_Gt, ZCS_Lt, ZCS_And, ZCS_Or:
//		v1 := zhen.GetGlobalVarValue(s.valueName1)
//		v2 := zhen.GetGlobalVarValue(s.valueName2)
//		var v ZhenValue
//		v, err = ZhenValueOperation(st, v1, v2)
//		if err != nil {
//			return
//		}
//		zhen.SetTempVarValue(s.tempValueNo1, v)
//
//	case ZCS_Not:
//		v1 := zhen.GetGlobalVarValue(s.valueName1)
//		var v2 ZhenValue
//		var v ZhenValue
//		v, err = ZhenValueOperation(st, v1, v2)
//		if err != nil {
//			return
//		}
//		zhen.SetTempVarValue(s.tempValueNo1, v)
//
//	case ZCS_TVar:
//		zhen.SetTempVarValue(s.tempValueNo1, s.value)
//	case ZCS_TFrom:
//		v := zhen.GetGlobalVarValue(s.valueName1)
//		zhen.SetTempVarValue(s.tempValueNo1, v)
//	case ZCS_TAs:
//		v := zhen.GetTempVarValue(s.tempValueNo1)
//		zhen.SetGlobalVarValue(s.valueName1, v)
//	case ZCS_TAdd, ZCS_TSub, ZCS_TMul, ZCS_TDiv, ZCS_TEq, ZCS_TNe, ZCS_TGt, ZCS_TLt, ZCS_TAnd, ZCS_TOr:
//		v1 := zhen.GetTempVarValue(s.tempValueNo1)
//		v2 := zhen.GetGlobalVarValue(s.valueName1)
//		var v ZhenValue
//
//		v, err = ZhenValueOperation(st, v1, v2)
//		if err != nil {
//			return
//		}
//		zhen.SetTempVarValue(s.tempValueNo1, v)
//
//	case ZCS_TNot:
//
//		v1 := zhen.GetTempVarValue(s.tempValueNo1)
//		var v2 ZhenValue
//		var v ZhenValue
//		v, err = ZhenValueOperation(st, v1, v2)
//		if err != nil {
//			return
//		}
//		zhen.SetTempVarValue(s.tempValueNo1, v)
//
//	case ZCS_TTAs:
//		v := zhen.GetTempVarValue(s.tempValueNo1)
//		zhen.SetTempVarValue(s.tempValueNo2, v)
//	case ZCS_TTAdd, ZCS_TTSub, ZCS_TTMul, ZCS_TTDiv, ZCS_TTEq, ZCS_TTNe, ZCS_TTGt, ZCS_TTLt, ZCS_TTAnd, ZCS_TTOr:
//		v1 := zhen.GetTempVarValue(s.tempValueNo1)
//		v2 := zhen.GetTempVarValue(s.tempValueNo2)
//		var v ZhenValue
//
//		v, err = ZhenValueOperation(st, v1, v2)
//		if err != nil {
//			return
//		}
//		zhen.SetTempVarValue(s.tempValueNo1, v)
//
//	case ZCS_If:
//		//todo 条件判断指令
//
//	case ZCS_For:
//		//todo 次数循环指令
//	case ZCS_While:
//		//todo 条件循环指令
//	case ZCS_Break:
//		//todo 跳出循环指令
//	case ZCS_Return:
//		//todo 返回指令
//	case ZCS_Call:
//		//todo 运行函数指令
//	case ZCS_PrintVar:
//		v := zhen.GetGlobalVarValue(s.valueName1)
//		fmt.Printf("变量：%s，值为：%s\n", s.valueName1, ZhenValueToString(v))
//
//	default:
//		fmt.Println(s.codeStepType)
//		err = errors.Copy("未知指令")
//
//	}
//	return
//}

//func (zhen *ZhenState) Run() (err error) {
//	zhen.runCodeStep = 0
//	var codeCountStep = len(zhen.allCodeStepPointers)
//	if codeCountStep == 0 {
//		return
//	}
//	for {
//
//		s := zhen.allCodeStepPointers[zhen.runCodeStep]
//		err = zhen.codeStepRun(s)
//		if err != nil {
//			return err
//		}
//
//		zhen.runCodeStep += 1
//		if zhen.runCodeStep >= codeCountStep {
//			break
//		}
//	}
//	return
//
//}

//
//func (zhen *ZhenState) getKeyWord(identifier string) (isKeyWord bool, keyWord ZhenValue, err error) {
//	value := zhen.GetGlobalVarValue("@关键字")
//	if value.valueType == ZvtObject {
//
//		keyWord, isKeyWord = value.valueTable[identifier]
//
//	}
//
//	return
//}
//
//func (zhen *ZhenState) AddTextKeyWord(keyWord string) (err error) {
//	value := zhen.GetGlobalVarValue("@文本型关键字")
//	if value.valueType != ZvtObject {
//		t := make(map[string]ZhenValue)
//		value = NewZhenValueTable(t)
//	}
//	value.valueTable[keyWord] = NewZhenValueFunction(DefineText)
//
//	zhen.SetGlobalVarValue("@关键字", value)
//	return
//}
//
//func (zhen *ZhenState) addKeyWord(keyWord KeyWord) (err error) {
//	globalVarName := "@关键字"
//
//	value := zhen.GetGlobalVarValue(globalVarName)
//	if value.valueType != ZvtObject {
//		t := make(map[string]ZhenValue)
//		value = NewZhenValueTable(t)
//	}
//
//	value.valueTable[keyWord.Names] = KeyWordToZhenValue(keyWord)
//
//	zhen.SetGlobalVarValue(globalVarName, value)
//	return
//}
//func (zhen *ZhenState) LoadBaseCodePre() (err error) {
//
//	//KwtDefineText
//	//KwtDefineConstant
//	//KwtDefineVar
//	//KwtDefineGlobalVar
//	//KwtDefineFun
//	//KwtDefineFunPara
//	//KwtDefineFunReturn
//	zhen.addKeyWord(KeyWord{Names: "程序名", TokenType: KwtDefineText, PreFun: DefineText})
//	zhen.addKeyWord(KeyWord{Names: "版本号", TokenType: KwtDefineText, PreFun: DefineText})
//
//	zhen.addKeyWord(KeyWord{Names: "定义", TokenType: KwtDefineVar, PreFun: DefineVar})
//
//	zhen.addKeyWord(KeyWord{Names: "常量", TokenType: KwtDefineConstant, PreFun: DefineConstantPreFun})
//	zhen.addKeyWord(KeyWord{Names: "定义常量", TokenType: KwtDefineConstant, PreFun: DefineConstantPreFun})
//
//	zhen.addKeyWord(KeyWord{Names: "变量", TokenType: KwtDefineVar, PreFun: DefineVar})
//	zhen.addKeyWord(KeyWord{Names: "定义变量", TokenType: KwtDefineVar, PreFun: DefineVar})
//
//	zhen.addKeyWord(KeyWord{Names: "全局变量", TokenType: KwtDefineGlobalVar, PreFun: DefineGlobalVarPreFun})
//	zhen.addKeyWord(KeyWord{Names: "定义全局变量", TokenType: KwtDefineGlobalVar, PreFun: DefineGlobalVarPreFun})
//
//	zhen.addKeyWord(KeyWord{Names: "如果", TokenType: KwtIf, PreFun: DefineVar})
//	zhen.addKeyWord(KeyWord{Names: "否则", TokenType: KwtElse, PreFun: DefineVar})
//
//	zhen.addKeyWord(KeyWord{Names: "循环", TokenType: KwtWhile, PreFun: DefineVar})
//	zhen.addKeyWord(KeyWord{Names: "按条件循环", TokenType: KwtWhile, PreFun: DefineVar})
//	zhen.addKeyWord(KeyWord{Names: "按次数循环", TokenType: KwtFor, PreFun: DefineVar})
//
//	zhen.addKeyWord(KeyWord{Names: "定义函数", TokenType: KwtDefineFun, PreFun: DefineFun})
//
//	zhen.addKeyWord(KeyWord{Names: "参数", TokenType: KwtDefineFunPara, PreFun: DefineFunParaPreFun})
//	zhen.addKeyWord(KeyWord{Names: "返回", TokenType: KwtDefineFunReturn, PreFun: DefineFunReturnPreFun})
//
//	zhen.addKeyWord(KeyWord{Names: "运行", TokenType: KwtCallFun, PreFun: DefineFun})
//
//	zhen.addKeyWord(KeyWord{Names: "显示", TokenType: KwtFun, PreFun: DefineFun})
//
//	return
//}

//func (zhen *ZhenState) Run2() (err error) {
//	zhen.runCodeStep = 0
//	var codeCountStep = len(zhen.allCodeStepPointers)
//	if codeCountStep == 0 {
//		return
//	}
//	for {
//
//		s := zhen.allCodeStepPointers[zhen.runCodeStep]
//		err = zhen.codeStepRun(s)
//		if err != nil {
//			return err
//		}
//
//		zhen.runCodeStep += 1
//		if zhen.runCodeStep >= codeCountStep {
//			break
//		}
//	}
//	return
//
//}
