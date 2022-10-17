package zhen_0_03

type ParseFun func(node *Parser) ParseFun
type ParseEndFun func(node *Parser)

type NodeParserType uint8

const (
	NptUnknown NodeParserType = iota
	NptRoot
	NptLine
	NptChildLine
	NptOperator
	NptBracket
	NptColon
)

var NodeParserTypeNames = [...]string{
	NptUnknown:   "未设置",
	NptRoot:      "根目录",
	NptLine:      "行",
	NptChildLine: "子行",
	NptOperator:  "运算",
	NptBracket:   "括号",
	NptColon:     "冒号",
}

func (npt NodeParserType) String() string {
	return NodeParserTypeNames[npt]
}

type NodeParser struct {
	parserType  NodeParserType
	mainNode    Node
	parseFun    ParseFun
	parseEndFun ParseEndFun
	stack       []Node
	//startTokenNo  int
	//endTokenNo    int
	childLineNode *ChildLineNode
}

func NewNodeParser(nodeParserType NodeParserType, mainNode Node, parseFun ParseFun, parseEndFun ParseEndFun) (nodeParser *NodeParser) {
	nodeParser = &NodeParser{}
	nodeParser.parserType = nodeParserType
	nodeParser.mainNode = mainNode
	nodeParser.parseFun = parseFun
	nodeParser.parseEndFun = parseEndFun

	nodeParser.stack = make([]Node, 0, 20)
	return
}

func NewRootNodeParser(mainNode Node) (nodeParser *NodeParser) {
	return NewNodeParser(NptRoot, mainNode, parseRoot, parseRootEnd)
}

func NewLineNodeParser(lineNode Node) (nodeParser *NodeParser) {
	return NewNodeParser(NptLine, lineNode, parseLine, parseLineEnd)
}

func NewChildLineNodeParser(childLineNode Node) (nodeParser *NodeParser) {
	return NewNodeParser(NptChildLine, childLineNode, parseToken, parseChildLineEnd)
}

func NewOperatorNodeParser() (nodeParser *NodeParser) {
	return NewNodeParser(NptOperator, nil, parseOperator, parseOperatorEnd)
}

func NewBracketNodeParser(bracketNode Node) (nodeParser *NodeParser) {
	return NewNodeParser(NptBracket, bracketNode, parseToken, parseBracketEnd)
}

func NewColonNodeParser(colonNode Node) (nodeParser *NodeParser) {
	return NewNodeParser(NptColon, colonNode, parseColon, parseColonEnd)
}

func parseRoot(p *Parser) (parseFun ParseFun) {
	startParseLine(p)
	parseFun = parseRoot
	return
}

func parseRootEnd(p *Parser) {
	rootParse := p.popNodeParser()
	p.code.RootNode.Items = rootParse.stack
	return
}

func startParseLine(p *Parser) {
	lineNode := NewLineNode()
	lineParse := NewLineNodeParser(lineNode)
	//lineParse := p.NewNodeParser(lineNode, parseLine, parseLineEnd)
	p.lineIndex = 0
	p.pushNodeParser(lineParse)
	return
}

func checkLineWithColonEnd(lineNode *LineNode) (isColonLine bool) {
	if len(lineNode.Items) > 0 {
		lastItem := lineNode.Items[len(lineNode.Items)-1]
		colonNode, ok := lastItem.(*ColonNode)
		if ok {
			if len(colonNode.Items) == 0 {
				isColonLine = true
			}
		}
	}
	return
}
func (p *Parser) checkChildLineNode(parse *NodeParser) (hasChildLine bool) {
	//parse : p.getTopNode()
	//lineNode, ok := parse.mainNode.(*LineNode)
	//if ok {
	if parse.childLineNode != nil {
		parse.childLineNode.Items = append(parse.childLineNode.Items, parse.stack...)
		parse.stack = parse.stack[0:0]
		//mainNode, ok := parse.mainNode.(*li)
		//if ok {
		//	mainNode.addItem(parse.childLineNode)
		//}
		parse.mainNode.addItem(parse.childLineNode)
		parse.childLineNode = nil
		p.pushNodeParser(parse)
		hasChildLine = true
	}
	//}
	return
}

