package zhen_0_01

import (
	"fmt"
	"github.com/beevik/etree"
	"strconv"
)

type CodeBlockToXml struct {
	MainCodeBlock *CodeBlock2
}

func NewCodeBlockToXml(codeBlock *CodeBlock2) (toXml CodeBlockToXml) {
	toXml.MainCodeBlock = codeBlock
	return
}

func (toXml *CodeBlockToXml) ToXmlElement(codeBlock *CodeBlock2, element *etree.Element) {
	showWords := false
	name := ""
	words := codeBlock.getChars()
	switch codeBlock.WordType {
	case CwtUnSet:
		switch codeBlock.Operator.Type {
		case OtUnSet:
			switch codeBlock.BlockType {
			case CbtFile:
				name = "程序"
			case CbtLine:
				name = fmt.Sprintf("代码行-%d", codeBlock.Pos.LineNo)
			case CbtOperator, CbtLetter, CbtNumber, CbtString, CbtComment:
				name = codeBlock.BlockType.String()
				showWords = true
			default:
				name = codeBlock.BlockType.String()
				showWords = false
			}
		default:
			name = codeBlock.Operator.Type.String()
			showWords = false
		}

	default:
		name = codeBlock.WordType.String()
		words = codeBlock.Word
		showWords = true
	}
	if name != "" {
		e := element.CreateElement(name)
		if showWords {
			e.SetText(words)
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
		//for _, c := range codeBlock.cod {
		//	toXml.codeBlockToXmlElement(c, e)
		//}

		for _, c := range codeBlock.items {
			toXml.ToXmlElement(c, e)
		}
	}

	return
}
func (toXml *CodeBlockToXml) ToXmlFile(path string) (err error) {

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	element := doc.CreateElement("程序")
	for _, codeBlock := range toXml.MainCodeBlock.items {
		toXml.ToXmlElement(codeBlock, element)
	}

	doc.Indent(4)
	err = doc.WriteToFile(path)

	if err != nil {
		return
	}

	return

}
