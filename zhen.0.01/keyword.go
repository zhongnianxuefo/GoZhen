package zhen_0_01

import (
	"errors"
)

type KeyWord struct {
	Name string
	//Type   KeyWordType
	PreFun func(*CodePre) (err error)
}

type KeyWordType int

const (
	KwtUnknown KeyWordType = iota
	KwtDefineText
	KwtDefineConstant
	KwtDefineVar
	KwtDefineGlobalVar
	KwtDefineFun
	KwtDefineFunPara
	KwtDefineFunReturn

	KwtText
	KwtConstant
	KwtVar
	KwtGlobalVar
	KwtFun
	KwtFunPara
	KwtFunReturn

	KwtIf
	KwtElse
	KwtElseIf

	KwtWhile
	KwtFor

	KwtCallFun
)

//func StringToKeyWordType(s string) (t KeyWordType) {
//	switch s {
//	case "文本":
//		t = KwtText
//	case "常量":
//		t = KwtConstant
//	default:
//		t = KwtUnknown
//	}
//	return
//}

//func KeyWordTypeToString(t KeyWordType) (s string) {
//	switch t {
//	case KwtText:
//		s = "文本"
//	case KwtConstant:
//		s = "常量"
//	default:
//		s = "未知"
//	}
//	return
//}

func KeyWordToZhenValue(keyWord KeyWord) (value ZValue) {

	//table := make(ZhenValueTable)
	//table["名称"] = StringToZhenValue(keyWord.Name)
	//table["类型"] = StringToZhenValue(KeyWordTypeToString(keyWord.Type))
	//table["预处理函数"] = NewZhenValueFunction(keyWord.PreFun)
	//value = NewZhenValueTable(table)
	return ZValue(keyWord)
}

func ZhenValueToKeyWord(value ZValue) (keyWord KeyWord) {
	k, ok := value.(KeyWord)
	if ok {
		keyWord = k
	}
	return
}
func ZValueToCodeBlock(value ZValue) (codeBlock *CodeBlock2, err error) {
	codeBlock, ok := value.(*CodeBlock2)
	if !ok {
		err = errors.New("预处理函数传入参数类似错误")
		return
	}
	return
}

//
//func ConstantStatement(zhen *ZhenState, block *CodeBlock2) (err error) {
//	switch block.BlockType {
//	case CbtLetter:
//		name := block.getChars()
//		value := NewZhenValueNil()
//		blnEnable := false
//		next, ok := block.getNext()
//		block.WordType = CwtConstant
//		block.Word = block.getChars()
//		if ok && next.BlockType == CbtOperator {
//			if next.getChars() == "=" {
//				next2, ok2 := next.getNext()
//				if ok2 {
//					switch next2.BlockType {
//					case CbtString:
//						value = StringToZhenValue(next2.getChars())
//						blnEnable = true
//					case CbtNumber:
//						n, err := strconv.ParseFloat(next2.getChars(), 64)
//						if err != nil {
//							return err
//						}
//						value = NewZhenValueNumber(ZFloat(n))
//						blnEnable = true
//					}
//				}
//			}
//		}
//		if blnEnable {
//			cs := ZhenCodeStep{codeStepType: ZCS_Var, valueName1: name, value: value}
//			block.codeSteps = append(block.codeSteps, cs)
//		}
//
//	case CbtLine, CbtChildLine:
//		for _, c := range block.items {
//			ConstantStatement(zhen, c)
//		}
//
//	}
//	return
//}
//

//func DefineVar(zhen *ZhenState) (err error) {
//	nowCodeBlock := zhen.NowCodeBlock
//	nowCodeBlock.Word = nowCodeBlock.getChars()
//	//var values []string
//	//isDef := false
//	next, ok := zhen.NowCodeBlock.getNext()
//	if ok {
//		if next.BlockType == CbtColon {
//			//isDef = true
//			for _, c := range next.items {
//				ConstantStatement(zhen, c)
//				//增加 类型文本（CbtTxt）
//				//c.WordType = CwtTxt
//				//c.Word = c.getChars()
//				//values = append(values, c.Word)
//				//if c.BlockType != CbtString {
//				//	c.w = CbtString
//				//}
//			}
//		}
//	}
//	nowCodeBlock.WordType = CwtKeyWord
//	return
//}
//func DefineConstantPreFun(zhen *ZhenState) (err error) {
//	nowCodeBlock := zhen.NowCodeBlock
//	nowCodeBlock.Word = nowCodeBlock.getChars()
//	//var values []string
//	//isDef := false
//	next, ok := zhen.NowCodeBlock.getNext()
//	if ok {
//		if next.BlockType == CbtColon {
//			//isDef = true
//			for _, c := range next.items {
//				ConstantStatement(zhen, c)
//				//增加 类型文本（CbtTxt）
//				//c.WordType = CwtTxt
//				//c.Word = c.getChars()
//				//values = append(values, c.Word)
//				//if c.BlockType != CbtString {
//				//	c.w = CbtString
//				//}
//			}
//		}
//	}
//	nowCodeBlock.WordType = CwtKeyWord
//	return
//}
//
//func DefineFun(zhen *ZhenState) (err error) {
//	nowCodeBlock := zhen.NowCodeBlock
//	nowCodeBlock.Word = nowCodeBlock.getChars()
//	//var values []string
//	//isDef := false
//	next, ok := zhen.NowCodeBlock.getNext()
//	if ok {
//		if next.BlockType == CbtLetter {
//			cs := ZhenCodeStep{codeStepType: ZCS_PrintVar, valueName1: next.getChars()}
//			nowCodeBlock.codeSteps = append(nowCodeBlock.codeSteps, cs)
//
//		}
//	}
//	nowCodeBlock.WordType = CwtKeyWord
//
//	return
//}
//func DefineFunParaPreFun(zhen *ZhenState) (err error) {
//	return
//}
//
//func DefineFunReturnPreFun(zhen *ZhenState) (err error) {
//	return
//}
