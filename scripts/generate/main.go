package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const oapiCodegenMod = "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.7.0"

func main() {
	root, err := filepath.Abs(".")
	if err != nil {
		fatal(err)
	}
	if len(os.Args) > 1 {
		root, err = filepath.Abs(os.Args[1])
		if err != nil {
			fatal(err)
		}
	}
	run := func(name string, args ...string) {
		cmd := exec.Command(name, args...)
		cmd.Dir = root
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			fatal(fmt.Errorf("%s %v: %w", name, args, err))
		}
	}
	run("go", "run", "./scripts/split-openapi", root)
	services := []struct{ pkg, slug string }{
		{"employer", "employer"},
		{"applicant", "applicant"},
		{"public", "public"},
		{"app", "app"},
	}
	for i, s := range services {
		out := filepath.Join("gen", s.pkg, "client.gen.go")
		spec := filepath.Join("api", fmt.Sprintf("openapi.%s.yaml", s.slug))
		run("go", "run", oapiCodegenMod,
			"-package", s.pkg,
			"-generate", "types,client",
			"-o", out,
			spec,
		)
		if i == 1 {
			applicantPath := filepath.Join(root, "gen", "applicant", "client.gen.go")
			patchApplicantGenFile(applicantPath)
		}
	}
	for _, s := range services {
		p := filepath.Join(root, "api", fmt.Sprintf("openapi.%s.yaml", s.slug))
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			fatal(fmt.Errorf("remove %s: %w", p, err))
		}
	}
}

func patchApplicantGenFile(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		fatal(err)
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
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
