package utils

import (
	"strings"
)

var source string

func SetSource(s string) {
	source = s
}

func ReadSourceAsLines(line int) string {
	return string(strings.Split(string(source), "\n")[line-1])
}
