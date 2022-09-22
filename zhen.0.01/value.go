package zhen_0_01

import (
	"fmt"
	"go/types"
	"strconv"
)

type ZValue interface{}

type ZValueType byte

const (
	ZvtNone ZValueType = iota
	ZvtNil
	ZvtBool
	ZvtInt
	ZvtFloat
	ZvtString
	ZvtFun
	ZvtObject
)

var ZValueNames = [...]string{
	ZvtNone:   "未定义",
	ZvtNil:    "空值",
	ZvtBool:   "布尔值",
	ZvtInt:    "整数",
	ZvtFloat:  "小数",
	ZvtString: "字符串",
	ZvtFun:    "函数",
	ZvtObject: "对象",
}

func (zt ZValueType) String() string {
	return ZValueNames[int(zt)]
}

type ZNone bool

type ZInt int
type ZFloat float64

//type ZArray []ValueString

//type ZObject []ValueString

//type ZhenValueTable map[string]ZhenValue

type ZFun func(state *ZhenState) error

//type ZDefFun func(codeBlock *CodeBlock2) error
type ZObject struct {
}

//type ZhenValue struct {
//	valueName     string
//	valueType     ValueType
//	valueBool     ZhenValueBoolean
//	valueNumber   ZFloat
//	valueString   ZhenValueString
//	valueArray    ZArray
//	valueTable    ZhenValueTable
//	valueFunction func(state *ZhenState) error
//}
const None = ZNone(false)

//
//func NewZhenValueNil() (v ZhenValue) {
//	v.valueType = ZvtNil
//	return
//}
//
//func NewZhenValueBoolean(valueBoolean ZhenValueBoolean) (v ZhenValue) {
//	v.valueType = ZvtBool
//	v.valueBool = valueBoolean
//	return
//}
//
//func NewZhenValueNumber(valueNumber ZFloat) (v ZhenValue) {
//	v.valueType = ZvtFloat
//	v.valueNumber = valueNumber
//	return
//}
//
//func NewZhenValueString(valueString ZhenValueString) (v ZhenValue) {
//	v.valueType = ZvtString
//	v.valueString = valueString
//	return
//}
//
//func NewZhenValueArray(valueArray ZArray) (v ZhenValue) {
//	v.valueType = VtArray
//	v.valueArray = valueArray
//	return
//}
//
//func NewZhenValueTable(valueTable ZhenValueTable) (v ZhenValue) {
//	v.valueType = ZvtObject
//	v.valueTable = valueTable
//	return
//}
//
//func NewZhenValueFunction(valueFunction ZFun) (v ZhenValue) {
//	v.valueType = ZvtFun
//	v.valueFunction = (*ZFun)(unsafe.Pointer(&valueFunction))
//	return
//}
//func StringToZhenValue(valueString string) (v ZhenValue) {
//	v.valueType = ZvtString
//	v.valueString = ZhenValueString(valueString)
//	return
//}
//
//func ZhenValueOperation(s ZhenCodeStepType, v1 ZhenValue, v2 ZhenValue) (v ZhenValue, err error) {
//	isArithmetic := false
//	isCompare := false
//	isBool := false
//	isSingleBool := false
//	switch s {
//	case ZCS_Add, ZCS_Sub, ZCS_Mul, ZCS_Div,
//		ZCS_TAdd, ZCS_TSub, ZCS_TMul, ZCS_TDiv,
//		ZCS_TTAdd, ZCS_TTSub, ZCS_TTMul, ZCS_TTDiv:
//		isArithmetic = true
//	case ZCS_Eq, ZCS_Ne, ZCS_Gt, ZCS_Lt,
//		ZCS_TEq, ZCS_TNe, ZCS_TGt, ZCS_TLt,
//		ZCS_TTEq, ZCS_TTNe, ZCS_TTGt, ZCS_TTLt:
//		isCompare = true
//	case ZCS_And, ZCS_Or, ZCS_TAnd, ZCS_TOr, ZCS_TTAnd, ZCS_TTOr:
//		isBool = true
//	case ZCS_Not, ZCS_TNot:
//		isSingleBool = true
//	}
//	vt1 := v1.valueType
//	vt2 := v2.valueType
//
//	if isArithmetic {
//		if vt1 == ZvtFloat && vt2 == ZvtFloat {
//			v.valueType = ZvtFloat
//			switch s {
//			case ZCS_Add, ZCS_TAdd, ZCS_TTAdd:
//				v.valueNumber = v1.valueNumber + v2.valueNumber
//			case ZCS_Sub, ZCS_TSub, ZCS_TTSub:
//				v.valueNumber = v1.valueNumber - v2.valueNumber
//			case ZCS_Mul, ZCS_TMul, ZCS_TTMul:
//				v.valueNumber = v1.valueNumber * v2.valueNumber
//			case ZCS_Div, ZCS_TDiv, ZCS_TTDiv:
//				v.valueNumber = v1.valueNumber / v2.valueNumber
//			}
//		} else {
//			err = errors.New("只有数字类型可以进行四则运算")
//			return
//		}
//
//	} else if isCompare {
//		if vt1 == ZvtBool && vt2 == ZvtBool {
//			v.valueType = ZvtBool
//			vv1 := v1.valueBool
//			vv2 := v2.valueBool
//
//			switch s {
//			case ZCS_Eq, ZCS_TEq, ZCS_TTEq:
//				v.valueBool = vv1 == vv2
//			case ZCS_Ne, ZCS_TNe, ZCS_TTNe:
//				v.valueBool = vv1 != vv2
//			default:
//				err = errors.New("布尔类型不能进行大小比较")
//				return
//			}
//
//		} else if vt1 == ZvtFloat && vt2 == ZvtFloat {
//			v.valueType = ZvtBool
//			vv1 := v1.valueNumber
//			vv2 := v2.valueNumber
//
//			switch s {
//			case ZCS_Eq, ZCS_TEq, ZCS_TTEq:
//				v.valueBool = vv1 == vv2
//			case ZCS_Ne, ZCS_TNe, ZCS_TTNe:
//				v.valueBool = vv1 != vv2
//
//			case ZCS_Gt, ZCS_TGt, ZCS_TTGt:
//				v.valueBool = vv1 > vv2
//			case ZCS_Lt, ZCS_TLt, ZCS_TTLt:
//				v.valueBool = vv1 < vv2
//
//			}
//		} else if vt1 == ZvtString && vt2 == ZvtString {
//			v.valueType = ZvtBool
//			vv1 := v1.valueString
//			vv2 := v2.valueString
//
//			switch s {
//			case ZCS_Eq, ZCS_TEq, ZCS_TTEq:
//				v.valueBool = vv1 == vv2
//			case ZCS_Ne, ZCS_TNe, ZCS_TTNe:
//				v.valueBool = vv1 != vv2
//			case ZCS_Gt, ZCS_TGt, ZCS_TTGt:
//				v.valueBool = vv1 > vv2
//			case ZCS_Lt, ZCS_TLt, ZCS_TTLt:
//				v.valueBool = vv1 < vv2
//			}
//		} else {
//			//todo 其他类型的比较待定
//			err = errors.New("比较类型未定义")
//			return
//		}
//
//	} else if isBool {
//		if vt1 == ZvtBool && vt2 == ZvtBool {
//			v.valueType = ZvtBool
//			vv1 := v1.valueBool
//			vv2 := v2.valueBool
//
//			switch s {
//			case ZCS_And, ZCS_TAnd, ZCS_TTAnd:
//				v.valueBool = vv1 && vv2
//			case ZCS_Or, ZCS_TOr, ZCS_TTOr:
//				v.valueBool = vv1 || vv2
//
//			}
//		} else {
//			err = errors.New("只有布尔类型可以进行比较运算")
//			return
//		}
//	} else if isSingleBool {
//		if vt1 == ZvtBool {
//			v.valueType = ZvtBool
//			v.valueBool = !v1.valueBool
//		} else {
//			err = errors.New("只有布尔类型可以进行非运算")
//			return
//		}
//	}
//	return
//}

