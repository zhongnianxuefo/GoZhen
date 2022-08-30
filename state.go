package main

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
)

const ZhenTempValuesLen = 10

type ZhenState struct {
	debug        bool
	globalValues ZhenValueTable
	localValues  ZhenValueTable
	tempValues   [ZhenTempValuesLen]ZhenValue

	allCodes []ZhenCode

	allCodeStepPointers []ZhenCodeStep
	runCodeStep         int
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
func (zhen *ZhenState) LoadString(code string) (err error) {
	doc := etree.NewDocument()
	err = doc.ReadFromString(code)
	if err != nil {
		return
	}

	program := doc.SelectElement("程序")

	for _, line := range program.ChildElements() {
		var code ZhenCode
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

func (zhen *ZhenState) AddCode(code ZhenCode) {
	zhen.allCodes = append(zhen.allCodes, code)
	if code.needRun {
		for i := range code.codeSteps {
			//fmt.Println(i, s, &code.codeSteps[i])
			zhen.allCodeStepPointers = append(zhen.allCodeStepPointers, code.codeSteps[i])
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

func (zhen *ZhenState) SetTempVarValue(tempValueNo int, value ZhenValue) {

	zhen.tempValues[tempValueNo] = value
}

func (zhen *ZhenState) GetTempVarValue(tempValueNo int) (value ZhenValue) {

	value = zhen.tempValues[tempValueNo]
	return
}

func (zhen *ZhenState) Var(s *ZhenCodeStep) {
	//todo 局部变量如何处理
	zhen.globalValues[s.valueName1] = s.value
	if zhen.debug {
		fmt.Printf("定义变量：%s，初始数值为：%s\n", s.valueName1, ZhenValueToString(s.value))
	}

}

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

func (zhen *ZhenState) codeStepRun(s ZhenCodeStep) (err error) {

	switch s.codeStepType {
	case ZCS_None:

	case ZCS_Var:
		zhen.SetVarValue(s.valueName1, s.value)
	case ZCS_As:
		v := zhen.GetVarValue(s.valueName1)
		zhen.SetVarValue(s.valueName2, v)
	case ZCS_Add:
		v1 := zhen.GetVarValue(s.valueName1)
		v2 := zhen.GetVarValue(s.valueName2)
		var v ZhenValue
		v, err = ZhenValueAdd(v1, v2)
		if err != nil {
			return
		}
		zhen.SetTempVarValue(s.tempValueNo1, v)
	case ZCS_Sub:
	case ZCS_Mul:
	case ZCS_Div:
	case ZCS_Eq:
	case ZCS_Ne:
	case ZCS_Gt:
	case ZCS_Lt:
	case ZCS_And:
	case ZCS_Or:
	case ZCS_Not:

	case ZCS_TVar:
		zhen.SetTempVarValue(s.tempValueNo1, s.value)
	case ZCS_TFrom:
		v := zhen.GetVarValue(s.valueName1)
		zhen.SetTempVarValue(s.tempValueNo1, v)
	case ZCS_TAs:
		v := zhen.GetTempVarValue(s.tempValueNo1)
		zhen.SetVarValue(s.valueName1, v)
	case ZCS_TAdd:
		v1 := zhen.GetTempVarValue(s.tempValueNo1)
		v2 := zhen.GetVarValue(s.valueName1)
		var v ZhenValue
		v, err = ZhenValueAdd(v1, v2)
		if err != nil {
			return
		}
		zhen.SetTempVarValue(s.tempValueNo1, v)
	case ZCS_TSub:
	case ZCS_TMul:
	case ZCS_TDiv:
	case ZCS_TEq:
	case ZCS_TNe:
	case ZCS_TGt:
	case ZCS_TLt:
	case ZCS_TAnd:
	case ZCS_TOr:
	case ZCS_TNot:
	case ZCS_TTAs:
	case ZCS_TTAdd:
	case ZCS_TTSub:
	case ZCS_TTMul:
	case ZCS_TTDiv:
	case ZCS_TTEq:
	case ZCS_TTNe:
	case ZCS_TTGt:
	case ZCS_TTLt:
	case ZCS_TTAnd:
	case ZCS_TTOr:
	case ZCS_If:
	case ZCS_For:
	case ZCS_While:
	case ZCS_Break:
	case ZCS_Return:
	case ZCS_Call:
	case ZCS_PrintVar:
		v := zhen.GetVarValue(s.valueName1)
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