func parseLineEnd(p *Parser) {
	lineParse := p.popNodeParser()
	lineNode, ok := lineParse.mainNode.(*LineNode)
	if ok {
		lineNode.LineIndent = p.lineIndex

		if !p.checkChildLineNode(lineParse) {

			lineNode.addChildLineNode(lineParse.stack)

			if len(lineNode.Items) > 0 {
				lineNode.arrangeAllChildLine()
				//addStackToNode(&lineNode.BaseNode, lineParse)
				p.pushLineNode(lineNode)
				if checkLineWithColonEnd(lineNode) {
					p.lineWithColonEnd = append(p.lineWithColonEnd, lineNode)
				}

			}
		}
	}
	return
}

func parseChildLineEnd(p *Parser) {
	lineParse := p.popNodeParser()
	lineNode, ok := lineParse.mainNode.(*ChildLineNode)
	if ok {
		lineNode.Items = lineParse.stack

		if len(lineNode.Items) > 0 {
			p.pushNode(lineNode)

		}
	}
	return
}

func startParseOperator(p *Parser) {
	top := p.getTopNode()
	if isValueNode(top) {
		p.popNode()
	}
	operatorParser := NewOperatorNodeParser()
	//operatorParser := p.NewNodeParser(nil, parseOperator, parseOperatorEnd)
	p.pushNodeParser(operatorParser)
	if isValueNode(top) {
		p.pushNode(top)
	}

	return
}
func parseOperatorEnd(p *Parser) {
	for true {
		node := p.popOperator()
		if node == nil {
			break
		}
	}
	operatorParse := p.popNodeParser()
	p.pushNodes(operatorParse.stack)
}

func startParseBracket(p *Parser) {
	bracketNode := NewBracketNode(p.getNowToken())
	p.nextToken()
	bracketParser := NewBracketNodeParser(bracketNode)
	//bracketParser := p.NewNodeParser(bracketNode, parseToken, parseBracketEnd)
	p.pushNodeParser(bracketParser)
	p.pushBracket(bracketNode)
	return
}

func parseBracketEnd(p *Parser) {

	bracketParse := p.popNodeParser()
	bracketNode, ok := bracketParse.mainNode.(*BracketNode)
	if ok {
		//bracketNode.Items = bracketParse.stack
		//p.pushValueNode(bracketNode)

		if !p.checkChildLineNode(bracketParse) {
			p.nextToken()
			p.popBracket()
			bracketNode.addChildLineNode(bracketParse.stack)
			bracketNode.arrangeAllChildLine()
			p.pushValueNode(bracketNode)

		}
	}

	return
}

func startParseColon(p *Parser, t *Token) {
	p.nextToken()
	colonNode := NewColonNode(t)
	colonParser := NewColonNodeParser(colonNode)
	//colonParser := p.NewNodeParser(colonNode, parseColon, parseColonEnd)
	p.pushNodeParser(colonParser)
	return
}

func parseColonEnd(p *Parser) {
	colonParser := p.popNodeParser()
	colonNode, ok := colonParser.mainNode.(*ColonNode)
	if ok {
		if !p.checkChildLineNode(colonParser) {
			colonNode.addChildLineNode(colonParser.stack)
			colonNode.arrangeAllChildLine()
			p.pushNode(colonNode)
		}

		//if len(lineNode.Items) > 0 {
		//	p.pushLineNode(lineNode)
		//	if checkLineWithColonEnd(lineNode) {
		//		p.lineWithColonEnd = append(p.lineWithColonEnd, lineNode)
		//	}
		//}

	}
	return
}

func parseLine(p *Parser) (parseFun ParseFun) {
	parseFun = parseLine
	indent := 0
	token := p.getNowToken()
	switch token.TokenType {
	case TtSpace:
		indent = 1
		p.nextToken()
	case TtFullWidthSpace:
		indent = 2
		p.nextToken()
	case TtTab:
		indent = 4
		p.nextToken()
	default:
		parseFun = parseToken
	}

	p.lineIndex += indent

	return
}

