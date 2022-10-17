package zhen_0_02

import (
	"fmt"
	"strconv"
)

type State struct {
	debug         bool
	codeFile      *CodeFile
	mainCodeBlock *CodeBlock
	nowCodeBlock  *CodeBlock

	//codeBlockArea []*CodeBlock
	globalValues *ValueGroup
	localValues  *ValueGroup
	areaValues   *ValueGroup
	tempValues   *TempValueGroup
}

func NewState(codeFile *CodeFile) (s State) {
	s.debug = false
	s.codeFile = codeFile
	s.mainCodeBlock = &s.codeFile.AllCodeBlock[0]

	return

}

func (state *State) Run() (err error) {
	err = state.runCodeBlock(state.mainCodeBlock)
	if err != nil {
		return
	}
	return
}

func (state *State) getCodeBlock(n int) (codeBlock *CodeBlock) {
	codeBlock = &state.codeFile.AllCodeBlock[n]

	return
}
func (codePre *State) getCodeBlockArea(codeBlock *CodeBlock) (codeBlockArea *CodeBlock) {
	if codeBlock.ParNo >= 0 {
		parCodeBlock := codePre.getCodeBlock(codeBlock.ParNo)
		switch parCodeBlock.BlockType {
		case CbtFile, CbtFun:
			codeBlockArea = parCodeBlock
		case CbtLine, CbtChildLine, CbtColon:
			codeBlockArea = parCodeBlock
		default:
			codeBlockArea = codePre.getCodeBlockArea(parCodeBlock)
		}
	}
	return
}

func (codePre *State) checkCodeBlockAreaIn(codeBlock *CodeBlock) (isArea bool) {

	switch codeBlock.BlockType {
	case CbtFile:
		codePre.globalValues = NewValueGroup(codePre.globalValues, codeBlock)
		codePre.localValues = NewValueGroup(codePre.localValues, codeBlock)
		codePre.areaValues = NewValueGroup(codePre.localValues, codeBlock)
		codePre.tempValues = NewTempValueGroup(codePre.tempValues, codeBlock)
		codePre.globalValues.Values["显示"] = StatePrint
	case CbtFun:
		codePre.globalValues = NewValueGroup(codePre.globalValues, codeBlock)
		codePre.localValues = NewValueGroup(codePre.localValues, codeBlock)
		codePre.areaValues = NewValueGroup(codePre.localValues, codeBlock)

		codePre.tempValues = NewTempValueGroup(codePre.tempValues, codeBlock)

	case CbtLine, CbtColon:
		codePre.areaValues = NewValueGroup(codePre.localValues, codeBlock)
		codePre.tempValues = NewTempValueGroup(codePre.tempValues, codeBlock)
	case CbtChildLine:
		codePre.tempValues = NewTempValueGroup(codePre.tempValues, codeBlock)
	}

	return
}
func (codePre *State) checkCodeBlockAreaOut(codeBlock *CodeBlock) (isArea bool) {

	switch codeBlock.BlockType {
	case CbtFile, CbtFun:
		if codePre.globalValues != nil {
			codePre.globalValues = codePre.globalValues.ParValueGroup
		}
		if codePre.localValues != nil {
			codePre.localValues = codePre.localValues.ParValueGroup
		}
		if codePre.tempValues != nil {
			codePre.tempValues = codePre.tempValues.ParValueGroup
		}

	case CbtLine, CbtChildLine, CbtColon:
		if codePre.tempValues != nil {
			codePre.tempValues = codePre.tempValues.ParValueGroup
		}
	}
	return
}

func (state *State) runCodeBlock(codeBlock *CodeBlock) (err error) {
	state.checkCodeBlockAreaIn(codeBlock)
	switch codeBlock.BlockType {
	case CbtFun:

	default:

		childNo := codeBlock.FirstChildNo
		for childNo >= 0 {
			childCodeBlock := state.getCodeBlock(childNo)
			err = state.runCodeBlock(childCodeBlock)
			if err != nil {
				return
			}

			childNo = childCodeBlock.NextNo
		}

	}

	for _, s := range codeBlock.Steps {
		err = state.runCodeStep(codeBlock, &s)
		if err != nil {
			return
		}
	}
	state.checkCodeBlockAreaOut(codeBlock)
	return
}
func (state *State) CodeStringToValue(s string) (v Value) {

	if s != "" {
		r := []rune(s)
		switch r[0] {
		case '\'', '‘', '’', '"', '“', '”':
			r = r[1:]
			if len(r) > 0 {
				switch r[len(r)-1] {
				case '\'', '‘', '’', '"', '“', '”':
					r = r[:len(r)-1]
				}
			}

			v = string(r)
			return
		}
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			v = f
			return
		}
		v = s
		return

	}
	return
}

