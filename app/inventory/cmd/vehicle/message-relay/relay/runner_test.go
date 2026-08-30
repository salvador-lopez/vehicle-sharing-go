//go:build unit

package relay_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/app/inventory/cmd/vehicle/message-relay/relay"
)

type runnerSuite struct {
	suite.Suite
}

func TestRunnerSuite(t *testing.T) {
	suite.Run(t, new(runnerSuite))
}

// stubSource is a handwritten test double for relay.Source.
type stubSource struct {
	runErr      error
	closeCalled bool
}

func (s *stubSource) Run() error { return s.runErr }
func (s *stubSource) Close()     { s.closeCalled = true }

func (s *runnerSuite) TestRunDelegatesToSource() {
	runner := relay.NewRunner(&stubSource{})
	s.NoError(runner.Run())
}

func (s *runnerSuite) TestRunReturnsSourceError() {
	expected := errors.New("source error")
	runner := relay.NewRunner(&stubSource{runErr: expected})
	s.ErrorIs(runner.Run(), expected)
}

func (s *runnerSuite) TestRunCallsSourceCloseOnCompletion() {
	source := &stubSource{}
	_ = relay.NewRunner(source).Run()
	s.True(source.closeCalled, "Close must always be called after Run completes")
}

func (s *runnerSuite) TestRunCallsSourceCloseOnError() {
	source := &stubSource{runErr: errors.New("source error")}
	_ = relay.NewRunner(source).Run()
	s.True(source.closeCalled, "Close must be called even when Run fails")
}
