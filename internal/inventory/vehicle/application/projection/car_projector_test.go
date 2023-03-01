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
	"vehicle-sharing-go/internal/inventory/vehicle/domain"
)

const (
	validVinNumber = "4Y1SL65848Z411439"
	color          = "Spectral Blue"
)

type carProjectorUnitSuite struct {
	suite.Suite
	ctx        context.Context
	mockCtrl   *gomock.Controller
	vinDecoder *mock.MockVINDecoder
	carRepo    *mock.MockCarRepository
	sut        *projection.CarProjector
}

func (s *carProjectorUnitSuite) SetupTest() {
	s.ctx = context.Background()

	s.mockCtrl = gomock.NewController(s.T())
	s.vinDecoder = mock.NewMockVINDecoder(s.mockCtrl)
	s.carRepo = mock.NewMockCarRepository(s.mockCtrl)

	s.sut = projection.NewCarProjector(s.vinDecoder, s.carRepo)
}

func (s *carProjectorUnitSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestVehicleUnitSuite(t *testing.T) {
	suite.Run(t, new(carProjectorUnitSuite))
}

func (s *carProjectorUnitSuite) TestProjectCarCreated() {
	carID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Second)

	vinData := &projection.VINData{
		VIN:           validVinNumber,
		Country:       "country",
		Manufacturer:  "manufacturer",
		Brand:         "brand",
		EngineSize:    "2000",
		FuelType:      "Gasoline",
		Model:         "Jazz",
		Year:          "2023",
		AssemblyPlant: "Barcelona",
		SN:            "411439",
	}
	s.vinDecoder.EXPECT().Decode(s.ctx, validVinNumber).Return(vinData, nil)
	carProjection := &projection.Car{
		ID:        carID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		VINData:   vinData,
		Color:     color,
	}
	s.carRepo.EXPECT().Create(s.ctx, carProjection).Return(nil)

	payload := &domain.CarCreatedEventPayloadDTO{
		VinNumber: validVinNumber,
		Color:     color,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	s.Require().NoError(s.sut.ProjectCarCreated(s.ctx, carID, payload))
}

func (s *carProjectorUnitSuite) TestProjectCarCreatedVINDecoderErr() {
	carID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Second)

	vinDecoderErr := errors.New("vin decoder error")
	s.vinDecoder.EXPECT().Decode(s.ctx, validVinNumber).Return(nil, vinDecoderErr)

	payload := &domain.CarCreatedEventPayloadDTO{
		VinNumber: validVinNumber,
		Color:     color,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	s.Require().EqualError(s.sut.ProjectCarCreated(s.ctx, carID, payload), vinDecoderErr.Error())
}

func (s *carProjectorUnitSuite) TestProjectCarCreatedCarRepoCreateErr() {
	carID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now().Add(time.Second)

	vinData := &projection.VINData{
		VIN:           validVinNumber,
		Country:       "country",
		Manufacturer:  "manufacturer",
		Brand:         "brand",
		EngineSize:    "2000",
		FuelType:      "Gasoline",
		Model:         "Jazz",
		Year:          "2023",
		AssemblyPlant: "Barcelona",
		SN:            "411439",
	}
	s.vinDecoder.EXPECT().Decode(s.ctx, validVinNumber).Return(vinData, nil)

	carRepoCreateErr := errors.New("car repo create error")
	s.carRepo.EXPECT().Create(s.ctx, gomock.Any()).Return(carRepoCreateErr)

	payload := &domain.CarCreatedEventPayloadDTO{
		VinNumber: validVinNumber,
		Color:     color,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	s.Require().EqualError(s.sut.ProjectCarCreated(s.ctx, carID, payload), carRepoCreateErr.Error())
}
