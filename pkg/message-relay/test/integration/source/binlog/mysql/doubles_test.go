//go:build integration

package mysql

import (
	"sync"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
)

// inMemoryStateHandler is a minimal, concurrency-safe StateHandler used to
// isolate this test from the real state persistence mechanism (Redis),
// which is covered by its own integration suite.
type inMemoryStateHandler struct {
	mu  sync.Mutex
	pos mysql.Position
}

func newInMemoryStateHandler(initial mysql.Position) *inMemoryStateHandler {
	return &inMemoryStateHandler{pos: initial}
}

func (h *inMemoryStateHandler) GetLastPosition() (mysql.Position, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.pos, nil
}

func (h *inMemoryStateHandler) SetLastPosition(pos mysql.Position) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.pos = pos
	return nil
}

func (h *inMemoryStateHandler) lastPosition() mysql.Position {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.pos
}

// capturingEventHandler is a canal.EventHandler test double that forwards
// captured row events and synced positions over channels, mirroring how
// relay.NewEventHandler wires production event handling.
type capturingEventHandler struct {
	canal.DummyEventHandler
	rows  chan *canal.RowsEvent
	posCh chan mysql.Position
}

func newCapturingEventHandler() (*capturingEventHandler, <-chan mysql.Position) {
	posCh := make(chan mysql.Position)
	return &capturingEventHandler{rows: make(chan *canal.RowsEvent, 10), posCh: posCh}, posCh
}

func (h *capturingEventHandler) OnRow(e *canal.RowsEvent) error {
	h.rows <- e
	return nil
}

func (h *capturingEventHandler) OnPosSynced(pos mysql.Position, _ mysql.GTIDSet, _ bool) error {
	h.posCh <- pos
	return nil
}
