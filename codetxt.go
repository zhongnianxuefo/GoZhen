package main

import (
	"fmt"
	"os"
	"strings"
)

type TxtCodeWordPos struct {
	//WordNo int

	StartNo  int
	BlockLen int
	LineNo   int
	ColNo    int
}

//type TxtCodeWord struct {
//	Pos      TxtCodeWordPos
//	Words    string
//	EndSpace bool
//	Type     TxtCodeWordType
//}

type CodeBlock struct {
	Pos          TxtCodeWordPos
	Type         TxtCodeRuneType
	NeedSpace    bool
	Words        string
	Items        []*CodeBlock
	ParCodeBlock *CodeBlock
}

//type TxtCodeWordGroup struct {
//	WordGroup []TxtCodeWordGroup
//}
//
//type TxtCodeLine struct {
//	Floor int
//	Words []TxtCodeWord
//}
//type TxtCodeType int
//
//const (
//	TCT_Block TxtCodeType = iota
//	TCT_Line
//)
//
//type TxtCodeBlock struct {
//	FirstLine  TxtCodeLine
//	ChildLines []TxtCodeBlock
//}
//
type TxtCode struct {
	Txt           string
	MainCodeBlock *CodeBlock
}
type TxtCodeRuneType int

const (
	TCRY_LeftQuotation TxtCodeRuneType = iota
	TCRY_RightQuotation
	TCRY_LeftBracket
	TCRY_RightBracket
	TCRY_Space
	TCRY_Tab
	TCRY_StartNote
	TCRY_EndLine
	TCRY_Colon
	TCRY_Comma
	TCRY_DH
	TCRY_Semicolon
	TCRY_Period
	TCRY_Operator
	TCRY_Mark
	TCRY_Other

	TCRY_String
	TCRY_NewLineTab
	TCRY_Note
	TCRY_Line
	TCRY_File
)

//type TxtCodeWordType int
//
//const (
//	TCWY_String TxtCodeWordType = iota
//	TCWY_Note
//	TCWY_Tab
//	TCWY_Space
//	TCWY_NewLine
//	TCWY_LeftBracket
//	TCWY_RightBracket
//	TCWY_Colon
//	TCWY_Comma
//	TCWY_DH
//	TCWY_Semicolon
//	TCWY_Period
//	TCWY_Operator
//	TCWY_Word
//	TCWY_Number
//	TCWY_KeyWord
//	TCWY_Other
//)

