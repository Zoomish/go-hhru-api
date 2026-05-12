package main

import (
	"bytes"
	"os"
	"path/filepath"
)

func main() {
	path := filepath.Join("gen", "applicant", "client.gen.go")
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	repls := []struct{ old, new string }{
		{"\tv.AuthType = \"application\"\n", "\ts := \"application\"\n\tv.AuthType = &s\n"},
		{"\tv.AuthType = \"applicant\"\n", "\ts := \"applicant\"\n\tv.AuthType = &s\n"},
		{"\tv.AuthType = \"employer\"\n", "\ts := \"employer\"\n\tv.AuthType = &s\n"},
		{"\tv.AuthType = \"employer_integration\"\n", "\ts := \"employer_integration\"\n\tv.AuthType = &s\n"},
	}
	out := b
	for _, r := range repls {
		out = bytes.ReplaceAll(out, []byte(r.old), []byte(r.new))
	}
	if err := os.WriteFile(path, out, 0644); err != nil {
		panic(err)
	}
}
