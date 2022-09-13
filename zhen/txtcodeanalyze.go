package zhen

import (
	"errors"
	"fmt"
	"unicode"
)

type TxtCodeAnalyze struct {
	AllCodeRunes []rune
	MainCode     *CodeBlock

	NowPos        CodeBlockPos
	LastCodeBlock *CodeBlock
	LastChar      rune
	NewLine       bool
	Nowrap        bool
	BracketCount  int
}

func NewTxtCodeAnalyze(codeTxt string) (analyze *TxtCodeAnalyze) {
	analyze = &TxtCodeAnalyze{}
	txt := codeTxt
	//txt = strings.ReplaceAll(codeTxt, "\r\n", "\n")
	//txt = strings.ReplaceAll(txt, "\r", "\n")
	analyze.AllCodeRunes = []rune(txt)

	analyze.NowPos = CodeBlockPos{StartNo: 0, BlockLen: 0, LineNo: 1, LineCount: 1, ColNo: 1}
	analyze.MainCode = analyze.newCodeBlock(analyze.NowPos, CbtFile)
	analyze.MainCode.ParCodeBlock = analyze.MainCode

	analyze.LastCodeBlock = analyze.MainCode
	analyze.LastChar = 0
	analyze.NewLine = true
	analyze.Nowrap = false
	analyze.BracketCount = 0
	return
}

func (analyze *TxtCodeAnalyze) AnalyzeCode() (err error) {

	for n, char := range analyze.AllCodeRunes {
		analyze.CheckChar(n, char)
		analyze.ChangePosLineNo(char)
	}
	analyze.CheckCharEnd()

	analyze.CheckChildLine(analyze.MainCode)
	err = analyze.CheckLineIndent(analyze.MainCode)
	if err != nil {
		return
	}
	analyze.ClearEmptyCodeBlock(analyze.MainCode)
	analyze.CheckLineColon(analyze.MainCode)

	return
}

func (analyze *TxtCodeAnalyze) newCodeBlock(pos CodeBlockPos, codeBlockType CodeBlockType) (codeBlock *CodeBlock) {
	codeBlock = NewCodeBlock(analyze.AllCodeRunes, pos, codeBlockType)
	return
}

func (analyze *TxtCodeAnalyze) appendNext(codeBlock *CodeBlock) {
	analyze.LastCodeBlock.ParCodeBlock.addItem(codeBlock)
	analyze.LastCodeBlock = codeBlock
}

func (analyze *TxtCodeAnalyze) appendChild(codeBlock *CodeBlock) {
	analyze.LastCodeBlock.addItem(codeBlock)
	//analyze.NewLine = false
	analyze.LastCodeBlock = codeBlock
}

func (analyze *TxtCodeAnalyze) CheckChar(n int, char rune) {
	charType := analyze.getRuneCodeBlockType(char)
	analyze.NowPos.StartNo = n
	analyze.NowPos.BlockLen = 1
	nowCodeBlock := analyze.newCodeBlock(analyze.NowPos, charType)
	check := false
	if !check {
		check = analyze.CheckNewLine(nowCodeBlock)
	}
	if !check {
		check = analyze.CheckComment(nowCodeBlock)
	}
	if !check {
		check = analyze.CheckString(nowCodeBlock)
	}
	//if !check {
	//	check = analyze.CheckCRLF(nowCodeBlock)
	//}
	if !check {
		check = analyze.CheckBracket(nowCodeBlock)
	}

	if !check {
		check = analyze.CheckOperator(nowCodeBlock)
	}
	if !check {
		check = analyze.CheckLetter(nowCodeBlock)
	}
	if !check {
		check = analyze.CheckNumber(nowCodeBlock)
	}
	if !check {
		check = analyze.CheckCRLF(nowCodeBlock)
	}
	if !check {
		check = analyze.CheckNowrap(nowCodeBlock)
	}
	if !check {
		analyze.appendNext(nowCodeBlock)
	}

}
func (analyze *TxtCodeAnalyze) CheckCharEnd() {
	switch analyze.LastCodeBlock.BlockType {
	case CbtPound:
		analyze.LastCodeBlock.BlockType = CbtComment

	case CbtLeftQuotation, CbtLeftApostrophe, CbtApostrophe, CbtQuotation:
		analyze.LastCodeBlock.BlockType = CbtString

	default:
		//todo 括号是否结束等问题
	}
}

