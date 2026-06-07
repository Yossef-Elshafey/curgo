package object

import (
	"fmt"
	"net/http"
)

type Response struct {
	Res *http.Response
	Opts map[string]bool
	Cache struct{
		Body string
	}
}

type Expect struct {
	Value Object
	ToBe  ToBe
}

type ToBe struct {
	Value Object
}

func (r *Response) Type() ObjectType { return RESPONSE }
func (r *Response) Visit() string {
	return fmt.Sprintf("%+v", r.Res)
}
