package db

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context) (*Store, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://metrics:metrics@localhost:5432/metricsdb?sslmode=disable"
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// init schema
	ctxInit, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = pool.Exec(ctxInit, `
CREATE TABLE IF NOT EXISTS metrics (
  id       BIGSERIAL PRIMARY KEY,
  hostname TEXT            NOT NULL,
  cpu      DOUBLE PRECISION NOT NULL,
  memory   DOUBLE PRECISION NOT NULL,
  ts       TIMESTAMPTZ      NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_metrics_host_ts ON metrics(hostname, ts DESC);
`)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &Store{Pool: pool}, nil
}

func (s *Store) Close() {
	if s != nil && s.Pool != nil {
		s.Pool.Close()
	}
}

func (s *Store) InsertMetric(ctx context.Context, hostname string, cpu, mem float64, at time.Time) error {
	if s == nil || s.Pool == nil {
		return errors.New("nil store")
	}
	_, err := s.Pool.Exec(ctx, `
INSERT INTO metrics (hostname, cpu, memory, ts)
VALUES ($1, $2, $3, $4)
`, hostname, cpu, mem, at)
	return err
}

// Optionnel: lecture simple pour tests / Grafana POC
type MetricRow struct {
	Hostname string
	CPU      float64
	Memory   float64
	TS       time.Time
}

func (s *Store) Recent(ctx context.Context, hostname string, limit int32) ([]MetricRow, error) {
	rows, err := s.Pool.Query(ctx, `
SELECT hostname, cpu, memory, ts
FROM metrics
WHERE ($1 = '' OR hostname = $1)
ORDER BY ts DESC
LIMIT $2
`, hostname, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MetricRow
	for rows.Next() {
		var r MetricRow
		if err := rows.Scan(&r.Hostname, &r.CPU, &r.Memory, &r.TS); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
