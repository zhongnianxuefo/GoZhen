package zhen_0_01

import (
	"encoding/json"
	"io/ioutil"
)

type CodeBlockPos struct {
	StartNo  int
	BlockLen int

	LineNo    int
	LineCount int
	ColNo     int
}

type CodeBlockType int

const (
	CbtUnSet CodeBlockType = iota
	CbtApostrophe
	CbtLeftApostrophe
	CbtRightApostrophe
	CbtQuotation
	CbtLeftQuotation
	CbtRightQuotation
	CbtLeftBracket
	CbtRightBracket
	CbtLeftSquareBracket
	CbtRightSquareBracket
	CbtLeftBigBracket
	CbtRightBigBracket
	CbtSpace
	CbtFullWidthSpace
	CbtTab
	CbtPound
	CbtCR
	CbtLF
	CbtCRLF
	CbtColon
	CbtComma
	CbtDunHao
	CbtSemicolon
	CbtPeriod
	CBtBackslash
	CbtOperator
	//CbtUnderscore
	CbtLetter
	CbtNumber
	CbtPoint
	CbtOtherChar

	CbtFile
	CbtLine
	CbtChildLine
	CbtClass
	CbtFun

	CbtString
	CbtComment
)

var CodeBlockTypeNames = [...]string{
	CbtUnSet:              "未设置",
	CbtApostrophe:         "单引号",
	CbtLeftApostrophe:     "左单引号",
	CbtRightApostrophe:    "右单引号",
	CbtQuotation:          "双引号",
	CbtLeftQuotation:      "左双引号",
	CbtRightQuotation:     "左双引号",
	CbtLeftBracket:        "左括号",
	CbtRightBracket:       "右括号",
	CbtLeftSquareBracket:  "左中括号",
	CbtRightSquareBracket: "右中括号",
	CbtLeftBigBracket:     "左大括号",
	CbtRightBigBracket:    "右大括号",
	CbtSpace:              "空格",
	CbtFullWidthSpace:     "全角空格",
	CbtTab:                "Tab",
	CbtPound:              "井号",
	CbtCR:                 "回车",
	CbtLF:                 "换行",
	CbtCRLF:               "回车换行",
	CbtColon:              "冒号",
	CbtComma:              "逗号",
	CbtDunHao:             "顿号",
	CbtSemicolon:          "分号",
	CbtPeriod:             "句号",
	CBtBackslash:          "反斜杠",
	CbtOperator:           "运算符",
	CbtLetter:             "标识符",
	CbtNumber:             "数字",
	CbtPoint:              "点",
	CbtOtherChar:          "其他字符",
	CbtFile:               "文件",
	CbtLine:               "代码行",
	CbtChildLine:          "子代码行",
	CbtClass:              "类",
	CbtFun:                "函数",
	CbtString:             "字符串",
	CbtComment:            "注释",
}

func (cbt CodeBlockType) String() string {
	return CodeBlockTypeNames[cbt]
}

type CodeWordType int

const (
	CwtUnSet CodeWordType = iota
	CwtKeyWord
	CwtTxt

	CwtConstant
	CwtGlobalVar
	CwtLocalVar

	CwtFunName
	CwtFunPara
	CwtFunReturn
)

var CodeWordTypeNames = [...]string{
	CwtUnSet:     "未设置",
	CwtKeyWord:   "关键字",
	CwtTxt:       "文本",
	CwtConstant:  "常量",
	CwtGlobalVar: "全局变量",
	CwtLocalVar:  "局部变量",
	CwtFunName:   "函数名",
	CwtFunPara:   "参数",
	CwtFunReturn: "返回值",
}

func (cwt CodeWordType) String() string {
	return CodeWordTypeNames[cwt]
}

