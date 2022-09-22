package zhen_0_02

type CodeBlockPos struct {
	startNo  int
	blockLen int

	LineNo    int
	LineCount int
	ColNo     int
}

type CodeBlockType uint8

const (
	CbtUnSet CodeBlockType = iota
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
	CbtClass
	CbtFun

	CbtString
	CbtComment
)

var CodeBlockTypeNames = [...]string{
	CbtUnSet:              "未设置",
	CbtApostrophe:         "单引号",
	CbtLeftApostrophe:     "左单引号",
	CbtRightApostrophe:    "右单引号",
	CbtQuotation:          "双引号",
	CbtLeftQuotation:      "左双引号",
	CbtRightQuotation:     "左双引号",
	CbtLeftBracket:        "左括号",
	CbtRightBracket:       "右括号",
	CbtLeftSquareBracket:  "左中括号",
	CbtRightSquareBracket: "右中括号",
	CbtLeftBigBracket:     "左大括号",
	CbtRightBigBracket:    "右大括号",
	CbtSpace:              "空格",
	CbtFullWidthSpace:     "全角空格",
	CbtTab:                "Tab",
	CbtPound:              "井号",
	CbtCR:                 "回车",
	CbtLF:                 "换行",
	CbtCRLF:               "回车换行",
	CbtColon:              "冒号",
	CbtComma:              "逗号",
	CbtDunHao:             "顿号",
	CbtSemicolon:          "分号",
	CbtPeriod:             "句号",
	CBtBackslash:          "反斜杠",
	CbtOperator:           "运算符",
	CbtLetter:             "标识符",
	CbtNumber:             "数字",
	CbtPoint:              "点",
	CbtOtherChar:          "其他字符",
	CbtFile:               "文件",
	CbtLine:               "代码行",
	CbtChildLine:          "子代码行",
	CbtClass:              "类",
	CbtFun:                "函数",
	CbtString:             "字符串",
	CbtComment:            "注释",
}

func (cbt CodeBlockType) String() string {
	return CodeBlockTypeNames[cbt]
}

type CodeWordType uint8

const (
	CwtUnSet CodeWordType = iota
	CwtKeyWord
	CwtTxt

	CwtConstant
	CwtVar
	CwtGlobalVar
	CwtLocalVar

	CwtFunName
	CwtFunPara
	CwtFunReturn
)

var CodeWordTypeNames = [...]string{
	CwtUnSet:     "未设置",
	CwtKeyWord:   "关键字",
	CwtTxt:       "文本",
	CwtConstant:  "常量",
	CwtVar:       "变量",
	CwtGlobalVar: "全局变量",
	CwtLocalVar:  "局部变量",
	CwtFunName:   "函数名",
	CwtFunPara:   "参数",
	CwtFunReturn: "返回值",
}

func (cwt CodeWordType) String() string {
	return CodeWordTypeNames[cwt]
}

var CharCodeBlockType = map[rune]CodeBlockType{
	'\'': CbtApostrophe,
	'‘':  CbtLeftApostrophe,
	'’':  CbtRightApostrophe,
	'"':  CbtQuotation,
	'“':  CbtLeftQuotation,
	'”':  CbtRightQuotation,
	'#':  CbtPound,
	' ':  CbtSpace,
	'　':  CbtFullWidthSpace,
	'\r': CbtCR,
	'\n': CbtLF,
	'(':  CbtLeftBracket, '（': CbtLeftBracket,
	')': CbtRightBracket, '）': CbtRightBracket,
	'[':  CbtLeftSquareBracket,
	']':  CbtRightSquareBracket,
	'{':  CbtLeftBigBracket,
	'}':  CbtRightBigBracket,
	'\t': CbtTab,
	':':  CbtColon, '：': CbtColon,
	'、': CbtDunHao,
	',': CbtComma, '，': CbtComma,
	'。': CbtPeriod,
	';': CbtSemicolon, '；': CbtSemicolon,
	'\\': CBtBackslash,
	'=':  CbtOperator,
	'+':  CbtOperator,
	'-':  CbtOperator,
	'*':  CbtOperator,
	'/':  CbtOperator,
	'>':  CbtOperator,
	'<':  CbtOperator,
	'&':  CbtOperator,
	'|':  CbtOperator,
	'!':  CbtOperator,
	'^':  CbtOperator,
	'?':  CbtOperator,
	'.':  CbtPoint,
	'_':  CbtLetter,
	'@':  CbtLetter,
}

type CodeBlock struct {
	Pos       CodeBlockPos
	BlockType CodeBlockType
	Chars     string

	//No           int
	ParNo        int
	NextNo       int
	FirstChildNo int
	LastChildNo  int

	Words    string
	WordType CodeWordType
	Operator Operator

	LineIndent int
	Comment    string

	//Vars  CodeVars
	Steps []CodeStep
}
type CodeStepType int

const (
	CstNone CodeStepType = iota
	CstDefineText
	CstDefineVar
	CstDefineLocalVar
	CstDefineGlobalVar
	CstDefineConstant
	CstDefineFun
	CstDefineFunPara
	CstDefineFunReturn

	CstAs
	CstAdd
	CstSub
	CstMul
	CstDiv
	CstPoint
)

var CodeStepTypeNames = [...]string{
	CstNone:            "未设置",
	CstDefineText:      "定义文本",
	CstDefineVar:       "定义变量",
	CstDefineLocalVar:  "定义局部变量",
	CstDefineGlobalVar: "定义全局变量",
	CstDefineConstant:  "定义常量",
	CstDefineFun:       "定义函数",
	CstDefineFunPara:   "定义函数参数",
	CstDefineFunReturn: "定义函数返回值",

	CstAs:    "=",
	CstAdd:   "+",
	CstSub:   "-",
	CstMul:   "*",
	CstDiv:   "/",
	CstPoint: ".",
}

func (cst CodeStepType) String() string {
	return CodeStepTypeNames[cst]
}

type CodeStep struct {
	CodeStepType CodeStepType
	TempVarNo1   int
	VarName1     string
	TempVarNo2   int
	VarName2     string
	ReturnVarNo  int
	ValueString  string
}
