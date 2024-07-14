package ticker_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/goaux/ticker"
)

func Example() {
	fn := func() error {
		fmt.Println("Executing task...")
		return nil
	}

	task := ticker.New(fn)
	ctx := context.Background()

	// Run the ticker every second, for a maximum of 5 executions
	err := task.Run(ctx, time.Second, ticker.WithLimit(5), ticker.WithImmediate(true))
	if err != nil {
		fmt.Printf("Task stopped with error: %v\n", err)
	}
	// Output:
	// Executing task...
	// Executing task...
	// Executing task...
	// Executing task...
	// Executing task...
}

// TestNew ensures that New creates a Task correctly
func TestNew(t *testing.T) {
	fn := func() error { return nil }
	task := ticker.New(fn)
	if task == nil {
		t.Error("New should return a non-nil Task")
	}

	task = ticker.New(nil)
	if task != nil {
		t.Error("New(nil) should return a nil Task")
	}
}

// TestTask_Run tests various scenarios of the Run method
func TestTask_Run(t *testing.T) {
	ErrTask := errors.New("task error")

	tests := []struct {
		name        string
		task        func() error
		duration    time.Duration
		options     []ticker.Option
		expectedErr error
		runTime     time.Duration
	}{
		{
			name:        "Invalid duration",
			task:        func() error { return nil },
			duration:    -time.Second,
			expectedErr: ticker.ErrNonPositiveInterval,
		},
		{
			name:        "Nil function",
			task:        nil,
			duration:    time.Second,
			expectedErr: ticker.ErrNilFunction,
		},
		{
			name:        "Run with limit",
			task:        func() error { return nil },
			duration:    100 * time.Millisecond,
			options:     []ticker.Option{ticker.WithLimit(3)},
			runTime:     350 * time.Millisecond,
			expectedErr: nil,
		},
		{
			name:        "Run with error",
			task:        func() error { return ErrTask },
			duration:    100 * time.Millisecond,
			options:     []ticker.Option{ticker.WithImmediate(true)},
			expectedErr: ErrTask,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := ticker.New(tt.task)
			ctx, cancel := context.WithTimeout(context.Background(), tt.runTime)
			defer cancel()

			err := task.Run(ctx, tt.duration, tt.options...)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

// TestWithImmediately tests the WithImmediately option
func TestWithImmediate(t *testing.T) {
	count := 0
	fn := func() error {
		count++
		return nil
	}

	task := ticker.New(fn)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	err := task.Run(ctx, 100*time.Millisecond, ticker.WithImmediate(true))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("unexpected error: %v", err)
	}

	if count < 2 {
		t.Errorf("expected at least 2 executions, got %d", count)
	}
}

// TestWithLimit tests the WithLimit option
func TestWithLimit(t *testing.T) {
	count := 0
	fn := func() error {
		count++
		return nil
	}

	task := ticker.New(fn)
	ctx := context.Background()

	err := task.Run(ctx, 10*time.Millisecond, ticker.WithLimit(5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if count != 5 {
		t.Errorf("expected 5 executions, got %d", count)
	}
}

// TestContextCancellation tests if the ticker stops when the context is canceled
func TestContextCancellation(t *testing.T) {
	count := 0
	fn := func() error {
		count++
		return nil
	}

	task := ticker.New(fn)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := task.Run(ctx, 10*time.Millisecond)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	if count == 0 {
		t.Error("expected at least one execution before cancellation")
	}
}
