package singleton

import (
	"context"

	"github.com/ianlopshire/go-async"
)

type Singleton[T any] interface {
	Instance() (T, error)
	InstanceCtx(ctx context.Context) (T, error)
}

type singleton[T any] struct {
	instance T
	err      error
}

func New[T any](new func() (T, error)) Singleton[T] {
	i, err := new()
	return singleton[T]{
		instance: i,
		err:      err,
	}
}

func (s singleton[T]) Instance() (T, error) {
	return s.InstanceCtx(context.Background())
}

func (s singleton[T]) InstanceCtx(ctx context.Context) (T, error) {
	return s.instance, s.err
}

// Lazy is a thread-safe lazy-initialized singleton value.
type lazy[T any] struct {
	new      func() (T, error)
	instance T
	err      error

	latch async.Latch
}

// Lazy is a thread-safe lazy-initialized singleton value.
func Lazy[T any](new func() (T, error)) Singleton[T] {
	return &lazy[T]{new: new}
}

// NewLazy is a thread-safe lazy-initialized singleton value.
//
// Deprecated: Use Lazy instead.
func NewLazy[T any](new func() (T, error)) Singleton[T] {
	return &lazy[T]{new: new}
}

func (l *lazy[T]) InstanceCtx(ctx context.Context) (T, error) {
	select {
	case <-l.latch.Done():
		return l.instance, l.err
	default:
		// Intentionally left blank.
	}

	go func() {
		async.Resolve(&l.latch, func() {
			l.instance, l.err = l.new()
		})
	}()

	select {
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	case <-l.latch.Done():
		return l.instance, l.err
	}
}

func (l *lazy[T]) Instance() (T, error) {
	return l.InstanceCtx(context.Background())
}