func (analyze *TxtCodeAnalyze) ChangePosLineNo(char rune) {
	if analyze.LastChar == '\r' && char == '\n' {
		analyze.NowPos.LineNo += 1
		analyze.NowPos.ColNo = 1
	} else if analyze.LastChar == '\r' && char == '\r' {
		analyze.NowPos.LineNo += 1
		analyze.NowPos.ColNo = 1
	} else if char == '\n' {
		analyze.NowPos.LineNo += 1
		analyze.NowPos.ColNo = 1
	} else {
		analyze.NowPos.ColNo += 1
	}
	analyze.LastChar = char
}

func (analyze *TxtCodeAnalyze) getRuneCodeBlockType(r rune) (t CodeBlockType) {
	switch r {
	case '\'':
		t = CbtApostrophe
	case '‘':
		t = CbtLeftApostrophe
	case '’':
		t = CbtRightApostrophe
	case '"':
		t = CbtQuotation
	case '“':
		t = CbtLeftQuotation
	case '”':
		t = CbtRightQuotation
	case '#':
		t = CbtPound
	case ' ':
		t = CbtSpace
	case '　':
		t = CbtFullWidthSpace
	case '\r':
		t = CbtCR
	case '\n':
		t = CbtLF
	case '(', '（':
		t = CbtLeftBracket
	case ')', '）':
		t = CbtRightBracket
	case '[':
		t = CbtLeftSquareBracket
	case ']':
		t = CbtRightSquareBracket
	case '{':
		t = CbtLeftBigBracket
	case '}':
		t = CbtRightBigBracket
	case '\t':
		t = CbtTab
	case ':', '：':
		t = CbtColon
	case '、':
		t = CbtDunHao
	case ',', '，':
		t = CbtComma
	case '。':
		t = CbtPeriod
	case ';', '；':
		t = CbtSemicolon
	case '\\':
		t = CBtBackslash
	case '=', '+', '-', '*', '/':
		t = CbtOperator
	case '>', '<', '&', '|', '!':
		t = CbtOperator
	case '^', '?':
		t = CbtOperator
	case '.':
		t = CbtPoint
	case '_', '@':
		t = CbtLetter
	default:
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

func (analyze *TxtCodeAnalyze) CheckNewLine(nowCodeBlock *CodeBlock) (check bool) {
	if analyze.NewLine {
		switch analyze.LastCodeBlock.BlockType {
		case CbtLeftApostrophe, CbtLeftQuotation:
			//多行字符串中换行符保留到字符串中
			analyze.NewLine = false
		default:
			if analyze.BracketCount > 0 {
				//括号中的代码当做一行处理
				analyze.NewLine = false
			} else if analyze.Nowrap == true {
				//行末尾有反斜杠，下行代码自动合并到当前行中一起处理
				analyze.Nowrap = false
				analyze.NewLine = false
			} else {
				lineCodeBlock := analyze.newCodeBlock(nowCodeBlock.Pos, CbtLine)
				analyze.LastCodeBlock = analyze.MainCode
				analyze.appendChild(lineCodeBlock)
				analyze.appendChild(nowCodeBlock)
				switch nowCodeBlock.BlockType {
				case CbtLF:
					analyze.NewLine = true
				default:
					analyze.NewLine = false
				}
				check = true
			}
		}
	}
	return
}

func (analyze *TxtCodeAnalyze) CheckComment(nowCodeBlock *CodeBlock) (check bool) {
	//单行注释#号键开始
	if analyze.LastCodeBlock.BlockType == CbtPound {
		switch nowCodeBlock.BlockType {
		case CbtCR, CbtLF, CbtCRLF:
			analyze.LastCodeBlock.BlockType = CbtComment
			//analyze.appendNext(nowCodeBlock)
		default:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
			check = true
		}

	}
	return
}

func (analyze *TxtCodeAnalyze) CheckString(nowCodeBlock *CodeBlock) (check bool) {
	//用中文单引号和双引号可以声明多行字符串，用英文单引号和双引号只是声明单行字符串
	switch analyze.LastCodeBlock.BlockType {
	case CbtLeftQuotation:
		switch nowCodeBlock.BlockType {
		case CbtRightQuotation:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
			analyze.LastCodeBlock.BlockType = CbtString
		default:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
		}
		check = true
	case CbtLeftApostrophe:
		switch nowCodeBlock.BlockType {
		case CbtRightApostrophe:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
			analyze.LastCodeBlock.BlockType = CbtString
		default:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
		}
		check = true
	case CbtApostrophe:
		switch nowCodeBlock.BlockType {
		case CbtApostrophe:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
			analyze.LastCodeBlock.BlockType = CbtString
		case CbtCR, CbtLF, CbtCRLF:
			analyze.LastCodeBlock.BlockType = CbtString
			analyze.appendNext(nowCodeBlock)
		default:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
		}
		check = true
	case CbtQuotation:
		switch nowCodeBlock.BlockType {
		case CbtQuotation:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
			analyze.LastCodeBlock.BlockType = CbtString
		case CbtCR, CbtLF, CbtCRLF:
			analyze.LastCodeBlock.BlockType = CbtString
			analyze.appendNext(nowCodeBlock)
		default:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
		}
		check = true
	}
	return
}

func (analyze *TxtCodeAnalyze) CheckCRLF(nowCodeBlock *CodeBlock) (check bool) {
	//震语言默认用\n换行，但是也支持windows中\n\r换行

	//if analyze.LastChar == '\r' && char == '\n' {
	//	analyze.NowPos.LineNo += 1
	//	analyze.NowPos.ColNo = 1
	//} else if analyze.LastChar == '\r' && char == '\r' {
	//	analyze.NowPos.LineNo += 1
	//	analyze.NowPos.ColNo = 1
	//} else if char == '\n' {
	//	analyze.NowPos.LineNo += 1
	//	analyze.NowPos.ColNo = 1
	//} else {
	//	analyze.NowPos.ColNo += 1
	//}
	switch analyze.LastCodeBlock.BlockType {
	case CbtCR:
		switch nowCodeBlock.BlockType {
		case CbtLF:
			analyze.LastCodeBlock.BlockType = CbtCRLF
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
			analyze.NewLine = true
			check = true
		case CbtCR:
			analyze.appendNext(nowCodeBlock)
			analyze.NewLine = true
			check = true
		default:
			//analyze.appendNext(nowCodeBlock)
			//analyze.NewLine = true
		}

	default:
		switch nowCodeBlock.BlockType {
		case CbtLF:
			analyze.appendNext(nowCodeBlock)
			analyze.NewLine = true
			check = true
		}

		//analyze.NewLine = true
	}
	//if !check {
	//
	//}
	//switch nowCodeBlock.BlockType {
	//case CbtLF:
	//	switch nowCodeBlock.BlockType {
	//	case CbtCR, CbtLF:
	//		analyze.appendNext(nowCodeBlock)
	//		analyze.NewLine = true
	//	default:
	//		analyze.appendNext(nowCodeBlock)
	//		analyze.NewLine = true
	//
	//	}
	//	analyze.NewLine = true
	//	//analyze.appendNext(nowCodeBlock)
	//	//analyze.NewLine = true
	//	check = true
	//default:
	//	//switch nowCodeBlock.BlockType {
	//	//case CbtCR, CbtLF:
	//	//	analyze.appendNext(nowCodeBlock)
	//	//	analyze.NewLine = true
	//	//	check = true
	//	//}
	//	//analyze.NewLine = true
	//}
	//analyze.NewLine = true
	return
}

func (analyze *TxtCodeAnalyze) CheckBracket(nowCodeBlock *CodeBlock) (check bool) {
	//代码分析的时候不检查括号匹配情况，留到预处理的时候检查
	switch analyze.LastCodeBlock.BlockType {
	case CbtLeftBracket, CbtLeftSquareBracket, CbtLeftBigBracket:
		analyze.BracketCount += 1
		analyze.appendChild(nowCodeBlock)
		check = true
	case CbtRightBracket, CbtRightSquareBracket, CbtRightBigBracket:
		analyze.BracketCount -= 1
		analyze.LastCodeBlock = analyze.LastCodeBlock.ParCodeBlock
		//analyze.appendNext(nowCodeBlock)
		//check = true
		//case CbtCR, CbtLF, CbtCRLF:
		//	if analyze.BracketCount > 0 {
		//		analyze.LastCodeBlock.BlockType = CbtEnter
		//		analyze.appendNext(nowCodeBlock)
		//		check = true
		//	}
	}
	return
}

func (analyze *TxtCodeAnalyze) CheckOperator(nowCodeBlock *CodeBlock) (check bool) {
	//震语言支持多种组合运算符
	switch analyze.LastCodeBlock.BlockType {
	case CbtOperator, CbtColon:
		switch nowCodeBlock.BlockType {
		case CbtOperator:
			analyze.LastCodeBlock.BlockType = CbtOperator
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
			check = true
		default:
			//analyze.appendNext(nowCodeBlock)
		}

	}
	return

}

func (analyze *TxtCodeAnalyze) CheckLetter(nowCodeBlock *CodeBlock) (check bool) {
	//标识符需要以下划线、字母或者汉字开头，然后标识符中可以包含数字
	switch analyze.LastCodeBlock.BlockType {
	case CbtLetter:
		switch nowCodeBlock.BlockType {
		case CbtLetter, CbtNumber:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
			check = true
		case CbtPoint:
			//todo 标识符要不要用点直接分割
			analyze.appendNext(nowCodeBlock)
			//analyze.LastCodeBlock.addLen(1)
			check = true
		default:
			//analyze.appendNext(nowCodeBlock)
		}

	}
	return
}

func (analyze *TxtCodeAnalyze) CheckNumber(nowCodeBlock *CodeBlock) (check bool) {
	//数字可以是十进制也可以16进制，此处不检查数字是否有问题，留到预处理中进行检查
	switch analyze.LastCodeBlock.BlockType {
	case CbtNumber:
		switch nowCodeBlock.BlockType {
		case CbtNumber, CbtPoint, CbtLetter:
			analyze.LastCodeBlock.setEndPos(nowCodeBlock)
			check = true
		default:
			//analyze.appendNext(nowCodeBlock)
		}

	}
	return
}

func (analyze *TxtCodeAnalyze) CheckNowrap(nowCodeBlock *CodeBlock) (check bool) {
	//当行末尾如果是反斜杠，下行代码当做一行处理
	switch analyze.LastCodeBlock.BlockType {
	case CBtBackslash:
		analyze.Nowrap = true
		analyze.appendNext(nowCodeBlock)
		check = true
	}
	if analyze.Nowrap == true {
		switch analyze.LastCodeBlock.BlockType {
		case CbtSpace, CbtFullWidthSpace, CbtTab:
			analyze.Nowrap = true
		case CbtLF, CbtCR, CbtCRLF:
			//analyze.LastCodeBlock.BlockType = CbtEnter
			analyze.Nowrap = true
		case CbtComment:
			//todo 反斜杠后面如果是 注释 需要再分析
		case CbtLine:
			analyze.Nowrap = false
		default:
			analyze.Nowrap = false
		}
	}
	return
}

func (analyze *TxtCodeAnalyze) CheckChildLine(codeBlock *CodeBlock) {
	if analyze.FindSeparator(codeBlock, CbtPeriod) {
		analyze.SeparatorLine(codeBlock, CbtPeriod)
	}
	if analyze.FindSeparator(codeBlock, CbtSemicolon) {
		analyze.SeparatorLine(codeBlock, CbtSemicolon)
	}
	if analyze.FindSeparator(codeBlock, CbtColon) {
		analyze.SeparatorLineWithColon(codeBlock)
	}
	if analyze.FindSeparator(codeBlock, CbtComma) {
		analyze.SeparatorLine(codeBlock, CbtComma)
	}
	if analyze.FindSeparator(codeBlock, CbtDunHao) {
		analyze.SeparatorLine(codeBlock, CbtDunHao)
	}
	for _, c := range codeBlock.Items {
		analyze.CheckChildLine(c)
	}
	return
}

func (analyze *TxtCodeAnalyze) FindSeparator(codeBlock *CodeBlock, separator CodeBlockType) (check bool) {
	hasSeparator := false
	for _, c := range codeBlock.Items {
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
	}
	return
}

func (analyze *TxtCodeAnalyze) SeparatorLine(codeBlock *CodeBlock, separator CodeBlockType) {
	oldItems := codeBlock.Items
	codeBlock.Items = []*CodeBlock{}
	codeGroup := analyze.newCodeBlock(codeBlock.Pos, CbtChildLine)
	codeBlock.addItem(codeGroup)
	for _, c := range oldItems {
		codeGroup.addItem(c)
		if c.BlockType == separator {
			codeGroup = analyze.newCodeBlock(codeBlock.Pos, CbtChildLine)
			codeBlock.addItem(codeGroup)
		}
	}
}

func (analyze *TxtCodeAnalyze) SeparatorLineWithColon(codeBlock *CodeBlock) {
	oldItems := codeBlock.Items
	codeBlock.Items = []*CodeBlock{}
	nowCode := codeBlock
	for n, c := range oldItems {
		if n == 0 {
			nowCode.addItem(c)
		} else {
			if nowCode.BlockType == CbtColon {
				nowCode.addItem(c)
			} else {
				nowCode.ParCodeBlock.addItem(c)
			}
		}
		nowCode = c
	}
}

func (analyze *TxtCodeAnalyze) CheckLineEmpty(codeBlock *CodeBlock, ignoreComment bool) (empty bool) {
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
		for _, c := range codeBlock.Items {
			empty = analyze.CheckLineEmpty(c, ignoreComment)
			if empty == false {
				return
			}
		}
	}
	return
}

func (analyze *TxtCodeAnalyze) getIndent(codeBlock *CodeBlock) (indent int) {
	indent = 0
	for _, c := range codeBlock.Items {
		switch c.BlockType {
		case CbtSpace:
			indent += 1
		case CbtFullWidthSpace:
			indent += 2
		case CbtTab:
			indent += 4
		case CbtChildLine:
			indent += analyze.getIndent(c)
			return
		default:
			return
		}
	}

	return
}
func (analyze *TxtCodeAnalyze) CheckLineColon(codeBlock *CodeBlock) {
	if analyze.FindSeparator(codeBlock, CbtColon) {
		analyze.SeparatorLineWithColon(codeBlock)
	}
	for _, c := range codeBlock.Items {
		//switch c.BlockType {
		//case CbtLine:
		analyze.CheckLineColon(c)
		//}
	}
}
func (analyze *TxtCodeAnalyze) CheckLineIndent(codeBlock *CodeBlock) (err error) {
	oldItems := codeBlock.Items
	codeBlock.Items = []*CodeBlock{}

	codeBlockIndents := make(map[int]*CodeBlock)
	var lastLine *CodeBlock
	firstLine := true
	for _, c := range oldItems {
		if analyze.CheckLineEmpty(c, true) {
			//codeBlock.addItem(c)
			//} else if analyze.CheckLineEmpty(c, false) {
			//	c.LineIndent = lastLine.LineIndent
			//	codeBlock.addItem(c)
		} else {
			c.LineIndent = analyze.getIndent(c)
			if firstLine {
				codeBlock.addItem(c)
				codeBlockIndents[c.LineIndent] = c
				firstLine = false
			} else {
				if lastLine.LineIndent < c.LineIndent {
					codeBlockIndents[lastLine.LineIndent].addItem(c)
					codeBlockIndents[c.LineIndent] = c
				} else {
					_, ok := codeBlockIndents[c.LineIndent]
					if !ok {
						//todo 错误格式需要统一
						err := fmt.Sprintf("代码缩进错误 %d，%d", c.Pos.LineNo, c.Pos.ColNo)
						return errors.New(err)
					}
					codeBlockIndents[c.LineIndent].appendNext(c)
					codeBlockIndents[c.LineIndent] = c
				}
			}
			lastLine = c
		}
	}
	return
}

func (analyze *TxtCodeAnalyze) ClearEmptyCodeBlock(codeBlock *CodeBlock) {
	oldItems := codeBlock.Items
	codeBlock.Items = []*CodeBlock{}

	beforeCode := codeBlock
	for _, c := range oldItems {
		switch c.BlockType {
		case CbtCR, CbtLF, CbtCRLF:
			//codeBlock.addItem(c)
		case CbtSpace, CbtFullWidthSpace, CbtTab:
		case CbtComment:
			beforeCode.Comment = c.getChars()
		default:
			codeBlock.addItem(c)
			beforeCode = c
		}
	}

	for _, c := range codeBlock.Items {
		analyze.ClearEmptyCodeBlock(c)
	}
}
