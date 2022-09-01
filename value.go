package main

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"strconv"
	"strings"
)

type ZhenValueType byte

const (
	ZhenValueTypeNone ZhenValueType = iota
	ZhenValueTypeNil
	ZhenValueTypeBoolean
	ZhenValueTypeNumber
	ZhenValueTypeString
	ZhenValueTypeArray
	ZhenValueTypeTable
	ZhenValueTypeFunction
)

type ZhenValueBoolean bool
type ZhenValueNumber float64
type ZhenValueString string
type ZhenValueArray []ZhenValue
type ZhenValueTable map[string]ZhenValue

type ZhenValueFunction struct {
	FunctionName string
	Function     func()
}

type ZhenValue struct {
	valueType     ZhenValueType
	valueBool     ZhenValueBoolean
	valueNumber   ZhenValueNumber
	valueString   ZhenValueString
	valueArray    ZhenValueArray
	valueTable    ZhenValueTable
	valueFunction ZhenValueFunction
}

func NewZhenValueNone() (v ZhenValue) {
	v.valueType = ZhenValueTypeNone
	return
}

func NewZhenValueNil() (v ZhenValue) {
	v.valueType = ZhenValueTypeNil
	return
}

func NewZhenValueBoolean(valueBoolean ZhenValueBoolean) (v ZhenValue) {
	v.valueType = ZhenValueTypeBoolean
	v.valueBool = valueBoolean
	return
}

func NewZhenValueNumber(valueNumber ZhenValueNumber) (v ZhenValue) {
	v.valueType = ZhenValueTypeNumber
	v.valueNumber = valueNumber
	return
}

func NewZhenValueString(valueString ZhenValueString) (v ZhenValue) {
	v.valueType = ZhenValueTypeString
	v.valueString = valueString
	return
}

func NewZhenValueArray(valueArray ZhenValueArray) (v ZhenValue) {
	v.valueType = ZhenValueTypeArray
	v.valueArray = valueArray
	return
}

func NewZhenValueTable(valueTable ZhenValueTable) (v ZhenValue) {
	v.valueType = ZhenValueTypeTable
	v.valueTable = valueTable
	return
}

func NewZhenValueFunction(valueFunction ZhenValueFunction) (v ZhenValue) {
	v.valueType = ZhenValueTypeFunction
	v.valueFunction = valueFunction
	return
}

