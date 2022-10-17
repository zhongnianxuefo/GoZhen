package zhen_0_03

import (
	"fmt"
	"github.com/beevik/etree"
	"strconv"
)

type ToXml struct {
	code *Code
}

func NewCodeToXml(code *Code) (toXml *ToXml) {
	toXml = &ToXml{}
	toXml.code = code
	return
}

func (toXml *ToXml) tokenToXmlElement(token *Token, element *etree.Element) {
	showWords := false
	name := token.TokenType.String()
	words := string(token.Chars)
	switch token.TokenType {
	case TtCR, TtLF, TtCRLF:
		showWords = false
	default:
		showWords = true
	}
	if name != "" {
		e := element.CreateElement(name)
		if showWords {
			e.SetText(words)
		}
		e.CreateAttr("line", strconv.Itoa(token.LineNo))
		e.CreateAttr("col", strconv.Itoa(token.ColNo))
	}
	return
}

func (toXml *ToXml) nodeToXmlElement(node Node, element *etree.Element) {
	var e *etree.Element

	var items []Node
	switch n := node.(type) {
	case *RootNode:
		e = element.CreateElement("主程序")
		items = n.Items
	case *LineNode:
		e = element.CreateElement("代码行")
		e.SetText(n.String())
		if n.LineIndent > 0 {
			e.CreateAttr("Indent", strconv.Itoa(n.LineIndent))
		}
		items = n.Items
	case *ChildLineNode:
		e = element.CreateElement("子行")
		e.CreateAttr("分隔符", n.Symbol.String())
		if n.LineNo > 0 {
			e.CreateAttr("line", strconv.Itoa(n.LineNo))
		}
		if n.ColNo > 0 {
			e.CreateAttr("col", strconv.Itoa(n.ColNo))
		}
		items = n.Items
	case *ColonNode:
		e = element.CreateElement("冒号")
		items = n.Items
	case *OperatorNode:
		e = element.CreateElement(n.OperatorType.String())
		e.CreateAttr("line", strconv.Itoa(n.LineNo))
		e.CreateAttr("col", strconv.Itoa(n.ColNo))
		items = n.Items
	case *IntNode:
		e = element.CreateElement("整数")
		e.SetText(strconv.Itoa(n.Value))
	case *FloatNode:
		e = element.CreateElement("小数")
		e.SetText(strconv.FormatFloat(n.Value, 'E', -1, 64))
		e.SetText(fmt.Sprintf("%f", n.Value))
	case *StringNode:
		e = element.CreateElement("字符串")
		e.SetText(n.Value)
	case *LetterNode:
		e = element.CreateElement("标识符")
		e.CreateAttr("类型", n.LetterType.String())
		e.SetText(n.Words)
	case *EmptyNode:
		e = element.CreateElement("空格")
		e.SetText(n.Words)
	case *BracketNode:
		e = element.CreateElement("括号")
		items = n.Items
	case *OtherNode:
		e = element.CreateElement("内容")
		e.SetText(n.String())
	}
	for _, item := range items {
		toXml.nodeToXmlElement(item, e)
	}
	return
}
func (toXml *ToXml) ToXmlFile(path string) (err error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	element := doc.CreateElement("程序")

	toXml.nodeToXmlElement(toXml.code.RootNode, element)

	tokenElement := element.CreateElement("程序代码")
	for _, t := range toXml.code.Tokens {
		toXml.tokenToXmlElement(t, tokenElement)
	}
	commentsElement := element.CreateElement("注释")
	for _, c := range toXml.code.Comments {
		toXml.tokenToXmlElement(c, commentsElement)
	}

	doc.Indent(4)
	err = doc.WriteToFile(path)
	if err != nil {
		return
	}
	return

}
