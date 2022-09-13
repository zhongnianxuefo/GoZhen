package zhen

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"io/ioutil"
	"os"
)

const ZhenTempValuesLen = 10

type ZhenState struct {
	debug        bool
	globalValues ZhenValueTable
	localValues  ZhenValueTable
	tempValues   [ZhenTempValuesLen]ZhenValue

	MainCodeBlock *CodeBlock
	//txtCode  TxtCode
	allCodes []ZhenCodeOld

	allCodeStepPointers []*ZhenCodeStep
	runCodeStep         int

	NowCodeBlock *CodeBlock
}

func NewZhenState() (zhen ZhenState) {
	zhen.globalValues = make(map[string]ZhenValue)
	zhen.localValues = make(map[string]ZhenValue)

	return

}
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
func (zhen *ZhenState) SetVarValue(valueName string, value ZhenValue) {
	//todo 局部变量如何处理
	zhen.globalValues[valueName] = value
}

func (zhen *ZhenState) GetVarValue(valueName string) (value ZhenValue) {
	//todo 局部变量如何处理
	value = zhen.globalValues[valueName]
	return
}

func (zhen *ZhenState) SetGlobalVarValue(valueName string, value ZhenValue) {
	zhen.globalValues[valueName] = value
}

func (zhen *ZhenState) GetGlobalVarValue(valueName string) (value ZhenValue) {
	value = zhen.globalValues[valueName]
	return
}

func (zhen *ZhenState) SetTempVarValue(tempValueNo int, value ZhenValue) {

	zhen.tempValues[tempValueNo] = value
}

func (zhen *ZhenState) GetTempVarValue(tempValueNo int) (value ZhenValue) {

	value = zhen.tempValues[tempValueNo]
	return
}

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

func (zhen *ZhenState) codeStepRun(s *ZhenCodeStep) (err error) {
	st := s.codeStepType

	switch st {
	case ZCS_None:

	case ZCS_Var:
		zhen.SetGlobalVarValue(s.valueName1, s.value)
	case ZCS_As:
		v := zhen.GetGlobalVarValue(s.valueName1)
		zhen.SetGlobalVarValue(s.valueName2, v)

	case ZCS_Add, ZCS_Sub, ZCS_Mul, ZCS_Div, ZCS_Eq, ZCS_Ne, ZCS_Gt, ZCS_Lt, ZCS_And, ZCS_Or:
		v1 := zhen.GetGlobalVarValue(s.valueName1)
		v2 := zhen.GetGlobalVarValue(s.valueName2)
		var v ZhenValue
		v, err = ZhenValueOperation(st, v1, v2)
		if err != nil {
			return
		}
		zhen.SetTempVarValue(s.tempValueNo1, v)

	case ZCS_Not:
		v1 := zhen.GetGlobalVarValue(s.valueName1)
		var v2 ZhenValue
		var v ZhenValue
		v, err = ZhenValueOperation(st, v1, v2)
		if err != nil {
			return
		}
		zhen.SetTempVarValue(s.tempValueNo1, v)

	case ZCS_TVar:
		zhen.SetTempVarValue(s.tempValueNo1, s.value)
	case ZCS_TFrom:
		v := zhen.GetGlobalVarValue(s.valueName1)
		zhen.SetTempVarValue(s.tempValueNo1, v)
	case ZCS_TAs:
		v := zhen.GetTempVarValue(s.tempValueNo1)
		zhen.SetGlobalVarValue(s.valueName1, v)
	case ZCS_TAdd, ZCS_TSub, ZCS_TMul, ZCS_TDiv, ZCS_TEq, ZCS_TNe, ZCS_TGt, ZCS_TLt, ZCS_TAnd, ZCS_TOr:
		v1 := zhen.GetTempVarValue(s.tempValueNo1)
		v2 := zhen.GetGlobalVarValue(s.valueName1)
		var v ZhenValue

		v, err = ZhenValueOperation(st, v1, v2)
		if err != nil {
			return
		}
		zhen.SetTempVarValue(s.tempValueNo1, v)

	case ZCS_TNot:

		v1 := zhen.GetTempVarValue(s.tempValueNo1)
		var v2 ZhenValue
		var v ZhenValue
		v, err = ZhenValueOperation(st, v1, v2)
		if err != nil {
			return
		}
		zhen.SetTempVarValue(s.tempValueNo1, v)

	case ZCS_TTAs:
		v := zhen.GetTempVarValue(s.tempValueNo1)
		zhen.SetTempVarValue(s.tempValueNo2, v)
	case ZCS_TTAdd, ZCS_TTSub, ZCS_TTMul, ZCS_TTDiv, ZCS_TTEq, ZCS_TTNe, ZCS_TTGt, ZCS_TTLt, ZCS_TTAnd, ZCS_TTOr:
		v1 := zhen.GetTempVarValue(s.tempValueNo1)
		v2 := zhen.GetTempVarValue(s.tempValueNo2)
		var v ZhenValue

		v, err = ZhenValueOperation(st, v1, v2)
		if err != nil {
			return
		}
		zhen.SetTempVarValue(s.tempValueNo1, v)

	case ZCS_If:
		//todo 条件判断指令

	case ZCS_For:
		//todo 次数循环指令
	case ZCS_While:
		//todo 条件循环指令
	case ZCS_Break:
		//todo 跳出循环指令
	case ZCS_Return:
		//todo 返回指令
	case ZCS_Call:
		//todo 运行函数指令
	case ZCS_PrintVar:
		v := zhen.GetGlobalVarValue(s.valueName1)
		fmt.Printf("变量：%s，值为：%s\n", s.valueName1, ZhenValueToString(v))

	default:
		fmt.Println(s.codeStepType)
		err = errors.New("未知指令")

	}
	return
}

