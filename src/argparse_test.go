package main

import (
	"strings"
	"testing"
)

// TestSanitizeFlagPresence locks in that sanitize mode is selected by the
// flag being present, not by its value. An empty sanitize value must still
// count as "sanitize was requested" and must never fall through to renaming.
func TestSanitizeFlagPresence(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantSet   bool
		wantValue string
	}{
		{"flag absent", nil, false, ""},
		{"short flag with value", []string{"-s", "Hello World"}, true, "Hello World"},
		{"short flag with empty value", []string{"-s", ""}, true, ""},
		{"long flag with empty value", []string{"--sanitize", ""}, true, ""},
		{"long flag with equals and empty value", []string{"--sanitize="}, true, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var cli tCLI
			parser, err := newArgParser(&cli)
			if err != nil {
				t.Fatalf("cannot create argument parser: %v", err)
			}
			if _, err := parser.Parse(c.args); err != nil {
				t.Fatalf("cannot parse %q: %v", c.args, err)
			}
			if got := cli.Sanitize != nil; got != c.wantSet {
				t.Fatalf("sanitize present = %v, want %v (args %q)", got, c.wantSet, c.args)
			}
			if c.wantSet && *cli.Sanitize != c.wantValue {
				t.Fatalf("sanitize = %q, want %q", *cli.Sanitize, c.wantValue)
			}
		})
	}
}

func TestSanitizeWithoutValue(t *testing.T) {
	// A missing value is a parse error. parseArgs reports it on stderr;
	// kong must not be allowed to print the help to stdout.
	for _, args := range [][]string{{"-s"}, {"--sanitize"}} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			var cli tCLI
			parser, err := newArgParser(&cli)
			if err != nil {
				t.Fatalf("cannot create argument parser: %v", err)
			}
			if _, err := parser.Parse(args); err == nil {
				t.Fatalf("expected a parse error for %q", args)
			}
			if cli.Sanitize != nil {
				t.Fatalf("sanitize must stay unset after a parse error")
			}
		})
	}
}

func TestSanitizeEmptyString(t *testing.T) {
	// This is why main prints nothing for `-s ""`: the sanitized value
	// itself is empty, so there is not even a newline to emit.
	rs := initReplacerSchemes()
	if got := normalize(tPath{Name: ""}, rs).Name; got != "" {
		t.Fatalf("normalize(\"\") = %q, want empty string", got)
	}
}
