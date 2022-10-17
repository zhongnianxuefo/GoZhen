package zhen_0_03

type Parser struct {
	code *Code

	tokenNo   int
	stack     []*NodeParser
	endParser *NodeParser

	lineIndex        int
	nowrap           bool
	brackets         []*BracketNode
	lineWithColonEnd []*LineNode
	childLineNode    *ChildLineNode
}

//type ParserError struct {
//	name string
//}
//
//func (e *ParserError) Error() string {
//	return "name 不能为空"
//}

func NewParser(code *Code) (p *Parser) {
	p = &Parser{}
	p.code = code

	p.code.RootNode = NewRootNode()

	rootParser := NewRootNodeParser(p.code.RootNode)
	p.pushNodeParser(rootParser)
	return
}

func (p *Parser) Parse() (err error) {
	for p.tokenNo < len(p.code.Tokens) {
		//fmt.Println(p.tokenNo, p.getNowToken())
		nodeParser := p.getTopNodeParser()
		if nodeParser == nil || nodeParser.parseFun == nil {
			//fmt.Println("分析函数为空")
			panic("分析函数为空")
		} else {
			nodeParser.parseFun = nodeParser.parseFun(p)
		}

		if nodeParser.parseFun == nil || p.endParser != nil {
			p.endNodeParse()
		}
	}
	for len(p.stack) > 0 {
		p.endParser = p.stack[0]
		p.endNodeParse()
	}

	return
}

func (p *Parser) endNodeParse() (err error) {
	nodeParser := p.getTopNodeParser()
	nodeParser.parseEndFun(p)
	if nodeParser == p.endParser || p.endParser == nil {
		p.endParser = nil
		return
	}
	for i := len(p.stack) - 1; i >= 0; i-- {
		np := p.stack[i]
		np.parseEndFun(p)
		if np == p.endParser {
			p.endParser = nil
			return
		}
	}
	return
}

func (p *Parser) getNodeParser(n int) (nodeParser *NodeParser) {
	stackLen := len(p.stack)
	if n >= 0 {
		if n < stackLen {
			nodeParser = p.stack[n]
		}
	} else {
		n = stackLen + n
		if n >= 0 {
			nodeParser = p.stack[n]
		}
	}
	return
}

func (p *Parser) getTopNodeParser() (nodeParser *NodeParser) {
	stackLen := len(p.stack)
	if stackLen > 0 {
		nodeParser = p.stack[stackLen-1]
	}
	return
}

func (p *Parser) findNodeParser(parserTypes ...NodeParserType) (nodeParser *NodeParser) {
	for i := len(p.stack) - 1; i >= 0; i-- {
		np := p.stack[i]
		blnFind := false
		for _, parserType := range parserTypes {
			if np.parserType == parserType {
				blnFind = true
			}
		}
		if blnFind {
			nodeParser = np
			return
		}
	}
	return
}

func (p *Parser) findNodeParserNo(parserTypes ...NodeParserType) (no int) {
	no = -1
	for i := len(p.stack) - 1; i >= 0; i-- {
		np := p.stack[i]
		blnFind := false
		for _, parserType := range parserTypes {
			if np.parserType == parserType {
				blnFind = true
			}
		}
		if blnFind {
			no = i
			return
		}
	}
	return
}

func (p *Parser) endNodeParser(parserTypes ...NodeParserType) {
	p.endParser = p.findNodeParser(parserTypes...)
	return
}

func (p *Parser) topNodeParserIsLine() (isLine bool) {
	nodeParser := p.getTopNodeParser()
	if nodeParser != nil {
		_, isLine = nodeParser.mainNode.(*LineNode)
	}
	return
}

func (p *Parser) popNodeParser() (nodeParser *NodeParser) {
	stackLen := len(p.stack)
	if stackLen > 0 {
		nodeParser = p.stack[stackLen-1]
		p.stack = p.stack[0 : stackLen-1]
	}

	return
}

func (p *Parser) pushNodeParser(nodeParser *NodeParser) {
	//nodeParser.startTokenNo = p.tokenNo
	p.stack = append(p.stack, nodeParser)
	return
}

