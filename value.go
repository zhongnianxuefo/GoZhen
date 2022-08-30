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
	ZhenValueTypeInt
	ZhenValueTypeFloat
	ZhenValueTypeString
	ZhenValueTypeArray
	ZhenValueTypeTable
	ZhenValueTypeFunction
)

type ZhenValueBoolean bool
type ZhenValueInt int64
type ZhenValueFloat float64
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
	valueInt      ZhenValueInt
	valueFloat    ZhenValueFloat
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

func NewZhenValueInt(valueInt ZhenValueInt) (v ZhenValue) {
	v.valueType = ZhenValueTypeInt
	v.valueInt = valueInt
	return
}

func NewZhenValueFloat(valueFloat ZhenValueFloat) (v ZhenValue) {
	v.valueType = ZhenValueTypeFloat
	v.valueFloat = valueFloat
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

func ZhenValueToInt(v ZhenValue) (valueInt ZhenValueInt, err error) {
	switch v.valueType {
	case ZhenValueTypeInt:
		valueInt = v.valueInt
	case ZhenValueTypeFloat:
		f := float64(v.valueFloat)
		valueInt = ZhenValueInt(int64(f))
	default:
		err = errors.New("只有整数和小数类型可以转换为整数")
	}
	return
}
func ZhenValueToFloat(v ZhenValue) (valueFloat ZhenValueFloat, err error) {
	switch v.valueType {
	case ZhenValueTypeInt:
		f := float64(v.valueInt)
		valueFloat = ZhenValueFloat(f)
	case ZhenValueTypeFloat:
		valueFloat = v.valueFloat
	default:
		err = errors.New("只有整数和小数类型可以转换为小数")
	}
	return
}

func ZhenValueAdd(v1 ZhenValue, v2 ZhenValue) (v ZhenValue, err error) {
	switch v1.valueType {
	case ZhenValueTypeInt:
		switch v2.valueType {
		case ZhenValueTypeInt:
			v.valueType = ZhenValueTypeInt
		case ZhenValueTypeFloat:
			v.valueType = ZhenValueTypeFloat
		default:
			err = errors.New("整数类型只能和整数或者小数相加")
			return
		}

	case ZhenValueTypeFloat:
		switch v2.valueType {
		case ZhenValueTypeInt:
			v.valueType = ZhenValueTypeFloat
		case ZhenValueTypeFloat:
			v.valueType = ZhenValueTypeFloat
		default:
			err = errors.New("整数类型只能和整数或者小数相加")
			return
		}

	case ZhenValueTypeString:
		switch v2.valueType {
		case ZhenValueTypeString:
			v.valueType = ZhenValueTypeString
		}
	}

	switch v.valueType {
	case ZhenValueTypeInt:
		var n1, n2 ZhenValueInt
		n1, err = ZhenValueToInt(v1)
		if err != nil {
			return
		}
		n2, err = ZhenValueToInt(v2)
		if err != nil {
			return
		}
		v.valueInt = n1 + n2
	case ZhenValueTypeFloat:
		var f1, f2 ZhenValueFloat
		f1, err = ZhenValueToFloat(v1)
		if err != nil {
			return
		}
		f2, err = ZhenValueToFloat(v2)
		if err != nil {
			return
		}
		v.valueFloat = f1 + f2
	case ZhenValueTypeString:
		v.valueString = v1.valueString + v2.valueString
	}
	return
}

func ZhenValueToString(v ZhenValue) (s string) {
	switch v.valueType {
	case ZhenValueTypeNone:
		s = "none"
	case ZhenValueTypeNil:
		s = "nil"
	case ZhenValueTypeBoolean:
		if v.valueBool == true {
			s = "true"
		} else {
			s = "false"
		}
	case ZhenValueTypeInt:
		s = strconv.FormatInt(int64(v.valueInt), 10)
	case ZhenValueTypeFloat:
		s = strconv.FormatFloat(float64(v.valueFloat), 'f', -1, 32)
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
		fmt.Println(v.valueType)
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
	} else if valueType == "整数" {
		n, e := strconv.Atoi(value)
		if e != nil {
			err = errors.New("解析错误:不能把文本转换为整数")
			return
		}
		v = NewZhenValueInt(ZhenValueInt(n))
	} else if valueType == "小数" {
		n, e := strconv.ParseFloat(value, 64)
		if e != nil {
			err = errors.New("解析错误:不能把文本转换为小数")
			return
		}
		v = NewZhenValueFloat(ZhenValueFloat(n))
	} else if valueType == "字符串" {
		v = NewZhenValueString(ZhenValueString(value))
	}

	return
}
