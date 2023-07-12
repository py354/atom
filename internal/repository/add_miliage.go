package repository

import "context"

func (r repository) AddMileage(ctx context.Context, registrationPlate string, mileage int) (err error) {
	const query = "update cars set Mileage=Mileage+$2 where RegistrationPlate = $1"

	result, err := r.pool.Exec(ctx, query, registrationPlate, mileage)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return
}
