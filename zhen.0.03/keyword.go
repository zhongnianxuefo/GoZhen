package zhen_0_03

type KeyWordType int8

const (
	KwtUnknown KeyWordType = iota
	KwtDefine
	KwtDefineConstant
	KwtDefineVar
	KwtDefineGlobalVar
	KwtDefineLocalVar
	KwtDefineFun
	KwtDefineFunPara
	KwtDefineFunReturn

	KwtIf
	KwtElse
	KwtElseIf

	KwtWhile
	KwtFor

	KwtCallFun
)

func NewKeyWord(name string, wordType KeyWordType) (keyWord *Letter) {
	keyWord = NewLetter(name, LtKeyWord)
	keyWord.Data = wordType
	return
}
