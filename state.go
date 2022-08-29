package main

import "fmt"

type ZhenState struct {
	globalValues ZhenValueTable
	localValues  ZhenValueTable
	stackValues  ZhenValueArray

	allCodes []ZhenCode

	allCodeStepPointers []*ZhenCodeStep
	runCodeStep         int
}

func NewZhenState() (zhen ZhenState) {
	zhen.globalValues = make(map[string]ZhenValue)
	zhen.localValues = make(map[string]ZhenValue)

	return

}
func (zhen *ZhenState) AddCodeString(code string) {

	return

}

func (zhen *ZhenState) AddCode(code ZhenCode) {
	zhen.allCodes = append(zhen.allCodes, code)
	if code.needRun {
		for i, _ := range code.codeSteps {
			//fmt.Println(i, s, &code.codeSteps[i])
			zhen.allCodeStepPointers = append(zhen.allCodeStepPointers, &code.codeSteps[i])
		}
	}

}
func (zhen *ZhenState) Var(s *ZhenCodeStep) {
	//todo 局部变量如何处理
	zhen.globalValues[s.valueName] = NewZhenValueNil()

	fmt.Printf("定义变量：%s，初始数值为：%s\n", s.valueName, ZhenValueToString(s.value))
}

func (zhen *ZhenState) Assign(s *ZhenCodeStep) {
	//todo 局部变量如何处理
	zhen.globalValues[s.valueName] = s.value

	fmt.Printf("变量：%s，赋值为：%s\n", s.valueName, ZhenValueToString(s.value))
}
func (zhen *ZhenState) PrintVar(s *ZhenCodeStep) {
	//todo 局部变量如何处理
	v := zhen.globalValues[s.valueName]

	fmt.Printf("变量：%s，值为：%s\n", s.valueName, ZhenValueToString(v))
}
