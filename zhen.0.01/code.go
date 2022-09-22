package zhen_0_01

import (
	"github.com/beevik/etree"
	"strconv"
)

type ZhenCodeStepType int

const (
	ZCS_None ZhenCodeStepType = iota

	ZCS_Var
	ZCS_As

	ZCS_Add
	ZCS_Sub
	ZCS_Mul
	ZCS_Div
	ZCS_Eq
	ZCS_Ne
	ZCS_Gt
	ZCS_Lt
	ZCS_And
	ZCS_Or
	ZCS_Not

	ZCS_TVar
	ZCS_TFrom
	ZCS_TAs
	ZCS_TAdd
	ZCS_TSub
	ZCS_TMul
	ZCS_TDiv
	ZCS_TEq
	ZCS_TNe
	ZCS_TGt
	ZCS_TLt
	ZCS_TAnd
	ZCS_TOr
	ZCS_TNot

	ZCS_TTAs
	ZCS_TTAdd
	ZCS_TTSub
	ZCS_TTMul
	ZCS_TTDiv
	ZCS_TTEq
	ZCS_TTNe
	ZCS_TTGt
	ZCS_TTLt
	ZCS_TTAnd
	ZCS_TTOr

	ZCS_If
	ZCS_For
	ZCS_While
	ZCS_Break
	ZCS_Return
	ZCS_Call

	ZCS_PrintVar
)

type ZhenCodeStep struct {
	codeStepType ZhenCodeStepType
	valueName1   string
	valueName2   string
	tempValueNo1 int
	tempValueNo2 int
	value        ZValue
}

type ZhenCodeWord struct {
	tag     string
	content string
}

type ZhenCodeOld struct {
	needRun   bool
	codeWords []ZhenCodeWord
	codeSteps []ZhenCodeStep
}

func getCodeStepFromElement(item *etree.Element) (s ZhenCodeStep, err error) {

	allComm := map[string]ZhenCodeStepType{
		"变量定义":    ZCS_Var,
		"变量赋值":    ZCS_As,
		"变量相加":    ZCS_Add,
		"变量相减":    ZCS_Sub,
		"变量相乘":    ZCS_Mul,
		"变量相除":    ZCS_Div,
		"变量相等判断":  ZCS_Eq,
		"变量不等于判断": ZCS_Ne,
		"变量大于判断":  ZCS_Gt,
		"变量小于判断":  ZCS_Lt,
		"变量且运算":   ZCS_And,
		"变量或运算":   ZCS_Or,
		"变量非运算":   ZCS_Not,

		"临时变量定义":         ZCS_TVar,
		"普通变量赋值给临时变量":    ZCS_TFrom,
		"临时变量赋值给普通变量":    ZCS_TAs,
		"临时变量和普通变量相加":    ZCS_TAdd,
		"临时变量和普通变量相减":    ZCS_TSub,
		"临时变量和普通变量相乘":    ZCS_TMul,
		"临时变量和普通变量相除":    ZCS_TDiv,
		"临时变量和普通变量相等判断":  ZCS_TEq,
		"临时变量和普通变量不等于判断": ZCS_TNe,
		"临时变量和普通变量大于判断":  ZCS_TGt,
		"临时变量和普通变量小于判断":  ZCS_TLt,
		"临时变量和普通变量且运算":   ZCS_TAnd,
		"临时变量和普通变量或运算":   ZCS_TOr,
		"临时变量非运算":        ZCS_TNot,

		"临时变量赋值给临时变量": ZCS_TTAs,
		"两个临时变量相加":    ZCS_TTAdd,
		"两个临时变量相减":    ZCS_TTSub,
		"两个临时变量相乘":    ZCS_TTMul,
		"两个临时变量相除":    ZCS_TTDiv,
		"两个临时变量相等判断":  ZCS_TTEq,
		"两个临时变量不等于判断": ZCS_TTNe,
		"两个临时变量大于判断":  ZCS_TTGt,
		"两个临时变量小于判断":  ZCS_TTLt,
		"两个临时变量且运算":   ZCS_TTAnd,
		"两个临时变量或运算":   ZCS_TTOr,

		"如果":   ZCS_If,
		"次数循环": ZCS_For,
		"条件循环": ZCS_While,
		"跳出循环": ZCS_Break,
		"函数返回": ZCS_Return,
		"运行函数": ZCS_Call,

		"显示变量": ZCS_PrintVar,
	}
	comm := item.SelectAttrValue("指令", "")
	s.codeStepType = allComm[comm]
	needValueName1 := false
	needValueName2 := false
	needTempValueName1 := false
	needTempValueName2 := false
	needValue := false

	switch s.codeStepType {
	case ZCS_Var:
		needValueName1 = true
		needValue = true
	case ZCS_As:
		needValueName1 = true
		needValueName2 = true
	case ZCS_Add, ZCS_Sub, ZCS_Mul, ZCS_Div, ZCS_Eq, ZCS_Ne, ZCS_Gt, ZCS_Lt, ZCS_And, ZCS_Or:
		needValueName1 = true
		needValueName2 = true
		needTempValueName1 = true
	case ZCS_Not:
		needValueName1 = true
		needTempValueName1 = true
	case ZCS_TVar:
		needTempValueName1 = true
		needValue = true
	case ZCS_TFrom, ZCS_TAs, ZCS_TAdd, ZCS_TSub, ZCS_TMul, ZCS_TDiv, ZCS_TEq, ZCS_TNe, ZCS_TGt, ZCS_TLt, ZCS_TAnd, ZCS_TOr:
		needValueName1 = true
		needTempValueName1 = true
	case ZCS_TNot:
		needTempValueName1 = true
	case ZCS_TTAs, ZCS_TTAdd, ZCS_TTSub, ZCS_TTMul, ZCS_TTDiv, ZCS_TTEq, ZCS_TTNe, ZCS_TTGt, ZCS_TTLt, ZCS_TTAnd, ZCS_TTOr:
		needTempValueName1 = true
		needTempValueName2 = true
	case ZCS_PrintVar:
		needValueName1 = true
	}

	if needValueName1 {
		s.valueName1 = item.SelectAttrValue("变量名", "")
		if s.valueName1 == "" {
			s.valueName1 = item.SelectAttrValue("变量名1", "")
		}
	}
	if needValueName2 {
		s.valueName2 = item.SelectAttrValue("变量名2", "")
	}
	if needTempValueName1 {
		s.tempValueNo1, err = strconv.Atoi(item.SelectAttrValue("临时变量1", ""))
		if err != nil {
			return
		}
	}
	if needTempValueName2 {
		s.tempValueNo2, err = strconv.Atoi(item.SelectAttrValue("临时变量2", ""))
		if err != nil {
			return
		}
	}
	if needValue {
		//s.value, err = GetZhenValueFromElement(item)
	}
	return
}
