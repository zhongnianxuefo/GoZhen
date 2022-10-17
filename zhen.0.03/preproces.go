package zhen_0_03

type Preprocess struct {
	code  *Code
	stack []Node

	keyWords    []string
	keyWordType KeyWordType
	defineType  KeyWordType
}

//type ParserError struct {
//	name string
//}
//
//func (e *ParserError) Error() string {
//	return "name 不能为空"
//}

func NewPreprocess(code *Code) (p *Preprocess) {
	p = &Preprocess{}
	p.code = code

	return
}

func (p *Preprocess) Preprocess() (err error) {
	rootNode := p.code.RootNode
	p.pushNode(rootNode)
	for _, item := range rootNode.Items {
		line, ok := item.(*LineNode)
		if ok {
			p.prepLineNode(line)
		}
	}
	p.popNode()
	return
}

func (p *Preprocess) pushNode(node Node) {
	p.stack = append(p.stack, node)
	return
}

func (p *Preprocess) popNode() (node Node) {
	stackLen := len(p.stack)
	if stackLen > 0 {
		node = p.stack[stackLen-1]
		p.stack = p.stack[0 : stackLen-1]
	}
	return
}

func (p *Preprocess) getStackNode(n int) (node Node) {
	stackLen := len(p.stack)
	if n >= 0 {
		if n < stackLen {
			node = p.stack[n]
		}
	} else {
		n = stackLen + n
		if n >= 0 {
			node = p.stack[n]
		}
	}
	return
}

func (p *Preprocess) getTopNode() (node Node) {
	stackLen := len(p.stack)
	if stackLen > 0 {
		node = p.stack[stackLen-1]
	}
	return
}

func (p *Preprocess) getLetterInfo(name string, n int) (letter *Letter, ok bool) {
	node := p.getStackNode(n)
	if node == nil {
		return
	}

	switch mainNode := node.(type) {
	case *RootNode:
		letter, ok = mainNode.LetterSet.getByName(name)
	case *LineNode:
		letter, ok = mainNode.LetterSet.getByName(name)
	}
	if ok == false {
		letter, ok = p.getLetterInfo(name, n-1)
	}
	return
}

func (p *Preprocess) prepLineNode(lineNode *LineNode) {
	p.pushNode(lineNode)
	p.keyWords = []string{}
	p.keyWordType = KwtUnknown
	p.prepNodes(lineNode.Items)
	p.popNode()
	return
}

func (p *Preprocess) prepChildLineNode(childLineNode *ChildLineNode) {
	p.pushNode(childLineNode)
	switch childLineNode.Symbol {
	case TtSemicolon, TtPeriod:
		p.keyWords = []string{}
		p.keyWordType = KwtUnknown
	}
	p.prepNodes(childLineNode.Items)
	p.popNode()
}

func (p *Preprocess) prepColonNode(n *ColonNode) {
	if p.keyWordType == KwtUnknown {
		p.keyWordType = KwtDefine
		//defineName = strings.Join(keyWords, " ")
	}
	switch p.keyWordType {
	case KwtDefineConstant, KwtDefineVar, KwtDefineLocalVar, KwtDefineGlobalVar:
		p.defineType = p.keyWordType
		p.defineNodes(n.Items)
	default:
		p.prepNodes(n.Items)
	}
	return
}

func (p *Preprocess) prepOperatorNode(operatorNode *OperatorNode) {
	p.prepNodes(operatorNode.Items)
	return
}

func (p *Preprocess) prepLetterNode(n *LetterNode) {
	p.keyWords = append(p.keyWords, n.Words)
	letter, ok := p.getLetterInfo(n.Words, -1)
	if ok {
		n.LetterType = letter.Type
		switch letter.Type {
		case LtKeyWord:
			p.keyWordType = letter.Data.(KeyWordType)
		}
	}
	return
}

func (p *Preprocess) prepNodes(nodes []Node) {
	for _, item := range nodes {
		switch node := item.(type) {
		case *LetterNode:
			p.prepLetterNode(node)
		case *ColonNode:
			p.prepColonNode(node)
		case *OperatorNode:
			p.prepOperatorNode(node)
		case *ChildLineNode:
			p.prepChildLineNode(node)
		case *LineNode:
			p.prepLineNode(node)
		}
	}
}

func (p *Preprocess) prepFunLineNode(lineNode *LineNode) {

}

func (p *Preprocess) defineNodes(nodes []Node) {
	for _, item := range nodes {
		switch node := item.(type) {
		case *LetterNode:
			p.defineLetterNode(node)
		case *OperatorNode:
			if node.OperatorType == OtEqual {
				if len(node.Items) > 0 {
					i, ok := node.Items[0].(*LetterNode)
					if ok {
						p.defineLetterNode(i)
					}
				}
			}
		//case *ColonNode:
		//	p.prepColonNode(node)
		case *ChildLineNode:
			p.defineNodes(node.Items)
		case *LineNode:
			p.defineNodes(node.Items)
		}
	}
}

func (p *Preprocess) defineLetterNode(node *LetterNode) {
	defineType := p.defineType
	switch defineType {
	case KwtDefineConstant:
		node.LetterType = LtConstant
		p.defineVar(node, defineType)
	case KwtDefineVar:
		node.LetterType = LtVar
		p.defineVar(node, defineType)
	case KwtDefineLocalVar:
		node.LetterType = LtLocalVar
		p.defineVar(node, defineType)
	case KwtDefineGlobalVar:
		node.LetterType = LtGlobalVar
		p.defineVar(node, defineType)
	}
	return
}

func (p *Preprocess) defineVar(letterNode *LetterNode, defineType KeyWordType) {
	switch defineType {
	case KwtDefineConstant:
		p.code.RootNode.LetterSet.addConstant(letterNode)
	case KwtDefineVar:
		p.code.RootNode.LetterSet.addVar(letterNode)
	case KwtDefineLocalVar:
		p.code.RootNode.LetterSet.addLocalVar(letterNode)
	case KwtDefineGlobalVar:
		p.code.RootNode.LetterSet.addGlobalVar(letterNode)
	}

	return
}
