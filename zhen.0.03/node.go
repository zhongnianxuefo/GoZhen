package zhen_0_03

import (
	"strconv"
	"strings"
)

type NodeType int

const (
	NtUnknown NodeType = iota
	NtRoot
	NtLine
	NtChildLine
	NtColon
	NtOperator
	NtInt
	NtFloat
	NtString
	NtLetter
	NtEmpty
	NtBracket
	NtOther
)

type Node interface {
	addItem(Node)
	String() string
}

type TokenNode struct {
	LineNo int
	ColNo  int
	Words  string
}

type BaseNode struct {
	TokenNode
	Items     []Node
	NeedItems int
	NeedSpace bool
}

type RootNode struct {
	BaseNode
	LetterSet *LetterSet
}

type LineNode struct {
	BaseNode
	LineIndent int
	LetterSet  *LetterSet
	//KeyWordType KeyWordType
}

type ChildLineNode struct {
	BaseNode
	Symbol TokenType
}

type ColonNode struct {
	BaseNode
}

type OperatorNode struct {
	BaseNode
	OperatorType  OperatorType
	LeftPriority  int
	RightPriority int
}

type IntNode struct {
	BaseNode
	Value int
}

type FloatNode struct {
	BaseNode
	Value float64
}

type StringNode struct {
	BaseNode
	Value string
}

type LetterNode struct {
	BaseNode
	LetterType LetterType
	//Value string
}

type EmptyNode struct {
	BaseNode
	//Value string
}

type BracketNode struct {
	BaseNode
	LeftBracket  TokenType
	RightBracket TokenType
	EndToken     TokenNode
}

type OtherNode struct {
	BaseNode
}

func isValueNode(node Node) (ok bool) {
	switch node.(type) {
	case *LetterNode, *StringNode, *IntNode, *FloatNode, *BracketNode, *OtherNode:
		ok = true
	}
	return
}
func getNodeType(node Node) (nodeType NodeType) {
	switch node.(type) {
	case *RootNode:
		nodeType = NtRoot
	case *LineNode:
		nodeType = NtLine
	case *ChildLineNode:
		nodeType = NtChildLine
	case *ColonNode:
		nodeType = NtColon
	case *OperatorNode:
		nodeType = NtOperator
	case *IntNode:
		nodeType = NtInt
	case *FloatNode:
		nodeType = NtFloat
	case *StringNode:
		nodeType = NtString
	case *LetterNode:
		nodeType = NtLetter
	case *EmptyNode:
		nodeType = NtEmpty
	case *BracketNode:
		nodeType = NtBracket
	case *OtherNode:
		nodeType = NtOther
	default:
		nodeType = NtUnknown
	}
	return
}

func NewRootNode() (root *RootNode) {
	root = &RootNode{}
	root.LetterSet = NewLetterSet()
	root.LoadKeyWords()
	return
}

func (root *RootNode) LoadKeyWords() {
	root.LetterSet.addKeyWord("定义", KwtDefine)
	root.LetterSet.addKeyWord("常量", KwtDefineConstant)
	root.LetterSet.addKeyWord("定义常量", KwtDefineConstant)
	root.LetterSet.addKeyWord("变量", KwtDefineVar)
	root.LetterSet.addKeyWord("定义变量", KwtDefineVar)
	root.LetterSet.addKeyWord("全局变量", KwtDefineGlobalVar)
	root.LetterSet.addKeyWord("定义全局变量", KwtDefineGlobalVar)
	root.LetterSet.addKeyWord("局部变量", KwtDefineLocalVar)
	root.LetterSet.addKeyWord("定义局部变量", KwtDefineLocalVar)

	root.LetterSet.addKeyWord("如果", KwtIf)
	root.LetterSet.addKeyWord("否则", KwtElse)
	root.LetterSet.addKeyWord("否则如果", KwtElseIf)

	root.LetterSet.addKeyWord("循环", KwtWhile)
	root.LetterSet.addKeyWord("按条件循环", KwtWhile)
	root.LetterSet.addKeyWord("按次数循环", KwtFor)

	root.LetterSet.addKeyWord("定义函数", KwtDefineFun)
	root.LetterSet.addKeyWord("函数", KwtDefineFun)
	root.LetterSet.addKeyWord("参数", KwtDefineFunPara)
	root.LetterSet.addKeyWord("返回", KwtDefineFunReturn)

	root.LetterSet.addKeyWord("运行", KwtCallFun)
	root.LetterSet.addKeyWord("显示", KwtCallFun)

	return
}

func NewLineNode() (line *LineNode) {
	line = &LineNode{}
	line.LetterSet = NewLetterSet()
	return
}

