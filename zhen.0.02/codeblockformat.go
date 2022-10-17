package zhen_0_02

import (
	"os"
	"strings"
)

type CodeBlockFormatParameter struct {
	MaxEmptyLine int
}
type CodeBlockFormat struct {
	CodeFile *CodeFile
	//MainCodeBlock *CodeBlock2
	Parameter CodeBlockFormatParameter

	AllWords []string
	NowLine  int

	Indent            int
	BaseIndent        int
	LastCodeLine      int
	LastLineIndent    int
	LastCodeBlockType CodeBlockType
}

func NewCodeBlockFormat(codeFile *CodeFile) (format CodeBlockFormat) {
	format.CodeFile = codeFile
	format.Parameter = CodeBlockFormatParameter{MaxEmptyLine: 3}
	format.NowLine = 1

	format.LastCodeLine = 1
	format.LastLineIndent = 0
	format.BaseIndent = 4
	format.AnalyseCodeBlock(0)
	return
}

func (format *CodeBlockFormat) NewLine(codeBlockNo int) {
	block := &format.CodeFile.AllCodeBlock[codeBlockNo]

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
	crlfCount := strings.Count(block.Chars, "\r\n")
	crCount := strings.Count(block.Chars, "\r")
	if crCount > crlfCount {
		crlfCount = crCount
	}
	lfCount := strings.Count(block.Chars, "\r")
	if lfCount > crlfCount {
		crlfCount = lfCount
	}
	//if block.Pos.lineCount != crlfCount+1 {
	//	fmt.Println(block.Pos.lineCount, crlfCount+1, format.parser.AllCodeBlock[codeBlockNo].Chars)
	//}

	//format.LastCodeLine = block.Pos.LineNo + block.Pos.lineCount - 1
	format.LastCodeLine = block.Pos.LineNo + crlfCount
	return
}

func (format *CodeBlockFormat) AnalyseCodeBlock(codeBlockNo int) {
	//todo 换行 行首加tab等

	format.NewLine(codeBlockNo)
	block := &format.CodeFile.AllCodeBlock[codeBlockNo]
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
		w = block.Chars
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

	n := block.FirstChildNo
	for n >= 0 {
		format.AnalyseCodeBlock(n)
		n = format.CodeFile.AllCodeBlock[n].NextNo
	}
	//for _, item := range block.items {
	//
	//}
	//if block.BlockType == CbtLine {
	//	format.allWords = append(format.allWords, "\n")
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
func (format *CodeBlockFormat) String() (code string) {

	code = strings.Join(format.AllWords, "")

	return
}

func (format *CodeBlockFormat) ToFile(path string) (err error) {
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
