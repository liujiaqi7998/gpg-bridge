package app

import (
	"errors"
	"testing"
	"time"
)

func TestRunResultAllowsMultipleWaiters(t *testing.T) {
	r := newRunResult()
	expected := errors.New("boom")
	go func() {
		time.Sleep(10 * time.Millisecond)
		r.finish(expected)
	}()

	if err := r.wait(); !errors.Is(err, expected) {
		t.Fatalf("first wait got %v, want %v", err, expected)
	}
	if err := r.wait(); !errors.Is(err, expected) {
		t.Fatalf("second wait got %v, want %v", err, expected)
	}
}

func TestRunResultKeepsFirstResult(t *testing.T) {
	r := newRunResult()
	first := errors.New("first")
	second := errors.New("second")

	r.finish(first)
	r.finish(second)

	if err := r.wait(); !errors.Is(err, first) {
		t.Fatalf("wait got %v, want %v", err, first)
	}
}
