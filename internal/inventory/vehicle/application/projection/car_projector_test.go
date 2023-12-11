//go:build unit || unit_inventory

package projection_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
	"vehicle-sharing-go/internal/inventory/vehicle/application/projection/mock"
	"vehicle-sharing-go/internal/inventory/vehicle/domain/event"
)

type carProjectorUnitSuite struct {
	suite.Suite
	ctx           context.Context
	mockCtrl      *gomock.Controller
	vinDecoder    *mock.MockVINDecoder
	carRepo       *mock.MockCarRepository
	carProjection *projection.Car
	sut           *projection.CarProjector
}

func (s *carProjectorUnitSuite) SetupTest() {
	s.ctx = context.Background()

	s.mockCtrl = gomock.NewController(s.T())
	s.vinDecoder = mock.NewMockVINDecoder(s.mockCtrl)
	s.carRepo = mock.NewMockCarRepository(s.mockCtrl)

	s.carProjection = s.buildCarProjection(
		uuid.New(),
		time.Now(),
		time.Now().Add(time.Second),
		s.buildVinDataProjection(
			"4Y1SL65848Z411439",
			"country",
			"manufacturer",
			"brand",
			"2000",
			"Gasoline",
			"Jazz",
			"2023",
			"Barcelona",
			"411439",
		),
		"Spectral Blue",
	)

	s.sut = projection.NewCarProjector(s.vinDecoder, s.carRepo)
}

func (s *carProjectorUnitSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestVehicleUnitSuite(t *testing.T) {
	suite.Run(t, new(carProjectorUnitSuite))
}

func (s *carProjectorUnitSuite) TestProjectCarCreated() {
	s.vinDecoder.EXPECT().Decode(s.ctx, s.carProjection.VIN).Return(s.carProjection.VINData, nil)
	s.carRepo.EXPECT().Create(s.ctx, s.carProjection).Return(nil)

	s.Require().NoError(s.runSut())
}

func (s *carProjectorUnitSuite) TestProjectCarCreatedVINDecoderErr() {
	vinDecoderErr := errors.New("vin decoder error")
	s.vinDecoder.EXPECT().Decode(s.ctx, s.carProjection.VIN).Return(nil, vinDecoderErr)

	s.Require().EqualError(s.runSut(), vinDecoderErr.Error())
}

func (s *carProjectorUnitSuite) TestProjectCarCreatedCarRepoCreateErr() {
	s.vinDecoder.EXPECT().Decode(s.ctx, s.carProjection.VIN).Return(s.carProjection.VINData, nil)

	carRepoCreateErr := errors.New("car repo create error")
	s.carRepo.EXPECT().Create(s.ctx, gomock.Any()).Return(carRepoCreateErr)

	s.Require().EqualError(s.runSut(), carRepoCreateErr.Error())
}

func (s *carProjectorUnitSuite) buildCarProjection(
	id uuid.UUID,
	createdAt,
	updatedAt time.Time,
	vinData *projection.VINData,
	color string,
) *projection.Car {
	return &projection.Car{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		VINData:   vinData,
		Color:     color,
	}
}

func (s *carProjectorUnitSuite) buildVinDataProjection(
	vinNumber,
	country,
	manufacturer,
	brand,
	engineSize,
	fuelType,
	model,
	year,
	assemblyPlant,
	sn string,
) *projection.VINData {
	return &projection.VINData{
		VIN:           vinNumber,
		Country:       &country,
		Manufacturer:  &manufacturer,
		Brand:         &brand,
		EngineSize:    &engineSize,
		FuelType:      &fuelType,
		Model:         &model,
		Year:          &year,
		AssemblyPlant: &assemblyPlant,
		SN:            &sn,
	}
}

func (s *carProjectorUnitSuite) runSut() error {
	payload := &event.CarCreatedPayload{
		VinNumber: s.carProjection.VIN,
		Color:     s.carProjection.Color,
		CreatedAt: s.carProjection.CreatedAt,
		UpdatedAt: s.carProjection.UpdatedAt,
	}

	return s.sut.ProjectCarCreated(s.ctx, s.carProjection.ID, payload)
}
