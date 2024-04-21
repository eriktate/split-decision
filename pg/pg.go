package pg

import (
	"context"
	"fmt"

	uuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PG struct {
	pool *pgxpool.Pool
}

func New(url string) (PG, error) {
	poolCfg, err := pgxpool.ParseConfig(url)
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		uuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return PG{}, fmt.Errorf("failed to create PG connection pool: %w", err)
	}

	return PG{pool}, nil
}

func (pg PG) Close() {
	pg.pool.Close()
}
