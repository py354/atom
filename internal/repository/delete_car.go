package repository

import (
	"context"
)

func (r repository) DeleteCar(ctx context.Context, registrationPlate string) (err error) {
	const query = "delete from cars where RegistrationPlate = $1"

	result, err := r.pool.Exec(ctx, query, registrationPlate)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return
}