func newTxtCode(codes string) (txtCode TxtCode, err error) {
	txtCode.Txt = codes //strings.ReplaceAll(codes, "\r\n", "\n")
	err = txtCode.AnalyzeWords()
	if err != nil {
		return
	}

	return
}
func (txt *TxtCode) AnalyzeWordsStep1() (err error) {
	txtCode := txt.Txt

	needEndSpace := func(t1 TxtCodeRuneType, t2 TxtCodeRuneType) (space bool) {
		switch t1 {
		case TCRY_Line, TCRY_NewLineTab, TCRY_EndLine, TCRY_Tab, TCRY_Space:
			space = false
		case TCRY_LeftBracket, TCRY_RightBracket:
			space = false
		case TCRY_Colon, TCRY_Comma, TCRY_DH, TCRY_Semicolon, TCRY_Period:
			space = false
		default:
			switch t2 {
			case TCRY_Line, TCRY_NewLineTab, TCRY_EndLine, TCRY_Tab, TCRY_LeftBracket, TCRY_RightBracket:
				space = false
			case TCRY_Colon, TCRY_Comma, TCRY_DH, TCRY_Semicolon, TCRY_Period:
				space = false
			default:
				space = true
			}
		}
		return
	}

	newCodeBlock := func(txtCode string, pos TxtCodeWordPos, codeBlockType TxtCodeRuneType) (codeBlock *CodeBlock) {
		codeBlock = &CodeBlock{}
		codeBlock.Pos = pos
		codeBlock.Type = codeBlockType
		codeBlock.Words = txtCode[pos.StartNo : pos.StartNo+pos.BlockLen]
		return
	}
	codeBlockAddBlockLen := func(codeBlock *CodeBlock, addBlockLen int) {
		codeBlock.Pos.BlockLen += addBlockLen
		codeBlock.Words = txtCode[codeBlock.Pos.StartNo : codeBlock.Pos.StartNo+codeBlock.Pos.BlockLen]
	}

	addItemCodeBlock := func(codeBlock *CodeBlock, item *CodeBlock) {
		codeBlock.Items = append(codeBlock.Items, item)
		item.ParCodeBlock = codeBlock
	}

	newLineCodeBlock := func(pos TxtCodeWordPos) (codeBlock *CodeBlock) {
		codeBlock = &CodeBlock{}
		codeBlock.Pos = pos
		codeBlock.Type = TCRY_Line
		codeBlock.Words = ""
		return
	}
	newFileCodeBlock := func() (codeBlock *CodeBlock) {
		codeBlock = &CodeBlock{}
		codeBlock.Pos = TxtCodeWordPos{StartNo: 0, BlockLen: 0, LineNo: 1, ColNo: 1}
		codeBlock.Type = TCRY_File
		codeBlock.Words = ""
		codeBlock.ParCodeBlock = codeBlock
		return
	}

	appendCodeBlock := func(beforeCodeBlock *CodeBlock, nowCodeBlock *CodeBlock) *CodeBlock {
		beforeCodeBlock.NeedSpace = needEndSpace(beforeCodeBlock.Type, nowCodeBlock.Type)
		addItemCodeBlock(beforeCodeBlock.ParCodeBlock, nowCodeBlock)
		return nowCodeBlock
	}
	appendChildCodeBlock := func(beforeCodeBlock *CodeBlock, nowCodeBlock *CodeBlock) *CodeBlock {
		addItemCodeBlock(beforeCodeBlock, nowCodeBlock)
		return nowCodeBlock
	}

	appendNewLineCodeBlock := func(beforeCodeBlock *CodeBlock, nowCodeBlock *CodeBlock,
		floorCodeBlocks map[int]*CodeBlock, beforeFloor int, nowFloor int) *CodeBlock {

		//newLineCodeBlock := func(pos TxtCodeWordPos) (codeBlock *CodeBlock) {
		//	codeBlock = &CodeBlock{}
		//	codeBlock.Pos = pos
		//	codeBlock.Type = TCRY_Line
		//	codeBlock.Words = ""
		//	return
		//}
		lineCodeBlock := newLineCodeBlock(nowCodeBlock.Pos)

		if beforeFloor == nowFloor {
			addItemCodeBlock(floorCodeBlocks[nowFloor].ParCodeBlock, lineCodeBlock)
			floorCodeBlocks[nowFloor] = lineCodeBlock
			addItemCodeBlock(lineCodeBlock, nowCodeBlock)
		} else if beforeFloor < nowFloor {
			addItemCodeBlock(floorCodeBlocks[beforeFloor], lineCodeBlock)
			floorCodeBlocks[nowFloor] = lineCodeBlock
			addItemCodeBlock(lineCodeBlock, nowCodeBlock)
		} else if beforeFloor > nowFloor {
			//beforeCodeBlock = floorCodeBlocks[nowFloor]
			addItemCodeBlock(floorCodeBlocks[nowFloor].ParCodeBlock, lineCodeBlock)
			floorCodeBlocks[nowFloor] = lineCodeBlock
			addItemCodeBlock(lineCodeBlock, nowCodeBlock)
		}
		return nowCodeBlock
	}

	var mainCodeBlock *CodeBlock
	mainCodeBlock = newFileCodeBlock()

	var lineCodeBlock *CodeBlock
	var beforeCodeBlock *CodeBlock
	var nowCodeBlock *CodeBlock
	var beforeChar int32
	var nowChar int32

	floorCodeBlocks := make(map[int]*CodeBlock)
	beforeFloor := 0
	nowFloor := 0

	nowPos := TxtCodeWordPos{StartNo: 0, BlockLen: 0, LineNo: 1, ColNo: 1}
	for n, r := range txtCode {
		//fmt.Println(n, r, string(r))
		nowChar = r
		nowPos.StartNo = n
		nowPos.BlockLen = len(string(r))

		var t TxtCodeRuneType
		switch r {
		case '“':
			t = TCRY_LeftQuotation
		case '”':
			t = TCRY_RightQuotation
		case '#':
			t = TCRY_StartNote
		case ' ', '　':
			t = TCRY_Space
		case '\n', '\r':
			t = TCRY_EndLine
		case '(', '（':
			t = TCRY_LeftBracket
		case ')', '）':
			t = TCRY_RightBracket
		case '\t':
			t = TCRY_Tab
		case ':', '：':
			t = TCRY_Colon
		case '、':
			t = TCRY_DH
		case ',', '，':
			t = TCRY_Comma
		case '。':
			t = TCRY_Period
		case ';', '；':
			t = TCRY_Semicolon
		case '=', '+', '-', '*', '/', '<', '>', '!':
			t = TCRY_Operator
			//t = TCRY_Mark
		default:
			t = TCRY_Other
		}
		nowCodeBlock = newCodeBlock(txtCode, nowPos, t)
		if n == 0 {
			lineCodeBlock = newLineCodeBlock(nowPos)
			addItemCodeBlock(mainCodeBlock, lineCodeBlock)
			addItemCodeBlock(lineCodeBlock, nowCodeBlock)
			beforeCodeBlock = nowCodeBlock
			floorCodeBlocks[nowFloor] = lineCodeBlock
			beforeFloor = nowFloor
		} else {
			switch beforeCodeBlock.Type {
			case TCRY_LeftQuotation:
				if t != TCRY_RightQuotation {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
				} else {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
					beforeCodeBlock.Type = TCRY_String
				}
			case TCRY_StartNote:
				if t != TCRY_EndLine {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
				} else {
					beforeCodeBlock.Type = TCRY_Note
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
					//beforeCodeBlock = appendNewLineCodeBlock(beforeCodeBlock, nowCodeBlock, floorCodeBlocks, nowFloor, nowFloor)
					//addItemCodeBlock(beforeCodeBlock.ParCodeBlock, nowCodeBlock)
					//beforeCodeBlock = nowCodeBlock
				}
			case TCRY_EndLine:
				if t == TCRY_EndLine {
					if beforeCodeBlock.Words == "\r" && nowCodeBlock.Words == "\n" {
						codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
					} else {
						beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
						//addItemCodeBlock(beforeCodeBlock.ParCodeBlock, nowCodeBlock)
						//beforeCodeBlock = nowCodeBlock
					}
				} else if t == TCRY_Tab {
					nowCodeBlock.Type = TCRY_NewLineTab
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
					//addItemCodeBlock(beforeCodeBlock.ParCodeBlock, nowCodeBlock)
					//beforeCodeBlock = nowCodeBlock
				} else {
					nowFloor = 0
					//fmt.Println(nowPos.LineNo, nowFloor, nowCodeBlock.Words)
					//addItemCodeBlock(beforeCodeBlock.ParCodeBlock, nowCodeBlock)
					//beforeCodeBlock = nowCodeBlock
					beforeCodeBlock = appendNewLineCodeBlock(beforeCodeBlock, nowCodeBlock, floorCodeBlocks, beforeFloor, nowFloor)
					beforeFloor = nowFloor
				}
			case TCRY_NewLineTab:
				if t == TCRY_Tab {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
				} else if t == TCRY_EndLine {
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
					//addItemCodeBlock(beforeCodeBlock.ParCodeBlock, nowCodeBlock)
					//beforeCodeBlock = nowCodeBlock
				} else {
					nowFloor = beforeCodeBlock.Pos.BlockLen
					//fmt.Println(nowPos.LineNo, nowFloor, nowCodeBlock.Words)
					//addItemCodeBlock(beforeCodeBlock.ParCodeBlock, nowCodeBlock)
					//beforeCodeBlock = nowCodeBlock
					beforeCodeBlock = appendNewLineCodeBlock(beforeCodeBlock, nowCodeBlock, floorCodeBlocks, beforeFloor, nowFloor)
					beforeFloor = nowFloor
				}
			case TCRY_Space:
				if t == TCRY_Space {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
				} else {
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
					//addItemCodeBlock(beforeCodeBlock.ParCodeBlock, nowCodeBlock)
					//
					//beforeCodeBlock = nowCodeBlock
					//beforeCodeBlock = appendNewLineCodeBlock(beforeCodeBlock, nowCodeBlock, floorCodeBlocks, nowFloor, nowFloor)
				}
			case TCRY_Tab:
				if t == TCRY_Tab {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
				} else {
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
					//addItemCodeBlock(beforeCodeBlock.ParCodeBlock, nowCodeBlock)
					//beforeCodeBlock = nowCodeBlock
					//beforeCodeBlock = appendNewLineCodeBlock(beforeCodeBlock, nowCodeBlock, floorCodeBlocks, nowFloor, nowFloor)
				}
			case TCRY_LeftBracket:
				//addItemCodeBlock(beforeCodeBlock, nowCodeBlock)
				beforeCodeBlock = appendChildCodeBlock(beforeCodeBlock, nowCodeBlock)
			case TCRY_RightBracket:
				beforeCodeBlock = beforeCodeBlock.ParCodeBlock
				beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
				//addItemCodeBlock(beforeCodeBlock.ParCodeBlock, nowCodeBlock)
				//beforeCodeBlock = nowCodeBlock
			case TCRY_Colon:
				beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
			case TCRY_Comma:
				beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
			case TCRY_DH:
				beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
			case TCRY_Semicolon:
				beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
			case TCRY_Period:
				beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
			case TCRY_Operator:
				if t == TCRY_Operator {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
				} else {
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
				}
			case TCRY_Other:
				if t == TCRY_Other {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
				} else {
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
				}
			default:
				beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
			}
		}
		nowPos.ColNo = nowPos.ColNo + nowPos.BlockLen
		if n > 0 {
			//fmt.Println(nowPos, "++", beforeChar, nowChar, string(nowChar))
			if (beforeChar == '\r') && (nowChar == '\n') {
				nowPos.ColNo = 1
				//fmt.Println(nowPos, "+1", beforeChar, nowChar, string(nowChar))
			} else if nowChar == '\n' || nowChar == '\r' {
				//fmt.Println(nowPos, "+2", beforeChar, nowChar)
				nowPos.LineNo = nowPos.LineNo + 1
				nowPos.ColNo = 1
			}
		}
		beforeChar = nowChar
	}
	switch beforeCodeBlock.Type {
	case TCRY_LeftQuotation:
		//codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
	case TCRY_StartNote:
		//codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
		//case TCRY_Other:
		//	codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
	}

	//var enableCodeBlock  *CodeBlock
	//var copyEnableCodeBlock func(block *CodeBlock)( enableCodeBlock  *CodeBlock)
	//copyEnableCodeBlock = func(block *CodeBlock)( enableCodeBlock  *CodeBlock){
	//
	//
	//
	//	return
	//}

	var printCodeBlock func(block *CodeBlock, floor int)
	printCodeBlock = func(block *CodeBlock, floor int) {
		if block.Type == TCRY_Line {
			tab := strings.Repeat("\t", block.Pos.ColNo)
			fmt.Print("\n")
			fmt.Printf("%s行%d，列%d, %d  ", tab, block.Pos.LineNo,
				block.Pos.ColNo, block.Type)
		}
		switch block.Type {
		case TCRY_Space, TCRY_Tab, TCRY_NewLineTab, TCRY_EndLine, TCRY_Line:
		default:
			fmt.Printf("<%s>", block.Words)

		}

		if len(block.Items) > 0 {
			fmt.Print("[")

			for _, item := range block.Items {
				//fmt.Println(n)
				printCodeBlock(item, floor+1)

			}
			fmt.Print("]")

		}
		if block.Type == TCRY_Line {
			//fmt.Print("\n")
		}

	}
	txt.MainCodeBlock = mainCodeBlock
	////printCodeBlock(mainCodeBlock, 0)
	//
	//var CodeBlockToString func(block *CodeBlock) (txt string)
	//
	//CodeBlockToString = func(block *CodeBlock) (txt string) {
	//
	//	var words []string
	//	var w string
	//	switch block.Type {
	//	case TCRY_Space, TCRY_Line:
	//
	//	case TCRY_NewLineTab, TCRY_EndLine:
	//		w = block.Words
	//	case TCRY_LeftBracket:
	//		w = "（"
	//	case TCRY_RightBracket:
	//		w = "）"
	//	case TCRY_Colon:
	//		w = "："
	//	case TCRY_DH:
	//		w = "、"
	//	case TCRY_Comma:
	//		w = "，"
	//	case TCRY_Period:
	//		w = "。"
	//	case TCRY_Semicolon:
	//		w = "；"
	//
	//	default:
	//		if block.NeedSpace {
	//			w = block.Words + " "
	//		} else {
	//			w = block.Words
	//		}
	//
	//	}
	//	words = append(words, w)
	//
	//	for _, item := range block.Items {
	//		w = CodeBlockToString(item)
	//		words = append(words, w)
	//	}
	//	return strings.Join(words, "")
	//}
	//
	//allTxtCode := CodeBlockToString(mainCodeBlock)
	//path := "Zhen\\格式化演示代码2.z1"
	//f, err := os.Create(path)
	//if err != nil {
	//	return
	//}
	//defer f.Close()
	//data := allTxtCode
	//_, err = f.Write([]byte(data))
	//if err != nil {
	//	return
	//}
	return

}

