package singleton

import (
	"context"
	"sync"

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

	once  sync.Once
	latch async.Latch
}

func NewLazy[T any](new func() (T, error)) Singleton[T] {
	return &lazy[T]{new: new}
}

func (l *lazy[T]) InstanceCtx(ctx context.Context) (T, error) {
	go func() {
		l.once.Do(func() {
			i, err := l.new()
			async.Resolve(&l.latch, func() {
				l.instance = i
				l.err = err
			})
		})
	}()

	select {
	case <-l.latch.Done():
		return l.instance, l.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

func (l *lazy[T]) Instance() (T, error) {
	return l.InstanceCtx(context.Background())
}
