package zhen_0_01

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"io/ioutil"
	"os"
	"strconv"
	"unicode"
)

type CodeBlock struct {
	Pos                 CodeBlockPos
	BlockType           CodeBlockType
	ParCodeBlock        int
	NextCodeBlock       int
	ChildCodeBlockFirst int
	ChildCodeBlockLast  int

	Words    string
	WordType CodeWordType
	Operator Operator

	LineIndent int
	Comment    string
}
type CodeFile struct {
	allCodeChars []rune
	NowPos       CodeBlockPos
	nowBlock     int
	nowLineBlock int

	LastChar rune

	lastBlockType    CodeBlockType
	lastBlockStartNo int
	lastBlockLen     int
	lastBlockLineNo  int

	NewLine      bool
	Nowrap       bool
	BracketCount int

	AllCodeBlock []CodeBlock
}

var CharCodeBlockType = map[rune]CodeBlockType{
	'\'': CbtApostrophe,
	'‘':  CbtLeftApostrophe,
	'’':  CbtRightApostrophe,
	'"':  CbtQuotation,
	'“':  CbtLeftQuotation,
	'”':  CbtRightQuotation,
	'#':  CbtPound,
	' ':  CbtSpace,
	'　':  CbtFullWidthSpace,
	'\r': CbtCR,
	'\n': CbtLF,
	'(':  CbtLeftBracket, '（': CbtLeftBracket,
	')': CbtRightBracket, '）': CbtRightBracket,
	'[':  CbtLeftSquareBracket,
	']':  CbtRightSquareBracket,
	'{':  CbtLeftBigBracket,
	'}':  CbtRightBigBracket,
	'\t': CbtTab,
	':':  CbtColon, '：': CbtColon,
	'、': CbtDunHao,
	',': CbtComma, '，': CbtComma,
	'。': CbtPeriod,
	';': CbtSemicolon, '；': CbtSemicolon,
	'\\': CBtBackslash,
	'=':  CbtOperator,
	'+':  CbtOperator,
	'-':  CbtOperator,
	'*':  CbtOperator,
	'/':  CbtOperator,
	'>':  CbtOperator,
	'<':  CbtOperator,
	'&':  CbtOperator,
	'|':  CbtOperator,
	'!':  CbtOperator,
	'^':  CbtOperator,
	'?':  CbtOperator,
	'.':  CbtPoint,
	'_':  CbtLetter,
	'@':  CbtLetter,
}

func NewCodeFile(codeTxt string) (codeFile *CodeFile) {
	codeFile = &CodeFile{}
	txt := codeTxt

	codeFile.allCodeChars = []rune(txt)
	codeFile.NowPos = CodeBlockPos{StartNo: 0, BlockLen: 0, LineNo: 1, LineCount: 1, ColNo: 1}

	codeFile.nowBlock = codeFile.newCodeBlock(codeFile.NowPos, CbtFile)

	//codeBlock.Pos = CodeBlockPos{startNo: 0, blockLen: 0, LineNo: 1, LineCount: 1, ColNo: 1}
	//codeBlock.BlockType = CbtFile
	//codeBlock.ParNo = 0
	//codeBlock.FirstChildNo = -1
	//codeBlock.LastChildNo = -1
	//
	//codeBlock.NextNo = -1

	codeFile.LastChar = 0
	codeFile.NewLine = true
	codeFile.Nowrap = false
	codeFile.BracketCount = 0

	//codeFile.AllCodeBlock = append(codeFile.AllCodeBlock, codeBlock)
	codeFile.nowBlock = 0
	return
}
func NewCodeFileFromJsonFile(jsonFile string) (codeFile *CodeFile, err error) {
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return
	}
	codeFile = &CodeFile{}
	err = json.Unmarshal(data, &codeFile)
	if err != nil {
		return
	}
	return
}
func (codeFile *CodeFile) ToJsonFile(jsonFile string) (err error) {
	data, err := json.MarshalIndent(codeFile, "", "\t")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(jsonFile, data, 0666)
	if err != nil {
		return
	}
	return
}
func NewCodeFileFromGobFile(gobFile string) (codeFile *CodeFile, err error) {

	file, err := os.Open(gobFile)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(file)
	codeFile = &CodeFile{}
	err = dec.Decode(&codeFile)
	if err != nil {
		return
	}
	return
	//data, err := ioutil.ReadFile(jsonFile)
	//if err != nil {
	//	return
	//}
	//codeFile = &CodeFile{}
	//err = json.Unmarshal(data, &codeFile)
	//if err != nil {
	//	return
	//}
	//return
}
func (codeFile *CodeFile) ToGobFile(gobFile string) (err error) {
	file, err := os.Create(gobFile)
	if err != nil {
		return
	}

	enc := gob.NewEncoder(file)
	err = enc.Encode(codeFile)
	if err != nil {
		return
	}

	return
}

func (codeFile *CodeFile) newCodeBlock(pos CodeBlockPos, charType CodeBlockType) (codeBlockNo int) {
	codeBlock := CodeBlock{}

	codeBlock.Pos = pos
	codeBlock.BlockType = charType

	codeBlock.ParCodeBlock = -1
	codeBlock.ChildCodeBlockFirst = -1
	codeBlock.ChildCodeBlockLast = -1
	codeBlock.NextCodeBlock = -1

	if codeBlock.Pos.BlockLen > 0 {
		codeBlock.Words = string(codeFile.allCodeChars[codeBlock.Pos.StartNo : codeBlock.Pos.StartNo+codeBlock.Pos.BlockLen])
	}

	codeBlockNo = len(codeFile.AllCodeBlock)
	codeFile.AllCodeBlock = append(codeFile.AllCodeBlock, codeBlock)

	return
}

