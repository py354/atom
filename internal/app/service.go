package app

import (
	"atom/internal/models"
	"atom/internal/repository"
	pb "atom/pkg/api"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type tserver struct {
	repo Repository
	pb.UnimplementedFleetServer
}

func New(repo Repository) *tserver {
	return &tserver{repo: repo}
}

var (
	errCarDoesntExists = status.Errorf(codes.NotFound, "car with this plate doesn't exists")
	errUnknown         = status.Errorf(codes.Unknown, "database error, maybe bad data")
)

func (t *tserver) RegisterCar(ctx context.Context, car *pb.Car) (*pb.Empty, error) {
	err := t.repo.CreateCar(ctx, models.Car{
		RegistrationPlate: car.RegistrationPlate,
		Model:             car.Model,
		Purpose:           car.Purpose.String(),
		ManufactureYear:   int(car.ManufactureYear),
		Mileage:           int(car.Mileage),
	})

	if err == nil {
		return &pb.Empty{}, nil
	}

	return nil, errUnknown
}

func (t *tserver) GetCar(ctx context.Context, rp *pb.RP) (*pb.Car, error) {
	car, err := t.repo.GetCar(ctx, rp.GetRegistrationPlate())

	if errors.Is(err, repository.ErrNotFound) {
		return nil, errCarDoesntExists
	} else if err != nil {
		return nil, errUnknown
	}

	return &pb.Car{
		RegistrationPlate: car.RegistrationPlate,
		Model:             car.Model,
		Purpose:           pb.Purpose(pb.Purpose_value[car.Purpose]),
		ManufactureYear:   int32(car.ManufactureYear),
		Mileage:           int32(car.Mileage),
	}, nil
}

func (t *tserver) GetAllCars(ctx context.Context, _ *pb.Empty) (*pb.CarList, error) {
	cars, err := t.repo.GetAllCars(ctx)
	if err != nil {
		return nil, errUnknown
	}

	result := make([]*pb.Car, 0)
	for _, car := range cars {
		result = append(result, &pb.Car{
			RegistrationPlate: car.RegistrationPlate,
			Model:             car.Model,
			Purpose:           pb.Purpose(pb.Purpose_value[car.Purpose]),
			ManufactureYear:   int32(car.ManufactureYear),
			Mileage:           int32(car.Mileage),
		})
	}

	return &pb.CarList{Cars: result}, nil
}

func (t *tserver) AddMileage(ctx context.Context, addMailReq *pb.AddMileageRequest) (*pb.Empty, error) {
	err := t.repo.AddMileage(ctx, addMailReq.GetRegistrationPlate(), int(addMailReq.GetMileage()))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, errCarDoesntExists
	} else if err != nil {
		return nil, errUnknown
	}
	return &pb.Empty{}, nil
}

func (t *tserver) DeleteCar(ctx context.Context, rp *pb.RP) (*pb.Empty, error) {
	err := t.repo.DeleteCar(ctx, rp.GetRegistrationPlate())
	if errors.Is(err, repository.ErrNotFound) {
		return nil, errCarDoesntExists
	} else if err != nil {
		return nil, errUnknown
	}

	return &pb.Empty{}, nil
}

func (t *tserver) GetEstimatedCost(ctx context.Context, rp *pb.RP) (*pb.EstimatedCostResp, error) {
	car, err := t.GetCar(ctx, rp)
	if err != nil {
		return nil, err
	}

	factor1 := float64(car.Mileage) * 0.1
	factor2 := float64(time.Now().Year()) - float64(car.ManufactureYear)
	var factor3 float64

	switch car.Purpose.String() {
	case "DELIVERY":
		factor3 = 100_000
	case "SHARING":
		factor3 = 1_000_000
	case "TAXI":
		factor3 = 2_000_000
	}

	cost := 10_000_000.0 - factor1 - factor2 - factor3
	if cost < 0 {
		cost = 0
	}

	return &pb.EstimatedCostResp{Cost: int64(cost)}, nil
}
