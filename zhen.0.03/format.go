package zhen_0_03

import (
	"os"
	"strings"
)

type Format struct {
	code     *Code
	allWords []string

	lastLine     int
	lastCol      int
	lastNodeType NodeType

	lineIndent int
	isPoint    bool

	comments    map[int]*Token
	backslashes map[int]*Token

	maxEmptyLine int
}

func NewCodeBlockFormat(code *Code) (format Format) {
	format.code = code
	format.maxEmptyLine = 1

	format.lastLine = 1

	format.comments = make(map[int]*Token)
	for _, c := range format.code.Comments {
		format.comments[c.LineNo] = c
	}

	format.backslashes = make(map[int]*Token)
	for _, b := range format.code.Backslashes {
		format.backslashes[b.LineNo] = b
	}

	format.AnalyseCodeBlock(format.code.RootNode)
	return
}

func (format *Format) checkBackslashes(line int) (nowrap bool) {
	b, ok := format.backslashes[line]
	if ok {
		if b.ColNo > 1 {
			format.allWords = append(format.allWords, " "+b.String())
		} else {
			format.allWords = append(format.allWords, b.String())
		}
		format.lineIndent += 1
		nowrap = true
	}
	return
}

func (format *Format) checkComments(line int) (hasComment bool) {
	c, ok := format.comments[line]
	if ok {
		if c.ColNo > 1 {
			format.allWords = append(format.allWords, " "+c.String())
		} else {
			format.allWords = append(format.allWords, c.String())
		}
		hasComment = true
	}
	return
}

func (format *Format) checkLineComments(line int) (hasComment bool) {
	c, ok := format.comments[line]
	if ok {
		format.allWords = append(format.allWords, "\n")
		if c.ColNo > 1 {
			format.allWords = append(format.allWords, strings.Repeat("\t", format.lineIndent))
		}
		format.allWords = append(format.allWords, c.String())
		hasComment = true
	}
	return
}

func (format *Format) getStringLine(words string) (line int) {
	line = 1
	if strings.Index(words, "\n") >= 0 || strings.Index(words, "\r") >= 0 {
		words = strings.Replace(words, "\r\n", "\n", -1)
		words = strings.Replace(words, "\r", "\n", -1)
		lines := strings.Split(words, "\n")
		line = len(lines)
	}
	return
}

func (format *Format) pushWords(line int, col int, words string, nodeType NodeType) {
	if words == "" || line < format.lastLine {
		return
	}
	nowrap := false
	//fmt.Println(line, format.lastLine, col, format.lastCol, words, nodeType)

	if line > format.lastLine {
		nowrap = format.checkBackslashes(format.lastLine)
		format.checkComments(format.lastLine)
		emptyLine := 0
		for i := format.lastLine + 1; i < line; i++ {
			if format.checkLineComments(i) {
				emptyLine = 0
			} else {
				emptyLine += 1
				if emptyLine <= format.maxEmptyLine && format.lastNodeType != NtColon {
					format.allWords = append(format.allWords, "\n")
				}
			}
		}
		format.allWords = append(format.allWords, "\n")
		format.allWords = append(format.allWords, strings.Repeat("\t", format.lineIndent))
		//format.LastLineIndent = format.lineIndent
		format.lastLine = line
		format.lastCol = 0
		format.lastNodeType = NtUnknown
	}

	if col > format.lastCol {
		if format.isPoint == false {
			switch format.lastNodeType {
			case NtLetter, NtInt, NtFloat, NtString, NtOperator:
				switch nodeType {
				case NtLetter, NtInt, NtFloat, NtString, NtOperator:
					format.allWords = append(format.allWords, " ")
				}
			}
		}

		format.allWords = append(format.allWords, words)
		wordsLine := 1
		if nodeType == NtString {
			wordsLine = format.getStringLine(words)
		}
		format.lastLine += wordsLine - 1
		if wordsLine > 1 {
			format.lastCol = 0
		} else {
			format.lastCol = col
		}
		format.lastNodeType = nodeType
	}
	if nowrap == true {
		format.lineIndent -= 1
	}
	return

}

func (format *Format) AnalyseCodeBlock(node Node) {
	var items []Node
	leftWord := ""
	rightWord := ""
	middleWord := ""
	//needItems := true
	isPoint := false
	leftLine := 0
	leftCol := 0
	rightLine := 0
	rightCol := 0
	nodeType := getNodeType(node)
	switch n := node.(type) {
	case *RootNode:
		leftLine = n.LineNo
		leftCol = n.ColNo
		items = n.Items
	case *LineNode:
		leftLine = n.LineNo
		leftCol = n.ColNo
		if n.LineIndent == 0 {
			format.lineIndent = 0
		}
		items = n.Items
	case *ChildLineNode:
		rightWord = n.Words
		rightLine = n.BaseNode.LineNo
		rightCol = n.BaseNode.ColNo
		rightWord = n.Words
		items = n.Items
	case *ColonNode:
		leftWord = n.Words
		leftLine = n.LineNo
		leftCol = n.ColNo
		items = n.Items
		format.lineIndent += 1
	case *OperatorNode:
		leftLine = n.LineNo
		leftCol = n.ColNo
		if n.OperatorType == OtNegative {
			leftWord = n.Words
		} else if n.OperatorType == OtPoint {
			isPoint = true
			middleWord = n.Words
		} else {
			middleWord = n.Words
		}
		items = n.Items
	case *IntNode:
		leftWord = n.String()
		leftLine = n.LineNo
		leftCol = n.ColNo
	case *FloatNode:
		leftWord = n.String()
		leftLine = n.LineNo
		leftCol = n.ColNo
	case *StringNode:
		leftWord = n.String()
		leftLine = n.LineNo
		leftCol = n.ColNo
	case *LetterNode:
		leftWord = n.String()
		leftLine = n.LineNo
		leftCol = n.ColNo
	case *EmptyNode:
		leftWord = n.String()
		leftLine = n.LineNo
		leftCol = n.ColNo
	case *BracketNode:
		switch n.LeftBracket {
		case TtLeftBracket:
			leftWord = "（"
			rightWord = "）"
		case TtLeftSquareBracket:
			leftWord = "["
			rightWord = "]"
		case TtLeftBigBracket:
			leftWord = "{"
			rightWord = "}"

		}
		leftLine = n.LineNo
		leftCol = n.ColNo
		rightLine = n.EndToken.LineNo
		rightCol = n.EndToken.ColNo
		items = n.Items
	case *OtherNode:
		leftWord = n.String()
		leftLine = n.LineNo
		leftCol = n.ColNo
	}
	format.pushWords(leftLine, leftCol, leftWord, nodeType)
	if nodeType == NtBracket {
		format.lineIndent += 1
	}

	for _, item := range items {
		format.AnalyseCodeBlock(item)
		if isPoint {
			format.isPoint = true
		}
		format.pushWords(leftLine, leftCol, middleWord, nodeType)
	}

	if isPoint {
		format.isPoint = false
	}
	format.pushWords(rightLine, rightCol, rightWord, nodeType)
	if nodeType == NtColon {
		format.lineIndent -= 1
	}
	if nodeType == NtBracket {
		format.lineIndent -= 1
	}
	return

}

func (format *Format) String() (code string) {
	code = strings.Join(format.allWords, "")
	return
}

func (format *Format) ToFile(path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	data := strings.Join(format.allWords, "")
	_, err = f.Write([]byte(data))
	if err != nil {
		return
	}
	return
}
