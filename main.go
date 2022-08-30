package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, World!")
	codes := ` 
<?xml version="1.0" encoding="UTF-8"?>
<程序>
<代码行><关键字>定义变量</关键字><变量名>A</变量名><代码 指令="变量定义" 变量名="A" /></代码行>
<代码行><关键字>变量</关键字><变量名>A</变量名><关键字>等于</关键字><数字>32</数字>
<代码 指令="临时变量定义" 临时变量1="1" 值类型="整数" 值="35" />
<代码 指令="临时变量赋值给普通变量" 临时变量1="1" 变量名="A"  />
</代码行>
<代码行><关键字>定义变量</关键字><变量名>B</变量名><代码 指令="变量定义" 变量名="B" /></代码行>
<代码行><关键字>变量</关键字><变量名>B</变量名><关键字>等于</关键字><数字>100</数字>
<代码 指令="临时变量定义" 临时变量1="1" 值类型="小数" 值="100.5" />
<代码 指令="临时变量赋值给普通变量" 临时变量1="1" 变量名="B"  />
</代码行>
<代码行><关键字>定义变量</关键字><变量名>C</变量名><代码 指令="变量定义" 变量名="C" /></代码行>
<代码行><关键字>变量</关键字><变量名>C</变量名><关键字>等于</关键字><关键字>变量</关键字><变量名>A</变量名><关键字>加</关键字><关键字>变量</关键字><变量名>B</变量名>
	<代码 指令="变量相加" 临时变量1="1" 变量名1="A" 变量名2="B"  />
	<代码 指令="临时变量赋值给普通变量" 临时变量1="1" 变量名="C"  />
</代码行>
<代码行><关键字>显示变量</关键字><变量名>C</变量名><代码 指令="显示变量" 变量名="C" /></代码行>
</程序>
`
	z := NewZhenState()
	z.debug = false
	err := z.LoadString(codes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("代码分析完成!")

	start := time.Now()
	for i := 1; i <= 1; i++ {
		err = z.Run()
	}

	cost := time.Since(start)
	if err != nil {
		fmt.Println(z.getRunInfo())
		fmt.Println(err)
	}
	fmt.Println("cost:", cost)
	start = time.Now()
	for i := 1; i <= 1000000; i++ {
		var a = 100
		var b = 64
		c := a + b
		_ = c
		//fmt.Println(c)
	}

	cost = time.Since(start)
	fmt.Println("cost:", cost)

	//var c ZhenCode
	//var cs ZhenCodeStep
	//cs.codeStepType = ZCS_Var
	//cs.valueName = "a"
	//c.codeSteps = append(c.codeSteps, cs)
	//var cs1 ZhenCodeStep
	//cs1.codeStepType = ZCS_As
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

	//z.Run()

	return
}
