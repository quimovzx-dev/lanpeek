package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestJSONCleanStdout(t *testing.T) {
	for _, args := range [][]string{
		{"info", "--json"},
		{"ports", "127.0.0.1", "--ports", "1", "--json"},
		{"discover", "--cidr", "127.0.0.1/32", "--no-dns", "--json"},
	} {
		root := NewRoot()
		var out, stderr bytes.Buffer
		root.SetOut(&out)
		root.SetErr(&stderr)
		root.SetArgs(args)
		if err := root.Execute(); err != nil {
			t.Fatal(err)
		}
		if !json.Valid(out.Bytes()) {
			t.Fatalf("not JSON: %s", out.String())
		}
		if args[0] != "info" && stderr.Len() == 0 {
			t.Fatal("missing progress on stderr")
		}
	}
}

func TestInvalidArguments(t *testing.T) {
	for _, args := range [][]string{
		{"discover"}, {"discover", "--cidr", "10.0.0.0/8"},
		{"ports"}, {"ports", "localhost", "--ports", "0"},
		{"ports", "https://example.com"}, {"ports", "localhost", "--concurrency", "0"},
		{"info", "--timeout", "0s"}, {"info", "--timeout", "31s"}, {"info", "unexpected"},
	} {
		root := NewRoot()
		root.SetArgs(args)
		var out bytes.Buffer
		root.SetOut(&out)
		root.SetErr(&out)
		if err := root.Execute(); err == nil {
			t.Errorf("accepted %v", args)
		}
	}
}
