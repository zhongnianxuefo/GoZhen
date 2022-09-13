package zhen

import (
	"strconv"
	"strings"
)

func KeyWordTextPreFun(zhen *ZhenState) (err error) {
	nowCodeBlock := zhen.NowCodeBlock
	nowCodeBlock.Word = nowCodeBlock.getChars()
	var values []string
	isDef := false
	next, ok := zhen.NowCodeBlock.getNext()
	if ok {
		if next.BlockType == CbtColon {
			isDef = true
			for _, c := range next.Items {
				//增加 类型文本（CbtTxt）
				c.WordType = CwtTxt
				c.Word = c.getChars()
				values = append(values, c.Word)
				//if c.BlockType != CbtString {
				//	c.w = CbtString
				//}
			}
		}
	}
	if isDef {

		nowCodeBlock.WordType = CwtKeyWord
		v := NewZhenValueString(ZhenValueString(strings.Join(values, " ")))
		zhen.SetGlobalVarValue(nowCodeBlock.Word, v)

	} else {

		nowCodeBlock.WordType = CwtConstant
	}

	return
}
func ConstantStatement(zhen *ZhenState, block *CodeBlock) (err error) {
	switch block.BlockType {
	case CbtLetter:
		name := block.getChars()
		value := NewZhenValueNil()
		blnEnable := false
		next, ok := block.getNext()
		block.WordType = CwtConstant
		block.Word = block.getChars()
		if ok && next.BlockType == CbtOperator {
			if next.getChars() == "=" {
				next2, ok2 := next.getNext()
				if ok2 {
					switch next2.BlockType {
					case CbtString:
						value = StringToZhenValue(next2.getChars())
						blnEnable = true
					case CbtNumber:
						n, err := strconv.ParseFloat(next2.getChars(), 64)
						if err != nil {
							return err
						}
						value = NewZhenValueNumber(ZhenValueNumber(n))
						blnEnable = true
					}
				}
			}
		}
		if blnEnable {
			cs := ZhenCodeStep{codeStepType: ZCS_Var, valueName1: name, value: value}
			block.codeSteps = append(block.codeSteps, cs)
		}

	case CbtLine, CbtChildLine:
		for _, c := range block.Items {
			ConstantStatement(zhen, c)
		}

	}
	return
}
func KeyWordConstantPreFun(zhen *ZhenState) (err error) {
	nowCodeBlock := zhen.NowCodeBlock
	nowCodeBlock.Word = nowCodeBlock.getChars()
	//var values []string
	//isDef := false
	next, ok := zhen.NowCodeBlock.getNext()
	if ok {
		if next.BlockType == CbtColon {
			//isDef = true
			for _, c := range next.Items {
				ConstantStatement(zhen, c)
				//增加 类型文本（CbtTxt）
				//c.WordType = CwtTxt
				//c.Word = c.getChars()
				//values = append(values, c.Word)
				//if c.BlockType != CbtString {
				//	c.w = CbtString
				//}
			}
		}
	}
	nowCodeBlock.WordType = CwtKeyWord
	return
}
func KeyWordFunPreFun(zhen *ZhenState) (err error) {
	nowCodeBlock := zhen.NowCodeBlock
	nowCodeBlock.Word = nowCodeBlock.getChars()
	//var values []string
	//isDef := false
	next, ok := zhen.NowCodeBlock.getNext()
	if ok {
		if next.BlockType == CbtLetter {
			cs := ZhenCodeStep{codeStepType: ZCS_PrintVar, valueName1: next.getChars()}
			nowCodeBlock.codeSteps = append(nowCodeBlock.codeSteps, cs)

		}
	}
	nowCodeBlock.WordType = CwtKeyWord

	return
}
