package main

import (
	"fmt"
	"time"
)

func test(txtCodeFile string, xmlFile string, formatFile string) (err error) {
	z := NewZhenState()
	z.debug = false
	err = z.LoadTxtFile(txtCodeFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	toXml := NewCodeBlockToXml(z.MainCodeBlock)
	err = toXml.ToXmlFile(xmlFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	format := NewCodeBlockFormat(z.MainCodeBlock)
	err = format.formatToFile(formatFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func main() {
	fmt.Println("Hello, World!")
	codes := ` 
<?xml version="1.0" encoding="UTF-8"?>
<程序>
<代码行><关键字>定义变量</关键字><变量名>A</变量名><代码 指令="变量定义" 变量名="A" /></代码行>
<代码行><关键字>变量</关键字><变量名>A</变量名><关键字>等于</关键字><数字>32</数字>
<代码 指令="临时变量定义" 临时变量1="1" 值类型="数字" 值="35" />
<代码 指令="临时变量赋值给普通变量" 临时变量1="1" 变量名="A"  />
</代码行>
<代码行><关键字>定义变量</关键字><变量名>B</变量名><代码 指令="变量定义" 变量名="B" /></代码行>
<代码行><关键字>变量</关键字><变量名>B</变量名><关键字>等于</关键字><数字>20</数字>
<代码 指令="临时变量定义" 临时变量1="1" 值类型="数字" 值="1023.7" />
<代码 指令="临时变量赋值给普通变量" 临时变量1="1" 变量名="B"  />
</代码行>
<代码行><关键字>定义变量</关键字><变量名>C</变量名><代码 指令="变量定义" 变量名="C" /></代码行>
<代码行><关键字>变量</关键字><变量名>C</变量名><关键字>等于</关键字><关键字>变量</关键字><变量名>A</变量名><关键字>乘</关键字><关键字>变量</关键字><变量名>B</变量名>
	<代码 指令="变量相乘" 临时变量1="1" 变量名1="A" 变量名2="B"  />
	<代码 指令="临时变量赋值给普通变量" 临时变量1="1" 变量名="C"  />
</代码行>
<代码行><关键字>显示变量</关键字><变量名>C</变量名><代码 指令="显示变量" 变量名="C" /></代码行>
</程序>
`

	err1 := test("Zhen/演示代码1.z1", "Zhen/格式化演示代码1.xml", "Zhen/格式化演示代码1.z1")
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	err2 := test("Zhen/格式化演示代码1.z1", "Zhen/格式化演示代码2.xml", "Zhen/格式化演示代码2.z1")
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	//err = z.LoadTxtFile("Zhen/格式化演示代码1.z1")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//z.txtCode.formatToFile("Zhen/格式化演示代码2.z1")
	//z.txtCode.ToXmlFile("Zhen/格式化演示代码2.xml")
	z := NewZhenState()
	z.debug = false
	err := z.LoadString(codes)
	//if err != nil {
	//	fmt.Println(err)
	//}
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

	return
}