func (p *Parser) addComment() {
	p.code.Comments = append(p.code.Comments, p.getNowToken())
	p.tokenNo += 1
	return
}

func (p *Parser) addBackslash() {
	p.code.Backslashes = append(p.code.Backslashes, p.getNowToken())
	p.nowrap = true
	p.tokenNo += 1
	return
}

func (p *Parser) getTopNode() (node Node) {
	ps := p.getTopNodeParser()
	stackLen := len(ps.stack)
	if stackLen > 0 {
		node = ps.stack[stackLen-1]
	}
	return
}

func (p *Parser) getNode(n int) (node Node) {
	ps := p.getTopNodeParser()
	stackLen := len(ps.stack)
	if n >= 0 {
		if n < stackLen {
			node = ps.stack[n]
		}
	} else {
		n = stackLen + n
		if n >= 0 {
			node = ps.stack[n]
		}
	}
	return
}

func (p *Parser) popNode() (node Node) {
	ps := p.getTopNodeParser()
	stackLen := len(ps.stack)
	if stackLen > 0 {
		node = ps.stack[stackLen-1]
		ps.stack = ps.stack[0 : stackLen-1]
	}
	return
}

func (p *Parser) pushNode(node Node) {
	if node == nil {
		return
	}
	ps := p.getTopNodeParser()
	ps.stack = append(ps.stack, node)

	return
}

func (p *Parser) pushNodes(nodes []Node) {
	for _, node := range nodes {
		p.pushNode(node)
	}
	return
}

func (p *Parser) getNowToken() (token *Token) {
	token = p.code.Tokens[p.tokenNo]
	return
}

func (p *Parser) nextToken() {
	p.tokenNo += 1
	return
}

func (p *Parser) pushNowToken() {
	t := p.code.Tokens[p.tokenNo]
	node := p.nodeFromToken(t)
	if node != nil {
		p.pushNode(node)
	}
	p.nextToken()
	return
}

func (p *Parser) nodeFromToken(t *Token) (node Node) {
	switch t.TokenType {
	case TtLetter:
		letterNode := NewLetterNode(t)
		//p.checkLetterNode(letterNode)
		node = letterNode
	case TtApostrophe, TtLeftApostrophe, TtQuotation, TtLeftQuotation:
		node = NewStringNode(t)
	case TtInt:
		node = NewIntNode(t)
	case TtFloat:
		node = NewFloatNode(t)
	case TtSpace, TtFullWidthSpace, TtTab:
		node = NewEmptyNode(t)
	case TtOtherChar:
		node = NewOtherNode(t)
	}

	return
}

//
//func (p *Parser) defineLetterNode(node *LetterNode, defineType KeyWordType) {
//	switch defineType {
//	case KwtDefineConstant:
//		node.LetterType = LtConstant
//		p.defineVar(node, defineType, -1)
//	}
//	return
//}

//func (p *Parser) setEqualOperatorDefine(operatorNode *OperatorNode, defineType KeyWordType) {
//
//	return
//}
//
//func (p *Parser) setChildLineDefine(childLineNode *ChildLineNode, defineType KeyWordType) {
//	for _, item := range childLineNode.Items {
//		switch node := item.(type) {
//		case *LetterNode:
//			p.defineLetterNode(node, defineType)
//		case *OperatorNode:
//			switch node.OperatorType {
//			case OtEqual:
//				p.setEqualOperatorDefine(node, defineType)
//			}
//		}
//	}
//	return
//}
//
//func (p *Parser) setNodeDefine(node Node, defineType KeyWordType) {
//	switch n := node.(type) {
//	case *LetterNode:
//		switch defineType {
//		case KwtDefineConstant:
//			n.LetterType = LtConstant
//			p.defineVar(n, defineType, -1)
//		}
//	case *OperatorNode:
//		if n.OperatorType == OtEqual {
//			if len(n.Items) > 0 {
//				i, ok := n.Items[0].(*LetterNode)
//				if ok {
//					p.defineLetterNode(i, defineType)
//				}
//			}
//		}
//	case *ChildLineNode:
//		for _, item := range n.Items {
//			p.setNodeDefine(item, defineType)
//		}
//	case *LineNode:
//		for _, item := range n.Items {
//			p.setNodeDefine(item, defineType)
//		}
//	case *ColonNode:
//		for _, item := range n.Items {
//			p.setNodeDefine(item, defineType)
//		}
//	}
//
//	return
//}

