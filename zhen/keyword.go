package zhen

type KeyWord struct {
	Name   string
	Type   KeyWordType
	PreFun ZhenValueFunction
}

type KeyWordType int

const (
	KwtUnknown KeyWordType = iota
	KwtText
	KwtConstant
	KwtFun
)

func StringToKeyWordType(s string) (t KeyWordType) {
	switch s {
	case "文本":
		t = KwtText
	case "常量":
		t = KwtConstant
	default:
		t = KwtUnknown
	}
	return
}

func KeyWordTypeToString(t KeyWordType) (s string) {
	switch t {
	case KwtText:
		s = "文本"
	case KwtConstant:
		s = "常量"
	default:
		s = "未知"
	}
	return
}

func KeyWordToZhenValue(keyWord KeyWord) (value ZhenValue) {
	table := make(ZhenValueTable)
	table["名称"] = StringToZhenValue(keyWord.Name)
	table["类型"] = StringToZhenValue(KeyWordTypeToString(keyWord.Type))
	table["预处理函数"] = NewZhenValueFunction(keyWord.PreFun)
	value = NewZhenValueTable(table)
	return
}
func ZhenValueToKeyWord(value ZhenValue) (keyWord KeyWord) {
	if value.valueType == ZhenValueTypeTable {
		table := value.valueTable
		keyWord.Name = string(table["名称"].valueString)
		keyWord.Type = StringToKeyWordType(string(table["类型"].valueString))
		keyWord.PreFun = table["预处理函数"].valueFunction
	}

	return
}
