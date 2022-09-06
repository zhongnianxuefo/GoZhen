package main

import (
	"fmt"
	"github.com/beevik/etree"
	"os"
	"strconv"
	"strings"
)

type CodeBlockPos struct {
	StartNo  int
	BlockLen int

	LineNo int
	ColNo  int
}

type CodeBlockType int

const (
	_ CodeBlockType = iota
	CbtLeftQuotation
	CbtRightQuotation
	CbtLeftBracket
	CbtRightBracket
	CbtSpace
	CbtTab
	CbtPound
	CbtEnter
	CbtColon
	CbtComma
	CbtDunHao
	CbtSemicolon
	CbtPeriod
	CbtOperator

	TCRY_Mark
	CbtOtherChar

	CbtFile
	CbtLine

	CbtString
	CbtNewLineTab
	CbtComment
)

type CodeBlock struct {
	Pos               CodeBlockPos
	BlockType         CodeBlockType
	NeedTrailingSpace bool
	Words             string
	Items             []*CodeBlock
	ParCodeBlock      *CodeBlock
}
type CodeBlockIndentation struct {
	CodeBlocks map[int]*CodeBlock
	Floor      int
}

type TxtCode struct {
	CodeTxt       string
	MainCodeBlock *CodeBlock
}

func getRuneCodeBlockType(r rune) (t CodeBlockType) {

	switch r {
	case '“':
		t = CbtLeftQuotation
	case '”':
		t = CbtRightQuotation
	case '#':
		t = CbtPound
	case ' ', '　':
		t = CbtSpace
	case '\n', '\r':
		t = CbtEnter
	case '(', '（':
		t = CbtLeftBracket
	case ')', '）':
		t = CbtRightBracket
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
	case '=', '+', '-', '*', '/', '<', '>', '!':
		t = CbtOperator
		//t = TCRY_Mark
	default:
		t = CbtOtherChar
	}
	return
}

func newCodeBlockPointer(txtCode string, pos CodeBlockPos, codeBlockType CodeBlockType) (codeBlock *CodeBlock) {
	codeBlock = &CodeBlock{}
	codeBlock.Pos = pos
	codeBlock.BlockType = codeBlockType
	codeBlock.Words = txtCode[pos.StartNo : pos.StartNo+pos.BlockLen]
	return
}

func newLineCodeBlock(pos CodeBlockPos) (codeBlock *CodeBlock) {
	codeBlock = &CodeBlock{}
	codeBlock.Pos = pos
	codeBlock.BlockType = CbtLine
	codeBlock.Words = ""
	return
}
func newFileCodeBlock() (codeBlock *CodeBlock) {
	codeBlock = &CodeBlock{}
	codeBlock.Pos = CodeBlockPos{StartNo: 0, BlockLen: 0, LineNo: 1, ColNo: 1}
	codeBlock.BlockType = CbtFile
	codeBlock.Words = ""
	codeBlock.ParCodeBlock = codeBlock
	return
}

func needTrailingSpace(t1 CodeBlockType, t2 CodeBlockType) (space bool) {
	switch t1 {
	case CbtLine, CbtNewLineTab, CbtEnter, CbtTab, CbtSpace:
		space = false
	case CbtLeftBracket, CbtRightBracket:
		space = false
	case CbtColon, CbtComma, CbtDunHao, CbtSemicolon, CbtPeriod:
		space = false
	default:
		switch t2 {
		case CbtLine, CbtNewLineTab, CbtEnter, CbtTab, CbtLeftBracket, CbtRightBracket:
			space = false
		case CbtColon, CbtComma, CbtDunHao, CbtSemicolon, CbtPeriod:
			space = false
		default:
			space = true
		}
	}
	return
}

func (codeBlock *CodeBlock) addBlockLen(txtCode string, addBlockLen int) {
	codeBlock.Pos.BlockLen += addBlockLen
	startNo := codeBlock.Pos.StartNo
	endNo := codeBlock.Pos.StartNo + codeBlock.Pos.BlockLen
	codeBlock.Words = txtCode[startNo:endNo]
}

func (codeBlock *CodeBlock) addItem(item *CodeBlock) {
	codeBlock.Items = append(codeBlock.Items, item)
	item.ParCodeBlock = codeBlock
}

