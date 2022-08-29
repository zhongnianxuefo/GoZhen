package main

import (
	"fmt"
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
