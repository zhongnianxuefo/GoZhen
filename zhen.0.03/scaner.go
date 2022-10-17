package zhen_0_03

import "strings"

type Scaner struct {
	code *Code

	chars    []rune
	charNo   int
	lineNo   int
	colNo    int
	char     rune
	charType TokenType
	scanFun  ScanFun
	nowToken *Token
}

func NewScan(code *Code) (scaner *Scaner) {
	scaner = &Scaner{}
	scaner.code = code
	scaner.chars = []rune(code.txt)
	scaner.lineNo = 1
	scaner.colNo = 1
	scaner.charNo = 0

	return
}
func (s *Scaner) Scan() (err error) {
	s.charNo = 0
	s.scanFun = scanNew
	for s.charNo < len(s.chars) {
		s.scanFun = s.scanFun(s)
	}

	return
}

func (s *Scaner) newToken() {
	s.nowToken = NewToken(s.lineNo, s.colNo, s.char, s.charType)
	s.code.Tokens = append(s.code.Tokens, s.nowToken)
	s.charNo += 1
	s.colNo += 1
	return
}

func (s *Scaner) addTokenChar() {
	s.nowToken.AddChar(s.char)
	s.charNo += 1
	s.colNo += 1
	return
}

func (s *Scaner) addStringTokenEndChar() {

	s.nowToken.AddChar(s.char)
	words := string(s.nowToken.Chars)
	if strings.Index(words, "\n") >= 0 || strings.Index(words, "\r") >= 0 {
		words = strings.Replace(words, "\r\n", "\n", -1)
		words = strings.Replace(words, "\r", "\n", -1)
		lines := strings.Split(words, "\n")

		s.lineNo = s.nowToken.LineNo + len(lines) - 1
		lastLine := lines[len(lines)-1]
		s.colNo = len([]rune(lastLine)) + 1
		//fmt.Println(len(lines), s.lineNo, s.colNo)
	} else {

		s.colNo += 1
	}
	s.charNo += 1
	return
}

func (s *Scaner) changeTokenType(tokenType TokenType) {
	s.nowToken.TokenType = tokenType
	return
}

func (s *Scaner) checkChar() {
	s.char = s.chars[s.charNo]
	s.charType = getCharType(s.char)
	return
}
