package zhen_0_02

type KeyWordType int

const (
	KwtUnknown KeyWordType = iota
	KwtDefineText
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
