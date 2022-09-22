package main

import (
	zhen_0_01 "GoZhen/zhen.0.01"
	zhen_0_02 "GoZhen/zhen.0.02"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func test(txtCodeFile string, jsonFile string, xmlFile string, formatFile string) (err error) {

	z := zhen_0_01.NewZhenState()
	//z.debug = false
	body, err := z.LoadTxtFile2(txtCodeFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	start := time.Now()
	err = z.LoadTxt(body)
	if err != nil {
		fmt.Println(err)
		return
	}

	cost := time.Since(start)
	fmt.Println("LoadTxtFile:", cost)

	start = time.Now()
	//err = z.LoadBaseCodePre()
	if err != nil {

		return
	}
	codePre := zhen_0_01.NewCodePre(z.MainCodeBlock)
	codePre.LoadBaseCodePre()
	//err = codePre.Preprocess()

	if err != nil {

		return
	}

	cost = time.Since(start)
	fmt.Println("Preprocess:", cost)

	start = time.Now()
	err = codePre.FileCodeBlock.ToJsonFile(jsonFile)
	if err != nil {

		return
	}
	cost = time.Since(start)
	fmt.Println("ToJsonFile:", cost)

	start = time.Now()

	toXml := zhen_0_01.NewCodeBlockToXml(z.MainCodeBlock)
	err = toXml.ToXmlFile(xmlFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	cost = time.Since(start)
	fmt.Println("ToXml:", cost)

	start = time.Now()
	format := zhen_0_01.NewCodeBlockFormat(z.MainCodeBlock)
	err = format.FormatToFile(formatFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	cost = time.Since(start)
	fmt.Println("Format:", cost)
	return
}
func test2(fileName string) (err error) {
	txtCodeFile := fmt.Sprintf("test/%s.z1", fileName)
	body, err := ioutil.ReadFile(txtCodeFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	start := time.Now()
	codeFile := zhen_0_02.NewCodeFile(string(body))
	err = codeFile.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	cost := time.Since(start)
	fmt.Println("Parse:", cost)

	start = time.Now()
	codePP := zhen_0_02.NewCodePP(codeFile)
	err = codePP.Preprocess()
	if err != nil {
		fmt.Println(err)
		return
	}
	cost = time.Since(start)
	fmt.Println("Preprocess:", cost)

	formatFile := fmt.Sprintf("test/%s-格式化.z1", fileName)
	start = time.Now()
	codeFile.FormatToFile(formatFile)
	cost = time.Since(start)
	fmt.Println("formatFile:", cost)

	xmlFile := fmt.Sprintf("test/%s.xml", fileName)
	start = time.Now()
	codeFile.ToXmlFile(xmlFile)
	cost = time.Since(start)
	fmt.Println("ToXmlFile:", cost)

	gobFile := fmt.Sprintf("test/%s.gob", fileName)
	start = time.Now()
	codeFile.ToGobFile(gobFile)
	cost = time.Since(start)
	fmt.Println("ToGobFile:", cost)

	//start = time.Now()
	//jsonFile := fmt.Sprintf("test/%s.json", fileName)
	//codeFile.ToJsonFile(jsonFile)
	//cost = time.Since(start)
	//fmt.Println("ToJsonFile:", cost)

	//start = time.Now()
	//codeFile2, err := zhen_0_02.NewCodeFileFromJsonFile(jsonFile)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//cost = time.Since(start)
	//fmt.Println("NewCodeFileFromJsonFile:", cost)
	//
	//start = time.Now()
	//jsonFile2 := fmt.Sprintf("test/%s2.json", fileName)
	//codeFile2.ToJsonFile(jsonFile2)
	//cost = time.Since(start)
	//fmt.Println("ToJsonFile 2:", cost)

	start = time.Now()
	codeFile2, err := zhen_0_02.NewCodeFileFromGobFile(gobFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	cost = time.Since(start)
	fmt.Println("NewCodeFileFromGobFile:", cost)

	formatFile2 := fmt.Sprintf("test/%s-格式化2.z1", fileName)
	start = time.Now()
	codeFile2.FormatToFile(formatFile2)
	cost = time.Since(start)
	fmt.Println("formatFile2:", cost)

	xmlFile2 := fmt.Sprintf("test/%s2.xml", fileName)
	start = time.Now()
	codeFile2.ToXmlFile(xmlFile2)
	cost = time.Since(start)
	fmt.Println("ToXmlFile2:", cost)

	//start = time.Now()
	//jsonFile3 := fmt.Sprintf("test/%s3.json", fileName)
	//codeFile3.ToJsonFile(jsonFile3)
	//cost = time.Since(start)
	//fmt.Println("ToJsonFile 3:", cost)

	//start = time.Now()
	//body2, err := ioutil.ReadFile(formatFile)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//codeFile4 := zhen_0_02.NewCodeFile(string(body2))
	//err = codeFile4.Parse()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//cost = time.Since(start)
	//fmt.Println("Parse2:", cost)

	//formatFile2 := fmt.Sprintf("test/%s-格式化2.z1", fileName)
	//start = time.Now()
	//codeFile4.FormatToFile(formatFile2)
	//cost = time.Since(start)
	//fmt.Println("formatFile2:", cost)
	return
}

type DoOperator struct {
	Bracket int

	Priority int
	Value    string
	Operator string
	Value2   string
	NextDo   *DoOperator
	t1       int
	t2       int
	t3       int
	leftvar  bool
	Rightvar bool
}

func test3() {
	words := strings.Split("1 + 2 + 3 * 4 / 5 + 6 * 7 + 8", " ")
	var ops []DoOperator
	//blnFirst := true
	//var NowOperator *DoOperator
	//t := 0
	for i := 0; i < len(words)-1; i += 2 {
		Value := words[i]
		Operator := words[i+1]
		value2 := words[i+2]
		Bracket := 0
		Priority := 0
		switch Operator {
		case "+", "-":
			Priority = 10
		case "*", "/":
			Priority = 20
		}
		op := &DoOperator{Bracket: Bracket, Operator: Operator, Priority: Priority, Value: Value, Value2: value2}
		//op.t1 = t + 1
		//op.t2 = t + 2
		ops = append(ops, *op)
		//t = op.t1
	}
	var opss []DoOperator
	var addOp func(nowNo int)
	//tNo := make(map[int]int)
	tt := 1
	lastTT := 0
	leftvar := false
	Rightvar := false
	appendOp := func(nowOp DoOperator) {
		//t1 := nowOp.t1
		//t2 := nowOp.t2
		//t, ok := tNo[t1]
		//if ok {
		//	nowOp.t1 = t
		//}
		//t, ok = tNo[t2]
		//if ok {
		//	nowOp.t2 = t
		//}
		if tt > lastTT {
			nowOp.t1 = tt
			nowOp.t2 = tt + 1
			nowOp.t3 = tt
		} else {
			nowOp.t1 = tt
			nowOp.t2 = tt + 1
			nowOp.t3 = tt
		}

		nowOp.leftvar = leftvar

		nowOp.Rightvar = Rightvar

		lastTT = tt
		opss = append(opss, nowOp)
		//tNo[t2] = nowOp.t1
	}

	addOp = func(nowNo int) {
		nowOp := ops[nowNo]
		if nowNo < len(ops)-1 {

			nextOp := ops[nowNo+1]
			if nextOp.Priority > nowOp.Priority {
				tt += 1
				leftvar = false
				Rightvar = false
				addOp(nowNo + 1)

				tt -= 1
				//nowOp.t1 = t + 1
				//nowOp.t2 = t + 2
				//opss = append(opss, nowOp)
				Rightvar = true
				appendOp(nowOp)
				leftvar = true
				//t = nowOp.t1
			} else {
				//nowOp.t1 = t + 1
				//nowOp.t2 = t + 2
				Rightvar = false
				appendOp(nowOp)
				leftvar = true
				//t = nowOp.t1
				addOp(nowNo + 1)
			}
		} else {
			//nowOp.t1 = t + 1
			//nowOp.t2 = t + 2
			appendOp(nowOp)
			//t = nowOp.t1
		}

	}
	addOp(0)
	for n, op := range ops {
		fmt.Println("a", n, op)
	}
	fmt.Println(words)
	for n, op := range opss {
		fmt.Println("b", n, op)
	}

}
func main() {
	fmt.Println("Hello, World!")

	err2 := test2("演示代码1")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
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
	_ = codes
	//start := time.Now()
	err1 := test("test/演示代码1.z1", "test/格式化演示代码1.json", "test/格式化演示代码1.xml", "test/格式化演示代码1.z1")
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	//cost := time.Since(start)
	//fmt.Println("cost:", cost)
	//
	//err2 := test("test/格式化演示代码1.z1", "test/格式化演示代码2.json", "test/格式化演示代码2.xml", "test/格式化演示代码2.z1")
	//if err2 != nil {
	//	fmt.Println(err2)
	//	return
	//}
	//
	////err = z.LoadTxtFile("zhen/格式化演示代码1.z1")
	////if err != nil {
	////	fmt.Println(err)
	////}
	////z.txtCode.formatToFile("zhen/格式化演示代码2.z1")
	////z.txtCode.ToXmlFile("zhen/格式化演示代码2.xml")
	//z := zhen.NewZhenState()
	//////z.debug = false
	//err := z.LoadString(codes)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("代码分析完成!")
	//
	//start = time.Now()
	//for i := 1; i <= 1; i++ {
	//	//err = z.Run()
	//}
	//
	//cost = time.Since(start)
	//if err != nil {
	//	//fmt.Println(z.getRunInfo())
	//	fmt.Println(err)
	//}
	//fmt.Println("cost:", cost)
	//start = time.Now()
	//for i := 1; i <= 1000000; i++ {
	//	var a = 100
	//	var b = 64
	//	c := a + b
	//	_ = c
	//	//fmt.Println(c)
	//}
	//
	//cost = time.Since(start)
	//fmt.Println("cost:", cost)
	//test3()
	return
}