func (codeFile *CodeFile) addNextCodeBlock(nowCodeBlockNo int, nextCodeBlockNo int) {
	nowCodeBlock := &codeFile.AllCodeBlock[nowCodeBlockNo]
	parCodeBlockNo := nowCodeBlock.ParCodeBlock
	parCodeBlock := &codeFile.AllCodeBlock[parCodeBlockNo]
	nextCodeBlock := &codeFile.AllCodeBlock[nextCodeBlockNo]
	parCodeBlock.ChildCodeBlockLast = nextCodeBlockNo
	nowCodeBlock.NextCodeBlock = nextCodeBlockNo
	nextCodeBlock.ParCodeBlock = parCodeBlockNo
	nextCodeBlock.NextCodeBlock = -1
	return
}
func (codeFile *CodeFile) addChildCodeBlock(nowCodeBlockNo int, childCodeBlockNo int) {

	nowCodeBlock := &codeFile.AllCodeBlock[nowCodeBlockNo]
	childCodeBlock := &codeFile.AllCodeBlock[childCodeBlockNo]

	if nowCodeBlock.ChildCodeBlockFirst == -1 {
		nowCodeBlock.ChildCodeBlockFirst = childCodeBlockNo
	}
	lastChildCodeBlockNo := nowCodeBlock.ChildCodeBlockLast
	if lastChildCodeBlockNo != -1 {
		lastChildCodeBlock := &codeFile.AllCodeBlock[lastChildCodeBlockNo]
		lastChildCodeBlock.NextCodeBlock = childCodeBlockNo
	}
	nowCodeBlock.ChildCodeBlockLast = childCodeBlockNo
	childCodeBlock.ParCodeBlock = nowCodeBlockNo
	childCodeBlock.NextCodeBlock = -1

	return
}
func (codeFile *CodeFile) addNewLineCodeBlock() {
	newLineCodeBlockNo := codeFile.newCodeBlock(codeFile.NowPos, CbtLine)
	codeFile.addChildCodeBlock(0, newLineCodeBlockNo)
	codeFile.nowBlock = newLineCodeBlockNo
	//codeFile.AllCodeBlock = append(codeFile.AllCodeBlock, codeBlock)
	//codeBlock := CodeBlock{}
	//
	//codeBlock.Pos = codeFile.nowPos
	//codeBlock.BlockType = CbtLine
	//
	//codeBlock.ParNo = 0
	//codeBlock.FirstChildNo = -1
	//codeBlock.LastChildNo = -1
	//
	//codeBlock.NextNo = -1

	//codeFile.nowBlock = len(codeFile.AllCodeBlock)
	//codeFile.AllCodeBlock = append(codeFile.AllCodeBlock, codeBlock)

	//if codeFile.nowLineBlock > 0 {
	//	codeFile.AllCodeBlock[codeFile.nowLineBlock].NextNo = codeFile.nowBlock
	//	codeFile.addNextCodeBlock(codeFile.nowLineBlock, codeFile.nowBlock)
	//} else {
	//	codeFile.AllCodeBlock[0].FirstChildNo = codeFile.nowBlock
	//	codeFile.addChildCodeBlock(0, codeFile.nowBlock)
	//}

	//codeFile.nowLineBlock = codeFile.nowBlock

	return
}
func (codeFile *CodeFile) addNewChildLineToChild(parNo int) (childLineCodeBlockNo int) {
	nowCodeBlock := codeFile.AllCodeBlock[parNo]

	codeBlock := CodeBlock{}
	codeBlock.Pos = nowCodeBlock.Pos
	codeBlock.Pos.BlockLen = 1
	codeBlock.BlockType = CbtChildLine

	codeBlock.ParCodeBlock = parNo
	codeBlock.ChildCodeBlockFirst = codeFile.AllCodeBlock[parNo].ChildCodeBlockFirst
	codeBlock.ChildCodeBlockLast = -1
	codeBlock.NextCodeBlock = -1

	childLineCodeBlockNo = len(codeFile.AllCodeBlock)

	codeFile.AllCodeBlock[parNo].ChildCodeBlockFirst = childLineCodeBlockNo

	codeFile.AllCodeBlock = append(codeFile.AllCodeBlock, codeBlock)
	codeFile.addChildCodeBlock(parNo, childLineCodeBlockNo)
	//codeFile.AllCodeBlock[nowCodeBlockNewNo].ParNo = parNo

	//codeFile.nowBlock = childLineCodeBlockNo
	return
}
func (codeFile *CodeFile) addNewChildLineToNext(beforeNo int, beforeLine int) (childLineCodeBlockNo int) {
	nowCodeBlock := codeFile.AllCodeBlock[beforeNo]
	par := codeFile.AllCodeBlock[nowCodeBlock.ParCodeBlock].ParCodeBlock
	codeBlock := CodeBlock{}
	codeBlock.Pos = nowCodeBlock.Pos
	codeBlock.Pos.BlockLen = 1
	codeBlock.BlockType = CbtChildLine

	codeBlock.ParCodeBlock = par
	codeBlock.ChildCodeBlockFirst = nowCodeBlock.NextCodeBlock
	codeBlock.ChildCodeBlockLast = -1
	codeBlock.NextCodeBlock = -1

	childLineCodeBlockNo = len(codeFile.AllCodeBlock)

	codeFile.AllCodeBlock[beforeNo].NextCodeBlock = -1
	codeFile.AllCodeBlock = append(codeFile.AllCodeBlock, codeBlock)
	codeFile.addNextCodeBlock(beforeNo, childLineCodeBlockNo)

	//bn := codeFile.AllCodeBlock[codeBlock.ParNo].FirstChildNo
	//for bn >= 0 {
	//	//fmt.Println(codeBlock.ParNo, bn, beforeLine, childLineCodeBlockNo)
	//	if codeFile.AllCodeBlock[bn].NextNo < 0 {
	//		//fmt.Println(bn, beforeLine, childLineCodeBlockNo)
	//		//codeFile.AllCodeBlock[bn].NextNo = childLineCodeBlockNo
	//	}
	//	bn = codeFile.AllCodeBlock[bn].NextNo
	//	//if bn < 0 {
	//	//	codeFile.AllCodeBlock[bn].NextNo = childLineCodeBlockNo
	//	//}
	//}
	codeFile.AllCodeBlock[beforeLine].NextCodeBlock = childLineCodeBlockNo
	//codeFile.AllCodeBlock[nowCodeBlockNewNo].ParNo = parNo

	//codeFile.nowBlock = childLineCodeBlockNo
	return
}

