package zhen

type CodePre struct {
	MainCodeBlock *CodeBlock
	state         *ZhenState
}

func NewCodePre(state *ZhenState, block *CodeBlock) (codePre CodePre) {
	codePre.MainCodeBlock = block
	codePre.state = state
	return
}

func (codePre *CodePre) Preprocess() (err error) {
	for _, code := range codePre.MainCodeBlock.Items {
		err = codePre.CodePreprocess(code)
		if err != nil {
			return
		}
	}
	return
}

func (codePre *CodePre) CodePreprocess(codeBlock *CodeBlock) (err error) {
	codePre.state.NowCodeBlock = codeBlock
	if codeBlock.WordType == CwtUnknown {
		switch codeBlock.BlockType {
		case CbtLetter:
			codePre.CheckLetter(codeBlock)
		case CbtNumber:
		case CbtString:

		}
	}

	for _, code := range codeBlock.Items {
		err = codePre.CodePreprocess(code)
		if err != nil {
			return
		}
	}
	return
}

func (codePre *CodePre) CheckLetter(codeBlock *CodeBlock) (err error) {
	isKeyWord, keyWord, err := codePre.CheckIsKeyWord(codeBlock)
	if err != nil {
		return
	}
	if isKeyWord {
		codeBlock.WordType = CwtKeyWord
		k := ZhenValueToKeyWord(keyWord)
		k.PreFun(codePre.state)
	}
	//_ = keyWord
	return
}

func (codePre *CodePre) CheckIsKeyWord(codeBlock *CodeBlock) (isKeyWord bool, keyWord ZhenValue, err error) {
	identifier := codeBlock.getChars()
	isKeyWord, keyWord, err = codePre.state.GetKeyWord(identifier)
	if err != nil {
		return
	}
	//if isKeyWord {
	//	codeBlock.KeyWords = identifier
	//	keyWord.valueFunction(codePre.state)
	//}

	return
}

//func (zhen *ZhenState) SetVarFunction(valueName string, value ZhenValue) {
//	//todo 局部变量如何处理
//	zhen.globalValues[valueName] = value
//}