//func (p *Parser) setColonNodeDefine(colonNode *ColonNode, defineType KeyWordType) {
//	for _, item := range colonNode.Items {
//		switch node := item.(type) {
//		case *ChildLineNode:
//			p.setChildLineDefine(node, defineType)
//		case *LetterNode:
//			p.defineLetterNode(node, defineType)
//		case *OperatorNode:
//			switch node.OperatorType {
//			case OtEqual:
//				p.setEqualOperatorDefine(node, defineType)
//			}
//		}
//	}
//
//	return
//}

//func (p *Parser) checkLineNodeItems(lineNode *LineNode) {
//	var keyWords []string
//	var keyWordType KeyWordType
//	//var defineName string
//
//	for _, item := range lineNode.Items {
//		switch node := item.(type) {
//		case *LetterNode:
//			keyWords = append(keyWords, node.Words)
//			letter, ok := p.getLetterInfo(node.Words, -1)
//			if ok {
//				node.LetterType = letter.Type
//				switch letter.Type {
//				case LtKeyWord:
//					keyWordType = letter.Data.(KeyWordType)
//
//				}
//			}
//		case *ColonNode:
//			if keyWordType == KwtUnknown {
//				keyWordType = KwtDefine
//				//defineName = strings.Join(keyWords, " ")
//			}
//			switch keyWordType {
//			case KwtDefineConstant, KwtDefineLocalVar, KwtDefineGlobalVar:
//				p.setNodeDefine(node, keyWordType)
//			}
//			break
//		}
//	}
//
//}

//func (p *Parser) checkLetterNode(node *LetterNode) {
//	letter, ok := p.getLetterInfo(node.Words, -1)
//	if ok {
//		node.LetterType = letter.Type
//		switch letter.Type {
//		case LtKeyWord:
//			lineNode := p.getLineNode(-1)
//			if lineNode != nil {
//				//lineNode.KeyWordType = letter.Data.(KeyWordType)
//			}
//		}
//	} else {
//		//lineNode := p.getLineNode(-1)
//		//if lineNode != nil {
//		//	switch lineNode.KeyWordType {
//		//	case KwtDefineConstant:
//		//		rootNode := p.RootNode
//		//		rootNode.LetterSet.addConstant(node)
//		//
//		//	}
//		//}
//	}
//
//	return
//}

//func (p *Parser) getLineNode(n int) (lineNode *LineNode) {
//	ps := p.getNodeParser(n)
//	if ps == nil {
//		return
//	}
//	node, ok := ps.mainNode.(*LineNode)
//	if ok {
//		lineNode = node
//		return
//	}
//
//	lineNode = p.getLineNode(n - 1)
//	return
//}

//func (p *Parser) getLetterInfo(name string, n int) (letter *Letter, ok bool) {
//	ps := p.getNodeParser(n)
//	if ps == nil {
//		return
//	}
//
//	switch mainNode := ps.mainNode.(type) {
//	case *RootNode:
//		letter, ok = mainNode.LetterSet.getByName(name)
//	case *LineNode:
//		letter, ok = mainNode.LetterSet.getByName(name)
//	}
//	if ok == false {
//		letter, ok = p.getLetterInfo(name, n-1)
//	}
//	return
//}

//func (p *Parser) defineVar(letterNode *LetterNode, defineType KeyWordType, n int) {
//	switch defineType {
//	case KwtDefineConstant:
//		p.code.RootNode.LetterSet.addConstant(letterNode)
//
//	}
//
//	return
//}

