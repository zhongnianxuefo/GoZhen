package main

import (
	"os"
	"strings"
)

type CodeBlockFormat struct {
	MainCodeBlock *CodeBlock
	AllWords      []string
}

func NewCodeBlockFormat(codeBlock *CodeBlock) (format CodeBlockFormat) {
	format.MainCodeBlock = codeBlock
	format.AnalyseCodeBlock(format.MainCodeBlock)
	return
}

func (format *CodeBlockFormat) AnalyseCodeBlock(block *CodeBlock) {
	//todo 换行 行首加tab等
	//var words []string
	var w string
	switch block.BlockType {
	case CbtSpace:
	case CbtLine:
		w = strings.Repeat(" ", block.LineIndent)
		//fmt.Sprintf("%d,%d", block.NowPos.LineNo, block.NowPos.Width)
	case CbtFile, CbtChildLine:

	case CbtEnter:

		//w = block.Words
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
		//if block.TrailingSpace {
		//	w = w + " "
		//}
	}
	if block.Comment != "" {
		if w != "" {
			w = w + " "
		}
		w = w + block.Comment
	}

	format.AllWords = append(format.AllWords, w)

	for _, item := range block.Items {
		format.AnalyseCodeBlock(item)

	}
	if block.BlockType == CbtLine {
		format.AllWords = append(format.AllWords, "\n")
	}
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

func (format *CodeBlockFormat) formatToFile(path string) (err error) {
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
