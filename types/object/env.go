package object

type Env struct {
	vars   map[string]Object
	outer  *Env
}

func NewEnvironment() *Env {
	vars := make(map[string]Object)
	return &Env{vars:vars}
}

func NewOuterEnv(outEnv *Env) *Env {
	env := NewEnvironment()
	env.outer = outEnv
	return env
}

func (e *Env) Set(k string, v Object) {
	e.vars[k] = v
}

func (e *Env) Get(k string) ( Object, bool ) {
	v, ok := e.vars[k] 
	if !ok && e.outer != nil {
		v, ok = e.outer.Get(k)
	}
	return v,ok
}