func NewChildLineNode(token *Token) (childLine *ChildLineNode) {
	childLine = &ChildLineNode{}
	childLine.setBaseInfo(token)
	childLine.Symbol = token.TokenType
	childLine.Words = getPunctuationWords(childLine.Symbol)

	return
}

func (node *ChildLineNode) Copy() (newNode *ChildLineNode) {
	newNode = &ChildLineNode{}
	//newNode.LineNo = node.LineNo
	//newNode.ColNo = node.ColNo
	newNode.Symbol = node.Symbol
	newNode.Words = getPunctuationWords(newNode.Symbol)
	return
}

func NewColonNode(token *Token) (colon *ColonNode) {
	colon = &ColonNode{}
	colon.setBaseInfo(token)
	colon.Words = "："
	return
}

func NewOperatorNode(token *Token) (operator *OperatorNode) {
	operator = &OperatorNode{}
	operator.setBaseInfo(token)
	operator.OperatorType = getOperatorType(token.TokenType)
	operator.Words = operator.OperatorType.Words()
	operator.LeftPriority = getOperatorLeftPriority(operator.OperatorType)
	operator.RightPriority = getOperatorRightPriority(operator.OperatorType)
	operator.NeedItems = getOperatorNeedItems(operator.OperatorType)

	return

}

func ChangeOperatorNodeType(operator *OperatorNode, newType OperatorType) (newOperator *OperatorNode) {
	newOperator = operator
	newOperator.OperatorType = newType
	newOperator.LeftPriority = getOperatorLeftPriority(newType)
	newOperator.RightPriority = getOperatorRightPriority(newType)
	newOperator.NeedItems = getOperatorNeedItems(newType)
	return
}

func NewIntNode(token *Token) (node *IntNode) {
	node = &IntNode{}
	node.setBaseInfo(token)
	i, err := strconv.Atoi(string(token.Chars))
	if err != nil {
		//todo  转换失败警告
	}
	node.Value = i
	//node.Words = string(token.Chars)
	return
}
func NewFloatNode(token *Token) (node *FloatNode) {
	node = &FloatNode{}
	node.setBaseInfo(token)
	f, err := strconv.ParseFloat(string(token.Chars), 64)
	if err != nil {
		//todo  转换失败警告
	}
	node.Value = f
	//node.Words = string(token.Chars)
	return
}

func NewStringNode(token *Token) (node *StringNode) {
	node = &StringNode{}
	node.setBaseInfo(token)
	s := token.Chars
	var leftQuotation rune
	var rightQuotation rune
	if len(s) > 0 {
		leftQuotation = s[0]
		switch leftQuotation {
		case '\'', '‘', '’', '"', '“', '”':
			s = s[1:]
			if len(s) > 0 {
				rightQuotation = s[len(s)-1]
				switch rightQuotation {
				case '\'', '‘', '’', '"', '“', '”':
					s = s[:len(s)-1]
				}
			}
		}
		if leftQuotation != rightQuotation {
			//todo  警告引号不匹配
		}
	}
	node.Value = string(s)
	//node.Words = string(token.Chars)
	return
}

func NewLetterNode(token *Token) (node *LetterNode) {
	node = &LetterNode{}
	node.setBaseInfo(token)
	//node.Value = string(token.Chars)
	//node.Words = string(token.Chars)
	return
}

func NewEmptyNode(token *Token) (node *EmptyNode) {
	node = &EmptyNode{}
	node.setBaseInfo(token)
	//node.Value = string(token.Chars)
	//node.Words = string(token.Chars)
	return
}

func NewOtherNode(token *Token) (node *OtherNode) {
	node = &OtherNode{}
	node.setBaseInfo(token)
	return
}

func NewBracketNode(token *Token) (node *BracketNode) {
	node = &BracketNode{}
	node.setBaseInfo(token)
	node.LeftBracket = token.TokenType
	switch node.LeftBracket {
	case TtLeftBracket:
		node.RightBracket = TtRightBracket
	case TtLeftSquareBracket:
		node.RightBracket = TtRightSquareBracket
	case TtLeftBigBracket:
		node.RightBracket = TtRightBigBracket
	default:
		panic("未知括号")
	}
	node.Words, node.EndToken.Words = getBracketSymbolWords(token.TokenType)
	return
}
func (node *BracketNode) SetEndToken(t *Token) {
	node.EndToken.LineNo = t.LineNo
	node.EndToken.ColNo = t.ColNo
}

func (node *BaseNode) setBaseInfo(token *Token) {
	node.LineNo = token.LineNo
	node.ColNo = token.ColNo
	node.Words = string(token.Chars)
	return
}

func (node *BaseNode) addItem(item Node) {
	if item != nil {
		node.Items = append(node.Items, item)
	}
	return
}

