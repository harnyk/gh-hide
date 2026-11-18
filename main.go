// Command gh-hide is a GitHub CLI extension that marks pull request files as viewed.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/harnyk/gh-hide/internal/github"
	"github.com/harnyk/gh-hide/internal/hide"
)

// Set by goreleaser via -ldflags.
var version = "dev"

// usageError signals invalid command-line usage (exit code 2).
type usageError struct{ error }

func newRootCmd(newClient func() (hide.Client, error)) *cobra.Command {
	var (
		opts      hide.Options
		countOnly bool
	)

	cmd := &cobra.Command{
		Use: "hide [--dry-run] [--unhide] [--all] [FILTER]",
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "gh hide",
		},
		Short: "Mark pull request files as viewed",
		Long: `Mark files whose pull request path contains FILTER as viewed.

FILTER is a literal substring of the file path, not a regex.
With --all, process every file in the pull request (FILTER must be omitted).
With --unhide, unmark matching files as viewed (show them again).
With --count, print how many files are viewed and not viewed.`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			switch {
			case countOnly, opts.All:
				if len(args) != 0 {
					return usageError{errors.New("FILTER must be omitted with --all or --count")}
				}
			case len(args) != 1:
				return usageError{errors.New("exactly one FILTER is required")}
			case args[0] == "":
				return usageError{errors.New("FILTER must not be empty")}
			}
			if opts.Concurrency < 1 {
				return usageError{errors.New("--concurrency must be at least 1")}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient()
			if err != nil {
				return err
			}
			r := &hide.Runner{Client: client, Out: cmd.OutOrStdout(), Err: cmd.ErrOrStderr()}
			if countOnly {
				return r.Count(cmd.Context())
			}
			if len(args) == 1 {
				opts.Filter = args[0]
			}
			return r.Run(cmd.Context(), opts)
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&opts.DryRun, "dry-run", "n", false, "show what would be changed without calling the API")
	f.BoolVarP(&opts.Unhide, "unhide", "u", false, "unmark matching files as viewed")
	f.BoolVarP(&opts.All, "all", "a", false, "process every file in the pull request")
	f.BoolVarP(&countOnly, "count", "c", false, "print viewed and not viewed file counts")
	f.IntVarP(&opts.Concurrency, "concurrency", "j", hide.DefaultConcurrency, "number of files to update in parallel")

	cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return usageError{err} })
	return cmd
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cmd := newRootCmd(func() (hide.Client, error) { return github.New() })
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		var ue usageError
		if errors.As(err, &ue) {
			fmt.Fprint(os.Stderr, "\n"+cmd.UsageString())
			os.Exit(2)
		}
		os.Exit(1)
	}
}