func (codeFile *CodeFile) appendChild(codeBlockType CodeBlockType) {

	newCodeBlockNo := codeFile.newCodeBlock(codeFile.NowPos, codeBlockType)
	codeFile.addChildCodeBlock(codeFile.nowBlock, newCodeBlockNo)
	codeFile.nowBlock = newCodeBlockNo

	//codeBlock := CodeBlock{}
	//
	//codeBlock.Pos = codeFile.nowPos
	//codeBlock.BlockType = codeBlockType
	//
	//codeBlock.ParNo = codeFile.nowBlock
	//codeBlock.FirstChildNo = -1
	//codeBlock.LastChildNo = -1
	//codeBlock.NextNo = -1
	//if codeBlock.Pos.blockLen > 0 {
	//	codeBlock.Words = string(codeFile.allCodeChars[codeBlock.Pos.startNo : codeBlock.Pos.startNo+codeBlock.Pos.blockLen])
	//}
	//
	//nowBlockNo := len(codeFile.AllCodeBlock)
	//codeFile.AllCodeBlock = append(codeFile.AllCodeBlock, codeBlock)
	//
	//codeFile.AllCodeBlock[codeFile.nowBlock].FirstChildNo = nowBlockNo
	//
	//codeFile.addChildCodeBlock(codeFile.nowBlock, nowBlockNo)
	//codeFile.nowBlock = nowBlockNo

	return
}

func (codeFile *CodeFile) appendNext(codeBlockType CodeBlockType) {

	newCodeBlockNo := codeFile.newCodeBlock(codeFile.NowPos, codeBlockType)
	codeFile.addNextCodeBlock(codeFile.nowBlock, newCodeBlockNo)
	codeFile.nowBlock = newCodeBlockNo

	//codeBlock := CodeBlock{}
	//
	//codeBlock.Pos = codeFile.nowPos
	//codeBlock.BlockType = codeBlockType
	//
	//codeBlock.ParNo = codeFile.AllCodeBlock[codeFile.nowBlock].ParNo
	//codeBlock.FirstChildNo = -1
	//codeBlock.LastChildNo = -1
	//codeBlock.NextNo = -1
	//if codeBlock.Pos.blockLen > 0 {
	//	codeBlock.Words = string(codeFile.allCodeChars[codeBlock.Pos.startNo : codeBlock.Pos.startNo+codeBlock.Pos.blockLen])
	//}
	//
	//nowBlockNo := len(codeFile.AllCodeBlock)
	//codeFile.AllCodeBlock = append(codeFile.AllCodeBlock, codeBlock)
	//
	//codeFile.AllCodeBlock[codeFile.nowBlock].NextNo = nowBlockNo
	//
	//codeFile.addNextCodeBlock(codeFile.nowBlock, nowBlockNo)
	//
	//codeFile.nowBlock = nowBlockNo

	return
}
func (codeFile *CodeFile) changeNowBlockType(blockType CodeBlockType) {
	codeFile.AllCodeBlock[codeFile.nowBlock].BlockType = blockType
}
func (codeFile *CodeFile) setEndPos() {
	endPos := codeFile.NowPos
	codeBlock := &codeFile.AllCodeBlock[codeFile.nowBlock]
	endNo := endPos.StartNo + endPos.BlockLen
	if endNo > codeBlock.Pos.StartNo {
		codeBlock.Pos.BlockLen = endNo - codeBlock.Pos.StartNo
	}
	endLineNo := endPos.LineNo
	if endLineNo > codeBlock.Pos.LineNo {
		codeBlock.Pos.LineCount = endLineNo - codeBlock.Pos.LineNo + 1
	}
	if codeBlock.Pos.BlockLen > 0 {
		codeBlock.Words = string(codeFile.allCodeChars[codeBlock.Pos.StartNo : codeBlock.Pos.StartNo+codeBlock.Pos.BlockLen])
	}

}