type CodeBlock2 struct {
	allCodeChars  []rune
	Pos           CodeBlockPos
	BlockType     CodeBlockType
	items         []*CodeBlock2
	parCodeBlock  *CodeBlock2
	nextCodeBlock *CodeBlock2

	LineIndent int
	Comment    string

	//isKeyWord bool

	Word         string
	WordType     CodeWordType
	codeSteps    []ZhenCodeStep
	nextStepCode *CodeBlock2

	keyWords   *VarGroup
	operators  *VarGroup
	constants  *VarGroup
	functions  *VarGroup
	globalVars *VarGroup
	localVars  *VarGroup

	Operator Operator
}

type CodeStep struct {
	codeStepType ZhenCodeStepType
	valueName1   string
	valueType1   CodeWordType

	valueName2 string
	valueType2 CodeWordType

	tempValueNo1 int
	tempValueNo2 int
	value        ZValue
}

func NewCodeBlock(codeChars []rune, pos CodeBlockPos, codeBlockType CodeBlockType) (codeBlock *CodeBlock2) {
	codeBlock = &CodeBlock2{}
	codeBlock.allCodeChars = codeChars
	codeBlock.Pos = pos
	codeBlock.BlockType = codeBlockType
	//codeBlock.lineIndent = 0

	codeBlock.keyWords = NewVarGroup()
	codeBlock.operators = NewVarGroup()
	codeBlock.constants = NewVarGroup()
	codeBlock.functions = NewVarGroup()
	codeBlock.globalVars = NewVarGroup()
	codeBlock.localVars = NewVarGroup()
	return
}

func (codeBlock *CodeBlock2) ToJsonFile(fileName string) (err error) {
	data, err := json.MarshalIndent(codeBlock, "", " ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		return
	}

	return
}

func (codeBlock *CodeBlock2) getChars() string {
	s := codeBlock.Pos.StartNo
	e := s + codeBlock.Pos.BlockLen
	return string(codeBlock.allCodeChars[s:e])
}

//func (codeBlock *CodeBlock2) addLen(addBlockLen int) {
//	codeBlock.Pos.blockLen += addBlockLen
//}

func (codeBlock *CodeBlock2) setEndPos(endItem *CodeBlock2) {
	endNo := endItem.Pos.StartNo + endItem.Pos.BlockLen
	if endNo > codeBlock.Pos.StartNo {
		codeBlock.Pos.BlockLen = endNo - codeBlock.Pos.StartNo
	}
	endLineNo := endItem.Pos.LineNo
	if endLineNo > codeBlock.Pos.LineNo {
		codeBlock.Pos.LineCount = endLineNo - codeBlock.Pos.LineNo + 1
	}
}

func (codeBlock *CodeBlock2) addItem(item *CodeBlock2) {
	codeBlock.items = append(codeBlock.items, item)
	item.parCodeBlock = codeBlock
}

func (codeBlock *CodeBlock2) appendNext(nextCodeBlock *CodeBlock2) *CodeBlock2 {
	codeBlock.parCodeBlock.addItem(nextCodeBlock)
	return nextCodeBlock
}

func (codeBlock *CodeBlock2) appendChild(child *CodeBlock2) *CodeBlock2 {
	codeBlock.addItem(child)
	return child
}

func (codeBlock *CodeBlock2) getNext() (next *CodeBlock2, isExist bool) {
	par := codeBlock.parCodeBlock
	parItemCount := len(par.items)
	for i, c := range par.items {
		if c == codeBlock {
			if i < parItemCount-1 {
				next = par.items[i+1]
				isExist = true
			}
		}
	}
	return
}

func (codeBlock *CodeBlock2) getCodeArea() (area *CodeBlock2) {
	par := codeBlock.parCodeBlock
	switch par.BlockType {
	case CbtFun, CbtClass, CbtFile:
		area = par
	default:
		area = par.getCodeArea()
	}
	return
}

func (codeBlock *CodeBlock2) getCodeFile() (area *CodeBlock2) {
	par := codeBlock.parCodeBlock
	switch par.BlockType {
	case CbtFile:
		area = par
	default:
		area = par.getCodeFile()
	}
	return
}
