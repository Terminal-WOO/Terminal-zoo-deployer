package client

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestErrInClusterConfig_Error(t *testing.T) {
	orig := errors.New("underlying failure")
	e := NewErrInClusterConfig(orig)

	got := e.Error()
	if got == "" {
		t.Fatal("Error() returned empty string")
	}

	// It should contain the original error text somewhere.
	if !strings.Contains(got, orig.Error()) {
		t.Errorf("Error() = %q, want it to contain %q", got, orig.Error())
	}
}

func TestErrInClusterConfig_Unwrap(t *testing.T) {
	orig := fmt.Errorf("wrapped: %w", errors.New("root cause"))
	e := NewErrInClusterConfig(orig)

	if unwrapped := e.Unwrap(); unwrapped == nil {
		t.Fatal("Unwrap() returned nil, expected non-nil error")
	}

	if !errors.Is(e, orig) {
		t.Errorf("e = %v, doesnt contain original %v", e, orig)
	}
}

func TestErrOutClusterConfig_Error(t *testing.T) {
	orig := errors.New("underlying failure")
	e := NewErrOutClusterConfig(orig)

	got := e.Error()
	if got == "" {
		t.Fatal("Error() returned empty string")
	}

	// It should contain the original error text somewhere.
	if !strings.Contains(got, orig.Error()) {
		t.Errorf("Error() = %q, want it to contain %q", got, orig.Error())
	}
}

func TestErrOutClusterConfig_Unwrap(t *testing.T) {
	orig := fmt.Errorf("wrapped: %w", errors.New("root cause"))
	e := NewErrOutClusterConfig(orig)

	if unwrapped := e.Unwrap(); unwrapped == nil {
		t.Fatal("Unwrap() returned nil, expected non-nil error")
	}

	if !errors.Is(e, orig) {
		t.Errorf("e = %v, doesnt contain original %v", e, orig)
	}
}

func TestErrMissingEnvVars_Error(t *testing.T) {
	e := new(ErrMissingEnvVars)

	got := e.Error()
	if got == "" {
		t.Fatal("Error() returned empty string")
	}
}

func TestErrClientCreation_Error(t *testing.T) {
	orig := errors.New("underlying failure")
	e := NewErrClientCreation(orig)

	got := e.Error()
	if got == "" {
		t.Fatal("Error() returned empty string")
	}

	// It should contain the original error text somewhere.
	if !strings.Contains(got, orig.Error()) {
		t.Errorf("Error() = %q, want it to contain %q", got, orig.Error())
	}
}

func TestErrClientCreation_Unwrap(t *testing.T) {
	orig := fmt.Errorf("wrapped: %w", errors.New("root cause"))
	e := NewErrClientCreation(orig)

	if unwrapped := e.Unwrap(); unwrapped == nil {
		t.Fatal("Unwrap() returned nil, expected non-nil error")
	}

	if !errors.Is(e, orig) {
		t.Errorf("e = %v, doesnt contain original %v", e, orig)
	}
}