func (codeFile *CodeFile) Parse() (err error) {

	for n, char := range codeFile.allCodeChars {
		codeFile.checkChar(n, char)
		codeFile.ChangePosLineNo(char)
	}
	codeFile.CheckCharEnd()
	fmt.Println("checkCharEnd")

	codeFile.CheckChildLine(0)
	err = codeFile.CheckLineIndent(0)
	if err != nil {
		return
	}
	codeFile.ClearEmptyCodeBlock(0)
	codeFile.CheckLineColon(0)

	return
}

func (codeFile *CodeFile) checkChar(n int, char rune) {
	charType := codeFile.getCharType(char)
	//fmt.Println(n, char, string(char), charType)
	codeFile.NowPos.StartNo = n
	codeFile.NowPos.BlockLen = 1
	//codeFile.appendChild(charType)
	//nowCodeBlock := codeFile.newCodeBlock(codeFile.nowPos, charType)
	check := false
	if !check {
		check = codeFile.CheckNewLine(charType)
	}
	if !check {
		check = codeFile.CheckComment(charType)
	}
	if !check {
		check = codeFile.CheckString(charType)
	}
	if !check {
		check = codeFile.CheckBracket(charType)
	}
	if !check {
		check = codeFile.CheckOperator(charType)
	}
	if !check {
		check = codeFile.CheckLetter(charType)
	}
	if !check {
		check = codeFile.CheckNumber(charType)
	}
	if !check {
		check = codeFile.CheckCRLF(charType)
	}
	if !check {
		check = codeFile.CheckNowrap(charType)
	}
	if !check {
		codeFile.appendNext(charType)
	}

}

func (codeFile *CodeFile) getCharType(r rune) (t CodeBlockType) {
	t, ok := CharCodeBlockType[r]
	if !ok {
		if unicode.IsLetter(r) {
			t = CbtLetter
		} else if unicode.IsNumber(r) {
			t = CbtNumber
		} else {
			t = CbtOtherChar
		}
	}
	return
}

func (codeFile *CodeFile) CheckNewLine(nowCharType CodeBlockType) (check bool) {

	if codeFile.NewLine {
		lastBlockType := codeFile.AllCodeBlock[codeFile.nowBlock].BlockType
		switch lastBlockType {
		case CbtLeftApostrophe, CbtLeftQuotation:
			//多行字符串中换行符保留到字符串中
			codeFile.NewLine = false
		default:
			if codeFile.BracketCount > 0 {
				//括号中的代码当做一行处理
				codeFile.NewLine = false
			} else if codeFile.Nowrap == true {
				//行末尾有反斜杠，下行代码自动合并到当前行中一起处理
				codeFile.Nowrap = false
				codeFile.NewLine = false
			} else {
				codeFile.addNewLineCodeBlock()
				codeFile.appendChild(nowCharType)
				//lineCodeBlock := codeFile.newCodeBlock(nowCodeBlock.Pos, CbtLine)
				//codeFile.LastCodeBlock = codeFile.MainCode
				//codeFile.appendChild(lineCodeBlock)
				//codeFile.appendChild(nowCodeBlock)
				switch nowCharType {
				case CbtLF:
					codeFile.NewLine = true
				default:
					codeFile.NewLine = false
				}
				check = true
			}
		}
	}
	return
}
func (codeFile *CodeFile) ChangePosLineNo(char rune) {
	if codeFile.LastChar == '\r' && char == '\n' {
		codeFile.NowPos.LineNo += 1
		codeFile.NowPos.ColNo = 1
	} else if codeFile.LastChar == '\r' && char == '\r' {
		codeFile.NowPos.LineNo += 1
		codeFile.NowPos.ColNo = 1
	} else if char == '\n' {
		codeFile.NowPos.LineNo += 1
		codeFile.NowPos.ColNo = 1
	} else {
		codeFile.NowPos.ColNo += 1
	}
	codeFile.LastChar = char
}

func (codeFile *CodeFile) CheckComment(nowCharType CodeBlockType) (check bool) {
	//单行注释#号键开始
	lastBlockType := codeFile.AllCodeBlock[codeFile.nowBlock].BlockType
	if lastBlockType == CbtPound {
		switch nowCharType {
		case CbtCR, CbtLF, CbtCRLF:
			codeFile.changeNowBlockType(CbtComment)
			//codeFile.AllCodeBlock[codeFile.nowBlock].BlockType = CbtComment
			//codeFile.appendNext(nowCodeBlock)
		default:
			codeFile.setEndPos()
			check = true
		}

	}
	return
}
func (codeFile *CodeFile) CheckString(nowCharType CodeBlockType) (check bool) {
	//用中文单引号和双引号可以声明多行字符串，用英文单引号和双引号只是声明单行字符串
	lastBlockType := codeFile.AllCodeBlock[codeFile.nowBlock].BlockType
	switch lastBlockType {
	case CbtLeftQuotation:
		switch nowCharType {
		case CbtRightQuotation:
			codeFile.setEndPos()
			codeFile.changeNowBlockType(CbtString)
			//codeFile.AllCodeBlock[codeFile.nowBlock].BlockType = CbtString
		default:
			codeFile.setEndPos()
		}
		check = true
	case CbtLeftApostrophe:
		switch nowCharType {
		case CbtRightApostrophe:
			codeFile.setEndPos()
			codeFile.changeNowBlockType(CbtString)
			//codeFile.AllCodeBlock[codeFile.nowBlock].BlockType = CbtString
		default:
			codeFile.setEndPos()
		}
		check = true
	case CbtApostrophe:
		switch nowCharType {
		case CbtApostrophe:
			codeFile.setEndPos()
			codeFile.changeNowBlockType(CbtString)
			//codeFile.AllCodeBlock[codeFile.nowBlock].BlockType = CbtString
		case CbtCR, CbtLF, CbtCRLF:
			codeFile.changeNowBlockType(CbtString)
			codeFile.appendNext(nowCharType)
		default:
			codeFile.setEndPos()
		}
		check = true
	case CbtQuotation:
		switch nowCharType {
		case CbtQuotation:
			codeFile.setEndPos()
			codeFile.changeNowBlockType(CbtString)
		case CbtCR, CbtLF, CbtCRLF:
			codeFile.changeNowBlockType(CbtString)
			codeFile.appendNext(nowCharType)
		default:
			codeFile.setEndPos()
		}
		check = true
	}
	return
}