func parseColon(p *Parser) (parseFun ParseFun) {
	parseFun = parseColon
	token := p.getNowToken()
	switch token.TokenType {
	case TtSpace, TtFullWidthSpace, TtTab:
		p.nextToken()
	case TtPound:
		p.addComment()
	case TtCR, TtLF, TtCRLF:
		parseFun = parseTokenCRLF(p, parseFun)
	default:
		parseFun = parseToken
	}
	return
}

func parseToken(ps *Parser) (parseFun ParseFun) {
	parseFun = parseToken
	token := ps.getNowToken()
	switch token.TokenType {
	case TtPoint, TtEqual, TtAdd, TtSub, TtMul, TtDiv, TtPower,
		TtEqualEqual, TtNotEqual, TtMoreThan,
		TtMoreThanEqual, TtLessThan, TtLessThanEqual:
		startParseOperator(ps)
	case TtLeftBracket, TtLeftSquareBracket, TtLeftBigBracket:
		startParseBracket(ps)
	case TtRightBracket, TtRightSquareBracket, TtRightBigBracket:
		bracketNode := ps.getTopBracket()
		if bracketNode != nil && bracketNode.RightBracket == token.TokenType {
			bracketNode.SetEndToken(token)
			//parseFun = nil
			ps.endNodeParser(NptBracket)
		} else {
			//todo 括号未匹配警告
			ps.nextToken()
		}
	case TtColon:
		startParseColon(ps, token)
	case TtPound:
		ps.addComment()
	case TtBackslash:
		ps.addBackslash()
		//ps.nextToken()
	case TtCR, TtLF, TtCRLF:
		parseFun = parseTokenCRLF(ps, parseFun)
	case TtSpace, TtFullWidthSpace, TtTab:
		parseFun = parseTokenSpace(ps, parseFun)
	case TtComma, TtDunHao, TtSemicolon, TtPeriod:
		parseFun = parseTokenChildLine(ps, parseFun)
	default:
		ps.pushNowToken()
	}

	return
}

func parseTokenChildLine(p *Parser, def ParseFun) (parseFun ParseFun) {
	parseFun = def

	var lineNodeParser *NodeParser
	childLineNode := NewChildLineNode(p.getNowToken())

	if len(p.brackets) > 0 {
		lineNodeParser = p.findNodeParser(NptBracket)
	} else {
		switch childLineNode.Symbol {
		case TtComma, TtDunHao:
			lineNodeParser = p.findNodeParser(NptColon, NptLine)
		case TtSemicolon, TtPeriod:
			lineNodeParser = p.findNodeParser(NptLine)
		}
	}

	if lineNodeParser != nil {
		lineNodeParser.childLineNode = childLineNode
		if len(p.brackets) > 0 {
			p.endNodeParser(NptBracket)
			parseFun = parseToken
		} else {
			switch childLineNode.Symbol {
			case TtComma, TtDunHao:
				//lineNodeParser = p.findNodeParser(NptColon, NptLine)
				p.endNodeParser(NptColon, NptLine)
				parseFun = parseToken
			case TtSemicolon, TtPeriod:
				//lineNodeParser = p.findNodeParser(NptLine)
				p.endNodeParser(NptLine)
				parseFun = parseToken
			}
		}

	}
	//p.pushNowToken()
	p.nextToken()

	////isEnd := false
	//nodeParser := p.getTopNode()
	//childLineNode, ok := nodeParser.mainNode.(*ChildLineNode)
	//if nodeParser.endTokenNo == p.tokenNo {
	//
	//	p.nextToken()
	//	parseFun = nil
	//} else {
	//	childLineNode := NewChildLineNode(p.getNowToken())
	//	childLineNodeParser := NewChildLineNodeParser(childLineNode)
	//	lineNodeParser := p.findNodeParser(NptLine)
	//	if lineNodeParser != nil {
	//		childLineNodeParser.startTokenNo = lineNodeParser.startTokenNo
	//	}
	//
	//	p.nextToken()
	//}

	return
}
func parseTokenCRLF(p *Parser, def ParseFun) (parseFun ParseFun) {
	parseFun = def
	if p.nowrap {
		p.nextToken()
		p.nowrap = false
		return
	}
	if len(p.brackets) == 0 {
		p.endNodeParser(NptLine)
		parseFun = parseToken
		//parseFun = nil
	}
	p.nextToken()
	return
}
func parseTokenSpace(p *Parser, def ParseFun) (parseFun ParseFun) {
	parseFun = def
	//todo 定义文本的界面保留空格
	p.nextToken()
	return
}

