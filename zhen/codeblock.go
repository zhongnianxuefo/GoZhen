package zhen

type CodeBlockPos struct {
	StartNo  int
	BlockLen int

	LineNo    int
	LineCount int
	ColNo     int
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

	CbtString
	CbtComment
)

type CodeWordType int

const (
	CwtUnknown CodeWordType = iota
	CwtKeyWord
	CwtTxt
	CwtConstant
)

type CodeBlock struct {
	AllCodeChars  []rune
	Pos           CodeBlockPos
	BlockType     CodeBlockType
	Items         []*CodeBlock
	ParCodeBlock  *CodeBlock
	NextCodeBlock *CodeBlock

	LineIndent int
	Comment    string

	//isKeyWord bool

	Word         string
	WordType     CodeWordType
	codeSteps    []ZhenCodeStep
	nextStepCode *CodeBlock
}

func NewCodeBlock(codeChars []rune, pos CodeBlockPos, codeBlockType CodeBlockType) (codeBlock *CodeBlock) {
	codeBlock = &CodeBlock{}
	codeBlock.AllCodeChars = codeChars
	codeBlock.Pos = pos
	codeBlock.BlockType = codeBlockType
	//codeBlock.LineIndent = 0

	return
}
func (codeBlock *CodeBlock) getChars() string {
	s := codeBlock.Pos.StartNo
	e := s + codeBlock.Pos.BlockLen
	return string(codeBlock.AllCodeChars[s:e])
}

//func (codeBlock *CodeBlock) addLen(addBlockLen int) {
//	codeBlock.Pos.BlockLen += addBlockLen
//}

func (codeBlock *CodeBlock) setEndPos(endItem *CodeBlock) {
	endNo := endItem.Pos.StartNo + endItem.Pos.BlockLen
	if endNo > codeBlock.Pos.StartNo {
		codeBlock.Pos.BlockLen = endNo - codeBlock.Pos.StartNo
	}
	endLineNo := endItem.Pos.LineNo
	if endLineNo > codeBlock.Pos.LineNo {
		codeBlock.Pos.LineCount = endLineNo - codeBlock.Pos.LineNo + 1
	}
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

func (codeBlock *CodeBlock) getNext() (next *CodeBlock, isExist bool) {
	par := codeBlock.ParCodeBlock
	parItemCount := len(par.Items)
	for i, c := range par.Items {
		if c == codeBlock {
			if i < parItemCount-1 {
				next = par.Items[i+1]
				isExist = true
			}
		}
	}
	return
}
