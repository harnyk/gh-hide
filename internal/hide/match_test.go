package hide

import "testing"

func TestMatcher(t *testing.T) {
	cases := []struct {
		name     string
		filters  []string
		wildcard bool
		path     string
		want     bool
	}{
		{"substring", []string{".spec."}, false, "src/a.spec.ts", true},
		{"substring miss", []string{".spec."}, false, "src/xspecy.ts", false},
		{"literal star", []string{"*.ts"}, false, "src/a.ts", false},
		{"any of several", []string{"nope", "lock"}, false, "pnpm-lock.yaml", true},
		{"none of several", []string{"nope", "nah"}, false, "pnpm-lock.yaml", false},

		{"basename glob in subdir", []string{"*.spec.ts"}, true, "src/deep/a.spec.ts", true},
		{"basename glob miss", []string{"*.spec.ts"}, true, "src/a.ts", false},
		{"path glob single level", []string{"src/*.ts"}, true, "src/a.ts", true},
		{"path glob no cross dir", []string{"src/*.ts"}, true, "src/deep/a.ts", false},
		{"doublestar", []string{"src/**/*.ts"}, true, "src/deep/er/a.ts", true},
		{"doublestar dir", []string{"generated/**"}, true, "generated/x/y.go", true},
		{"anchored to full path", []string{"generated/**"}, true, "pkg/generated/y.go", false},
		{"question mark", []string{"?.go"}, true, "x/a.go", true},
		{"braces", []string{"*.{yaml,yml}"}, true, "ci/a.yml", true},
		{"class", []string{"[ab].go"}, true, "b.go", true},
		{"wildcard is not substring", []string{"spec"}, true, "a.spec.ts", false},
		{"several globs", []string{"*.md", "*.lock"}, true, "Cargo.lock", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewMatcher(tc.filters, tc.wildcard)
			if err != nil {
				t.Fatal(err)
			}
			if got := m(tc.path); got != tc.want {
				t.Errorf("match(%q) with %q = %v, want %v", tc.path, tc.filters, got, tc.want)
			}
		})
	}
}

func TestMatcherErrors(t *testing.T) {
	cases := []struct {
		name     string
		filters  []string
		wildcard bool
	}{
		{"no filters", nil, false},
		{"empty filter", []string{"a", ""}, false},
		{"bad pattern", []string{"[a"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := NewMatcher(tc.filters, tc.wildcard); err == nil {
				t.Fatal("expected error")
			}
		})
	}
	if _, err := NewMatcher([]string{"[a"}, false); err != nil {
		t.Errorf("bracket must be literal without --wildcard: %v", err)
	}
}
