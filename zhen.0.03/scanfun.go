package zhen_0_03

type ScanFun func(scaner *Scaner) ScanFun

func scanNew(p *Scaner) (scanFun ScanFun) {
	scanFun = scanNew
	p.checkChar()
	p.newToken()
	switch p.charType {
	case TtPound:
		scanFun = scanComment
	case TtLeftQuotation, TtLeftApostrophe, TtApostrophe, TtQuotation:
		scanFun = scanString
	case TtLetter:
		scanFun = scanLetter
	case TtEqual, TtMoreThan, TtLessThan:
		scanFun = scanOperator
	case TtInt:
		scanFun = scanNumber
	case TtCR, TtLF:
		scanFun = scanNewLine
	default:
		scanFun = scanNew
	}

	return
}

func scanComment(p *Scaner) (scanFun ScanFun) {
	scanFun = scanComment
	p.checkChar()
	switch p.charType {
	case TtCR, TtLF, TtCRLF:
		scanFun = scanNew
	default:
		p.addTokenChar()
	}
	return
}

func scanString(p *Scaner) (scanFun ScanFun) {
	scanFun = scanString
	p.checkChar()
	switch p.nowToken.TokenType {
	case TtLeftQuotation:
		switch p.charType {
		case TtRightQuotation:
			p.addStringTokenEndChar()
			scanFun = scanNew
		default:
			p.addTokenChar()
		}
	case TtLeftApostrophe:
		switch p.charType {
		case TtRightApostrophe:
			p.addStringTokenEndChar()
			scanFun = scanNew
		default:
			p.addTokenChar()
		}
	case TtApostrophe:
		switch p.charType {
		case TtApostrophe:
			p.addTokenChar()
			scanFun = scanNew
		case TtCR, TtLF, TtCRLF:
			scanFun = scanNew
		default:
			p.addTokenChar()
		}
	case TtQuotation:
		switch p.charType {
		case TtQuotation:
			p.addTokenChar()
			scanFun = scanNew
		case TtCR, TtLF, TtCRLF:
			scanFun = scanNew
		default:
			p.addTokenChar()
		}
	}
	return
}

func scanLetter(p *Scaner) (scanFun ScanFun) {
	scanFun = scanLetter
	p.checkChar()
	switch p.charType {
	case TtLetter, TtInt:
		p.addTokenChar()
	default:
		scanFun = scanNew
	}
	return
}

func scanNumber(p *Scaner) (scanFun ScanFun) {
	scanFun = scanNumber
	p.checkChar()
	switch p.charType {
	case TtInt:
		p.addTokenChar()
	case TtPoint:
		if p.nowToken.TokenType == TtInt {
			p.addTokenChar()
			p.changeTokenType(TtFloat)
		} else {
			scanFun = scanNew
		}
	default:
		scanFun = scanNew
	}
	return
}

func scanOperator(p *Scaner) (scanFun ScanFun) {
	scanFun = scanNew
	p.checkChar()
	switch p.charType {
	case TtEqual:
		switch p.nowToken.TokenType {
		case TtEqual:
			p.addTokenChar()
			p.changeTokenType(TtEqualEqual)
		case TtMoreThan:
			p.addTokenChar()
			p.changeTokenType(TtMoreThanEqual)
		case TtLessThan:
			p.addTokenChar()
			p.changeTokenType(TtLessThanEqual)
		}
	case TtMoreThan:
		if p.nowToken.TokenType == TtLessThan {
			p.addTokenChar()
			p.changeTokenType(TtNotEqual)
		}
	case TtLessThan:
		if p.nowToken.TokenType == TtMoreThan {
			p.addTokenChar()
			p.changeTokenType(TtNotEqual)
		}
	}
	return
}

func scanNewLine(p *Scaner) (scanFun ScanFun) {
	scanFun = scanNewLine
	p.checkChar()
	switch p.charType {
	case TtLF:
		if p.nowToken.TokenType == TtCR {
			p.addTokenChar()
			p.changeTokenType(TtCRLF)
			scanFun = scanNew
		} else {
			scanFun = scanNew
		}
	default:
		scanFun = scanNew
	}
	p.lineNo += 1
	p.colNo = 1
	return
}
