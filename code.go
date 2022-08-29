package main

type ZhenCodeStepType byte

const (
	ZCS_Var ZhenCodeStepType = iota
	ZCS_Add
	ZCS_Sub
	ZCS_Mul
	ZCS_Div
	ZCS_Eq
	LUA_LT
	LUA_LE
	ZCS_Assign
	ZCS_If
	ZCS_For
	ZCS_While
	ZCS_Break
	ZCS_Return
	ZCS_Call

	ZCS_PrintVar
)

type ZhenCodeStep struct {
	codeStepType ZhenCodeStepType
	valueName    string
	value        ZhenValue
}

type ZhenCode struct {
	needRun   bool
	codeSteps []ZhenCodeStep
}
