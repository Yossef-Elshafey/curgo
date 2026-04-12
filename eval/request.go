package eval

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ParamType int

const (
	HEADER ParamType = iota
	HOST
	METHOD
	NONE
	BODY
)

type CurgoRequest struct {
	params              map[string]ParamType
	helpMessagePadding  int
	req                 *http.Request
}

func (cr *CurgoRequest) build() {
	cr.set("data",          BODY)
	cr.set("method",        METHOD)
	cr.set("host",          HOST)
	cr.set("header",        HEADER)
}

func New() *CurgoRequest {
	cr := &CurgoRequest{params: make(map[string]ParamType)}
	cr.build()
	return cr
}

func (cr *CurgoRequest) set(from string, to ParamType) {
	if len(from) > cr.helpMessagePadding {
		cr.helpMessagePadding = len(from)
	}

	if _, ok := cr.params[from]; ok {
		log.Fatalf("Curl Transpiler: cant set %s already exists",  from)
	}
	cr.params[from] = to
}

func (cr *CurgoRequest) Get(k string) ParamType {
	v, ok := cr.params[k]
	if !ok {
		return NONE
	}
	return v
}

func (cr *CurgoRequest) buildRequest(k,v string) error {
	if cr.req == nil {
		cr.req, _ = http.NewRequest("GET", "http://placeholder", nil)
	}
	switch cr.Get(k) {
		case HEADER: 
			header := strings.SplitN(v, ":", 2)
			cr.req.Header.Add(header[0], header[1])
		case HOST:
			url, err := url.Parse(v)
			if err != nil { return fmt.Errorf("cannot parse request url")}
			cr.req.URL = url
		case METHOD: cr.req.Method = v
		case BODY: cr.req.Body = io.NopCloser(strings.NewReader(v))
		case NONE: return fmt.Errorf("cant transpile %s", k)
	}
	return nil
}

func (cr *CurgoRequest) fire() ( *http.Response, error ) {
	resp, err := http.DefaultClient.Do(cr.req)
	if err != nil { return nil, err }
	return resp, nil
}
