//go:build unit

package event_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/pkg/domain/event"
	"vehicle-sharing-go/pkg/domain/event/mock"
)

func TestPublishRecordedEvents(t *testing.T) {
	recordedEvents := []*event.Event{
		{
			ID:            uuid.New(),
			AggregateID:   uuid.New(),
			AggregateType: "Car",
			EventType:     "CarCreatedEvent",
			Payload:       map[string]string{"vin": "vehicle-123"},
			Timestamp:     time.Now(),
		},
	}

	tests := []struct {
		name            string
		recordedEvents  []*event.Event
		evtPublisherErr error
		expectedErrMsg  string
	}{
		{
			name:            "preserves publisher error cause",
			recordedEvents:  recordedEvents,
			evtPublisherErr: errors.New("outbox unavailable"),
			expectedErrMsg:  "failed to publish events recorded in the aggregate root, cause: outbox unavailable",
		},
		{
			name:            "publish fails",
			recordedEvents:  recordedEvents,
			evtPublisherErr: errors.New("publish error"),
			expectedErrMsg:  "failed to publish events recorded in the aggregate root, cause: publish error",
		},
		{
			name:           "publish successfully",
			recordedEvents: recordedEvents,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			ctx := context.Background()
			mockCtrl := gomock.NewController(t)
			evtRecorder := mock.NewMockRecorder(mockCtrl)
			publisher := mock.NewMockPublisher(mockCtrl)
			sut := event.NewAgRootEventPublisher(publisher)

			evtRecorder.EXPECT().RecordedEvents().Return(tt.recordedEvents)
			publisher.EXPECT().Publish(ctx, tt.recordedEvents).Return(tt.evtPublisherErr)
			err := sut.Publish(ctx, evtRecorder)
			if tt.expectedErrMsg != "" {
				r.ErrorIs(err, tt.evtPublisherErr)
				serializedEvents, marshalErr := json.Marshal(tt.recordedEvents)
				r.NoError(marshalErr)
				r.EqualError(err, tt.expectedErrMsg+", events: "+string(serializedEvents))
				return
			}

			r.NoError(err)
		})
	}
}
