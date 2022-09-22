package zhen_0_02

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

type CodeFile struct {
	allCodeChars []rune

	nowPos     CodeBlockPos
	lastCodeNo int
	lastChar   rune

	newLine      bool
	nowrap       bool
	bracketCount int

	AllCodeBlock []CodeBlock
	//AllCodeArea   []CodeArea
	//CodeBlockArea map[int]int
	//AllVarNames   *CodeVars
	//MainCodeBlock *CodeBlock
}

func NewCodeFile(codeTxt string) (codeFile *CodeFile) {
	codeFile = &CodeFile{}
	codeFile.allCodeChars = []rune(codeTxt)

	codeFile.nowPos = CodeBlockPos{startNo: 0, blockLen: 0, LineNo: 1, LineCount: 1, ColNo: 1}
	codeFile.lastCodeNo = codeFile.newCodeBlock(codeFile.nowPos, CbtFile)

	codeFile.lastChar = 0
	codeFile.newLine = true
	codeFile.nowrap = false
	codeFile.bracketCount = 0

	//a := NewVarNames()
	//codeFile.AllVarNames = &a
	//codeFile.AllVarNames.AddVar(CodeVarKey{Name: "test", Type: CvtKeyWords})
	//codeFile.MainCodeBlock = &codeFile.AllCodeBlock[0]

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

func (code *CodeFile) ToJsonFile(jsonFile string) (err error) {
	data, err := json.MarshalIndent(code, "", "\t")
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
	defer file.Close()
	dec := gob.NewDecoder(file)
	codeFile = &CodeFile{}
	err = dec.Decode(&codeFile)
	if err != nil {
		return
	}
	return

}

func (code *CodeFile) ToGobFile(gobFile string) (err error) {
	file, err := os.Create(gobFile)
	if err != nil {
		return
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(code)
	if err != nil {
		return
	}
	return
}

func (code *CodeFile) newCodeBlock(pos CodeBlockPos, charType CodeBlockType) (codeBlockNo int) {
	codeBlockNo = len(code.AllCodeBlock)

	codeBlock := CodeBlock{}
	codeBlock.Pos = pos
	codeBlock.BlockType = charType
	switch charType {
	case CbtLine, CbtChildLine:
		codeBlock.Pos.blockLen = 0
	}

	//codeBlock.No = codeBlockNo
	codeBlock.ParNo = -1
	codeBlock.FirstChildNo = -1
	codeBlock.LastChildNo = -1
	codeBlock.NextNo = -1

	if codeBlock.Pos.blockLen > 0 {
		s := codeBlock.Pos.startNo
		e := codeBlock.Pos.startNo + codeBlock.Pos.blockLen
		codeBlock.Chars = string(code.allCodeChars[s:e])
	}

	code.AllCodeBlock = append(code.AllCodeBlock, codeBlock)
	return
}

func (code *CodeFile) newCodeBlock2(sourceCodeBlaock *CodeBlock) (codeBlock CodeBlock) {
	codeBlock = *sourceCodeBlaock
	code.AllCodeBlock = append(code.AllCodeBlock, codeBlock)
	return
}
func (code *CodeFile) getChildCodeBlock(codeBlockNo int) (childItems []int) {
	n := code.AllCodeBlock[codeBlockNo].FirstChildNo
	for n >= 0 {
		childItems = append(childItems, n)
		n = code.AllCodeBlock[n].NextNo
	}
	return
}

func (code *CodeFile) addNextCodeBlock(nowCodeBlockNo int, nextCodeBlockNo int) {
	nowCodeBlock := &code.AllCodeBlock[nowCodeBlockNo]
	parCodeBlockNo := nowCodeBlock.ParNo
	parCodeBlock := &code.AllCodeBlock[parCodeBlockNo]
	nextCodeBlock := &code.AllCodeBlock[nextCodeBlockNo]
	parCodeBlock.LastChildNo = nextCodeBlockNo
	nowCodeBlock.NextNo = nextCodeBlockNo
	nextCodeBlock.ParNo = parCodeBlockNo
	nextCodeBlock.NextNo = -1
	return
}

func (code *CodeFile) addChildCodeBlock(nowCodeBlockNo int, childCodeBlockNo int) {
	nowCodeBlock := &code.AllCodeBlock[nowCodeBlockNo]
	childCodeBlock := &code.AllCodeBlock[childCodeBlockNo]
	if nowCodeBlock.FirstChildNo == -1 {
		nowCodeBlock.FirstChildNo = childCodeBlockNo
	}
	lastChildCodeBlockNo := nowCodeBlock.LastChildNo
	if lastChildCodeBlockNo != -1 {
		lastChildCodeBlock := &code.AllCodeBlock[lastChildCodeBlockNo]
		lastChildCodeBlock.NextNo = childCodeBlockNo
	}
	nowCodeBlock.LastChildNo = childCodeBlockNo
	childCodeBlock.ParNo = nowCodeBlockNo
	childCodeBlock.NextNo = -1
	return
}

func (code *CodeFile) addNewLineCodeBlock() {
	newLineCodeBlockNo := code.newCodeBlock(code.nowPos, CbtLine)
	code.addChildCodeBlock(0, newLineCodeBlockNo)
	code.lastCodeNo = newLineCodeBlockNo

	return
}

func (code *CodeFile) appendChildCodeBlock(codeBlockType CodeBlockType) {
	newCodeBlockNo := code.newCodeBlock(code.nowPos, codeBlockType)
	code.addChildCodeBlock(code.lastCodeNo, newCodeBlockNo)
	code.lastCodeNo = newCodeBlockNo

	return
}

func (code *CodeFile) appendNextCodeBlock(codeBlockType CodeBlockType) {
	newCodeBlockNo := code.newCodeBlock(code.nowPos, codeBlockType)
	code.addNextCodeBlock(code.lastCodeNo, newCodeBlockNo)
	code.lastCodeNo = newCodeBlockNo
	return
}

func (code *CodeFile) changeNowCodeBlockType(blockType CodeBlockType) {
	code.AllCodeBlock[code.lastCodeNo].BlockType = blockType
}

func (code *CodeFile) setNowCodeBlockEndPos() {
	endPos := code.nowPos
	codeBlock := &code.AllCodeBlock[code.lastCodeNo]
	endNo := endPos.startNo + endPos.blockLen
	if endNo > codeBlock.Pos.startNo {
		codeBlock.Pos.blockLen = endNo - codeBlock.Pos.startNo
	}
	endLineNo := endPos.LineNo
	if endLineNo > codeBlock.Pos.LineNo {
		codeBlock.Pos.LineCount = endLineNo - codeBlock.Pos.LineNo + 1
	}
	if codeBlock.Pos.blockLen > 0 {
		s := codeBlock.Pos.startNo
		e := codeBlock.Pos.startNo + codeBlock.Pos.blockLen
		codeBlock.Chars = string(code.allCodeChars[s:e])
	}
}

func (code *CodeFile) Parse() (err error) {

	for n, char := range code.allCodeChars {
		code.checkChar(n, char)
		code.changePosLineNo(char)
	}
	code.checkCharEnd()

	code.checkChildLine(0)
	err = code.checkLineIndent(0)
	if err != nil {
		return
	}
	code.clearEmptyCodeBlock(0)
	code.checkLineColon(0)
	code.RefreshAllCodeBlock()
	return
}

func (code *CodeFile) checkChar(n int, char rune) {
	charType := code.getCharType(char)
	code.nowPos.startNo = n
	code.nowPos.blockLen = 1
	check := false
	if !check {
		check = code.checkNewLine(charType)
	}
	if !check {
		check = code.checkComment(charType)
	}
	if !check {
		check = code.checkString(charType)
	}
	if !check {
		check = code.checkBracket(charType)
	}
	if !check {
		check = code.checkOperator(charType)
	}
	if !check {
		check = code.checkLetter(charType)
	}
	if !check {
		check = code.checkNumber(charType)
	}
	if !check {
		check = code.checkCRLF(charType)
	}
	if !check {
		check = code.checkNowrap(charType)
	}
	if !check {
		code.appendNextCodeBlock(charType)
	}

}

func (code *CodeFile) getCharType(r rune) (t CodeBlockType) {
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

func (code *CodeFile) checkNewLine(nowCharType CodeBlockType) (check bool) {
	if code.newLine {
		lastBlockType := code.AllCodeBlock[code.lastCodeNo].BlockType
		switch lastBlockType {
		case CbtLeftApostrophe, CbtLeftQuotation:
			//多行字符串中换行符保留到字符串中
			code.newLine = false
		default:
			if code.bracketCount > 0 {
				//括号中的代码当做一行处理
				code.newLine = false
			} else if code.nowrap == true {
				//行末尾有反斜杠，下行代码自动合并到当前行中一起处理
				code.nowrap = false
				code.newLine = false
			} else {
				code.addNewLineCodeBlock()
				code.appendChildCodeBlock(nowCharType)
				switch nowCharType {
				case CbtLF:
					code.newLine = true
				default:
					code.newLine = false
				}
				check = true
			}
		}
	}
	return
}
func (code *CodeFile) changePosLineNo(char rune) {
	if code.lastChar == '\r' && char == '\n' {
		code.nowPos.LineNo += 1
		code.nowPos.ColNo = 1
	} else if code.lastChar == '\r' && char == '\r' {
		code.nowPos.LineNo += 1
		code.nowPos.ColNo = 1
	} else if char == '\n' {
		code.nowPos.LineNo += 1
		code.nowPos.ColNo = 1
	} else {
		code.nowPos.ColNo += 1
	}
	code.lastChar = char
}

func (code *CodeFile) checkComment(nowCharType CodeBlockType) (check bool) {
	//单行注释#号键开始
	lastBlockType := code.AllCodeBlock[code.lastCodeNo].BlockType
	if lastBlockType == CbtPound {
		switch nowCharType {
		case CbtCR, CbtLF, CbtCRLF:
			code.changeNowCodeBlockType(CbtComment)
		default:
			code.setNowCodeBlockEndPos()
			check = true
		}

	}
	return
}
func (code *CodeFile) checkString(nowCharType CodeBlockType) (check bool) {
	//用中文单引号和双引号可以声明多行字符串，用英文单引号和双引号只是声明单行字符串
	lastBlockType := code.AllCodeBlock[code.lastCodeNo].BlockType
	switch lastBlockType {
	case CbtLeftQuotation:
		switch nowCharType {
		case CbtRightQuotation:
			code.setNowCodeBlockEndPos()
			code.changeNowCodeBlockType(CbtString)
		default:
			code.setNowCodeBlockEndPos()
		}
		check = true
	case CbtLeftApostrophe:
		switch nowCharType {
		case CbtRightApostrophe:
			code.setNowCodeBlockEndPos()
			code.changeNowCodeBlockType(CbtString)
		default:
			code.setNowCodeBlockEndPos()
		}
		check = true
	case CbtApostrophe:
		switch nowCharType {
		case CbtApostrophe:
			code.setNowCodeBlockEndPos()
			code.changeNowCodeBlockType(CbtString)
		case CbtCR, CbtLF, CbtCRLF:
			code.changeNowCodeBlockType(CbtString)
			code.appendNextCodeBlock(nowCharType)
		default:
			code.setNowCodeBlockEndPos()
		}
		check = true
	case CbtQuotation:
		switch nowCharType {
		case CbtQuotation:
			code.setNowCodeBlockEndPos()
			code.changeNowCodeBlockType(CbtString)
		case CbtCR, CbtLF, CbtCRLF:
			code.changeNowCodeBlockType(CbtString)
			code.appendNextCodeBlock(nowCharType)
		default:
			code.setNowCodeBlockEndPos()
		}
		check = true
	}
	return
}

func (code *CodeFile) checkBracket(nowCharType CodeBlockType) (check bool) {
	//代码分析的时候不检查括号匹配情况，留到预处理的时候检查
	lastBlockType := code.AllCodeBlock[code.lastCodeNo].BlockType
	switch lastBlockType {
	case CbtLeftBracket, CbtLeftSquareBracket, CbtLeftBigBracket:
		code.bracketCount += 1
		code.appendChildCodeBlock(nowCharType)
		check = true
	case CbtRightBracket, CbtRightSquareBracket, CbtRightBigBracket:
		code.bracketCount -= 1
		code.lastCodeNo = code.AllCodeBlock[code.lastCodeNo].ParNo
	}
	return
}

func (code *CodeFile) checkOperator(nowCharType CodeBlockType) (check bool) {
	//震语言支持多种组合运算符
	lastBlockType := code.AllCodeBlock[code.lastCodeNo].BlockType
	switch lastBlockType {
	case CbtOperator, CbtColon:
		switch nowCharType {
		case CbtOperator:
			code.changeNowCodeBlockType(CbtOperator)
			code.setNowCodeBlockEndPos()
			check = true
		default:

		}

	}
	return

}

func (code *CodeFile) checkLetter(nowCharType CodeBlockType) (check bool) {
	//标识符需要以下划线、字母或者汉字开头，然后标识符中可以包含数字
	lastBlockType := code.AllCodeBlock[code.lastCodeNo].BlockType
	switch lastBlockType {
	case CbtLetter:
		switch nowCharType {
		case CbtLetter, CbtNumber:
			code.setNowCodeBlockEndPos()
			check = true
		case CbtPoint:
			//todo 标识符要不要用点直接分割
			code.appendNextCodeBlock(nowCharType)
			check = true
		default:

		}
	}
	return
}

func (code *CodeFile) checkNumber(nowCharType CodeBlockType) (check bool) {
	//数字可以是十进制也可以16进制，此处不检查数字是否有问题，留到预处理中进行检查
	lastBlockType := code.AllCodeBlock[code.lastCodeNo].BlockType
	switch lastBlockType {
	case CbtNumber:
		switch nowCharType {
		case CbtNumber, CbtPoint, CbtLetter:
			code.setNowCodeBlockEndPos()
			check = true
		default:

		}
	}
	return
}

func (code *CodeFile) checkCRLF(nowCharType CodeBlockType) (check bool) {
	//震语言默认用\n换行，但是也支持windows中\r\n换行
	lastBlockType := code.AllCodeBlock[code.lastCodeNo].BlockType
	switch lastBlockType {
	case CbtCR:
		switch nowCharType {
		case CbtLF:
			code.changeNowCodeBlockType(CbtCRLF)
			code.setNowCodeBlockEndPos()
			code.newLine = true
			check = true
		case CbtCR:
			code.appendNextCodeBlock(nowCharType)
			code.newLine = true
			check = true
		default:

		}
	default:
		switch nowCharType {
		case CbtLF:
			code.appendNextCodeBlock(nowCharType)
			code.newLine = true
			check = true
		}
	}

	return
}

func (code *CodeFile) checkNowrap(nowCharType CodeBlockType) (check bool) {
	//当行末尾如果是反斜杠，下行代码当做一行处理
	lastBlockType := code.AllCodeBlock[code.lastCodeNo].BlockType
	switch lastBlockType {
	case CBtBackslash:
		code.nowrap = true
		code.appendNextCodeBlock(nowCharType)
		check = true
	}
	if code.nowrap == true {
		switch nowCharType {
		case CbtSpace, CbtFullWidthSpace, CbtTab:
			code.nowrap = true
		case CbtLF, CbtCR, CbtCRLF:
			code.nowrap = true
		case CbtComment:
			//todo 反斜杠后面如果是 注释 需要再分析
		case CbtLine:
			code.nowrap = false
		default:
			code.nowrap = false
		}
	}
	return
}

func (code *CodeFile) checkCharEnd() {
	lastBlockType := code.AllCodeBlock[code.lastCodeNo].BlockType
	switch lastBlockType {
	case CbtPound:
		code.changeNowCodeBlockType(CbtComment)
	case CbtLeftQuotation, CbtLeftApostrophe, CbtApostrophe, CbtQuotation:
		code.changeNowCodeBlockType(CbtString)
	default:
		//todo 括号是否结束等问题
	}
}

func (code *CodeFile) checkChildLine(codeBlock int) {
	if code.findSeparator(codeBlock, CbtPeriod) {
		code.separatorLine(codeBlock, CbtPeriod)
	}
	if code.findSeparator(codeBlock, CbtSemicolon) {
		code.separatorLine(codeBlock, CbtSemicolon)
	}
	if code.findSeparator(codeBlock, CbtColon) {
		code.separatorLineWithColon(codeBlock)
	}
	if code.findSeparator(codeBlock, CbtComma) {
		code.separatorLine(codeBlock, CbtComma)
	}
	if code.findSeparator(codeBlock, CbtDunHao) {
		code.separatorLine(codeBlock, CbtDunHao)
	}

	n := code.AllCodeBlock[codeBlock].FirstChildNo
	for n >= 0 {
		code.checkChildLine(n)
		n = code.AllCodeBlock[n].NextNo
	}
	return
}

func (code *CodeFile) findSeparator(codeBlockNo int, separator CodeBlockType) (check bool) {
	hasSeparator := false
	codeBlock := &code.AllCodeBlock[codeBlockNo]
	n := codeBlock.FirstChildNo
	for n >= 0 {
		c := &code.AllCodeBlock[n]
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
		n = c.NextNo
	}
	return
}

func (code *CodeFile) separatorLine(codeBlockNo int, separator CodeBlockType) {
	childItems := code.getChildCodeBlock(codeBlockNo)
	codeBlock := &code.AllCodeBlock[codeBlockNo]
	codeBlock.FirstChildNo = -1
	codeBlock.LastChildNo = -1

	childLine := code.newCodeBlock(codeBlock.Pos, CbtChildLine)
	code.addChildCodeBlock(codeBlockNo, childLine)
	for _, n := range childItems {
		code.addChildCodeBlock(childLine, n)
		c := &code.AllCodeBlock[n]
		if c.BlockType == separator {
			childLine = code.newCodeBlock(c.Pos, CbtChildLine)
			code.addChildCodeBlock(codeBlockNo, childLine)
		}
	}
}

func (code *CodeFile) separatorLineWithColon(codeBlockNo int) {
	childItems := code.getChildCodeBlock(codeBlockNo)
	codeBlock := &code.AllCodeBlock[codeBlockNo]
	codeBlock.FirstChildNo = -1
	codeBlock.LastChildNo = -1

	nowCodeNo := codeBlockNo
	nowCode := &code.AllCodeBlock[nowCodeNo]
	for n, c := range childItems {
		if n == 0 {
			code.addChildCodeBlock(nowCodeNo, c)
		} else {

			if nowCode.BlockType == CbtColon {
				code.addChildCodeBlock(nowCodeNo, c)
			} else {
				code.addChildCodeBlock(nowCode.ParNo, c)
			}
		}

		nowCodeNo = c
		nowCode = &code.AllCodeBlock[nowCodeNo]
	}

}

func (code *CodeFile) checkLineEmpty(codeBlockNo int, ignoreComment bool) (empty bool) {
	codeBlock := &code.AllCodeBlock[codeBlockNo]

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
		c := codeBlock.FirstChildNo
		for c >= 0 {
			empty = code.checkLineEmpty(c, ignoreComment)
			if empty == false {
				return
			}
			c = code.AllCodeBlock[c].NextNo
		}
	}
	return
}
func (code *CodeFile) getIndent(codeBlockNo int) (indent int) {
	indent = 0
	n := code.AllCodeBlock[codeBlockNo].FirstChildNo
	for n >= 0 {
		c := &code.AllCodeBlock[n]
		switch c.BlockType {
		case CbtSpace:
			indent += 1
		case CbtFullWidthSpace:
			indent += 2
		case CbtTab:
			indent += 4
		case CbtChildLine:
			indent += code.getIndent(n)
			return
		default:
			return
		}
		n = c.NextNo
	}

	return
}

func (code *CodeFile) checkLineIndent(codeBlockNo int) (err error) {
	childItems := code.getChildCodeBlock(codeBlockNo)
	codeBlock := &code.AllCodeBlock[codeBlockNo]
	codeBlock.LastChildNo = -1
	codeBlock.FirstChildNo = -1

	codeBlockIndents := make(map[int]int)
	var lastLine *CodeBlock
	firstLine := true
	for _, n := range childItems {
		c := &code.AllCodeBlock[n]
		if code.checkLineEmpty(n, true) {

		} else {
			c.LineIndent = code.getIndent(n)
			if firstLine {
				code.addChildCodeBlock(codeBlockNo, n)
				codeBlockIndents[c.LineIndent] = n
				firstLine = false
			} else {
				if lastLine.LineIndent < c.LineIndent {
					code.addChildCodeBlock(codeBlockIndents[lastLine.LineIndent], n)
					codeBlockIndents[c.LineIndent] = n
				} else {
					_, ok := codeBlockIndents[c.LineIndent]
					if !ok {
						//todo 错误格式需要统一
						err := fmt.Sprintf("代码缩进错误 %d，%d", c.Pos.LineNo, c.Pos.ColNo)
						return errors.New(err)
					}
					code.addNextCodeBlock(codeBlockIndents[c.LineIndent], n)
					codeBlockIndents[c.LineIndent] = n
				}
			}
			lastLine = c
		}
	}
	return
}

func (code *CodeFile) clearEmptyCodeBlock(codeBlockNo int) {
	childItems := code.getChildCodeBlock(codeBlockNo)
	codeBlock := &code.AllCodeBlock[codeBlockNo]
	codeBlock.LastChildNo = -1
	codeBlock.FirstChildNo = -1

	beforeCode := codeBlock
	for _, n := range childItems {
		c := &code.AllCodeBlock[n]
		switch c.BlockType {
		case CbtCR, CbtLF, CbtCRLF:

		case CbtSpace, CbtFullWidthSpace, CbtTab:
		case CbtComment:
			beforeCode.Comment = c.Chars
		default:
			code.addChildCodeBlock(codeBlockNo, n)

			beforeCode = c
		}
	}
	n := codeBlock.FirstChildNo
	for n >= 0 {
		code.clearEmptyCodeBlock(n)
		n = code.AllCodeBlock[n].NextNo
	}

}
func (code *CodeFile) checkLineColon(codeBlockNo int) {
	if code.findSeparator(codeBlockNo, CbtColon) {
		code.separatorLineWithColon(codeBlockNo)
	}

	n := code.AllCodeBlock[codeBlockNo].FirstChildNo
	for n >= 0 {
		code.checkLineColon(n)
		n = code.AllCodeBlock[n].NextNo
	}

}
func (code *CodeFile) addCodeBlock(allCodeBlock []CodeBlock, no int, newNo int) {
	codeBlock := &allCodeBlock[no]
	n := codeBlock.FirstChildNo
	if n >= 0 {
		c := &allCodeBlock[n]
		newChildNo := len(code.AllCodeBlock)
		code.AllCodeBlock = append(code.AllCodeBlock, *c)
		newCodeBlock := &code.AllCodeBlock[newChildNo]
		newCodeBlock.NextNo = -1
		newCodeBlock.FirstChildNo = -1
		newCodeBlock.LastChildNo = -1
		code.addChildCodeBlock(newNo, newChildNo)
		code.addCodeBlock(allCodeBlock, n, newChildNo)
		n = c.NextNo
	}
	if codeBlock.NextNo >= 0 {
		c := &allCodeBlock[codeBlock.NextNo]
		newNextNo := len(code.AllCodeBlock)
		code.AllCodeBlock = append(code.AllCodeBlock, *c)
		newCodeBlock := &code.AllCodeBlock[newNextNo]
		newCodeBlock.NextNo = -1
		newCodeBlock.FirstChildNo = -1
		newCodeBlock.LastChildNo = -1
		code.addNextCodeBlock(newNo, newNextNo)
		code.addCodeBlock(allCodeBlock, codeBlock.NextNo, newNextNo)
	}

}
func (code *CodeFile) RefreshAllCodeBlock() {
	allCodeBlock := code.AllCodeBlock
	code.AllCodeBlock = make([]CodeBlock, 0)
	if len(allCodeBlock) >= 0 {
		c := &allCodeBlock[0]
		newNo := len(code.AllCodeBlock)
		code.AllCodeBlock = append(code.AllCodeBlock, *c)
		newCodeBlock := &code.AllCodeBlock[newNo]
		newCodeBlock.NextNo = -1
		newCodeBlock.FirstChildNo = -1
		newCodeBlock.LastChildNo = -1

		code.addCodeBlock(allCodeBlock, 0, newNo)
	}

}

func (code *CodeFile) codeBlockToXmlElement(codeBlockNo int, element *etree.Element) {
	showWords := false
	codeBlock := code.AllCodeBlock[codeBlockNo]
	name := ""
	words := codeBlock.Chars
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
		case OtPoint:
			name = "点"
			showWords = false
		default:
			name = codeBlock.Operator.Type.String()
			showWords = false
		}
	default:
		name = codeBlock.WordType.String()
		words = codeBlock.Chars
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
		//e.CreateAttr("len", strconv.Itoa(codeBlock.Pos.blockLen))

		if codeBlock.LineIndent > 0 {
			e.CreateAttr("LineIndent", strconv.Itoa(codeBlock.LineIndent))
		}

		c := codeBlock.FirstChildNo
		for c >= 0 {
			code.codeBlockToXmlElement(c, e)
			c = code.AllCodeBlock[c].NextNo
		}
		for _, s := range codeBlock.Steps {
			code.codeStepToXmlElement(s, e)
		}

	}

	return
}
func (code *CodeFile) codeStepToXmlElement(codeSetp CodeStep, element *etree.Element) {
	e := element.CreateElement("操作代码")
	if codeSetp.CodeStepType != CstNone {
		e.CreateAttr("指令", codeSetp.CodeStepType.String())
	}
	if codeSetp.VarName1 != "" {
		e.CreateAttr("变量1", codeSetp.VarName1)
	}
	if codeSetp.TempVarNo1 != 0 {
		e.CreateAttr("临时变量1", strconv.Itoa(codeSetp.TempVarNo1))
	}
	if codeSetp.VarName2 != "" {
		e.CreateAttr("变量2", codeSetp.VarName2)
	}
	if codeSetp.TempVarNo2 != 0 {
		e.CreateAttr("临时变量2", strconv.Itoa(codeSetp.TempVarNo2))
	}
	if codeSetp.ReturnVarNo != 0 {
		e.CreateAttr("计算结果临时变量", strconv.Itoa(codeSetp.ReturnVarNo))
	}

	if codeSetp.ValueString != "" {
		e.CreateAttr("值", codeSetp.ValueString)
	}
	//e.CreateAttr("line", strconv.Itoa(codeBlock.Pos.LineNo))
	//e.CreateAttr("col", strconv.Itoa(codeBlock.Pos.ColNo))
}
func (code *CodeFile) ToXmlFile(path string) (err error) {

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	element := doc.CreateElement("程序")
	c := code.AllCodeBlock[0].FirstChildNo
	for c >= 0 {
		code.codeBlockToXmlElement(c, element)
		c = code.AllCodeBlock[c].NextNo
	}

	doc.Indent(4)
	err = doc.WriteToFile(path)
	if err != nil {
		return
	}
	return

}

func (code *CodeFile) FormatToFile(path string) (err error) {
	format := NewCodeBlockFormat(code)
	return format.ToFile(path)
}
