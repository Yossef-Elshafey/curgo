package environment

import (
	"curgo/types/object"
	"testing"
)

func TestEnvSet(t *testing.T) {
	ff := &object.FetchFunction{}
	env := New()
	env.Set("id", ff)
}
