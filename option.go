package ticker

// Option represents a configuration option for the ticker.
// It is used to modify the behavior of a Task when running.
type Option interface {
	apply(*config)
}

// config holds the configuration for a ticker.
type config struct {
	Immediate bool
	Limit     int
}

// WithImmediate returns an Option to set whether the task should be executed immediately
// before starting the ticker.
func WithImmediate(v bool) Option {
	return immediate(v)
}

type immediate bool

func (o immediate) apply(c *config) {
	c.Immediate = bool(o)
}

// WithLimit returns an Option to set the maximum number of times the task should be executed.
//
// A value of 0 means no execution.
// A negative value means no limit (infinite executions).
func WithLimit(v int) Option {
	return limit(v)
}

type limit int

func (o limit) apply(c *config) {
	c.Limit = int(o)
}