func (zhen *ZhenState) Run() (err error) {
	zhen.runCodeStep = 0
	var codeCountStep = len(zhen.allCodeStepPointers)
	if codeCountStep == 0 {
		return
	}
	for {

		s := zhen.allCodeStepPointers[zhen.runCodeStep]
		err = zhen.codeStepRun(s)
		if err != nil {
			return err
		}

		zhen.runCodeStep += 1
		if zhen.runCodeStep >= codeCountStep {
			break
		}
	}
	return

}

func (zhen *ZhenState) GetKeyWord(identifier string) (isKeyWord bool, keyWord ZhenValue, err error) {
	value := zhen.GetGlobalVarValue("@关键字")
	if value.valueType == ZhenValueTypeTable {

		keyWord, isKeyWord = value.valueTable[identifier]

	}

	return
}

func (zhen *ZhenState) AddTextKeyWord(keyWord string) (err error) {
	value := zhen.GetGlobalVarValue("@文本型关键字")
	if value.valueType != ZhenValueTypeTable {
		t := make(map[string]ZhenValue)
		value = NewZhenValueTable(t)
	}
	value.valueTable[keyWord] = NewZhenValueFunction(KeyWordTextPreFun)

	zhen.SetGlobalVarValue("@关键字", value)
	return
}

func (zhen *ZhenState) AddKeyWord(keyWord KeyWord) (err error) {
	globalVarName := "@关键字"

	value := zhen.GetGlobalVarValue(globalVarName)
	if value.valueType != ZhenValueTypeTable {
		t := make(map[string]ZhenValue)
		value = NewZhenValueTable(t)
	}

	value.valueTable[keyWord.Name] = KeyWordToZhenValue(keyWord)

	zhen.SetGlobalVarValue(globalVarName, value)
	return
}
func (zhen *ZhenState) LoadBaseCodePre() (err error) {

	zhen.AddKeyWord(KeyWord{Name: "程序名", Type: KwtText, PreFun: KeyWordTextPreFun})
	zhen.AddKeyWord(KeyWord{Name: "版本号", Type: KwtText, PreFun: KeyWordTextPreFun})
	zhen.AddKeyWord(KeyWord{Name: "常量", Type: KwtConstant, PreFun: KeyWordConstantPreFun})
	zhen.AddKeyWord(KeyWord{Name: "显示", Type: KwtFun, PreFun: KeyWordFunPreFun})

	//zhen.AddTextKeyWord("版本号")
	//zhen.AddKeyWordConstant("")

	//value := zhen.GetGlobalVarValue("@文本型关键字")
	//if value.valueType != ZhenValueTypeTable {
	//	t := make(map[string]ZhenValue)
	//	value = NewZhenValueTable(t)
	//}
	//value.valueTable["程序名"] = NewZhenValueFunction(KeyWordTextPreFun)
	////key := "程序名"
	////value := NewZhenValueFunction(ZhenStateSetGlobalValuesByChildCode)
	//
	//zhen.SetGlobalVarValue("@关键字", value)
	return
}

func (zhen *ZhenState) Run2() (err error) {
	zhen.runCodeStep = 0
	var codeCountStep = len(zhen.allCodeStepPointers)
	if codeCountStep == 0 {
		return
	}
	for {

		s := zhen.allCodeStepPointers[zhen.runCodeStep]
		err = zhen.codeStepRun(s)
		if err != nil {
			return err
		}

		zhen.runCodeStep += 1
		if zhen.runCodeStep >= codeCountStep {
			break
		}
	}
	return

}