func parseOperator(ps *Parser) (parseFun ParseFun) {
	parseFun = parseOperator

	token := ps.getNowToken()
	switch token.TokenType {
	case TtLetter, TtInt, TtFloat,
		TtApostrophe, TtLeftApostrophe, TtQuotation, TtLeftQuotation:
		ok := ps.pushOperatorValue()
		if !ok {
			//parseFun = nil

			ps.endNodeParser(NptOperator)
		}
	case TtPoint, TtEqual, TtAdd, TtSub, TtMul, TtDiv, TtPower,
		TtEqualEqual, TtNotEqual, TtMoreThan, TtMoreThanEqual, TtLessThan, TtLessThanEqual:
		op := NewOperatorNode(ps.getNowToken())
		if op.OperatorType == OtSub {
			if ps.checkNegativeOperator() {
				op = ChangeOperatorNodeType(op, OtNegative)

			}
		}

		if !ps.pushOperator(op) {
			//todo  出错
			ps.nextToken()
		}

	case TtLeftBracket, TtLeftSquareBracket, TtLeftBigBracket:
		startParseBracket(ps)

	case TtRightBracket, TtRightSquareBracket, TtRightBigBracket:
		bracketNode := ps.getTopBracket()
		if bracketNode != nil && bracketNode.RightBracket == token.TokenType {
			bracketNode.SetEndToken(token)

			//parseFun = nil
			ps.endNodeParser(NptBracket)
		} else {
			//todo 括号未匹配警告
			ps.nextToken()
		}

	case TtPound:
		ps.addComment()
	case TtBackslash:
		ps.addBackslash()
		//ps.nowrap = true
		//ps.nextToken()
	case TtCR, TtLF, TtCRLF:
		parseFun = parseTokenCRLF(ps, parseFun)
	case TtSpace, TtFullWidthSpace, TtTab:
		parseFun = parseTokenSpace(ps, parseFun)
	case TtComma, TtDunHao, TtSemicolon, TtPeriod:
		parseFun = parseTokenChildLine(ps, parseFun)
	default:
		ps.endNodeParser(NptOperator)
		//ps.popNode()
		//parseFun = nil

	}

	return
}

//func startParseLine(ps *NodeParser) (parseFun ParseFun) {
//	lineNode := NewLineNode()
//	lineParse := ps.NewParserState(lineNode)
//	lineParse.StartParse(parseLine, nil)
//
//	lineNode.Items = lineParse.getStack()
//	ps.pushNode(lineNode)
//	parseFun = startParseLine
//	return
//}
//
//

//func checkBracketEnd(ps *NodeParser) (end bool) {
//	bracketNode, ok := ps.mainNode.(*BracketNode)
//	if ok {
//		token := ps.getNowToken()
//
//		if token.TokenType == bracketNode.RightBracket {
//			ps.nextToken()
//			end = true
//		} else {
//			switch token.TokenType {
//			case TtCR, TtLF, TtCRLF:
//				ps.nextToken()
//			}
//		}
//	}
//
//	return
//}
//

//
////func parseLine(ps *NodeParser) (parseFun ParseFun) {
////	parseFun = parseLineIndent
////	return
////}
//
//func parseBracket(ps *NodeParser) (parseFun ParseFun) {
//	parseFun = parseToken
//	return
//}
//
//
