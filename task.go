// Package ticker provides a flexible and customizable ticker implementation
// for executing periodic tasks.
//
// It allows users to create tickers that execute a given task at specified intervals,
// with options for immediate execution and execution limits. This package is useful
// for scenarios requiring repeated task execution, such as polling, scheduled jobs,
// or periodic health checks.
package ticker

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Task represents a function that can be executed periodically.
type Task func() error

// New creates a new Task from the given task function.
// It returns a Task type that can be used with the Run method for periodic execution.
// If a nil function is provided, New returns nil.
func New(task func() error) Task {
	return Task(task)
}

// Run executes the task periodically according to the specified duration and options.
//
// It returns an error if the task encounters an error or if the context is canceled.
// The duration d must be greater than zero; if not, Run returns ErrNonPositiveInterval.
//
// Options can be used to customize the behavior:
//   - WithImmediate: Execute the task immediately before starting the ticker.
//   - WithLimit: Limit the number of executions.
//
// If no error occurs, Run will continue until the context is canceled or, if specified,
// the execution limit is reached.
func (task Task) Run(ctx context.Context, d time.Duration, options ...Option) error {
	if d <= 0 {
		return ErrNonPositiveInterval
	}

	if task == nil {
		return ErrNilFunction
	}

	c := &config{
		Limit: -1,
	}
	for _, opt := range options {
		opt.apply(c)
	}

	if c.Limit == 0 {
		return nil
	}
	if c.Limit > 0 {
		return task.runLimit(ctx, d, c)
	}
	return task.run(ctx, d, c)
}

// runLimit executes the task for a limited number of times or until the context is canceled.
// It respects the immediate execution option and returns early if the limit is reached.
func (task Task) runLimit(ctx context.Context, d time.Duration, c *config) error {
	limit := c.Limit
	if c.Immediate {
		if err := task(); err != nil {
			return err
		}
		limit--
		if limit == 0 {
			return nil
		}
	}
	t := time.NewTicker(d)
	defer t.Stop()
	for ; limit > 0; limit-- {
		select {
		case <-t.C:
			if err := task(); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// run executes the task indefinitely or until the context is canceled.
// It respects the immediate execution option.
func (task Task) run(ctx context.Context, d time.Duration, c *config) error {
	if c.Immediate {
		if err := task(); err != nil {
			return err
		}
	}
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := task(); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

var (
	// ErrInvalidArgument is the base error indicating that an invalid argument was provided.
	// It can be used to check if an error is related to invalid arguments:
	//
	//     err := task.Run(ctx, -1 * time.Second)
	//     if errors.Is(err, ErrInvalidArgument) {
	//         fmt.Println("An invalid argument was provided")
	//     }
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrNonPositiveInterval indicates that a non-positive interval was provided.
	// This error wraps ErrInvalidArgument, so errors.Is(ErrNonPositiveInterval, ErrInvalidArgument) will return true.
	ErrNonPositiveInterval = fmt.Errorf("%w: non-positive interval", ErrInvalidArgument)

	// ErrNilFunction indicates that a nil function was provided.
	// This error wraps ErrInvalidArgument, so errors.Is(ErrNilFunction, ErrInvalidArgument) will return true.
	ErrNilFunction = fmt.Errorf("%w: function must not be nil", ErrInvalidArgument)
)
