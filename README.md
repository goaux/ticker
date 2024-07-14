# ticker
The `ticker` package provides a flexible and customizable ticker implementation for executing periodic tasks in Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/goaux/ticker.svg)](https://pkg.go.dev/github.com/goaux/ticker)
[![Go Report Card](https://goreportcard.com/badge/github.com/goaux/ticker)](https://goreportcard.com/report/github.com/goaux/ticker)

## Features

- Execute tasks at specified intervals
- Option for immediate execution before starting the ticker
- Limit the number of executions
- Context-aware for easy cancellation and timeout handling
- Customizable through functional options

## Installation

To install the `ticker` package, use `go get`:

    go get github.com/goaux/ticker

## Usage

Here's a basic example of how to use the `ticker` package:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/goaux/ticker"
)

func main() {
    task := ticker.New(func() error {
        fmt.Println("Executing task...")
        return nil
    })

    ctx := context.Background()
    err := task.Run(ctx, time.Second, ticker.WithLimit(5), ticker.WithImmediate(true))
    if err != nil {
        fmt.Printf("Task stopped with error: %v\n", err)
    }
}
```

This example creates a task that prints "Executing task..." every second, for a maximum of 5 executions, starting immediately.

## Best Practices

1. Always use a context to manage the lifecycle of your ticker:

    ```go
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    err := task.Run(ctx, time.Second)
    ```

2. Handle errors returned by the `Run` method:

    ```go
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            fmt.Println("Ticker stopped due to timeout")
        } else {
            fmt.Printf("Ticker stopped with error: %v\n", err)
        }
    }
    ```

3. Use the `WithImmediate` option when you want the task to execute right away:

    ```go
    task.Run(ctx, time.Minute, ticker.WithImmediate(true))
    ```

4. Use the `WithLimit` option to control the number of executions:

    ```go
    task.Run(ctx, time.Hour, ticker.WithLimit(24)) // Run 24 times (once per hour for a day)
    ```

5. For tasks that may fail, implement proper error handling in your task function:

    ```go
    task := ticker.New(func() error {
        if err := doSomething(); err != nil {
            log.Printf("Task error: %v", err)
            return nil // Continue running despite errors
        }
        return nil
    })
    ```

## Error Handling

The module defines several error types:

- `ErrInvalidArgument`: Base error for invalid arguments.
- `ErrNonPositiveInterval`: Indicates that a non-positive interval was provided.
- `ErrNilFunction`: Indicates that a nil function was provided.

These errors can be checked using `errors.Is()`.
