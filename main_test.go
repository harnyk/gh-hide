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
		"two filters":         {"a", "b"},
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
			var ue usageError
			if !errors.As(err, &ue) {
				t.Fatalf("expected usage error, got %v", err)
			}
			if called {
				t.Fatal("client was created for invalid arguments")
			}
		})
	}
}
