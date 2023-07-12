package app

import (
	"atom/internal/models"
	"context"
)

type Repository interface {
	CreateCar(context.Context, models.Car) error
	GetCar(context.Context, string) (models.Car, error)
	GetAllCars(context.Context) ([]models.Car, error)
	DeleteCar(context.Context, string) error
	AddMileage(context.Context, string, int) error
}
