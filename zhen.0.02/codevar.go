package zhen_0_02

type CodeVarType byte

const (
	CvtUnknown CodeVarType = iota
	CvtKeyWords
	CvtOperator
	CvtConstant
	CvtFunction
	CvtGlobalVar
	CvtLocalVar
)

var CodeVarTypeNames = [...]string{
	CvtUnknown:   "未知",
	CvtKeyWords:  "关键字",
	CvtOperator:  "运算符",
	CvtConstant:  "常量",
	CvtFunction:  "函数",
	CvtGlobalVar: "全局变量",
	CvtLocalVar:  "本地变量",
}

func (cvt CodeVarType) String() string {
	return CodeVarTypeNames[cvt]
}

type CodeVarKey struct {
	Name string
	Type CodeVarType
}

type CodeVars struct {
	Names map[CodeVarKey]int
	Count int
}

func NewVarNames() (varNames CodeVars) {
	varNames.Names = make(map[CodeVarKey]int)
	return
}

func (vars *CodeVars) AddVar(varName string, varType CodeVarType) {
	key := CodeVarKey{Name: varName, Type: varType}
	_, ok := vars.Names[key]
	if !ok {
		n := vars.Count
		vars.Names[key] = n
		vars.Count += 1
	}
	return
}

func (vars *CodeVars) GetVarNo(varName string, varType CodeVarType) (no int, ok bool) {
	key := CodeVarKey{Name: varName, Type: varType}
	no, ok = vars.Names[key]
	return
}

//
//type CodeArea struct {
//	Names map[CodeVarKey]int
//	Count int
//}
