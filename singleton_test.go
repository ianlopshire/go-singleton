package singleton_test

import (
	"testing"

	"github.com/ianlopshire/go-singleton"
)

func TestNew(t *testing.T) {
	s := singleton.New(func() (*string, error) {
		foo := "foo"
		return &foo, nil
	})

	v1, err := s.Instance()
	if err != nil {
		t.Fatalf("Instance() unexpected error: %v", err)
	}
	v2, err := s.Instance()
	if err != nil {
		t.Fatalf("Instance() unexpected error: %v", err)
	}

	if v1 != v2 {
		t.Fatalf("Instance() expected %v, got %v", v1, v2)
	}
}

func TestLazy(t *testing.T) {
	var callCount int

	s := singleton.Lazy(func() (*string, error) {
		callCount++
		foo := "foo"
		return &foo, nil
	})

	if callCount != 0 {
		t.Fatalf("Lazy() expected callCount to be 0, got %d", callCount)
	}

	v1, err := s.Instance()
	if err != nil {
		t.Fatalf("Instance() unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Fatalf("Lazy() expected callCount to be 1, got %d", callCount)
	}

	v2, err := s.Instance()
	if err != nil {
		t.Fatalf("Instance() unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Fatalf("Lazy() expected callCount to be 1, got %d", callCount)
	}

	if v1 != v2 {
		t.Fatalf("Instance() expected %v, got %v", v1, v2)
	}
}
