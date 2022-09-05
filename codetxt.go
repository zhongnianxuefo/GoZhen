package main

import (
	"fmt"
	"os"
	"strings"
)

type TxtCodeWordPos struct {
	StartNo  int
	BlockLen int
	LineNo   int
	ColNo    int
}

type CodeBlock struct {
	Pos          TxtCodeWordPos
	Type         TxtCodeRuneType
	NeedSpace    bool
	Words        string
	Items        []*CodeBlock
	ParCodeBlock *CodeBlock
}

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
				}
			case TCRY_EndLine:
				if t == TCRY_EndLine {
					if beforeCodeBlock.Words == "\r" && nowCodeBlock.Words == "\n" {
						codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
					} else {
						beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
					}
				} else if t == TCRY_Tab {
					nowCodeBlock.Type = TCRY_NewLineTab
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
				} else {
					nowFloor = 0
					beforeCodeBlock = appendNewLineCodeBlock(beforeCodeBlock, nowCodeBlock, floorCodeBlocks, beforeFloor, nowFloor)
					beforeFloor = nowFloor
				}
			case TCRY_NewLineTab:
				if t == TCRY_Tab {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
				} else if t == TCRY_EndLine {
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
				} else {
					nowFloor = beforeCodeBlock.Pos.BlockLen
					beforeCodeBlock = appendNewLineCodeBlock(beforeCodeBlock, nowCodeBlock, floorCodeBlocks, beforeFloor, nowFloor)
					beforeFloor = nowFloor
				}
			case TCRY_Space:
				if t == TCRY_Space {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
				} else {
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
				}
			case TCRY_Tab:
				if t == TCRY_Tab {
					codeBlockAddBlockLen(beforeCodeBlock, nowPos.BlockLen)
				} else {
					beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
				}
			case TCRY_LeftBracket:
				beforeCodeBlock = appendChildCodeBlock(beforeCodeBlock, nowCodeBlock)
			case TCRY_RightBracket:
				beforeCodeBlock = beforeCodeBlock.ParCodeBlock
				beforeCodeBlock = appendCodeBlock(beforeCodeBlock, nowCodeBlock)
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
			if (beforeChar == '\r') && (nowChar == '\n') {
				nowPos.ColNo = 1
			} else if nowChar == '\n' || nowChar == '\r' {
				nowPos.LineNo = nowPos.LineNo + 1
				nowPos.ColNo = 1
			}
		}
		beforeChar = nowChar
	}

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

	return

}

func (txt *TxtCode) AnalyzeWords() (err error) {

	err = txt.AnalyzeWordsStep1()
	if err != nil {
		return
	}
	return

}
func (txt *TxtCode) CodeBlockToString(block *CodeBlock) string {

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
