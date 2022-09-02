package main

type ZhenCodeType int

const (
	ZCT_Block ZhenCodeType = iota
	ZCT_Line
	ZCT_Words
)

type ZhenCodeLine struct {
	Words []TxtCodeWord
}

type ZhenCodeBlock struct {
	Lines  []ZhenCodeLine
	Blocks []ZhenCodeBlock
}
