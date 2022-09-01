package main

import (
	"errors"
	"os"
	"strings"
)

type TxtCodeWordPos struct {
	WordNo int
	LineNo int
	ColNo  int
}

type TxtCodeWord struct {
	Pos   TxtCodeWordPos
	Words string
}

type TxtCodeBlock struct {
	LineNo int
	Blocks []TxtCodeBlock
}
type TxtCode struct {
	Txt   string
	Words []TxtCodeWord
}

func newTxtCode(codes string) (txtCode TxtCode, err error) {
	txtCode.Txt = strings.ReplaceAll(codes, "\r\n", "\n")
	err = txtCode.AnalyzeWords()
	if err != nil {
		return
	}

	return
}

func (txt *TxtCode) addWord(pos TxtCodeWordPos, endWordNo int) {
	var word TxtCodeWord
	word.Pos = pos
	word.Words = txt.Txt[pos.WordNo:endWordNo]
	txt.Words = append(txt.Words, word)

}
func (txt *TxtCode) AnalyzeWords() (err error) {
	t := txt.Txt

	quotationPos := TxtCodeWordPos{WordNo: -1, LineNo: -1, ColNo: -1}
	bracketsPos := TxtCodeWordPos{WordNo: -1, LineNo: -1, ColNo: -1}
	wordPos := TxtCodeWordPos{WordNo: -1, LineNo: -1, ColNo: -1}
	bracketsCount := 0

	pos := TxtCodeWordPos{WordNo: 0, LineNo: 0, ColNo: 0}

	for n, w := range t {
		pos.WordNo = n
		l := len(string(w))

		var newLine = false
		var startQuotation = false
		var endQuotation = false
		var startBrackets = false
		var endBrackets = false
		var startWord = false
		var endWord = false
		var isMark = false

		switch w {
		case '“':
			startQuotation = true
			endWord = true
		case '”':
			endQuotation = true
			endWord = true
		case '(', '（':
			startBrackets = true
			endWord = true
		case ')', '）':
			endBrackets = true
			endWord = true
		case ' ':
			endWord = true
		case '=', '+', '-', '*', '/', '<', '>', '\\',
			',', ';', ':', '，', '。', '、', '；', '：',
			'\t':
			endWord = true
			isMark = true
		case '\n':
			newLine = true
			endWord = true
			isMark = true
		default:
			startWord = true
		}
		if startQuotation {
			if endWord && wordPos.WordNo >= 0 {
				txt.addWord(wordPos, n)
				wordPos.WordNo = -1
			}
			quotationPos = pos
		} else if endQuotation {
			if quotationPos.WordNo >= 0 {
				txt.addWord(quotationPos, n+l)
				quotationPos.WordNo = -1
			} else {
				return errors.New("未找到匹配的引号")
			}
		} else if quotationPos.WordNo == -1 {
			if startBrackets {
				if bracketsCount == 0 {
					if endWord && wordPos.WordNo >= 0 {
						txt.addWord(wordPos, n)
						wordPos.WordNo = -1
					}
					bracketsPos = pos
				}
				bracketsCount = bracketsCount + 1
			} else if endBrackets {
				if bracketsCount == 1 {
					txt.addWord(bracketsPos, n+l)
					bracketsPos.WordNo = -1
				}
				bracketsCount = bracketsCount - 1
			} else if bracketsCount == 0 {
				if startWord {
					if wordPos.WordNo == -1 {
						wordPos = pos
					}
				} else if endWord {
					if wordPos.WordNo >= 0 {
						txt.addWord(wordPos, n)
						wordPos.WordNo = -1
					}
				}
				if isMark {
					txt.addWord(pos, n+l)
				}
			}
		}
		pos.ColNo = pos.ColNo + l
		if newLine {
			pos.LineNo = pos.LineNo + 1
			pos.ColNo = 0
		}
	}

	if quotationPos.WordNo >= 0 {
		return errors.New("未找到匹配的引号")
	} else if bracketsCount > 0 {
		return errors.New("未找到匹配的括号")
	} else if wordPos.WordNo >= 0 {
		txt.addWord(wordPos, len(t))
	}

	return
}
func (txt *TxtCode) formatCode() string {
	var lines []string
	var words []string

	for _, word := range txt.Words {
		w := word.Words
		switch w {
		case "\n":
			line := strings.TrimRight(strings.Join(words, ""), " \t")
			lines = append(lines, line)
			words = []string{}
		case "\t":
			words = append(words, w)
		default:
			words = append(words, w+" ")

		}

	}
	return strings.Join(lines, "\n")
}

//保存内容到文件
func (txt *TxtCode) formatToFile(path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	data := txt.formatCode()
	_, err = f.Write([]byte(data))
	if err != nil {
		return
	}
	return
}