func (codeBlock *CodeBlock) appendNext(nextCodeBlock *CodeBlock) *CodeBlock {

	codeBlock.NeedTrailingSpace = needTrailingSpace(codeBlock.BlockType, nextCodeBlock.BlockType)
	codeBlock.ParCodeBlock.addItem(nextCodeBlock)

	return nextCodeBlock
}

func (codeBlock *CodeBlock) appendChild(child *CodeBlock) *CodeBlock {
	codeBlock.addItem(child)
	return child
}

func newCodeBlockIndentation(mainCodeBlock *CodeBlock, codeBlock *CodeBlock) (codeBlockIndentation CodeBlockIndentation) {

	lineCodeBlock := newLineCodeBlock(codeBlock.Pos)
	mainCodeBlock.addItem(lineCodeBlock)
	lineCodeBlock.addItem(codeBlock)
	codeBlockIndentation.CodeBlocks = make(map[int]*CodeBlock)
	codeBlockIndentation.CodeBlocks[0] = lineCodeBlock

	return
}
func (codeBlockLines *CodeBlockIndentation) appendNewLine(floor int, codeBlock *CodeBlock) *CodeBlock {
	lineCodeBlock := newLineCodeBlock(codeBlock.Pos)
	if floor > codeBlockLines.Floor {
		codeBlockLines.CodeBlocks[codeBlockLines.Floor].addItem(lineCodeBlock)
		codeBlockLines.CodeBlocks[floor] = lineCodeBlock
		lineCodeBlock.addItem(codeBlock)
	} else {
		codeBlockLines.CodeBlocks[floor].ParCodeBlock.addItem(lineCodeBlock)
		codeBlockLines.CodeBlocks[floor] = lineCodeBlock
		lineCodeBlock.addItem(codeBlock)
	}
	return codeBlock
}

