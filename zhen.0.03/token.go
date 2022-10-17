package zhen_0_03

import "unicode"

type TokenType uint8

const (
	TtUnSet TokenType = iota
	TtApostrophe
	TtLeftApostrophe
	TtRightApostrophe
	TtQuotation
	TtLeftQuotation
	TtRightQuotation
	TtLeftBracket
	TtRightBracket
	TtLeftSquareBracket
	TtRightSquareBracket
	TtLeftBigBracket
	TtRightBigBracket
	TtSpace
	TtFullWidthSpace
	TtTab
	TtPound
	TtCR
	TtLF
	TtCRLF
	TtColon
	TtDunHao
	TtComma
	TtSemicolon
	TtPeriod
	TtBackslash
	TtEqual
	TtAdd
	TtSub
	TtMul
	TtDiv
	TtPower
	TtNegative
	TtEqualEqual
	TtNotEqual
	TtMoreThan
	TtMoreThanEqual
	TtLessThan
	TtLessThanEqual
	TtLetter
	TtInt
	TtFloat
	TtPoint
	TtOtherChar
)

var TokenTypeNames = [...]string{
	TtUnSet:              "未设置",
	TtApostrophe:         "单引号",
	TtLeftApostrophe:     "左单引号",
	TtRightApostrophe:    "右单引号",
	TtQuotation:          "双引号",
	TtLeftQuotation:      "左双引号",
	TtRightQuotation:     "左双引号",
	TtLeftBracket:        "左括号",
	TtRightBracket:       "右括号",
	TtLeftSquareBracket:  "左中括号",
	TtRightSquareBracket: "右中括号",
	TtLeftBigBracket:     "左大括号",
	TtRightBigBracket:    "右大括号",
	TtSpace:              "空格",
	TtFullWidthSpace:     "全角空格",
	TtTab:                "Tab",
	TtPound:              "井号",
	TtCR:                 "回车",
	TtLF:                 "换行",
	TtCRLF:               "回车换行",
	TtColon:              "冒号",
	TtDunHao:             "顿号",
	TtComma:              "逗号",
	TtSemicolon:          "分号",
	TtPeriod:             "句号",
	TtBackslash:          "反斜杠",
	//CbtOperator:           "运算符",
	TtEqual:         "等于",
	TtAdd:           "加",
	TtSub:           "减",
	TtMul:           "乘",
	TtDiv:           "除",
	TtPower:         "幂",
	TtNegative:      "负",
	TtEqualEqual:    "相等",
	TtNotEqual:      "不等于",
	TtMoreThan:      "大于",
	TtMoreThanEqual: "大于等于",
	TtLessThan:      "小于",
	TtLessThanEqual: "小于等于",

	TtLetter:    "标识符",
	TtInt:       "整数",
	TtFloat:     "小数",
	TtPoint:     "点",
	TtOtherChar: "其他字符",
}

func (tt TokenType) String() string {
	return TokenTypeNames[tt]
}

type Token struct {
	LineNo int
	ColNo  int

	TokenType TokenType
	Chars     []rune
}

//type Tokens []*Token

var CharCodeBlockType = map[rune]TokenType{
	'\'': TtApostrophe,
	'‘':  TtLeftApostrophe,
	'’':  TtRightApostrophe,
	'"':  TtQuotation,
	'“':  TtLeftQuotation,
	'”':  TtRightQuotation,
	'#':  TtPound,
	' ':  TtSpace,
	'　':  TtFullWidthSpace,
	'\r': TtCR,
	'\n': TtLF,
	'(':  TtLeftBracket, '（': TtLeftBracket,
	')': TtRightBracket, '）': TtRightBracket,
	'[':  TtLeftSquareBracket,
	']':  TtRightSquareBracket,
	'{':  TtLeftBigBracket,
	'}':  TtRightBigBracket,
	'\t': TtTab,
	':':  TtColon, '：': TtColon,
	'、': TtDunHao,
	',': TtComma, '，': TtComma,
	';': TtSemicolon, '；': TtSemicolon,
	'。':  TtPeriod,
	'\\': TtBackslash,
	'=':  TtEqual,
	'+':  TtAdd,
	'-':  TtSub,
	'*':  TtMul,
	'/':  TtDiv,
	'^':  TtPower,
	'>':  TtMoreThan,
	'≥':  TtMoreThanEqual,
	'<':  TtLessThan,
	'≤':  TtLessThanEqual,
	'≠':  TtNotEqual,
	//'&':  CbtOperator,
	//'|':  CbtOperator,
	//'!':  CbtOperator,
	//'?':  CbtOperator,
	'.': TtPoint,
	'_': TtLetter,
	'@': TtLetter,
}

func getCharType(r rune) (t TokenType) {
	t, ok := CharCodeBlockType[r]
	if !ok {
		if unicode.IsLetter(r) {
			t = TtLetter
		} else if unicode.IsNumber(r) {
			t = TtInt
		} else {
			t = TtOtherChar
		}
	}
	return
}

func NewToken(lineNo int, colNo int, char rune, tokenType TokenType) (t *Token) {
	t = &Token{}
	t.LineNo = lineNo
	t.ColNo = colNo
	t.TokenType = tokenType
	t.Chars = append(t.Chars, char)
	return
}

func (t *Token) AddChar(char rune) {
	t.Chars = append(t.Chars, char)
	return
}

func (t *Token) String() string {
	return string(t.Chars)
}

func getPunctuationWords(tokenType TokenType) (words string) {

	switch tokenType {
	case TtDunHao:
		words = "、"
	case TtComma:
		words = "，"
	case TtSemicolon:
		words = ";"
	case TtPeriod:
		words = "。"
	}
	return
}

func getBracketSymbolWords(tokenType TokenType) (leftBracket string, rightBracket string) {
	switch tokenType {
	case TtLeftBracket:
		leftBracket = "（"
		rightBracket = "）"
	case TtLeftSquareBracket:
		leftBracket = "["
		rightBracket = "]"
	case TtLeftBigBracket:
		leftBracket = "{"
		rightBracket = "}"

	}
	return
}

//func (t *Token) addItem(node Node)   {
//	panic("Token 类型不能增加 item")
//	return
//}

//func (t *Token) String() {
//	t.Chars = append(t.Chars, char)
//	return
//}
