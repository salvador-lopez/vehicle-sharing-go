package event_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/pkg/domain/event"
	"vehicle-sharing-go/pkg/domain/event/mock"
)

func TestPublishRecordedEvents(t *testing.T) {
	var (
		r           *require.Assertions
		ctx         context.Context
		mockCtrl    *gomock.Controller
		evtRecorder *mock.MockRecorder
		publisher   *mock.MockPublisher
		sut         *event.AgRootEventPublisher
	)

	setup := func() {
		r = require.New(t)
		ctx = context.Background()
		mockCtrl = gomock.NewController(t)
		evtRecorder = mock.NewMockRecorder(mockCtrl)
		publisher = mock.NewMockPublisher(mockCtrl)
		sut = event.NewAgRootEventPublisher(publisher)
	}

	teardown := func() {
		mockCtrl.Finish()
	}

	tests := []struct {
		name            string
		recordedEvents  []*event.Event
		evtPublisherErr error
		expectedErrMsg  string
	}{
		{
			name: "publish successfully",
			recordedEvents: []*event.Event{
				{
					uuid.New(),
					uuid.New(),
					"Car",
					"CarCreatedEvent",
					nil,
					time.Now(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()
			evtRecorder.EXPECT().RecordedEvents().Return(tt.recordedEvents)
			publisher.EXPECT().Publish(ctx, tt.recordedEvents).Return(tt.evtPublisherErr)
			err := sut.Publish(ctx, evtRecorder)
			if tt.expectedErrMsg != "" {
				r.EqualError(err, tt.expectedErrMsg)
				return
			}

			r.NoError(err)
		})
	}
}
