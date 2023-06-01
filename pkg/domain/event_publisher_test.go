package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/pkg/domain"
	"vehicle-sharing-go/pkg/domain/mock"
)

func TestPublishRecordedEvents(t *testing.T) {
	var (
		r           *require.Assertions
		ctx         context.Context
		mockCtrl    *gomock.Controller
		evtRecorder *mock.MockEventRecorder
		publisher   *mock.MockEventPublisher
		sut         *domain.AgRootEventPublisher
	)

	setup := func() {
		r = require.New(t)
		ctx = context.Background()
		mockCtrl = gomock.NewController(t)
		evtRecorder = mock.NewMockEventRecorder(mockCtrl)
		publisher = mock.NewMockEventPublisher(mockCtrl)
		sut = domain.NewAgRootEventPublisher(publisher)
	}

	teardown := func() {
		mockCtrl.Finish()
	}

	tests := []struct {
		name            string
		topic           string
		recordedEvents  []*domain.Event
		evtPublisherErr error
		expectedErrMsg  string
	}{
		{
			name:  "publish successfully",
			topic: "inventory",
			recordedEvents: []*domain.Event{
				domain.NewEvent(
					uuid.New(),
					uuid.New(),
					"Car",
					"CarCreatedEvent",
					nil,
					time.Now(),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			defer teardown()
			evtRecorder.EXPECT().RecordedEvents().Return(tt.recordedEvents)
			publisher.EXPECT().Publish(ctx, tt.topic, tt.recordedEvents).Return(tt.evtPublisherErr)
			err := sut.Publish(ctx, tt.topic, evtRecorder)
			if tt.expectedErrMsg != "" {
				r.EqualError(err, tt.expectedErrMsg)
				return
			}

			r.NoError(err)
		})
	}
}
