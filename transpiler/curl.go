package transpiler

import (
	"fmt"
	"log"
)

type valueWrapper struct {
	argument string
	valueHandler  func(v string) string
}

type CurlTranspiler struct {
	params map[string]*valueWrapper
}

func doubleQuoteWrapper(v string) string {
	return fmt.Sprintf("\"%s\"", v)
}

func singleQuoteWrapper(v string) string {
	return fmt.Sprintf("'%s'", v)
}

func nullWrapper(v string) string { return v }

func (ct *CurlTranspiler) build() {
	ct.set("header",  &valueWrapper{argument:  "-H",  valueHandler:  doubleQuoteWrapper})
	ct.set("data",    &valueWrapper{argument:  "-d",  valueHandler:  singleQuoteWrapper})
	ct.set("method",  &valueWrapper{argument:  "-X",  valueHandler:  nullWrapper})
	ct.set("host",    &valueWrapper{argument:  "",    valueHandler:  nullWrapper})
}

func New() *CurlTranspiler {
	ct := &CurlTranspiler{params: make(map[string]*valueWrapper)}
	ct.build()
	return ct
}

func (ct *CurlTranspiler) set(from string, to *valueWrapper) {
	if _, ok := ct.params[from]; ok {
		log.Fatalf("Curl Transpiler: cant set %s already exists",  from)
	}
	ct.params[from] = to
}

func (ct *CurlTranspiler) Get(key, value string) (string, string) {
	wrapper, ok := ct.params[key]
	if !ok {
		log.Fatalf("Transpiler: can't find value for key %s", key)
	}
	value = wrapper.valueHandler(value)
	return wrapper.argument, value
}

