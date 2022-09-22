package zhen_0_01

type VarGroupWithDef struct {
	NewVarGroup  *VarGroup
	DefVarGroup  *VarGroup
	DefVarCount  int
	DefValues    []ZValue
	DefValueSets []bool
}

func NewVarGroupWithDefault(defGroup *VarGroup) (varGroup *VarGroupWithDef) {
	varGroup = &VarGroupWithDef{}
	varGroup.NewVarGroup = NewVarGroup()
	varGroup.DefVarGroup = defGroup
	varGroup.DefVarCount = len(defGroup.AllVars)
	if varGroup.DefVarCount > 0 {
		varGroup.DefValues = make([]ZValue, varGroup.DefVarCount)
		varGroup.DefValueSets = make([]bool, varGroup.DefVarCount)
	}

	return
}

func (varGroup *VarGroupWithDef) FindByName(name string) (n int) {
	n = -1
	if varGroup.DefVarCount > 0 {
		n = varGroup.DefVarGroup.FindByName(name)
	}
	if n < 0 {
		n = varGroup.DefVarCount + varGroup.NewVarGroup.FindByName(name)
	}
	return
}

func (varGroup *VarGroupWithDef) Get(n int) (v ZValue) {
	if n >= 0 && n < varGroup.DefVarCount {
		if varGroup.DefValueSets[n] {
			v = varGroup.DefValues[n]
		} else {
			v = varGroup.DefVarGroup.Get(n)
		}
	} else {
		n -= varGroup.DefVarCount
		v = varGroup.NewVarGroup.Get(n)
	}

	return
}
func (varGroup *VarGroupWithDef) Set(n int, v ZValue) (ok bool) {
	if n >= 0 && n < varGroup.DefVarCount {
		varGroup.DefValues[n] = v
		varGroup.DefValueSets[n] = true
		ok = true
	} else {
		n -= varGroup.DefVarCount
		ok = varGroup.NewVarGroup.Set(n, v)
	}
	return
}

func (varGroup *VarGroupWithDef) GetByName(name string) (v ZValue) {
	n := varGroup.FindByName(name)
	if n >= 0 {
		v = varGroup.Get(n)
	}

	return
}
func (varGroup *VarGroupWithDef) SetByName(name string, v ZValue) {
	n := varGroup.FindByName(name)
	if n >= 0 {
		_ = varGroup.Set(n, v)
	} else {
		varGroup.NewVarGroup.SetByName(name, v)
	}

	return
}
