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
