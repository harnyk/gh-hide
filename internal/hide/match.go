package hide

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// Matcher reports whether a pull request file path is selected by filters.
type Matcher func(filePath string) bool

// NewMatcher builds a Matcher that selects a path when it matches any filter.
//
// Without wildcard, a filter matches when the path contains it literally.
// With wildcard, a filter is a glob pattern (*, ?, [...], {a,b}, **) matched
// against the whole path; a pattern without "/" is also matched against the
// file's base name, so "*.spec.ts" selects such files in any directory.
func NewMatcher(filters []string, wildcard bool) (Matcher, error) {
	if len(filters) == 0 {
		return nil, errors.New("at least one filter is required")
	}
	for _, f := range filters {
		if f == "" {
			return nil, errors.New("filter must not be empty")
		}
		if wildcard && !doublestar.ValidatePattern(f) {
			return nil, fmt.Errorf("invalid wildcard pattern %q", f)
		}
	}
	patterns := append([]string(nil), filters...)

	if !wildcard {
		return func(p string) bool {
			for _, f := range patterns {
				if strings.Contains(p, f) {
					return true
				}
			}
			return false
		}, nil
	}

	return func(p string) bool {
		base := path.Base(p)
		for _, f := range patterns {
			if doublestar.MatchUnvalidated(f, p) {
				return true
			}
			if !strings.Contains(f, "/") && doublestar.MatchUnvalidated(f, base) {
				return true
			}
		}
		return false
	}, nil
}
