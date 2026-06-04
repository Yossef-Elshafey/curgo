package utils

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
	BODY
	NONE
)

type ReqBuilder struct {
	params              map[string]ParamType
	helpMessagePadding  int
	req                 *http.Request
}

func NewRequestBuilder() *ReqBuilder {
	b := &ReqBuilder{}
	b.params = make(map[string]ParamType)
	b.build()
	return b
}

func (b *ReqBuilder) build() {
	b.set("data",          BODY)
	b.set("method",        METHOD)
	b.set("host",          HOST)
	b.set("header",        HEADER)
}

func (b *ReqBuilder) set(from string, to ParamType) {
	if len(from) > b.helpMessagePadding {
		b.helpMessagePadding = len(from)
	}

	if _, ok := b.params[from]; ok {
		log.Fatalf("Curl Transpiler: cant set %s already exists",  from)
	}
	b.params[from] = to
}

func (b *ReqBuilder) Get(k string) ParamType {
	v, ok := b.params[k]
	if !ok {
		return NONE
	}
	return v
}

func (b *ReqBuilder) BuildRequest(k,v string) error {
	if b.req == nil {
		b.req, _ = http.NewRequest("GET", "", nil)
	}
	switch b.Get(k) {
		case HEADER: 
			header := strings.SplitN(v, ":", 2)
			if len(header) < 2 {
				return fmt.Errorf("missing ':' between header key, value: %s", v)
			}
			b.req.Header.Add(header[0], header[1])
		case HOST:
			u, err := url.Parse(v)
			if err != nil {
				return fmt.Errorf("cannot parse request url: %s", err)
			}
			if u.Scheme == "" || u.Host == "" {
				return fmt.Errorf("invalid URL: missing scheme or host: %s", v)
			}
			b.req.URL = u
		case METHOD: b.req.Method = v
		case BODY: b.req.Body = io.NopCloser(strings.NewReader(v))
		case NONE : return fmt.Errorf("cant transpile %s", k)
	}
	return nil
}

func (b *ReqBuilder) normalization() {
	// req.url errors is handled by url.Parse function
	if b.req.Method == "" {
		b.req.Method = "GET"
	}
}

func (b *ReqBuilder) Fire() ( *http.Response, error ) {
	b.normalization()
	resp, err := http.DefaultClient.Do(b.req)
	if err != nil { return nil, err }
	return resp, nil
}