func (p *Parser) pushValueNode(node Node) (pushOk bool) {
	//node = p.checkNegativeNode(node)
	pushOk = false
	top := p.getTopNode()
	op, ok := top.(*OperatorNode)
	if ok {
		if op.needItem() {
			op.addItem(node)
			pushOk = true
		}
	}
	if !pushOk {
		p.pushNode(node)
	}
	return
}

func (p *Parser) pushOperatorValue() (pushOk bool) {
	node := p.nodeFromToken(p.getNowToken())

	//node = p.checkNegativeNode(node)

	top := p.getTopNode()
	op, ok := top.(*OperatorNode)
	if ok {
		if op.needItem() {
			op.addItem(node)
			p.nextToken()
			pushOk = true
		}
	} else {
		panic("当前堆中最上面一个节点不是运算节点")
	}
	return
}

func (p *Parser) replaceOperator(operator *OperatorNode) {
	operator.addItem(p.popNode())
	p.pushNode(operator)
	p.nextToken()
	return
}

func (p *Parser) pushOperator(operator *OperatorNode) (pushOk bool) {
	top := p.getTopNode()
	if isValueNode(top) {
		p.replaceOperator(operator)
		pushOk = true
		return
	}
	topOperator, ok := top.(*OperatorNode)
	if !ok {
		p.pushNode(operator)
		pushOk = true
		return
	}
	if topOperator.LeftPriority == operator.RightPriority {
		p.replaceOperator(operator)
		pushOk = true
	} else if topOperator.LeftPriority < operator.RightPriority {
		if operator.NeedItems > 1 {
			l := len(topOperator.Items)
			if l > 0 {
				operator.addItem(topOperator.Items[l-1])
				topOperator.Items = topOperator.Items[0 : l-1]
			}
		}
		p.pushNode(operator)
		p.nextToken()
		pushOk = true
	} else {
		if p.popOperator() != nil {
			pushOk = true
		}
	}

	return
}

func (p *Parser) popOperator() (node Node) {
	top1, ok1 := p.getNode(-1).(*OperatorNode)
	if !ok1 {
		return
	}
	top2, ok2 := p.getNode(-2).(*OperatorNode)
	if !ok2 {
		return
	}
	node = p.popNode()
	top2.addItem(top1)

	return
}

func (p *Parser) pushLineNode(lineNode *LineNode) {
	nowLineIndent := lineNode.LineIndent
	for i := len(p.lineWithColonEnd) - 1; i >= 0; i-- {
		colonLine := p.lineWithColonEnd[i]
		if nowLineIndent > colonLine.LineIndent {
			colonNode := colonLine.Items[len(colonLine.Items)-1]
			//cn, ok := colonNode.(*ColonNode)
			//if ok {
			//	cn.addItem(lineNode)
			//}

			colonNode.addItem(lineNode)
			//p.prepLineNode(lineNode)
			return
		} else {
			p.lineWithColonEnd = p.lineWithColonEnd[0:i]
		}
	}
	//p.prepLineNode(lineNode)
	p.pushNode(lineNode)
	return
}

func (p *Parser) getTopBracket() (bracket *BracketNode) {
	bracketLen := len(p.brackets)
	if bracketLen > 0 {
		bracket = p.brackets[bracketLen-1]
	}
	return
}

func (p *Parser) pushBracket(bracket *BracketNode) {
	p.brackets = append(p.brackets, bracket)
	return
}

func (p *Parser) popBracket() (bracket *BracketNode) {
	bracketLen := len(p.brackets)
	if bracketLen > 0 {
		bracket = p.brackets[bracketLen-1]
		p.brackets = p.brackets[0 : bracketLen-1]
	}

	return
}

func (p *Parser) checkNegativeOperator() (isNegative bool) {
	isNegative = true
	node := p.getTopNode()
	if node == nil {
		isNegative = true
	} else if isValueNode(node) {
		isNegative = false
	} else {
		o, ok := node.(*OperatorNode)
		if ok {
			if o.needItem() == false {
				isNegative = false
			}
		}
	}
	return

	//if isNegative {
	//	op = ChangeOperatorNodeType(op, OtNegative)
	//}
}
