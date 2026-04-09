package environment

import "curgo/types/object"

type Env struct {
	vars map[string]object.Object
}

func New() *Env {
	vars := make(map[string]object.Object)
	return &Env{vars:vars}
}

func (e *Env) Set(k string, v object.Object) {
	e.vars[k] = v
}

func (e *Env) Get(k string) ( object.Object, bool ) {
	v, ok := e.vars[k] 
	return v,ok
}
