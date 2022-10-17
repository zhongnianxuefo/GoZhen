package zhen_0_02

type OperatorType int8

const (
	OtUnSet OperatorType = iota
	OtAs
	OtAdd
	OtSub
	OtMul
	OtDiv
	OtPoint
)

var OperatorTypeNames = [...]string{

	OtUnSet: "未设置",
	OtAs:    "赋值",
	OtAdd:   "加",
	OtSub:   "减",
	OtMul:   "乘",
	OtDiv:   "除",
	OtPoint: ".",
}

func (ot OperatorType) String() string {
	return OperatorTypeNames[ot]
}

type Operator struct {
	//Names     string
	Type     OperatorType
	Priority int
}
