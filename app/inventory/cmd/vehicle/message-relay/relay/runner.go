package relay

import okrun "github.com/oklog/run"

// Source is the relay backend interface. Implementations stream events from a
// source (e.g. MySQL binlog, database polling) and handle their own
// checkpointing. Swap the implementation to change the relay strategy without
// touching Runner.
type Source interface {
	Run() error
	Close()
}

// Runner manages the lifecycle of a Source. It guarantees Close is always
// called when Run stops, regardless of whether it succeeds or fails.
type Runner struct {
	source Source
}

// NewRunner constructs a Runner for the given Source.
func NewRunner(source Source) *Runner {
	return &Runner{source: source}
}

// Run starts the source and blocks until it stops, then calls Close.
func (r *Runner) Run() error {
	var g okrun.Group
	g.Add(
		func() error { return r.source.Run() },
		func(error) { r.source.Close() },
	)
	return g.Run()
}