func (state *State) Arithmetic(codeStep *CodeStep) (err error) {
	var aa Value
	var bb Value

	if codeStep.ValueString1 != "" {
		aa = state.CodeStringToValue(codeStep.ValueString1)
	} else if codeStep.VarName1 != "" {
		aa = state.localValues.Values[codeStep.VarName1]
	} else if codeStep.TempVarNo1 != 0 {
		aa = state.tempValues.Values[codeStep.TempVarNo1]
	}

	if codeStep.ValueString2 != "" {
		bb = state.CodeStringToValue(codeStep.ValueString2)
	} else if codeStep.VarName2 != "" {
		bb = state.localValues.Values[codeStep.VarName2]
	} else if codeStep.TempVarNo2 != 0 {
		bb = state.tempValues.Values[codeStep.TempVarNo2]
	}

	switch codeStep.CodeStepType {
	case CstAs:
		if codeStep.VarName1 != "" {
			state.localValues.Values[codeStep.VarName1] = bb
		} else if codeStep.TempVarNo1 != 0 {
			state.tempValues.Values[codeStep.TempVarNo1] = bb
		}
		state.tempValues.Values[codeStep.ReturnVarNo] = bb
	case CstAdd, CstSub, CstMul, CstDiv:
		var v Value
		v, err = Arithmetic(aa, bb, codeStep.CodeStepType)
		if err != nil {
			return
		}
		state.tempValues.Values[codeStep.ReturnVarNo] = v
	}
	return
}
func (state *State) runCodeStep(codBlock *CodeBlock, codeStep *CodeStep) (err error) {
	switch codeStep.CodeStepType {
	case CstDefineText, CstDefineVar, CstDefineLocalVar, CstDefineGlobalVar, CstDefineConstant:
	case CstDefineFun, CstDefineFunPara, CstDefineFunReturn:

	case CstAs, CstAdd, CstSub, CstMul, CstDiv:
		err = state.Arithmetic(codeStep)
		if err != nil {
			return
		}
	case CstPoint:

	case CstCall:
	case CstTryCall:
		var fun Value
		fun, err = state.getCodeStepValue1(codeStep)
		if err == nil {
			switch f := fun.(type) {
			case func(*State, *CodeStep) error:
				err = f(state, codeStep)
				if err != nil {
					return
				}
			}
		}
	}
	fmt.Println(*codeStep)
	return
}

func (state *State) getCodeStepValue1(codeStep *CodeStep) (aa Value, err error) {
	if codeStep.ValueString1 != "" {
		aa = state.CodeStringToValue(codeStep.ValueString1)
	} else if codeStep.VarName1 != "" {
		var ok bool
		aa, ok = state.GetVarValue(codeStep.VarName1)
		if !ok {
			aa = nil
		}
	} else if codeStep.TempVarNo1 != 0 {
		aa = state.tempValues.Values[codeStep.TempVarNo1]
	}

	return
}

func (state *State) getCodeStepValue2(codeStep *CodeStep) (aa Value, err error) {
	if codeStep.ValueString2 != "" {
		aa = state.CodeStringToValue(codeStep.ValueString2)
	} else if codeStep.VarName2 != "" {
		var ok bool
		aa, ok = state.GetVarValue(codeStep.VarName2)
		if !ok {
			aa = nil
		}
	} else if codeStep.TempVarNo2 != 0 {
		aa = state.tempValues.Values[codeStep.TempVarNo2]
	}

	return
}
func (state *State) GetVarValue(varName string) (v Value, ok bool) {

	areaValues := state.areaValues
	for areaValues != nil {
		switch areaValues.CodeBlock.BlockType {
		case CbtFun, CbtFile:
			v, ok = state.areaValues.Values[varName]
			if ok {
				return
			}
		}
		areaValues = areaValues.ParValueGroup
	}

	localValues := state.localValues
	for localValues != nil {
		switch localValues.CodeBlock.BlockType {
		case CbtFun, CbtFile:
			v, ok = state.localValues.Values[varName]
			if ok {
				return
			}
		}
		localValues = localValues.ParValueGroup
	}

	globalValues := state.globalValues
	for globalValues != nil {
		switch globalValues.CodeBlock.BlockType {
		case CbtFun, CbtFile:
			v, ok = state.globalValues.Values[varName]
			if ok {
				return
			}
		}
		globalValues = globalValues.ParValueGroup
	}

	return
}
func StatePrint(state *State, codeStep *CodeStep) (err error) {
	var aa Value
	aa, err = state.getCodeStepValue2(codeStep)
	if err != nil {
		return nil
	}
	PrintValue(aa)
	return
}