//func (txt *TxtCode) AnalyzeWordsStepA() (err error) {
//	txtCode := txt.Txt
//	startPos := TxtCodeWordPos{WordNo: -1, LineNo: -1, ColNo: -1}
//	isQuotation := false
//	isNote := false
//	isWord := false
//	var allWords []TxtCodeWord
//	addWord := func(pos TxtCodeWordPos, endWordNo int, wordType TxtCodeWordType) {
//		var word TxtCodeWord
//		word.Pos = pos
//		word.Words = txt.Txt[pos.WordNo:endWordNo]
//		word.Type = wordType
//		allWords = append(allWords, word)
//	}
//
//	nowPos := TxtCodeWordPos{WordNo: 0, LineNo: 0, ColNo: 0}
//	for n, r := range txtCode {
//		nowPos.WordNo = n
//		l := len(string(r))
//		var t TxtCodeRuneType
//		switch r {
//		case '“':
//			t = TCRY_LeftQuotation
//		case '”':
//			t = TCRY_RightQuotation
//		case '#':
//			t = TCRY_StartNote
//		case ' ', '　':
//			t = TCRY_Space
//		case '\n', '\r':
//			t = TCRY_EndLine
//		case '(', '（':
//			t = TCRY_LeftBracket
//		case ')', '）':
//			t = TCRY_RightBracket
//		case '\t':
//			t = TCRY_Tab
//		case ':', '：':
//			t = TCRY_Colon
//		case '、':
//			t = TCRY_DH
//		case ',', '，':
//			t = TCRY_Comma
//		case '。':
//			t = TCRY_Period
//		case ';', '；':
//			t = TCRY_Semicolon
//		case '=', '+', '-', '*', '/', '<', '>':
//			t = TCRY_Operator
//			//t = TCRY_Mark
//		default:
//			t = TCRY_Other
//		}
//
//		if isQuotation == false && isNote == false {
//			if isWord == false {
//				if t == TCRY_Other {
//					isWord = true
//					startPos = nowPos
//				}
//			} else if t != TCRY_Other {
//				isWord = false
//				addWord(startPos, n, TCWY_Word)
//			}
//		}
//
//		if t == TCRY_LeftQuotation {
//			isQuotation = true
//			startPos = nowPos
//		} else if t == TCRY_RightQuotation {
//			if isQuotation {
//				isQuotation = false
//				addWord(startPos, n+l, TCWY_String)
//			} else {
//				return errors.New("未找到匹配的引号")
//			}
//		}
//
//		if isQuotation == false {
//			if t == TCRY_StartNote {
//				isNote = true
//				startPos = nowPos
//			} else if t == TCRY_EndLine {
//				if isNote {
//					isNote = false
//					addWord(startPos, n, TCWY_Note)
//				}
//			}
//		}
//
//		if isQuotation == false && isNote == false {
//			switch t {
//			case TCRY_Space:
//				addWord(nowPos, n+l, TCWY_Space)
//			case TCRY_Tab:
//				addWord(nowPos, n+l, TCWY_Tab)
//			case TCRY_EndLine:
//				addWord(nowPos, n+l, TCWY_NewLine)
//			case TCRY_LeftBracket:
//				addWord(nowPos, n+l, TCWY_LeftBracket)
//			case TCRY_RightBracket:
//				addWord(nowPos, n+l, TCWY_RightBracket)
//			case TCRY_Colon:
//				addWord(nowPos, n+l, TCWY_Colon)
//			case TCRY_Comma:
//				addWord(nowPos, n+l, TCWY_Comma)
//			case TCRY_DH:
//				addWord(nowPos, n+l, TCWY_DH)
//			case TCRY_Semicolon:
//				addWord(nowPos, n+l, TCWY_Semicolon)
//			case TCRY_Period:
//				addWord(nowPos, n+l, TCWY_Period)
//			case TCRY_Operator:
//				addWord(nowPos, n+l, TCWY_Operator)
//			}
//
//		}
//
//		nowPos.ColNo = nowPos.ColNo + l
//		if t == TCRY_EndLine {
//			nowPos.LineNo = nowPos.LineNo + 1
//			nowPos.ColNo = 0
//		}
//	}
//	if isQuotation {
//		addWord(startPos, len(txtCode), TCWY_String)
//	} else if isNote {
//		addWord(startPos, len(txtCode), TCWY_Note)
//	} else if isWord {
//		addWord(startPos, len(txtCode), TCWY_Word)
//	}
//
//	txt.Words = allWords
//	return
//}
//
//func (txt *TxtCode) AnalyzeWordsStepB() (err error) {
//	allWords := txt.Words
//	needEndSpace := func(t1 TxtCodeWordType, t2 TxtCodeWordType) (space bool) {
//		switch t1 {
//		case TCWY_NewLine, TCWY_Tab, TCWY_Space:
//			space = false
//		case TCWY_LeftBracket, TCWY_RightBracket:
//			space = false
//		case TCWY_Colon, TCWY_Comma, TCWY_DH, TCWY_Semicolon, TCWY_Period:
//			space = false
//		default:
//			switch t2 {
//			case TCWY_NewLine, TCWY_Tab, TCWY_LeftBracket, TCWY_RightBracket:
//				space = false
//			case TCWY_Colon, TCWY_Comma, TCWY_DH, TCWY_Semicolon, TCWY_Period:
//				space = false
//			default:
//				space = true
//			}
//		}
//		return
//	}
//
//	var words []TxtCodeWord
//	checkAddWord := func(word TxtCodeWord, nextWordType TxtCodeWordType) {
//		word.EndSpace = needEndSpace(word.Type, nextWordType)
//		switch word.Type {
//		case TCWY_Space:
//
//		default:
//			words = append(words, word)
//		}
//
//	}
//
//	allWordsLen := len(allWords)
//	var ww TxtCodeWord
//	if allWordsLen > 0 {
//		ww = allWords[0]
//		for n := 1; n < allWordsLen; n = n + 1 {
//			w := allWords[n]
//
//			switch w.Type {
//			case TCWY_String, TCWY_Note, TCWY_LeftBracket, TCWY_RightBracket:
//
//				checkAddWord(ww, w.Type)
//				ww = w
//			case TCWY_Colon, TCWY_Comma, TCWY_DH, TCWY_Semicolon, TCWY_Period:
//				checkAddWord(ww, w.Type)
//				ww = w
//			case TCWY_Operator:
//				if ww.Type == TCWY_Operator {
//					ww.Words = ww.Words + w.Words
//				} else {
//					checkAddWord(ww, w.Type)
//					ww = w
//				}
//			case TCWY_Tab:
//				if ww.Type == TCWY_Tab {
//					ww.Words = ww.Words + w.Words
//				} else {
//					checkAddWord(ww, w.Type)
//					ww = w
//				}
//			case TCWY_Space:
//				if ww.Type == TCWY_Space {
//					ww.Words = ww.Words + w.Words
//				} else {
//					checkAddWord(ww, w.Type)
//					ww = w
//				}
//			case TCWY_NewLine:
//				if ww.Type == TCWY_NewLine && ww.Words == "\r" && w.Words == "\n" {
//
//					ww.Words = ww.Words + w.Words
//				} else {
//					checkAddWord(ww, w.Type)
//					ww = w
//				}
//			case TCWY_Word:
//				checkAddWord(ww, w.Type)
//				ww = w
//
//				//,TCWY_Word,TCWY_Number,TCWY_KeyWord,TCWY_Other
//
//			}
//
//		}
//
//		checkAddWord(ww, TCWY_NewLine)
//	}
//	txt.Words = words
//	return
//}
//
//func (txt *TxtCode) AnalyzeWordsStepC() (err error) {
//	allWords := txt.Words
//
//	var codeLines []TxtCodeLine
//	var codeLine TxtCodeLine
//
//	floor := 0
//	newLine := true
//	bracketCount := 0
//	for _, word := range allWords {
//
//		switch word.Type {
//		case TCWY_Note:
//			newLine = false
//		case TCWY_LeftBracket:
//			codeLine.Words = append(codeLine.Words, word)
//			newLine = false
//			bracketCount = bracketCount + 1
//		case TCWY_RightBracket:
//			codeLine.Words = append(codeLine.Words, word)
//			newLine = false
//			bracketCount = bracketCount - 1
//		case TCWY_Tab:
//			if newLine {
//				floor = floor + len(word.Words)
//			}
//		case TCWY_NewLine:
//			if bracketCount == 0 {
//				if len(codeLine.Words) > 0 {
//					codeLine.Floor = floor
//					codeLines = append(codeLines, codeLine)
//				}
//				codeLine = TxtCodeLine{}
//				floor = 0
//				newLine = true
//			}
//
//		default:
//			codeLine.Words = append(codeLine.Words, word)
//			newLine = false
//		}
//	}
//	//fmt.Println("************")
//	//for _, line := range codeLines {
//	//	PrintCodeLine(line)
//	//}
//	//fmt.Println("************")
//	//var codeBlocks  map[int]TxtCodeBlock
//	codeBlocks := make(map[int]TxtCodeBlock)
//	parFloors := make(map[int]int)
//	newCodeBlock := func(codeLine TxtCodeLine) (codeBlock TxtCodeBlock) {
//
//		codeBlock.FirstLine = codeLine
//		return
//	}
//	addCodeBlock := func(floor int, codeBlock TxtCodeBlock) {
//		parCodeBlock := codeBlocks[floor]
//		parCodeBlock.ChildLines = append(parCodeBlock.ChildLines, codeBlock)
//		codeBlocks[floor] = parCodeBlock
//	}
//	setCodeBlock := func(floor int, codeBlock TxtCodeBlock) {
//		codeBlocks[floor] = codeBlock
//	}
//	var closeCodeBlock func(floor int, nextFloor int)
//	closeCodeBlock = func(floor int, nextFloor int) {
//		if floor > nextFloor {
//
//			codeBlock := codeBlocks[floor]
//			parFloor := parFloors[floor]
//			addCodeBlock(parFloor, codeBlock)
//			closeCodeBlock(parFloor, nextFloor)
//		}
//
//	}
//
//	if len(codeLines) > 0 {
//		//var codeBlock TxtCodeBlock
//		//codeBlocks[0] = TxtCodeBlock{}
//		for i := 0; i < len(codeLines)-1; i++ {
//			codeBlock := newCodeBlock(codeLines[i])
//			lineFloor := codeLines[i].Floor
//			nextFloor := codeLines[i+1].Floor
//			//var codeBlock TxtCodeBlock
//
//			//fmt.Println(i, lineFloor)
//			if lineFloor == nextFloor {
//				addCodeBlock(lineFloor, codeBlock)
//				//codeBlock := codeBlocks[lineFloor]
//				//codeBlock.ChildLines = append(codeBlock.ChildLines, codeBlock)
//				//codeBlocks[lineFloor] = codeBlock
//			} else if nextFloor > lineFloor {
//				parFloors[nextFloor] = lineFloor
//				setCodeBlock(nextFloor, codeBlock)
//				//codeBlock := TxtCodeBlock{}
//				//codeBlock.FirstLine = codeLine
//				//codeBlocks[nextFloor] = codeBlock
//			} else if nextFloor < lineFloor {
//				addCodeBlock(lineFloor, codeBlock)
//				closeCodeBlock(lineFloor, nextFloor)
//				//codeBlock := codeBlocks[lineFloor]
//				//codeBlock.ChildLines = append(codeBlock.ChildLines, codeBlock)
//				//codeBlocks[lineFloor] = codeBlock
//				//
//				//for f := lineFloor - 1; f >= nextFloor; f-- {
//				//
//				//	parFloor := parFloors[f]
//				//	parCodeBlock := codeBlocks[parFloor]
//				//	parCodeBlock.ChildLines = append(parCodeBlock.ChildLines, codeBlock)
//				//	codeBlocks[parFloor] = parCodeBlock
//				//	codeBlock = codeBlocks[parFloor]
//				//}
//
//				//f := parFloors[floor]
//				//codeBlock.ChildLines = append(codeBlock.ChildLines, codeBlock)
//				//c := codeBlocks[f]
//				//c.ChildLines = append(c.ChildLines, codeBlock)
//				//codeBlocks[nextFloor] = c
//				//
//				//codeBlock = codeBlocks[nextFloor]
//
//			}
//
//		}
//
//		i := len(codeLines) - 1
//		//var b TxtCodeBlock
//		//b.FirstLine = codeLines[i]
//		codeBlock := newCodeBlock(codeLines[i])
//		lineFloor := codeLines[i].Floor
//		//codeBlock := codeBlocks[lineFloor]
//
//		addCodeBlock(lineFloor, codeBlock)
//		closeCodeBlock(lineFloor, 0)
//
//		//codeBlock.ChildLines = append(codeBlock.ChildLines, b)
//		//codeBlocks[lineFloor] = codeBlock
//		//
//		//parFloor := parFloors[lineFloor]
//		//parCodeBlock := codeBlocks[parFloor]
//		//parCodeBlock.ChildLines = append(parCodeBlock.ChildLines, codeBlock)
//		//codeBlocks[parFloor] = parCodeBlock
//
//		//codeBlock.ChildLines = append(codeBlock.ChildLines, b)
//		//codeBlocks[floor] = codeBlock
//
//	}
//	for i, c := range codeBlocks {
//		if i == 0 {
//			fmt.Println(i)
//			PrintCodeBlocksToString(c)
//		}
//
//	}
//	//PrintCodeBlocksToString(codeBlocks[0])
//	return
//}
//func PrintCodeLine(codeline TxtCodeLine) {
//	tab := strings.Repeat("\t", codeline.Floor)
//	var words []string
//	for _, w := range codeline.Words {
//		words = append(words, w.Words)
//	}
//	line := strings.Join(words, " ")
//	fmt.Println(tab, "[", line, "]")
//}
//func PrintCodeBlocksToString(codeBlock TxtCodeBlock) {
//	tab := strings.Repeat("\t", codeBlock.FirstLine.Floor)
//	var words []string
//	for _, w := range codeBlock.FirstLine.Words {
//		words = append(words, w.Words)
//	}
//
//	line := strings.Join(words, " ")
//	fmt.Println(tab, "[", line, "]")
//	if len(codeBlock.ChildLines) > 0 {
//		fmt.Println(tab, "{")
//		for _, c := range codeBlock.ChildLines {
//			PrintCodeBlocksToString(c)
//		}
//		fmt.Println(tab, "}")
//	}
//
//}
func (txt *TxtCode) AnalyzeWords() (err error) {

	err = txt.AnalyzeWordsStep1()
	if err != nil {
		return
	}
	return
	//err = txt.AnalyzeWordsStepA()
	//if err != nil {
	//	return
	//}
	//err = txt.AnalyzeWordsStepB()
	//if err != nil {
	//	return
	//}
	//err = txt.AnalyzeWordsStepC()
	//if err != nil {
	//	return
	//}
	//
	////txt.Words = words
	//return
}
func (txt *TxtCode) CodeBlockToString(block *CodeBlock) string {

	//CodeBlockToString = func(block *CodeBlock) (txt string) {

	var words []string
	var w string
	switch block.Type {
	case TCRY_Space, TCRY_Line:

	case TCRY_NewLineTab, TCRY_EndLine:
		w = block.Words
	case TCRY_LeftBracket:
		w = "（"
	case TCRY_RightBracket:
		w = "）"
	case TCRY_Colon:
		w = "："
	case TCRY_DH:
		w = "、"
	case TCRY_Comma:
		w = "，"
	case TCRY_Period:
		w = "。"
	case TCRY_Semicolon:
		w = "；"

	default:
		if block.NeedSpace {
			w = block.Words + " "
		} else {
			w = block.Words
		}

	}
	words = append(words, w)

	for _, item := range block.Items {
		w = txt.CodeBlockToString(item)
		words = append(words, w)
	}
	return strings.Join(words, "")

	//var words []string
	//for _, word := range txt.Words {
	//	w := word.Words
	//	switch word.Type {
	//	case TCWY_LeftBracket:
	//		w = "（"
	//	case TCWY_RightBracket:
	//		w = "）"
	//	case TCWY_Colon:
	//		w = "："
	//	case TCWY_DH:
	//		w = "、"
	//	case TCWY_Comma:
	//		w = "，"
	//	case TCWY_Period:
	//		w = "。"
	//	case TCWY_Semicolon:
	//		w = "；"
	//	}
	//	if word.Type != TCWY_Space {
	//		if word.EndSpace {
	//			words = append(words, w+" ")
	//		} else {
	//			words = append(words, w)
	//		}
	//	}
	//}
	//
	//return strings.Join(words, "")
}

//保存内容到文件
func (txt *TxtCode) formatToFile(path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	data := txt.CodeBlockToString(txt.MainCodeBlock)
	_, err = f.Write([]byte(data))
	if err != nil {
		return
	}
	return
}
