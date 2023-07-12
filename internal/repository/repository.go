package repository

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) repository {
	return repository{pool: pool}
}
