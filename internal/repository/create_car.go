package repository

import (
	"atom/internal/models"
	"context"
)

func (r repository) CreateCar(ctx context.Context, car models.Car) (err error) {
	const query = `
		insert into cars 
		    (RegistrationPlate, Model, Purpose, ManufactureYear, Mileage) 
		VALUES 
		    ($1, $2, $3, $4, $5)
`

	_, err = r.pool.Exec(ctx, query,
		car.RegistrationPlate,
		car.Model,
		car.Purpose,
		car.ManufactureYear,
		car.Mileage,
	)

	return
}
