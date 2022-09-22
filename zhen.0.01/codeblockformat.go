package zhen_0_01

import (
	"os"
	"strings"
)

type CodeBlockFormatParameter struct {
	MaxEmptyLine int
}
type CodeBlockFormat struct {
	MainCodeBlock *CodeBlock2
	Parameter     CodeBlockFormatParameter

	AllWords []string
	NowLine  int

	Indent            int
	BaseIndent        int
	LastCodeLine      int
	LastLineIndent    int
	LastCodeBlockType CodeBlockType
}

func NewCodeBlockFormat(codeBlock *CodeBlock2) (format CodeBlockFormat) {
	format.MainCodeBlock = codeBlock
	format.Parameter = CodeBlockFormatParameter{MaxEmptyLine: 3}
	format.NowLine = 1

	format.LastCodeLine = 1
	format.LastLineIndent = 0
	format.BaseIndent = 4
	format.AnalyseCodeBlock(format.MainCodeBlock)
	return
}

func (format *CodeBlockFormat) NewLine(block *CodeBlock2) {
	if block.Pos.LineNo > format.LastCodeLine {
		line := block.Pos.LineNo - format.LastCodeLine
		if line > format.Parameter.MaxEmptyLine {
			line = format.Parameter.MaxEmptyLine
		}
		indent := 0
		if block.BlockType == CbtLine {
			if block.LineIndent > format.LastLineIndent {
				if format.LastLineIndent == 0 {
					format.BaseIndent = block.LineIndent - format.LastLineIndent
				}

				line = 1
				format.Indent += 1
				format.LastLineIndent = block.LineIndent
			} else if block.LineIndent < format.LastLineIndent {
				format.LastLineIndent = block.LineIndent
				format.Indent = format.LastLineIndent / format.BaseIndent
				format.LastLineIndent = block.LineIndent
			}
			indent = format.Indent
		} else {
			indent = format.Indent + 1
		}

		for i := 0; i < line; i++ {
			w := "\n" + strings.Repeat("\t", indent)
			format.AllWords = append(format.AllWords, w)
			format.LastCodeBlockType = CbtLine
		}
	}
	format.LastCodeLine = block.Pos.LineNo + block.Pos.LineCount - 1
	return
}

func (format *CodeBlockFormat) AnalyseCodeBlock(block *CodeBlock2) {
	//todo 换行 行首加tab等

	format.NewLine(block)

	var w string
	switch block.BlockType {
	case CbtSpace, CbtFullWidthSpace, CbtTab:
	case CbtFile, CbtLine, CbtChildLine:
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
		w = block.getChars()
	}

	switch format.LastCodeBlockType {
	case CbtLetter, CbtNumber, CbtString, CbtOperator:
		switch block.BlockType {
		case CbtLetter, CbtNumber, CbtString, CbtOperator, CBtBackslash:
			w = " " + w
		}

	}
	switch block.BlockType {
	case CbtSpace, CbtFullWidthSpace, CbtTab:
	case CbtCR, CbtLF, CbtCRLF, CbtChildLine:
	case CbtComment:
	default:

		format.LastCodeBlockType = block.BlockType
	}
	//format.LastCodeBlockType = block.BlockType
	if block.Comment != "" {
		if w != "" {
			w = w + " "
		}
		w = w + block.Comment
	}

	format.AllWords = append(format.AllWords, w)

	for _, item := range block.items {
		format.AnalyseCodeBlock(item)

	}
	//if block.BlockType == CbtLine {
	//	format.AllWords = append(format.AllWords, "\n")
	//}
	//n	else if block.TrailingEnter > 0 {
	//		words = append(words, "\n")
	//		if block.NextLineIndent > 0 {
	//			w := strings.Repeat(" ", block.NextLineIndent)
	//			words = append(words, w)
	//		}
	//
	//	}
	return

}
func (format *CodeBlockFormat) formatToString() (code string) {

	code = strings.Join(format.AllWords, "")

	return
}

func (format *CodeBlockFormat) FormatToFile(path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	data := strings.Join(format.AllWords, "")
	_, err = f.Write([]byte(data))
	if err != nil {
		return
	}
	return
}