func (codeFile *CodeFile) CheckBracket(nowCharType CodeBlockType) (check bool) {
	//代码分析的时候不检查括号匹配情况，留到预处理的时候检查
	lastBlockType := codeFile.AllCodeBlock[codeFile.nowBlock].BlockType
	switch lastBlockType {
	case CbtLeftBracket, CbtLeftSquareBracket, CbtLeftBigBracket:
		codeFile.BracketCount += 1
		codeFile.appendChild(nowCharType)
		check = true
	case CbtRightBracket, CbtRightSquareBracket, CbtRightBigBracket:
		codeFile.BracketCount -= 1
		codeFile.nowBlock = codeFile.AllCodeBlock[codeFile.nowBlock].ParCodeBlock
		//codeFile.appendNext(nowCodeBlock)
		//check = true
		//case CbtCR, CbtLF, CbtCRLF:
		//	if codeFile.bracketCount > 0 {
		//		codeFile.LastCodeBlock.BlockType = CbtEnter
		//		codeFile.appendNext(nowCodeBlock)
		//		check = true
		//	}
	}
	return
}

func (codeFile *CodeFile) CheckOperator(nowCharType CodeBlockType) (check bool) {
	//震语言支持多种组合运算符
	lastBlockType := codeFile.AllCodeBlock[codeFile.nowBlock].BlockType
	switch lastBlockType {
	case CbtOperator, CbtColon:
		switch nowCharType {
		case CbtOperator:
			codeFile.changeNowBlockType(CbtOperator)
			codeFile.setEndPos()
			check = true
		default:
			//codeFile.appendNext(nowCodeBlock)
		}

	}
	return

}

func (codeFile *CodeFile) CheckLetter(nowCharType CodeBlockType) (check bool) {
	//标识符需要以下划线、字母或者汉字开头，然后标识符中可以包含数字
	lastBlockType := codeFile.AllCodeBlock[codeFile.nowBlock].BlockType
	switch lastBlockType {
	case CbtLetter:
		switch nowCharType {
		case CbtLetter, CbtNumber:
			codeFile.setEndPos()
			check = true
		case CbtPoint:
			//todo 标识符要不要用点直接分割
			codeFile.appendNext(nowCharType)
			//codeFile.LastCodeBlock.addLen(1)
			check = true
		default:
			//codeFile.appendNext(nowCodeBlock)
		}

	}
	return
}

func (codeFile *CodeFile) CheckNumber(nowCharType CodeBlockType) (check bool) {
	//数字可以是十进制也可以16进制，此处不检查数字是否有问题，留到预处理中进行检查
	lastBlockType := codeFile.AllCodeBlock[codeFile.nowBlock].BlockType
	switch lastBlockType {
	case CbtNumber:
		switch nowCharType {
		case CbtNumber, CbtPoint, CbtLetter:
			codeFile.setEndPos()
			check = true
		default:
			//codeFile.appendNext(nowCodeBlock)
		}

	}
	return
}

func (codeFile *CodeFile) CheckCRLF(nowCharType CodeBlockType) (check bool) {
	//震语言默认用\n换行，但是也支持windows中\r\n换行
	lastBlockType := codeFile.AllCodeBlock[codeFile.nowBlock].BlockType
	switch lastBlockType {
	case CbtCR:
		switch nowCharType {
		case CbtLF:
			codeFile.changeNowBlockType(CbtCRLF)
			codeFile.setEndPos()
			codeFile.NewLine = true
			check = true
		case CbtCR:
			codeFile.appendNext(nowCharType)
			codeFile.NewLine = true
			check = true
		default:
			//codeFile.appendNext(nowCodeBlock)
			//codeFile.newLine = true
		}
	default:
		switch nowCharType {
		case CbtLF:
			codeFile.appendNext(nowCharType)
			codeFile.NewLine = true
			check = true
		}
	}

	return
}

func (codeFile *CodeFile) CheckNowrap(nowCharType CodeBlockType) (check bool) {
	//当行末尾如果是反斜杠，下行代码当做一行处理
	lastBlockType := codeFile.AllCodeBlock[codeFile.nowBlock].BlockType
	switch lastBlockType {
	case CBtBackslash:
		codeFile.Nowrap = true
		codeFile.appendNext(nowCharType)
		check = true
	}
	if codeFile.Nowrap == true {
		switch nowCharType {
		case CbtSpace, CbtFullWidthSpace, CbtTab:
			codeFile.Nowrap = true
		case CbtLF, CbtCR, CbtCRLF:
			//codeFile.LastCodeBlock.BlockType = CbtEnter
			codeFile.Nowrap = true
		case CbtComment:
			//todo 反斜杠后面如果是 注释 需要再分析
		case CbtLine:
			codeFile.Nowrap = false
		default:
			codeFile.Nowrap = false
		}
	}
	return
}