func (node *BaseNode) needItem() (need bool) {
	need = true
	if node.NeedItems == 0 {
		need = false
	} else if len(node.Items) >= node.NeedItems {
		need = false
	}
	return
}

func (node *BaseNode) String() (s string) {
	var words []string
	words = append(words, node.Words)
	if node.NeedSpace {
		words = append(words, " ")
	}
	for _, items := range node.Items {
		words = append(words, items.String())
	}
	s = strings.Join(words, "")
	return
}

func (node *ChildLineNode) String() (s string) {
	var words []string
	hasChild := false
	for _, items := range node.Items {
		_, ok := items.(*ChildLineNode)
		if ok {
			words = append(words, items.String())
			hasChild = true
		} else {
			words = append(words, items.String())
		}
	}

	s = strings.Join(words, "")
	if hasChild {
		ss := []rune(s)
		ss = ss[0 : len(ss)-1]
		s = string(ss) + node.Words
	} else {
		s = s + node.Words
	}

	return
}

func (node *OperatorNode) String() (s string) {
	var words []string
	for _, item := range node.Items {
		leftWord := " " + node.Words

		//nodeType :=getNodeType(item)

		//words = append(words, node.Words+"")
		words = append(words, leftWord)
		words = append(words, item.String())
	}
	if node.NeedItems > 1 && len(words) > 0 {
		words = words[1:]
	}
	s = strings.Join(words, "")
	return
}

func (node *BracketNode) String() (s string) {
	var words []string

	leftBracket := ""
	rightBracket := ""
	switch node.LeftBracket {
	case TtLeftBracket:
		leftBracket = "（"
		rightBracket = "）"
	case TtLeftSquareBracket:
		leftBracket = "["
		rightBracket = "]"
	case TtLeftBigBracket:
		leftBracket = "{"
		rightBracket = "}"

	}
	words = append(words, leftBracket)
	for _, items := range node.Items {
		words = append(words, items.String())
	}
	words = append(words, rightBracket)
	s = strings.Join(words, "")
	return
}

func (node *BaseNode) addChildLineNode(items []Node) {
	hasChildLineNode := false
	if len(node.Items) > 0 {
		lastItem := node.Items[len(node.Items)-1]
		childLineNode, ok := lastItem.(*ChildLineNode)
		if ok && len(items) > 0 {
			var c *ChildLineNode
			c = childLineNode.Copy()
			c.Items = append(c.Items, items...)
			node.addItem(c)
			hasChildLineNode = true
		}
	}
	if hasChildLineNode == false {
		node.Items = append(node.Items, items...)
	}
	node.arrangeAllChildLine()
	return
}

func getSymbolGrade(symbol TokenType) (grade int) {
	switch symbol {
	case TtComma:
		grade = 20
	case TtDunHao:
		grade = 10
	case TtSemicolon:
		grade = 30
	case TtPeriod:
		grade = 40
	default:
		grade = 0
	}
	return
}

func (node *BaseNode) arrangeAllChildLine() {
	//return
	for true {
		if node.arrangeChildLine() == false {
			break
		}
	}
	return
}

func (node *BaseNode) arrangeChildLine() (isChange bool) {
	var maxGrad int
	var maxGradeNode *ChildLineNode

	for _, item := range node.Items {
		grade := 0
		c, ok := item.(*ChildLineNode)
		if ok {
			grade = getSymbolGrade(c.Symbol)
			if grade > maxGrad {
				maxGrad = grade
				maxGradeNode = c.Copy()
			}
		}
	}
	if maxGrad == 0 {
		return
	}
	var newItems []Node
	for _, item := range node.Items {
		grade := 0
		c, ok := item.(*ChildLineNode)
		if ok {
			grade = getSymbolGrade(c.Symbol)
		}
		if grade < maxGrad {
			maxGradeNode.addItem(item)
			isChange = true
		} else if grade == maxGrad {
			maxGradeNode.Items = append(maxGradeNode.Items, c.Items...)
			maxGradeNode.LineNo = c.LineNo
			maxGradeNode.ColNo = c.ColNo
			newItems = append(newItems, maxGradeNode)
			maxGradeNode = c.Copy()
		}
	}
	if len(maxGradeNode.Items) > 0 {
		newItems = append(newItems, maxGradeNode)
	}

	if isChange {
		node.Items = newItems
	}
	for _, item := range node.Items {
		c, ok := item.(*ChildLineNode)
		if ok {
			if c.arrangeChildLine() {
				isChange = true
			}
		}
	}
	return

}

//func isOperatorNode(node Node) (ok bool) {
//	switch node.(type) {
//	case *OperatorNode:
//		ok = true
//	}
//	return
//}
