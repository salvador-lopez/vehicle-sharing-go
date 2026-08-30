//go:build unit

package binlog_test

import (
	"errors"
	"testing"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/app/inventory/cmd/vehicle/message-relay/source/binlog"
	"vehicle-sharing-go/app/inventory/cmd/vehicle/message-relay/source/binlog/mock"
)

// stubEventHandler is a minimal canal.EventHandler test double.
// DummyEventHandler provides no-op implementations of every interface method;
// we only override OnPosSynced so the position channel works as in production.
type stubEventHandler struct {
	canal.DummyEventHandler
	posCh chan mysql.Position
}

func newStubEventHandler() (*stubEventHandler, <-chan mysql.Position) {
	posCh := make(chan mysql.Position)
	return &stubEventHandler{posCh: posCh}, posCh
}

func (h *stubEventHandler) OnPosSynced(p mysql.Position, _ mysql.GTIDSet, _ bool) error {
	h.posCh <- p
	return nil
}

// --- suite ---

type mySQLSourceSuite struct {
	suite.Suite
	mockCtrl         *gomock.Controller
	mockCanal        *mock.MockCanal
	mockStateHandler *mock.MockStateHandler
	handler          *stubEventHandler
	posCh            <-chan mysql.Position
}

func TestMySQLSourceSuite(t *testing.T) {
	suite.Run(t, new(mySQLSourceSuite))
}

func (s *mySQLSourceSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockCanal = mock.NewMockCanal(s.mockCtrl)
	s.mockStateHandler = mock.NewMockStateHandler(s.mockCtrl)
	s.handler, s.posCh = newStubEventHandler()
}

func (s *mySQLSourceSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

// --- constructor ---

func (s *mySQLSourceSuite) TestNewSourceRegistersEventHandlerOnCanal() {
	s.mockCanal.EXPECT().SetEventHandler(s.handler)

	binlog.NewSource(s.mockCanal, s.handler, s.posCh, s.mockStateHandler, time.Second)
}

// --- Run: error paths ---

func (s *mySQLSourceSuite) TestRunReturnsGetLastPositionError() {
	s.mockCanal.EXPECT().SetEventHandler(s.handler)
	s.mockStateHandler.EXPECT().GetLastPosition().Return(mysql.Position{}, errors.New("redis unavailable"))

	sut := binlog.NewSource(s.mockCanal, s.handler, s.posCh, s.mockStateHandler, time.Hour)
	s.EqualError(sut.Run(), "redis unavailable")
}

func (s *mySQLSourceSuite) TestRunReturnsRunFromError() {
	s.mockCanal.EXPECT().SetEventHandler(s.handler)
	s.mockStateHandler.EXPECT().GetLastPosition().Return(mysql.Position{}, nil)
	s.mockCanal.EXPECT().RunFrom(mysql.Position{}).Return(errors.New("binlog stream broken"))
	s.mockCanal.EXPECT().Close()
	s.mockStateHandler.EXPECT().SetLastPosition(mysql.Position{}).Return(nil)

	sut := binlog.NewSource(s.mockCanal, s.handler, s.posCh, s.mockStateHandler, time.Hour)
	s.EqualError(sut.Run(), "binlog stream broken")
}

// --- Run: happy path ---

func (s *mySQLSourceSuite) TestRunCallsRunFromWithLastCheckpointedPosition() {
	lastPos := mysql.Position{Name: "binlog.000003", Pos: 1024}
	s.mockCanal.EXPECT().SetEventHandler(s.handler)
	s.mockStateHandler.EXPECT().GetLastPosition().Return(lastPos, nil)
	s.mockCanal.EXPECT().RunFrom(lastPos).Return(nil)
	s.mockCanal.EXPECT().Close()
	s.mockStateHandler.EXPECT().SetLastPosition(lastPos).Return(nil)

	sut := binlog.NewSource(s.mockCanal, s.handler, s.posCh, s.mockStateHandler, time.Hour)
	s.NoError(sut.Run())
}

// --- Close ---

func (s *mySQLSourceSuite) TestCloseCallsCanalClose() {
	s.mockCanal.EXPECT().SetEventHandler(s.handler)
	s.mockCanal.EXPECT().Close()

	sut := binlog.NewSource(s.mockCanal, s.handler, s.posCh, s.mockStateHandler, time.Second)
	sut.Close()
}

// --- position channel update ---

func (s *mySQLSourceSuite) TestRunUpdatesLastPositionFromPositionChan() {
	newPos := mysql.Position{Name: "binlog.000005", Pos: 4096}
	canClose := make(chan struct{})

	s.mockCanal.EXPECT().SetEventHandler(s.handler)
	s.mockStateHandler.EXPECT().GetLastPosition().Return(mysql.Position{}, nil)
	s.mockCanal.EXPECT().RunFrom(mysql.Position{}).DoAndReturn(func(mysql.Position) error {
		<-canClose
		return nil
	})
	s.mockCanal.EXPECT().Close()
	// The ctx.Done final checkpoint must use the updated position, not the initial empty one.
	s.mockStateHandler.EXPECT().SetLastPosition(newPos).Return(nil)

	sut := binlog.NewSource(s.mockCanal, s.handler, s.posCh, s.mockStateHandler, time.Hour)

	done := make(chan error)
	go func() { done <- sut.Run() }()

	// OnPosSynced does a blocking channel send; it returns only after Actor 2 has read the
	// value, which guarantees lastPosition is updated before we trigger shutdown.
	posDelivered := make(chan struct{})
	go func() {
		_ = s.handler.OnPosSynced(newPos, nil, false)
		close(posDelivered)
	}()
	<-posDelivered

	close(canClose) // shut down RunFrom; ctx.Done final checkpoint will use newPos
	s.NoError(<-done)
}

// --- ticker checkpoint error ---

func (s *mySQLSourceSuite) TestRunReturnsCheckpointErrorOnTick() {
	canClose := make(chan struct{})
	checkpointErr := errors.New("checkpoint failed")

	s.mockCanal.EXPECT().SetEventHandler(s.handler)
	s.mockStateHandler.EXPECT().GetLastPosition().Return(mysql.Position{}, nil)
	s.mockCanal.EXPECT().RunFrom(mysql.Position{}).DoAndReturn(func(mysql.Position) error {
		<-canClose
		return nil
	})
	// When Actor 2 returns an error, okrun calls canal.Close() as the Actor 1 interrupt.
	// That interrupt must unblock RunFrom so the goroutine can exit cleanly.
	s.mockCanal.EXPECT().Close().DoAndReturn(func() { close(canClose) })
	s.mockStateHandler.EXPECT().SetLastPosition(mysql.Position{}).Return(checkpointErr)

	sut := binlog.NewSource(s.mockCanal, s.handler, s.posCh, s.mockStateHandler, 5*time.Millisecond)
	s.ErrorIs(sut.Run(), checkpointErr)
}
