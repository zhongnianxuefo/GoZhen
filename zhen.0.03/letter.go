package zhen_0_03

type LetterType int

const (
	LtUnknown LetterType = iota
	LtKeyWord
	LtDefine
	LtConstant
	LtVar
	LtGlobalVar
	LtLocalVar
	LtFun
	LtFunPara
	LtFunReturn

	LtCallFun
)

var LetterTypeNames = [...]string{
	LtUnknown:   "未设置",
	LtKeyWord:   "关键字",
	LtDefine:    "定义",
	LtConstant:  "常量",
	LtVar:       "变量",
	LtGlobalVar: "全局变量",
	LtLocalVar:  "本地变量",
	LtFun:       "函数",
	LtFunPara:   "参数",
	LtFunReturn: "返回值",

	LtCallFun: "运行",
}

func (lt LetterType) String() string {
	return LetterTypeNames[lt]
}

type Letter struct {
	Name   string
	Type   LetterType
	Words  []string
	Format []string
	Node   Node
	Data   interface{}
}

type LetterWords struct {
	Word  string
	Names map[string]bool
}

type LetterSet struct {
	Names map[string]*Letter
	Words map[string]*LetterWords
}

func NewLetter(name string, letterType LetterType) (letter *Letter) {
	letter = &Letter{Name: name}
	letter.Name = name
	letter.Type = letterType
	return
}

func NewLetterWords(word string) (letterWords *LetterWords) {
	letterWords = &LetterWords{}
	letterWords.Word = word
	letterWords.Names = make(map[string]bool)
	return
}

func NewLetterSet() (set *LetterSet) {
	set = &LetterSet{}
	set.Names = make(map[string]*Letter)
	set.Words = make(map[string]*LetterWords)
	return
}

func (s *LetterSet) del(keyWord *Letter) {
	name := keyWord.Name
	delete(s.Names, name)
	for _, word := range keyWord.Words {
		words, ok := s.Words[word]
		if ok {
			delete(words.Names, name)
		}
	}
	return
}

func (s *LetterSet) add(keyWord *Letter) {
	name := keyWord.Name
	k, ok := s.getByName(name)
	if ok {
		s.del(k)
	}
	s.Names[name] = keyWord
	for _, word := range keyWord.Words {
		words, wordOk := s.Words[word]
		if wordOk {
			words.Names[name] = true
		} else {
			words = NewLetterWords(word)
			words.Names[name] = true
			s.Words[word] = words
		}
	}
	return
}

func (s *LetterSet) getByName(name string) (keyWord *Letter, ok bool) {
	keyWord, ok = s.Names[name]
	return
}

func (s *LetterSet) addKeyWord(name string, wordType KeyWordType) {
	keyWord := NewKeyWord(name, wordType)
	s.add(keyWord)
	return
}

func (s *LetterSet) addConstant(letterNode *LetterNode) {
	letter := &Letter{}
	letter.Name = letterNode.Words
	letter.Type = LtConstant
	letter.Node = letterNode

	s.add(letter)
	return
}

func (s *LetterSet) addVar(letterNode *LetterNode) {
	letter := &Letter{}
	letter.Name = letterNode.Words
	letter.Type = LtVar
	letter.Node = letterNode

	s.add(letter)
	return
}

func (s *LetterSet) addLocalVar(letterNode *LetterNode) {
	letter := &Letter{}
	letter.Name = letterNode.Words
	letter.Type = LtLocalVar
	letter.Node = letterNode

	s.add(letter)
	return
}

func (s *LetterSet) addGlobalVar(letterNode *LetterNode) {
	letter := &Letter{}
	letter.Name = letterNode.Words
	letter.Type = LtGlobalVar
	letter.Node = letterNode

	s.add(letter)
	return
}
