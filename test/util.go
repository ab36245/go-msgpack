package test

import "testing"

func report(t *testing.T, a, e any) {
	t.Fatalf("\nexpected: %v\nactual:   %v\n", e, a)
}