//func debugValue(v value) string {
//	switch v := v.(type) {
//	case *table:
//		entry := func(x value) string {
//			if t, ok := x.(*table); ok {
//				return fmt.Sprintf("table %#v", t)
//			}
//			return debugValue(x)
//		}
//		s := fmt.Sprintf("table %#v {[", v)
//		for _, x := range v.array {
//			s += entry(x) + ", "
//		}
//		s += "], {"
//		for k, x := range v.hash {
//			s += entry(k) + ": " + entry(x) + ", "
//		}
//		return s + "}}"
//	case string:
//		return "'" + v + "'"
//	case float64:
//		return fmt.Sprintf("%f", v)
//	case *luaClosure:
//		return fmt.Sprintf("closure %s:%d %v", v.prototype.source, v.prototype.lineDefined, v)
//	case *goClosure:
//		return fmt.Sprintf("go closure %#v", v)
//	case *goFunction:
//		pc := reflect.ValueOf(v.Function).Pointer()
//		f := runtime.FuncForPC(pc)
//		file, line := f.FileLine(pc)
//		return fmt.Sprintf("go function %s %s:%d", f.Name(), file, line)
//	case *userData:
//		return fmt.Sprintf("userdata %#v", v)
//	case nil:
//		return "nil"
//	case bool:
//		return fmt.Sprintf("%#v", v)
//	}
//	return fmt.Sprintf("unknown %#v %s", v, reflect.TypeOf(v).Name())
//}
func ZhenValueToString(v interface{}) (s string) {

	switch value := v.(type) {
	case ZNone:
		s = "未定义"
	case types.Nil:
		s = "空值"

	case bool:
		if bool(value) == true {
			s = "真"
		} else {
			s = "假"
		}
	case ZInt:
		s = strconv.Itoa(int(value))
	case ZFloat:
		s = strconv.FormatFloat(float64(value), 'f', -1, 32)
	case string:
		s = string(value)
	case ZFun:
		//todo 函数转换为字符串
		s = "函数"
		//fmt.Println("未知类型。", v.valueType)
		//var ws []string
		//for _, a := range v.valueArray {
		//	ws = append(ws, ZhenValueToString(a))
		//}
		//s = strings.Join(ws, ", ")
	case *ZObject:
		//todo 对象转换为字符串
		s = "对象"
		//var ws []string
		//for k, a := range v.valueTable {
		//	w := fmt.Sprintf("%s = %s", k, ZhenValueToString(a))
		//	ws = append(ws, w)
		//}
		//s = strings.Join(ws, ", ")
	//case ZvtFun:
	//	s = v.valueName
	default:
		fmt.Println("未知类型。", value)
	}
	return
}

//
//func GetZhenValueFromElement(item *etree.Element) (v ZhenValue, err error) {
//	valueType := item.SelectAttrValue("值类型", "")
//	value := item.SelectAttrValue("值", "")
//
//	if valueType == "未定义" {
//		v = NewValueNone()
//	} else if valueType == "空值" {
//		v = NewZhenValueNil()
//	} else if valueType == "数字" {
//		n, e := strconv.ParseFloat(value, 64)
//		if e != nil {
//			err = errors.New("解析错误:不能把文本转换为数字")
//			return
//		}
//		v = NewZhenValueNumber(ZFloat(n))
//	} else if valueType == "字符串" {
//		v = NewZhenValueString(ZhenValueString(value))
//	}
//
//	return
//}
