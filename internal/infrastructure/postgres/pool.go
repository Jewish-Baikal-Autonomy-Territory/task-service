package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pgxgeom "github.com/twpayne/pgx-geom"
)

type PoolOptions struct {
	ConnectionString      string
	MinConnections        int32
	MaxConnections        int32
	MaxIdleConnectionTime time.Duration
}

func NewPool(ctx context.Context, opts PoolOptions) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(opts.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	config.AfterConnect = pgxgeom.Register
	config.MinConns = opts.MinConnections
	config.MaxConns = opts.MaxConnections
	config.MaxConnIdleTime = opts.MaxIdleConnectionTime

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	return pool, nil
}
