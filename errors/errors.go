package errors

import (
	"strings"
)

type cache struct {
	source string
}

var c cache = cache{}

func SetSourceName(s string) {
	c.source = s
}

func CaptureErrorLine(ln int) string {
	lines := strings.Split(c.source, "\n")
	return lines[ln]
}