func (codeFile *CodeFile) CheckCharEnd() {
	lastBlockType := codeFile.AllCodeBlock[codeFile.nowBlock].BlockType
	switch lastBlockType {
	case CbtPound:
		codeFile.changeNowBlockType(CbtComment)
	case CbtLeftQuotation, CbtLeftApostrophe, CbtApostrophe, CbtQuotation:
		codeFile.changeNowBlockType(CbtString)
	default:
		//todo 括号是否结束等问题
	}
}

func (codeFile *CodeFile) CheckChildLine(codeBlock int) {
	//var childItems []int
	//
	//c := codeFile.AllCodeBlock[codeBlock].FirstChildNo
	//for c >= 0 {
	//	childItems = append(childItems, c)
	//	c = codeFile.AllCodeBlock[c].NextNo
	//}

	if codeFile.FindSeparator(codeBlock, CbtPeriod) {
		//fmt.Println(codeBlock)
		codeFile.SeparatorLine(codeBlock, CbtPeriod)
	}

	//return
	if codeFile.FindSeparator(codeBlock, CbtSemicolon) {
		codeFile.SeparatorLine(codeBlock, CbtSemicolon)
	}
	if codeFile.FindSeparator(codeBlock, CbtColon) {
		codeFile.SeparatorLineWithColon(codeBlock)
	}
	if codeFile.FindSeparator(codeBlock, CbtComma) {
		codeFile.SeparatorLine(codeBlock, CbtComma)
	}
	if codeFile.FindSeparator(codeBlock, CbtDunHao) {
		codeFile.SeparatorLine(codeBlock, CbtDunHao)
	}

	cc := codeFile.AllCodeBlock[codeBlock].ChildCodeBlockFirst
	for cc >= 0 {
		codeFile.CheckChildLine(cc)
		cc = codeFile.AllCodeBlock[cc].NextCodeBlock
	}
	//for _, c := range childItems {
	//	codeFile.checkChildLine(c)
	//	c = codeFile.AllCodeBlock[c].NextNo
	//}
	return
}

func (codeFile *CodeFile) FindSeparator(codeBlockNo int, separator CodeBlockType) (check bool) {
	hasSeparator := false
	codeBlock := codeFile.AllCodeBlock[codeBlockNo]
	n := codeBlock.ChildCodeBlockFirst

	for n >= 0 {
		c := codeFile.AllCodeBlock[n]
		enableCode := false
		switch c.BlockType {
		case separator:
			hasSeparator = true
		case CbtLetter, CbtNumber, CbtOperator, CbtString:
			enableCode = true
		case CbtApostrophe, CbtLeftApostrophe, CbtRightApostrophe,
			CbtQuotation, CbtLeftQuotation, CbtRightQuotation,
			CbtLeftBracket, CbtRightBracket,
			CbtLeftSquareBracket, CbtRightSquareBracket,
			CbtLeftBigBracket, CbtRightBigBracket:
			enableCode = true
		case CbtLine, CbtChildLine:
			enableCode = true
		}

		if hasSeparator && enableCode {
			check = true
			return
		}
		n = c.NextCodeBlock
	}
	return
}

//func (analyze *TxtCodeAnalyze) separatorLine(codeBlock *CodeBlock2, separator CodeBlockType) {
//	oldItems := codeBlock.items
//	codeBlock.items = []*CodeBlock2{}
//	codeGroup := analyze.newCodeBlock(codeBlock.Pos, CbtChildLine)
//	codeBlock.addItem(codeGroup)
//	for _, c := range oldItems {
//		codeGroup.addItem(c)
//		if c.BlockType == separator {
//			codeGroup = analyze.newCodeBlock(codeBlock.Pos, CbtChildLine)
//			codeBlock.addItem(codeGroup)
//		}
//	}
//}

func (codeFile *CodeFile) SeparatorLine(codeBlockNo int, separator CodeBlockType) {

	var childItems []int
	c := codeFile.AllCodeBlock[codeBlockNo].ChildCodeBlockFirst
	for c >= 0 {
		childItems = append(childItems, c)
		c = codeFile.AllCodeBlock[c].NextCodeBlock
	}
	codeBlock := &codeFile.AllCodeBlock[codeBlockNo]
	//codeBlock.NextNo = -1
	codeBlock.ChildCodeBlockFirst = -1
	codeBlock.ChildCodeBlockLast = -1

	codeGroup := codeFile.newCodeBlock(codeBlock.Pos, CbtChildLine)
	codeFile.addChildCodeBlock(codeBlockNo, codeGroup)
	//codeBlock.addItem(codeGroup)

	//clno := codeFile.addNewChildLineToChild(codeBlockNo)

	for _, n := range childItems {
		codeFile.addChildCodeBlock(codeGroup, n)
		//codeGroup.addItem(c)
		c := codeFile.AllCodeBlock[n]

		if c.BlockType == separator {
			codeGroup = codeFile.newCodeBlock(c.Pos, CbtChildLine)
			codeFile.addChildCodeBlock(codeBlockNo, codeGroup)
			//clno = codeFile.addNewChildLineToNext(n, clno)
			//codeFile.AllCodeBlock[n].NextNo = -1
			//codeBlock.addItem(codeGroup)codeGroup
		}
	}

	//for c.NextNo >= 0 {
	//	if c.BlockType == separator {
	//		codeFile.AllCodeBlock[c.ParNo].NextNo = c.NextNo
	//		codeFile.addNewChildLineCodeBlock(c.NextNo)
	//		c = codeFile.AllCodeBlock[codeFile.nowBlock]
	//	} else {
	//		c = codeFile.AllCodeBlock[c.NextNo]
	//	}
	//
	//}
}

