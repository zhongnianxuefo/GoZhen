package zhen_0_03

type OperatorType uint8

const (
	OtUnSet OperatorType = iota
	OtPoint
	OtPower
	OtEqual
	OtMul
	OtDiv
	OtAdd
	OtSub
	OtNegative
	OtEqualEqual
	OtNotEqual
	OtMoreThan
	OtMoreThanEqual
	OtLessThan
	OtLessThanEqual
)

var OperatorTypeNames = [...]string{
	OtUnSet:         "未设置",
	OtPoint:         "点",
	OtPower:         "幂",
	OtNegative:      "负",
	OtMul:           "乘",
	OtDiv:           "除",
	OtAdd:           "加",
	OtSub:           "减",
	OtEqual:         "等于",
	OtEqualEqual:    "相等",
	OtNotEqual:      "不等于",
	OtMoreThan:      "大于",
	OtMoreThanEqual: "大于等于",
	OtLessThan:      "小于",
	OtLessThanEqual: "小于等于",
}

func (ot OperatorType) String() string {
	return OperatorTypeNames[ot]
}

var OperatorTypeWords = [...]string{
	OtUnSet:         "",
	OtPoint:         ".",
	OtPower:         "^",
	OtNegative:      "-",
	OtMul:           "*",
	OtDiv:           "/",
	OtAdd:           "+",
	OtSub:           "-",
	OtEqual:         "=",
	OtEqualEqual:    "==",
	OtNotEqual:      "≠",
	OtMoreThan:      ">",
	OtMoreThanEqual: "≥",
	OtLessThan:      "<",
	OtLessThanEqual: "≤",
}

func (ot OperatorType) Words() string {
	return OperatorTypeWords[ot]
}

var OperatorLeftPriority = map[OperatorType]int{
	OtUnSet:         0,
	OtPoint:         200,
	OtPower:         150,
	OtNegative:      140,
	OtMul:           130,
	OtDiv:           130,
	OtAdd:           120,
	OtSub:           120,
	OtEqual:         100,
	OtEqualEqual:    50,
	OtNotEqual:      50,
	OtMoreThan:      50,
	OtMoreThanEqual: 50,
	OtLessThan:      50,
	OtLessThanEqual: 50,
}

func getOperatorLeftPriority(o OperatorType) (p int) {
	p, _ = OperatorLeftPriority[o]

	return
}

var OperatorRightPriority = map[OperatorType]int{

	OtNegative: 160,
}

func getOperatorRightPriority(o OperatorType) (p int) {
	p, ok := OperatorRightPriority[o]
	if !ok {
		p = getOperatorLeftPriority(o)
	}
	return
}

func getOperatorNeedItems(o OperatorType) (items int) {
	if o == OtNegative {
		items = 1
	} else if o == OtUnSet {
		//todo 警告
	} else {
		items = 2
	}

	return
}

var TokenTypeOperatorType = map[TokenType]OperatorType{
	TtEqual:         OtEqual,
	TtAdd:           OtAdd,
	TtSub:           OtSub,
	TtMul:           OtMul,
	TtDiv:           OtDiv,
	TtPower:         OtPower,
	TtNegative:      OtNegative,
	TtEqualEqual:    OtEqualEqual,
	TtNotEqual:      OtNotEqual,
	TtMoreThan:      OtMoreThan,
	TtMoreThanEqual: OtMoreThanEqual,
	TtLessThan:      OtLessThan,
	TtLessThanEqual: OtLessThanEqual,
	TtPoint:         OtPoint,
}

func getOperatorType(t TokenType) (p OperatorType) {
	p, _ = TokenTypeOperatorType[t]

	return
}
