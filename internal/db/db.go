package db

import (
	"atom/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context) (*pgxpool.Pool, error) {
	template := "host=%s port=%s user=%s dbname=%s password=%s sslmode=disable"

	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(template,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Dbname,
		c.Postgres.Password)

	return pgxpool.New(ctx, dsn)
}
