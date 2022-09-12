package main

import (
	"fmt"
	"github.com/beevik/etree"
	"strconv"
)

type CodeBlockToXml struct {
	MainCodeBlock *CodeBlock
}

func NewCodeBlockToXml(codeBlock *CodeBlock) (toXml CodeBlockToXml) {
	toXml.MainCodeBlock = codeBlock
	return
}

func (toXml *CodeBlockToXml) ToXmlElement(codeBlock *CodeBlock, element *etree.Element) {
	showWords := false
	name := ""
	switch codeBlock.BlockType {
	case CbtFile:
		name = "程序"
	case CbtLine:
		name = fmt.Sprintf("代码行-%d", codeBlock.Pos.LineNo)
	case CbtChildLine:
		name = "代码子行"
	case CbtCR:
		name = "回车"
	case CbtLF:
		name = "换行"
	case CbtCRLF:
		name = "回车换行"
	case CbtSpace:
		name = "空格"
	case CbtTab:
		name = "Tab"
	case CbtLeftBracket:
		name = "左括号"
	case CbtRightBracket:
		name = "右括号"
	case CbtLeftSquareBracket:
		name = "左中括号"
	case CbtRightSquareBracket:
		name = "右中括号"
	case CbtLeftBigBracket:
		name = "左大括号"
	case CbtRightBigBracket:
		name = "右大括号"
	case CbtColon:
		name = "冒号"
	case CbtComma:
		name = "逗号"
	case CbtDunHao:
		name = "顿号"
	case CbtSemicolon:
		name = "分号"
	case CbtPeriod:
		name = "句号"
	case CBtBackslash:
		name = "反斜杠"
	case CbtPoint:
		name = "点"
	case CbtOperator:
		name = "运算符"
		showWords = true
	case CbtLetter:
		name = "标识符"
		showWords = true
	case CbtNumber:
		name = "数字"
		showWords = true
	case CbtString:
		name = "字符串"
		showWords = true
	case CbtComment:
		name = "注释"
		showWords = true
	default:
		name = "未知" + strconv.Itoa(int(codeBlock.BlockType))
		showWords = true
	}
	if name != "" {
		e := element.CreateElement(name)
		if showWords {
			e.SetText(codeBlock.getChars())
		}
		if codeBlock.Comment != "" {
			e.CreateAttr("注释", codeBlock.Comment)
		}
		e.CreateAttr("line", strconv.Itoa(codeBlock.Pos.LineNo))
		e.CreateAttr("col", strconv.Itoa(codeBlock.Pos.ColNo))
		e.CreateAttr("len", strconv.Itoa(codeBlock.Pos.BlockLen))

		if codeBlock.LineIndent > 0 {
			e.CreateAttr("LineIndent", strconv.Itoa(codeBlock.LineIndent))
		}

		for _, c := range codeBlock.Items {
			toXml.ToXmlElement(c, e)
		}
	}

	return
}
func (toXml *CodeBlockToXml) ToXmlFile(path string) (err error) {

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	element := doc.CreateElement("程序")
	for _, codeBlock := range toXml.MainCodeBlock.Items {
		toXml.ToXmlElement(codeBlock, element)
	}

	doc.Indent(4)
	err = doc.WriteToFile(path)

	if err != nil {
		return
	}

	return

}
