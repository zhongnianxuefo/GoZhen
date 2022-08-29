package main

import (
	"fmt"
)

func (zhen *ZhenState) Run() {
	zhen.runCodeStep = 0
	var codeCountStep = len(zhen.allCodeStepPointers)
	if codeCountStep == 0 {
		return
	}
	for {
		s := zhen.allCodeStepPointers[zhen.runCodeStep]

		switch s.codeStepType {
		case ZCS_Var:
			zhen.Var(s)
		case ZCS_Assign:
			zhen.Assign(s)
		case ZCS_PrintVar:
			zhen.PrintVar(s)

		default:
			fmt.Println(s.codeStepType)
		}

		zhen.runCodeStep += 1
		if zhen.runCodeStep >= codeCountStep {
			break
		}
	}
	return

}

func main() {
	fmt.Println("Hello, World!")
	codes := ` 
<代码行><关键字>定义</关键字><变量名>变量A</变量名><代码 操作="定义变量" 变量名="变量A" /></代码行>
<代码行><变量名>变量A</变量名><关键字>等于</关键字><数字>32</数字><代码 操作="变量赋值" 变量名="变量A" 值类型="数字" 值="32" /></代码行>
<代码行><关键字>显示</关键字><变量名>变量A</变量名><代码 操作="显示" 变量名="变量A" /></代码行>
`
	z := NewZhenState()
	z.AddCodeString(codes)
	//var c ZhenCode
	//var cs ZhenCodeStep
	//cs.codeStepType = ZCS_Var
	//cs.valueName = "a"
	//c.codeSteps = append(c.codeSteps, cs)
	//var cs1 ZhenCodeStep
	//cs1.codeStepType = ZCS_Assign
	//cs1.valueName = "a"
	//cs1.value = NewZhenValueInt(100)
	//c.codeSteps = append(c.codeSteps, cs1)
	//c.needRun = true
	//z.AddCode(c)
	//
	//var c2 ZhenCode
	//var cs2 ZhenCodeStep
	//cs2.codeStepType = ZCS_PrintVar
	//cs2.valueName = "a"
	//c2.codeSteps = append(c2.codeSteps, cs2)
	//c2.needRun = true
	//z.AddCode(c2)

	z.Run()

	return
}