//func (analyze *TxtCodeAnalyze) SeparatorLineWithColon(codeBlock *CodeBlock2) {
//	oldItems := codeBlock.items
//	codeBlock.items = []*CodeBlock2{}
//	nowCode := codeBlock
//	for n, c := range oldItems {
//		if n == 0 {
//			nowCode.addItem(c)
//		} else {
//			if nowCode.BlockType == CbtColon {
//				nowCode.addItem(c)
//			} else {
//				nowCode.parCodeBlock.addItem(c)
//			}
//		}
//		nowCode = c
//	}
//}
func (codeFile *CodeFile) SeparatorLineWithColon(codeBlockNo int) {

	var childItems []int
	c := codeFile.AllCodeBlock[codeBlockNo].ChildCodeBlockFirst
	for c >= 0 {
		childItems = append(childItems, c)
		c = codeFile.AllCodeBlock[c].NextCodeBlock
	}
	codeBlock := &codeFile.AllCodeBlock[codeBlockNo]
	//codeBlock.NextNo = -1
	codeBlock.ChildCodeBlockFirst = -1
	codeBlock.ChildCodeBlockLast = -1
	nowCodeNo := codeBlockNo
	nowCode := &codeFile.AllCodeBlock[nowCodeNo]
	//codeGroup := codeFile.newCodeBlock(codeBlock.Pos, CbtChildLine)
	//codeFile.addChildCodeBlock(codeBlockNo, codeGroup)
	//clno := codeFile.addNewChildLineToChild(codeBlockNo)
	for n, c := range childItems {
		if n == 0 {
			codeFile.addChildCodeBlock(nowCodeNo, c)
			//nowCode.addItem(c)
		} else {

			if nowCode.BlockType == CbtColon {
				codeFile.addChildCodeBlock(nowCodeNo, c)
			} else {
				codeFile.addChildCodeBlock(nowCode.ParCodeBlock, c)
				//.addItem(c)
			}
		}

		nowCodeNo = c
		nowCode = &codeFile.AllCodeBlock[nowCodeNo]
	}
	//for _, n := range childItems {
	//	//codeGroup.addItem(c)
	//	c := codeFile.AllCodeBlock[n]
	//	if c.BlockType == CbtColon {
	//		c.FirstChildNo = c.NextNo
	//		c.NextNo = -1
	//		//clno = codeFile.addNewChildLineToChild(n)
	//		//codeFile.AllCodeBlock[n].NextNo = -1
	//		//codeBlock.addItem(codeGroup)codeGroup
	//	}
	//}
	//oldItems := codeBlock.items
	//codeBlock.items = []*CodeBlock2{}
	//nowCode := codeBlock
	//for n, c := range oldItems {
	//	if n == 0 {
	//		nowCode.addItem(c)
	//	} else {
	//		if nowCode.BlockType == CbtColon {
	//			nowCode.addItem(c)
	//		} else {
	//			nowCode.parCodeBlock.addItem(c)
	//		}
	//	}
	//	nowCode = c
	//}
}
func (codeFile *CodeFile) ToXmlElement(codeBlockNo int, element *etree.Element) {
	showWords := false
	codeBlock := codeFile.AllCodeBlock[codeBlockNo]
	name := ""
	words := codeBlock.Words
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
		words = codeBlock.Words
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
		//	codeFile.codeBlockToXmlElement(c, e)
		//}
		c := codeBlock.ChildCodeBlockFirst
		for c >= 0 {
			codeFile.ToXmlElement(c, e)
			c = codeFile.AllCodeBlock[c].NextCodeBlock
		}
		//for _, c := range codeBlock.items {
		//
		//}
	}

	return
}
func (codeFile *CodeFile) ToXmlFile(path string) (err error) {

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	element := doc.CreateElement("程序")
	c := codeFile.AllCodeBlock[0].ChildCodeBlockFirst
	for c >= 0 {
		codeFile.ToXmlElement(c, element)
		c = codeFile.AllCodeBlock[c].NextCodeBlock
	}

	doc.Indent(4)
	err = doc.WriteToFile(path)

	if err != nil {
		return
	}

	return

}

