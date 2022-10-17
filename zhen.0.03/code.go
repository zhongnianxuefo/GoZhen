package zhen_0_03

import (
	"encoding/gob"
	"os"
)

type Code struct {
	Tokens      []*Token
	RootNode    *RootNode
	Comments    []*Token
	Backslashes []*Token

	txt string
}

func NewCode(codeTxt string) (code *Code) {
	code = &Code{}
	code.txt = codeTxt
	return
}

func (code *Code) Scan() (err error) {
	scaner := NewScan(code)
	return scaner.Scan()
}

func (code *Code) Parse() (err error) {
	parser := NewParser(code)
	return parser.Parse()
}

func (code *Code) Preprocess() (err error) {
	preprocess := NewPreprocess(code)
	return preprocess.Preprocess()
}

func (code *Code) ToXmlFile(path string) (err error) {
	toXml := NewCodeToXml(code)
	return toXml.ToXmlFile(path)
}

func (code *Code) FormatToFile(path string) (err error) {
	format := NewCodeBlockFormat(code)
	return format.ToFile(path)
}

func (code *Code) ToGobFile(gobFile string) (err error) {
	file, err := os.Create(gobFile)
	if err != nil {
		return
	}
	defer file.Close()

	gob.Register(&BaseNode{})
	gob.Register(&RootNode{})
	gob.Register(&LineNode{})
	gob.Register(&ChildLineNode{})
	gob.Register(&ColonNode{})
	gob.Register(&OperatorNode{})
	gob.Register(&IntNode{})
	gob.Register(&FloatNode{})
	gob.Register(&StringNode{})
	gob.Register(&LetterNode{})
	gob.Register(&EmptyNode{})
	gob.Register(&BracketNode{})
	gob.Register(&OtherNode{})

	gob.Register(&Token{})
	gob.Register(&Letter{})
	gob.Register(&LetterWords{})
	gob.Register(&LetterSet{})

	gob.Register(KwtUnknown)

	enc := gob.NewEncoder(file)
	err = enc.Encode(code)
	if err != nil {
		return
	}
	return
}

func NewCodeFromGobFile(gobFile string) (code *Code, err error) {
	file, err := os.Open(gobFile)
	if err != nil {
		return
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	code = &Code{}
	err = dec.Decode(&code)
	if err != nil {
		return
	}
	return

}
