package lexer_test

import (
	"regexp"
	"testing"
)

func TestAssignmentWithAndWithoutValue(t *testing.T) {
	regex := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_].*)(=(\s.*)?"[^"]*")?;`)
	content := `domain = "localhost:3000";`
	if !regex.MatchString(content) {
		t.Log("should match but fail")
	}
	content = `id;`
	if !regex.MatchString(content) {
		t.Log("should match but fail")
	}
	content = `transer`
	if regex.MatchString(content) {
		t.Logf("Expect %s to fail", content)
	}
}
