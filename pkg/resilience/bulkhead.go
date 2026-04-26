package resilience

import (
	"context"
	"errors"
)

var (
	ErrBulkheadFull = errors.New("bulkhead is full")
)

// Bulkhead defines the interface for bulkhead operations
type Bulkhead interface {
	Execute(ctx context.Context, fn func() error) error
}

type semaphoreBulkhead struct {
	sem chan struct{}
}

// NewBulkhead creates a new bulkhead with specified max concurrency
func NewBulkhead(maxConcurrency int) Bulkhead {
	return &semaphoreBulkhead{
		sem: make(chan struct{}, maxConcurrency),
	}
}

func (b *semaphoreBulkhead) Execute(ctx context.Context, fn func() error) error {
	select {
	case b.sem <- struct{}{}:
		defer func() { <-b.sem }()
		return fn()
	case <-ctx.Done():
		return ctx.Err()
	default:
		return ErrBulkheadFull
	}
}
