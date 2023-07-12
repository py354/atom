package repository

import (
	"atom/internal/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (r repository) GetCar(ctx context.Context, registrationPlate string) (models.Car, error) {
	const query = `
		select  RegistrationPlate,
				Model,
				Purpose,
				ManufactureYear,
				Mileage
		from cars where RegistrationPlate = $1
	`

	car := models.Car{}
	row := r.pool.QueryRow(ctx, query, registrationPlate)
	err := row.Scan(&car.RegistrationPlate, &car.Model, &car.Purpose, &car.ManufactureYear, &car.Mileage)
	if errors.Is(err, pgx.ErrNoRows) {
		return car, ErrNotFound
	}

	return car, err
}

func (r repository) GetAllCars(ctx context.Context) ([]models.Car, error) {
	const query = `
		select  RegistrationPlate,
				Model,
				Purpose,
				ManufactureYear,
				Mileage
		from cars
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.Car, 0)
	for rows.Next() {
		car := models.Car{}
		err = rows.Scan(&car.RegistrationPlate, &car.Model, &car.Purpose, &car.ManufactureYear, &car.Mileage)
		if err != nil {
			return nil, err
		}

		result = append(result, car)
	}

	return result, err
}
