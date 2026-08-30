package binlog

import (
	"context"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	okrun "github.com/oklog/run"
	"github.com/sirupsen/logrus"
)

// Canal is the subset of canal.Canal used by Source, allowing it to be
// replaced by a test double.
//go:generate mockgen -destination=mock/canal_mock.go -package=mock . Canal
type Canal interface {
	RunFrom(mysql.Position) error
	SetEventHandler(handler canal.EventHandler)
	Close()
}

//go:generate mockgen -destination=mock/state_handler_mock.go -package=mock . StateHandler
type StateHandler interface {
	GetLastPosition() (mysql.Position, error)
	SetLastPosition(position mysql.Position) error
}

// Source implements relay.Source by streaming MySQL binlog changes via a Canal.
// On each tick the most recently seen binlog position is persisted via
// StateHandler so the relay can resume after a restart.
type Source struct {
	canal                Canal
	stateHandler         StateHandler
	positionChan         <-chan mysql.Position
	stateUpdateFrequency time.Duration
}

// NewSource constructs a Source, registering handler on the canal.
// positionChan is the channel returned by relay.NewEventHandler; it carries
// the latest binlog position as the canal advances.
func NewSource(
	c Canal,
	handler canal.EventHandler,
	positionChan <-chan mysql.Position,
	stateHandler StateHandler,
	stateUpdateFrequency time.Duration,
) *Source {
	c.SetEventHandler(handler)

	return &Source{
		canal:                c,
		stateHandler:         stateHandler,
		positionChan:         positionChan,
		stateUpdateFrequency: stateUpdateFrequency,
	}
}

// Run starts streaming from the last checkpointed binlog position. It blocks
// until either the canal stream or the checkpoint ticker exits, then stops
// both and returns the first error.
func (s *Source) Run() error {
	lastPosition, err := s.stateHandler.GetLastPosition()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(s.stateUpdateFrequency)

	// startPosition is the fixed resume point passed to RunFrom. lastPosition
	// is updated by Actor 2 as the canal advances and is only used for checkpointing.
	startPosition := lastPosition

	var g okrun.Group

	// Actor 1: stream binlog changes from MySQL.
	g.Add(
		func() error { return s.canal.RunFrom(startPosition) },
		func(error) { s.canal.Close() },
	)

	// Actor 2: periodically checkpoint the latest binlog position.
	g.Add(
		func() error {
			for {
				select {
				case pos := <-s.positionChan:
					lastPosition = pos
				case <-ticker.C:
					if err := s.checkpoint(lastPosition); err != nil {
						return err
					}
				case <-ctx.Done():
					// Final checkpoint before exit.
					_ = s.checkpoint(lastPosition)
					return nil
				}
			}
		},
		func(error) {
			ticker.Stop()
			cancel()
		},
	)

	return g.Run()
}

// Close closes the underlying canal.
func (s *Source) Close() {
	s.canal.Close()
}

func (s *Source) checkpoint(p mysql.Position) error {
	if err := s.stateHandler.SetLastPosition(p); err != nil {
		return err
	}
	logrus.WithField("position", p).Debug("binlog position checkpointed")
	return nil
}
