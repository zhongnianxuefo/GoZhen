package main

type CodeBlockPos struct {
	StartNo  int
	BlockLen int

	LineNo int
	ColNo  int
}

type CodeBlockType int

const (
	_ CodeBlockType = iota
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
	CbtEnter
	CbtColon
	CbtComma
	CbtDunHao
	CbtSemicolon
	CbtPeriod
	CBtBackslash
	CbtOperator
	CbtUnderscore
	CbtLetter
	CbtNumber
	CbtPoint
	CbtOtherChar

	CbtFile
	CbtLine
	CbtChildLine

	CbtString
	CbtComment
)

type CodeBlock struct {
	AllCodeChars []rune
	Pos          CodeBlockPos
	BlockType    CodeBlockType
	Items        []*CodeBlock
	ParCodeBlock *CodeBlock

	LineIndent int
	Comment    string
}

func NewCodeBlock(codeChars []rune, pos CodeBlockPos, codeBlockType CodeBlockType) (codeBlock *CodeBlock) {
	codeBlock = &CodeBlock{}
	codeBlock.AllCodeChars = codeChars
	codeBlock.Pos = pos
	codeBlock.BlockType = codeBlockType
	return
}
func (codeBlock *CodeBlock) getChars() string {
	s := codeBlock.Pos.StartNo
	e := s + codeBlock.Pos.BlockLen
	return string(codeBlock.AllCodeChars[s:e])
}

func (codeBlock *CodeBlock) addLen(addBlockLen int) {
	codeBlock.Pos.BlockLen += addBlockLen
}

func (codeBlock *CodeBlock) addItem(item *CodeBlock) {
	codeBlock.Items = append(codeBlock.Items, item)
	item.ParCodeBlock = codeBlock
}

func (codeBlock *CodeBlock) appendNext(nextCodeBlock *CodeBlock) *CodeBlock {
	codeBlock.ParCodeBlock.addItem(nextCodeBlock)
	return nextCodeBlock
}

func (codeBlock *CodeBlock) appendChild(child *CodeBlock) *CodeBlock {
	codeBlock.addItem(child)
	return child
}
