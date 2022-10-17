package zhen_0_02

import (
	"errors"
	"fmt"
)

type Value interface{}

type ValueType byte

const (
	VtNone ValueType = iota
	VtNil
	VtBool
	VtInt
	VtFloat
	VtString
	VtFun
	VtObject
)

var ValueNames = [...]string{
	VtNone:   "未定义",
	VtNil:    "空值",
	VtBool:   "布尔值",
	VtInt:    "整数",
	VtFloat:  "小数",
	VtString: "字符串",
	VtFun:    "函数",
	VtObject: "对象",
}

func (vt ValueType) String() string {
	return ValueNames[int(vt)]
}

type ValueGroup struct {
	Values        map[string]Value
	ParValueGroup *ValueGroup
	CodeBlock     *CodeBlock
}

type TempValueGroup struct {
	Values        map[int]Value
	ParValueGroup *TempValueGroup
	CodeBlock     *CodeBlock
}

func NewValueGroup(par *ValueGroup, codeBlock *CodeBlock) (vg *ValueGroup) {
	vg = &ValueGroup{}
	vg.ParValueGroup = par
	vg.CodeBlock = codeBlock
	vg.Values = make(map[string]Value)
	return
}
func NewTempValueGroup(par *TempValueGroup, codeBlock *CodeBlock) (vg *TempValueGroup) {
	vg = &TempValueGroup{}
	vg.ParValueGroup = par
	vg.CodeBlock = codeBlock
	vg.Values = make(map[int]Value)
	return
}

func ValueToInt(v Value) (i int, ok bool) {
	a, ok := v.(int)
	if ok {
		i = a
		return
	}
	aa, ok := v.(*int)
	if ok {
		i = *aa
		return
	}
	return
}
func ValueToFloat(v Value) (f float64, ok bool) {
	i, ok := ValueToInt(v)
	if ok {
		f = float64(i)
		return
	}
	a, ok := v.(float64)
	if ok {
		f = a
		return
	}
	aa, ok := v.(*float64)
	if ok {
		f = *aa
		return
	}
	return
}
func ValueToString(v Value) (s string, ok bool) {
	a, ok := v.(string)
	if ok {
		s = a
		return
	}
	aa, ok := v.(*string)
	if ok {
		s = *aa
		return
	}

	return
}
func Arithmetic(aa Value, bb Value, t CodeStepType) (v Value, err error) {
	blnDo := false
	if !blnDo {
		a1, ok1 := ValueToInt(aa)
		a2, ok2 := ValueToInt(bb)
		if ok1 && ok2 {
			switch t {
			case CstAdd:
				v = a1 + a2
			case CstSub:
				v = a1 - a2
			case CstMul:
				v = a1 * a2
			case CstDiv:
				if a2 == 0 {
					err = errors.New("被除数不能为0")
				} else {
					v = float64(a1) / float64(a2)
				}
			}
			blnDo = true
		}
	}

	if !blnDo {
		a1, ok1 := ValueToFloat(aa)
		a2, ok2 := ValueToFloat(bb)
		if ok1 && ok2 {
			switch t {
			case CstAdd:
				v = a1 + a2
			case CstSub:
				v = a1 - a2
			case CstMul:
				v = a1 * a2
			case CstDiv:
				if a2 == 0 {
					err = errors.New("被除数不能为0")
				} else {
					v = a1 / a2
				}
			}
			blnDo = true
		}
	}
	if !blnDo {
		a1, ok1 := ValueToString(aa)
		a2, ok2 := ValueToString(bb)
		if ok1 && ok2 {
			switch t {
			case CstAdd:
				v = a1 + a2
			}
			blnDo = true
		}
	}
	return
}

func PrintValue(v Value) {
	switch a := v.(type) {
	case int:
		fmt.Println(a)
	case float64:
		fmt.Println(a)
	case string:
		fmt.Println(a)
	case *int:
		fmt.Println(*a)
	case *float64:
		fmt.Println(*a)
	case *string:
		fmt.Println(*a)
	}

	return
}

//type None bool
//
//type ZInt int
//type ZFloat float64

//type ZArray []ValueString1

//type ZObject []ValueString1

//type ZhenValueTable map[string]ZhenValue

//type ZFun func(state *State) error
//
////type ZDefFun func(codeBlock *CodeBlock2) error
//type ZObject struct {
//}