func (codeFile *CodeFile) CheckLineEmpty(codeBlockNo int, ignoreComment bool) (empty bool) {
	codeBlock := codeFile.AllCodeBlock[codeBlockNo]

	empty = true
	switch codeBlock.BlockType {
	case CbtComment:
		if ignoreComment {
			empty = false
		}
	case CbtLetter, CbtNumber, CbtString, CbtOperator:
		empty = false
	}
	if empty == true {

		c := codeFile.AllCodeBlock[codeBlockNo].ChildCodeBlockFirst
		for c >= 0 {
			empty = codeFile.CheckLineEmpty(c, ignoreComment)
			if empty == false {
				return
			}

			c = codeFile.AllCodeBlock[c].NextCodeBlock
		}

		//c codeBlock.FirstChildNo
		//for _, c := range codeBlock.items {
		//	empty = codeFile.checkLineEmpty(c, ignoreComment)
		//	if empty == false {
		//		return
		//	}
		//}
	}
	return
}
func (codeFile *CodeFile) getIndent(codeBlockNo int) (indent int) {
	//codeBlock := codeFile.AllCodeBlock[codeBlockNo]
	indent = 0
	c := codeFile.AllCodeBlock[codeBlockNo].ChildCodeBlockFirst
	for c >= 0 {
		//for _, c := range codeBlock.items {
		switch codeFile.AllCodeBlock[c].BlockType {
		case CbtSpace:
			indent += 1
		case CbtFullWidthSpace:
			indent += 2
		case CbtTab:
			indent += 4
		case CbtChildLine:
			indent += codeFile.getIndent(c)
			return
		default:
			return
		}
		c = codeFile.AllCodeBlock[c].NextCodeBlock
	}

	return
}

func (codeFile *CodeFile) CheckLineIndent(codeBlockNo int) (err error) {
	//oldItems := codeBlock.items

	var childItems []int
	c := codeFile.AllCodeBlock[codeBlockNo].ChildCodeBlockFirst
	for c >= 0 {
		childItems = append(childItems, c)
		c = codeFile.AllCodeBlock[c].NextCodeBlock
	}
	codeBlock := &codeFile.AllCodeBlock[codeBlockNo]

	codeBlock.ChildCodeBlockLast = -1
	codeBlock.ChildCodeBlockFirst = -1

	codeBlockIndents := make(map[int]int)
	//var lastLineNo int
	var lastLine *CodeBlock
	firstLine := true
	for _, n := range childItems {
		c := &codeFile.AllCodeBlock[n]
		if codeFile.CheckLineEmpty(n, true) {
			//codeBlock.addItem(c)
			//} else if analyze.checkLineEmpty(c, false) {
			//	c.LineIndent = lastLine.LineIndent
			//	codeBlock.addItem(c)
		} else {
			c.LineIndent = codeFile.getIndent(n)
			if firstLine {
				codeFile.addChildCodeBlock(codeBlockNo, n)
				//codeBlock.addItem(c)
				codeBlockIndents[c.LineIndent] = n
				firstLine = false
			} else {
				if lastLine.LineIndent < c.LineIndent {
					codeFile.addChildCodeBlock(codeBlockIndents[lastLine.LineIndent], n)
					//codeBlockIndents[lastLine.LineIndent].addItem(c)
					codeBlockIndents[c.LineIndent] = n
				} else {
					_, ok := codeBlockIndents[c.LineIndent]
					if !ok {
						//todo 错误格式需要统一
						err := fmt.Sprintf("代码缩进错误 %d，%d", c.Pos.LineNo, c.Pos.ColNo)
						return errors.New(err)
					}
					codeFile.addNextCodeBlock(codeBlockIndents[c.LineIndent], n)
					//codeBlockIndents[c.LineIndent].appendNext(c)
					codeBlockIndents[c.LineIndent] = n
				}
			}
			lastLine = c
			//lastLineNo = n
		}
	}
	return
}

func (codeFile *CodeFile) ClearEmptyCodeBlock(codeBlockNo int) {

	var childItems []int
	c := codeFile.AllCodeBlock[codeBlockNo].ChildCodeBlockFirst
	for c >= 0 {
		childItems = append(childItems, c)
		c = codeFile.AllCodeBlock[c].NextCodeBlock
	}
	codeBlock := &codeFile.AllCodeBlock[codeBlockNo]
	codeBlock.ChildCodeBlockLast = -1
	codeBlock.ChildCodeBlockFirst = -1
	//oldItems := codeBlock.items
	//codeBlock.items = []*CodeBlock2{}

	beforeCode := codeBlock
	for _, n := range childItems {
		c := &codeFile.AllCodeBlock[n]
		switch c.BlockType {
		case CbtCR, CbtLF, CbtCRLF:
			//codeBlock.addItem(c)
		case CbtSpace, CbtFullWidthSpace, CbtTab:
		case CbtComment:
			beforeCode.Comment = c.Words
		default:
			codeFile.addChildCodeBlock(codeBlockNo, n)
			//codeBlock.addItem(c)
			beforeCode = c
		}
	}
	c = codeBlock.ChildCodeBlockFirst
	for c >= 0 {
		codeFile.ClearEmptyCodeBlock(c)
		c = codeFile.AllCodeBlock[c].NextCodeBlock
	}
	//for _, c := range codeBlock.items {
	//	analyze.clearEmptyCodeBlock(c)
	//}
}
func (codeFile *CodeFile) CheckLineColon(codeBlockNo int) {
	if codeFile.FindSeparator(codeBlockNo, CbtColon) {
		codeFile.SeparatorLineWithColon(codeBlockNo)
	}
	codeBlock := &codeFile.AllCodeBlock[codeBlockNo]
	c := codeBlock.ChildCodeBlockFirst
	for c >= 0 {
		codeFile.CheckLineColon(c)
		c = codeFile.AllCodeBlock[c].NextCodeBlock
	}
	//for _, c := range codeBlock.items {
	//	//switch c.BlockType {
	//	//case CbtLine:
	//	codeFile.checkLineColon(c)
	//	//}
	//}
}