func newTxtCode(codes string) (txtCode TxtCode, err error) {
	txtCode.CodeTxt = codes //strings.ReplaceAll(codes, "\r\n", "\n")
	err = txtCode.AnalyzeWords()
	if err != nil {
		return
	}

	return
}
func getCodeBlockFromTxt(txtCode string) (mainCodeBlock *CodeBlock, err error) {

	mainCodeBlock = newFileCodeBlock()

	//var lineCodeBlock *CodeBlock
	var beforeCodeBlock *CodeBlock
	var nowCodeBlock *CodeBlock
	var beforeChar int32
	var nowChar int32

	var codeBlockIndentation CodeBlockIndentation

	//beforeFloor := 0
	//nowFloor := 0

	nowPos := CodeBlockPos{StartNo: 0, BlockLen: 0, LineNo: 1, ColNo: 1}
	for n, r := range txtCode {
		nowChar = r
		nowPos.StartNo = n
		nowPos.BlockLen = len(string(r))

		t := getRuneCodeBlockType(r)

		nowCodeBlock = newCodeBlockPointer(txtCode, nowPos, t)
		if n == 0 {

			codeBlockIndentation = newCodeBlockIndentation(mainCodeBlock, nowCodeBlock)
			beforeCodeBlock = nowCodeBlock

			//codeBlockIndentation[nowFloor] = lineCodeBlock
			//beforeFloor = nowFloor
		} else {
			switch beforeCodeBlock.BlockType {
			case CbtLeftQuotation:
				if t != CbtRightQuotation {
					beforeCodeBlock.addBlockLen(txtCode, nowPos.BlockLen)

				} else {
					beforeCodeBlock.addBlockLen(txtCode, nowPos.BlockLen)
					beforeCodeBlock.BlockType = CbtString
				}
			case CbtPound:
				if t != CbtEnter {
					beforeCodeBlock.addBlockLen(txtCode, nowPos.BlockLen)
				} else {
					beforeCodeBlock.BlockType = CbtComment
					beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
				}
			case CbtEnter:
				if t == CbtEnter {
					if beforeCodeBlock.Words == "\r" && nowCodeBlock.Words == "\n" {
						beforeCodeBlock.addBlockLen(txtCode, nowPos.BlockLen)
					} else {
						beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
					}
				} else if t == CbtTab {
					nowCodeBlock.BlockType = CbtNewLineTab
					beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
				} else {
					beforeCodeBlock = codeBlockIndentation.appendNewLine(0, nowCodeBlock)
				}
			case CbtNewLineTab:
				if t == CbtTab {
					beforeCodeBlock.addBlockLen(txtCode, nowPos.BlockLen)
				} else if t == CbtEnter {
					beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
				} else {
					beforeCodeBlock = codeBlockIndentation.appendNewLine(beforeCodeBlock.Pos.BlockLen, nowCodeBlock)
					//nowFloor = beforeCodeBlock.Pos.BlockLen
					//beforeCodeBlock = appendNewLineCodeBlock(beforeCodeBlock, nowCodeBlock, codeBlockIndentation, beforeFloor, nowFloor)
					//beforeFloor = nowFloor
				}
			case CbtSpace:
				if t == CbtSpace {
					beforeCodeBlock.addBlockLen(txtCode, nowPos.BlockLen)
				} else {
					beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
				}
			case CbtTab:
				if t == CbtTab {
					beforeCodeBlock.addBlockLen(txtCode, nowPos.BlockLen)
				} else {
					beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
				}
			case CbtLeftBracket:
				beforeCodeBlock = beforeCodeBlock.appendChild(nowCodeBlock)
			case CbtRightBracket:
				beforeCodeBlock = beforeCodeBlock.ParCodeBlock
				beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
			case CbtColon:
				beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
			case CbtComma:
				beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
			case CbtDunHao:
				beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
			case CbtSemicolon:
				beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
			case CbtPeriod:
				beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
			case CbtOperator:
				if t == CbtOperator {
					beforeCodeBlock.addBlockLen(txtCode, nowPos.BlockLen)
				} else {
					beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
				}
			case CbtOtherChar:
				if t == CbtOtherChar {
					beforeCodeBlock.addBlockLen(txtCode, nowPos.BlockLen)
				} else {
					beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
				}
			default:
				beforeCodeBlock = beforeCodeBlock.appendNext(nowCodeBlock)
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
		if block.BlockType == CbtLine {
			tab := strings.Repeat("\t", block.Pos.ColNo)
			fmt.Print("\n")
			fmt.Printf("%s行%d，列%d, %d  ", tab, block.Pos.LineNo,
				block.Pos.ColNo, block.BlockType)
		}
		switch block.BlockType {
		case CbtSpace, CbtTab, CbtNewLineTab, CbtEnter, CbtLine:
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
		if block.BlockType == CbtLine {
			//fmt.Print("\n")
		}

	}

	return
}

func (code *TxtCode) AnalyzeWords() (err error) {
	code.MainCodeBlock, err = getCodeBlockFromTxt(code.CodeTxt)
	if err != nil {
		return
	}

	return

}

func XmlElementAddCodeBlock(parElement *etree.Element, codeBlock *CodeBlock) {

	//doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)
	name := ""
	switch codeBlock.BlockType {
	case CbtFile:
		name = "程序"
	case CbtLine:
		name = "代码"
	default:
		name = "未知" + strconv.Itoa(int(codeBlock.BlockType))

	}
	element := parElement.CreateElement(name)
	element.SetText(codeBlock.Words)
	for _, c := range codeBlock.Items {
		XmlElementAddCodeBlock(element, c)
	}

	return
}
func (code *TxtCode) ToXmlFile(path string) (err error) {

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	element := doc.CreateElement("程序")
	for _, codeBlock := range code.MainCodeBlock.Items {
		XmlElementAddCodeBlock(element, codeBlock)
	}

	doc.Indent(2)
	err = doc.WriteToFile(path)

	if err != nil {
		return
	}

	return

}

func CodeBlockToString(block *CodeBlock) string {

	var words []string
	var w string
	switch block.BlockType {
	case CbtSpace, CbtLine:

	case CbtNewLineTab, CbtEnter:
		w = block.Words
	case CbtLeftBracket:
		w = "（"
	case CbtRightBracket:
		w = "）"
	case CbtColon:
		w = "："
	case CbtDunHao:
		w = "、"
	case CbtComma:
		w = "，"
	case CbtPeriod:
		w = "。"
	case CbtSemicolon:
		w = "；"

	default:
		if block.NeedTrailingSpace {
			w = block.Words + " "
		} else {
			w = block.Words
		}

	}
	words = append(words, w)

	for _, item := range block.Items {
		w = CodeBlockToString(item)
		words = append(words, w)
	}
	return strings.Join(words, "")

}

func (code *TxtCode) formatToFile(path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	data := CodeBlockToString(code.MainCodeBlock)
	_, err = f.Write([]byte(data))
	if err != nil {
		return
	}
	return
}