func ZhenValueOperation(s ZhenCodeStepType, v1 ZhenValue, v2 ZhenValue) (v ZhenValue, err error) {
	isArithmetic := false
	isCompare := false
	isBool := false
	isSingleBool := false
	switch s {
	case ZCS_Add, ZCS_Sub, ZCS_Mul, ZCS_Div,
		ZCS_TAdd, ZCS_TSub, ZCS_TMul, ZCS_TDiv,
		ZCS_TTAdd, ZCS_TTSub, ZCS_TTMul, ZCS_TTDiv:
		isArithmetic = true
	case ZCS_Eq, ZCS_Ne, ZCS_Gt, ZCS_Lt,
		ZCS_TEq, ZCS_TNe, ZCS_TGt, ZCS_TLt,
		ZCS_TTEq, ZCS_TTNe, ZCS_TTGt, ZCS_TTLt:
		isCompare = true
	case ZCS_And, ZCS_Or, ZCS_TAnd, ZCS_TOr, ZCS_TTAnd, ZCS_TTOr:
		isBool = true
	case ZCS_Not, ZCS_TNot:
		isSingleBool = true
	}
	vt1 := v1.valueType
	vt2 := v2.valueType

	if isArithmetic {
		if vt1 == ZhenValueTypeNumber && vt2 == ZhenValueTypeNumber {
			v.valueType = ZhenValueTypeNumber
			switch s {
			case ZCS_Add, ZCS_TAdd, ZCS_TTAdd:
				v.valueNumber = v1.valueNumber + v2.valueNumber
			case ZCS_Sub, ZCS_TSub, ZCS_TTSub:
				v.valueNumber = v1.valueNumber - v2.valueNumber
			case ZCS_Mul, ZCS_TMul, ZCS_TTMul:
				v.valueNumber = v1.valueNumber * v2.valueNumber
			case ZCS_Div, ZCS_TDiv, ZCS_TTDiv:
				v.valueNumber = v1.valueNumber / v2.valueNumber
			}
		} else {
			err = errors.New("只有数字类型可以进行四则运算")
			return
		}

	} else if isCompare {
		if vt1 == ZhenValueTypeBoolean && vt2 == ZhenValueTypeBoolean {
			v.valueType = ZhenValueTypeBoolean
			vv1 := v1.valueBool
			vv2 := v2.valueBool

			switch s {
			case ZCS_Eq, ZCS_TEq, ZCS_TTEq:
				v.valueBool = vv1 == vv2
			case ZCS_Ne, ZCS_TNe, ZCS_TTNe:
				v.valueBool = vv1 != vv2
			default:
				err = errors.New("布尔类型不能进行大小比较")
				return
			}

		} else if vt1 == ZhenValueTypeNumber && vt2 == ZhenValueTypeNumber {
			v.valueType = ZhenValueTypeBoolean
			vv1 := v1.valueNumber
			vv2 := v2.valueNumber

			switch s {
			case ZCS_Eq, ZCS_TEq, ZCS_TTEq:
				v.valueBool = vv1 == vv2
			case ZCS_Ne, ZCS_TNe, ZCS_TTNe:
				v.valueBool = vv1 != vv2

			case ZCS_Gt, ZCS_TGt, ZCS_TTGt:
				v.valueBool = vv1 > vv2
			case ZCS_Lt, ZCS_TLt, ZCS_TTLt:
				v.valueBool = vv1 < vv2

			}
		} else if vt1 == ZhenValueTypeString && vt2 == ZhenValueTypeString {
			v.valueType = ZhenValueTypeBoolean
			vv1 := v1.valueString
			vv2 := v2.valueString

			switch s {
			case ZCS_Eq, ZCS_TEq, ZCS_TTEq:
				v.valueBool = vv1 == vv2
			case ZCS_Ne, ZCS_TNe, ZCS_TTNe:
				v.valueBool = vv1 != vv2
			case ZCS_Gt, ZCS_TGt, ZCS_TTGt:
				v.valueBool = vv1 > vv2
			case ZCS_Lt, ZCS_TLt, ZCS_TTLt:
				v.valueBool = vv1 < vv2
			}
		} else {
			//todo 其他类型的比较待定
			err = errors.New("比较类型未定义")
			return
		}

	} else if isBool {
		if vt1 == ZhenValueTypeBoolean && vt2 == ZhenValueTypeBoolean {
			v.valueType = ZhenValueTypeBoolean
			vv1 := v1.valueBool
			vv2 := v2.valueBool

			switch s {
			case ZCS_And, ZCS_TAnd, ZCS_TTAnd:
				v.valueBool = vv1 && vv2
			case ZCS_Or, ZCS_TOr, ZCS_TTOr:
				v.valueBool = vv1 || vv2

			}
		} else {
			err = errors.New("只有布尔类型可以进行比较运算")
			return
		}
	} else if isSingleBool {
		if vt1 == ZhenValueTypeBoolean {
			v.valueType = ZhenValueTypeBoolean
			v.valueBool = !v1.valueBool
		} else {
			err = errors.New("只有布尔类型可以进行非运算")
			return
		}
	}
	return
}

func ZhenValueToString(v ZhenValue) (s string) {
	switch v.valueType {
	case ZhenValueTypeNone:
		s = "未定义"
	case ZhenValueTypeNil:
		s = "空值"
	case ZhenValueTypeBoolean:
		if v.valueBool == true {
			s = "真"
		} else {
			s = "假"
		}
	case ZhenValueTypeNumber:
		s = strconv.FormatFloat(float64(v.valueNumber), 'f', -1, 32)
	case ZhenValueTypeString:
		s = string(v.valueString)
	case ZhenValueTypeArray:
		var ws []string
		for _, a := range v.valueArray {
			ws = append(ws, ZhenValueToString(a))
		}
		s = strings.Join(ws, ", ")
	case ZhenValueTypeTable:
		var ws []string
		for k, a := range v.valueTable {
			w := fmt.Sprintf("%s = %s", k, ZhenValueToString(a))
			ws = append(ws, w)
		}
		s = strings.Join(ws, ", ")
	case ZhenValueTypeFunction:
		s = v.valueFunction.FunctionName
	default:
		fmt.Println("未知类型。", v.valueType)
	}
	return
}

func getZhenValueFromElement(item *etree.Element) (v ZhenValue, err error) {
	valueType := item.SelectAttrValue("值类型", "")
	value := item.SelectAttrValue("值", "")

	if valueType == "未定义" {
		v = NewZhenValueNone()
	} else if valueType == "空值" {
		v = NewZhenValueNil()
	} else if valueType == "数字" {
		n, e := strconv.ParseFloat(value, 64)
		if e != nil {
			err = errors.New("解析错误:不能把文本转换为数字")
			return
		}
		v = NewZhenValueNumber(ZhenValueNumber(n))
	} else if valueType == "字符串" {
		v = NewZhenValueString(ZhenValueString(value))
	}

	return
}
