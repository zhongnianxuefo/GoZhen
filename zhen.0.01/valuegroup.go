package zhen_0_01

type VarGroup struct {
	AllVars  []ZValue
	VarNames map[string]int
}

func NewVarGroup() (valueGroup *VarGroup) {
	valueGroup = &VarGroup{}
	valueGroup.VarNames = make(map[string]int)
	return
}
func (vars *VarGroup) FindByName(name string) (n int) {
	x, ok := vars.VarNames[name]
	if ok {
		n = x
	} else {
		n = -1
	}
	return
}

func (vars *VarGroup) Get(n int) (v ZValue) {
	if n >= 0 && n < len(vars.AllVars) {
		v = vars.AllVars[n]
	}
	return
}
func (vars *VarGroup) Set(n int, v ZValue) (ok bool) {
	if n >= 0 && n < len(vars.AllVars) {
		vars.AllVars[n] = v
		ok = true
	}
	return
}

func (vars *VarGroup) GetByName(name string) (v ZValue) {
	x, ok := vars.VarNames[name]
	if ok {
		v = vars.Get(x)
	}
	return
}
func (vars *VarGroup) SetByName(name string, v ZValue) {
	x, ok := vars.VarNames[name]
	if ok {
		_ = vars.Set(x, v)
	} else {
		vars.AllVars = append(vars.AllVars, v)
		vars.VarNames[name] = len(vars.AllVars) - 1
	}
	return
}
