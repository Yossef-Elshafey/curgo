package transpiler

import "log"

func (ct *CurlTranspiler) build() {
	ct.set("header", "-H")
	ct.set("data", "-d")
	ct.set("method", "-X")
	ct.set("host", "")
}

type CurlTranspiler struct {
	params map[string]string
	command string
}

func New() *CurlTranspiler {
	ct := &CurlTranspiler{params: make(map[string]string)}
	ct.build()
	return ct
}

func (ct *CurlTranspiler) set(from, to string) {
	if _, ok := ct.params[from]; ok {
		log.Fatalf("Curl Transpiler: cant set %s, %s already exists", from, from)
	}
	ct.params[from] = to
}

func (ct *CurlTranspiler) Get(key string) (string, bool) {
	v, ok := ct.params[key]
	return v, ok
}
