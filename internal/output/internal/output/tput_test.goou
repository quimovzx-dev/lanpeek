package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestJSON(t *testing.T) {
	var b bytes.Buffer
	if err := JSON(&b, []string{"open", "refused"}); err != nil {
		t.Fatal(err)
	}
	if !json.Valid(b.Bytes()) {
		t.Fatal(b.String())
	}
}

func TestTable(t *testing.T) {
	var b bytes.Buffer
	if err := Table(&b, []string{"PORT", "STATUS"}, [][]string{{"22", "open"}, {"80", "refused"}}); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"PORT", "STATUS", "22", "open", "80", "refused"} {
		if !strings.Contains(b.String(), s) {
			t.Errorf("missing %s in %s", s, b.String())
		}
	}
}
