package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/harnyk/gh-hide/internal/hide"
)

func TestInvalidArgsAreUsageErrors(t *testing.T) {
	cases := map[string][]string{
		"missing filter":      {},
		"empty filter":        {""},
		"all with filter":     {"--all", ".spec."},
		"count with filter":   {"--count", ".spec."},
		"empty among filters": {"a", ""},
		"bad wildcard":        {"-w", "[a"},
		"wildcard with all":   {"-w", "--all"},
		"wildcard with count": {"-w", "--count"},
		"unknown flag":        {"--nope", "a"},
		"zero concurrency":    {"-j", "0", "a"},
		"invalid concurrency": {"--concurrency", "x", "a"},
	}
	for name, args := range cases {
		t.Run(name, func(t *testing.T) {
			called := false
			cmd := newRootCmd(func() (hide.Client, error) {
				called = true
				return nil, errors.New("must not be called")
			})
			cmd.SetArgs(args)
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})
			err := cmd.Execute()
			if _, ok := errors.AsType[usageError](err); !ok {
				t.Fatalf("expected usage error, got %v", err)
			}
			if called {
				t.Fatal("client was created for invalid arguments")
			}
		})
	}
}
