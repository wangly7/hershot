package team

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	List(ctx context.Context) ([]Team, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) List(ctx context.Context) ([]Team, error) {
	sql := `
	SELECT
		id::text,
		name,
		city,
		abbreviation
	FROM teams
	ORDER BY city, name
	`

	row, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("query teams: %w", err)
	}
	defer row.Close()

	teams := make([]Team, 0)

	for row.Next() {
		var t Team

		if err := row.Scan(
			&t.ID,
			&t.Name,
			&t.City,
			&t.Abbreviation,
		); err != nil {
			return nil, fmt.Errorf("scan team: %w", err)
		}

		teams = append(teams, t)
	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("iterate teams: %w", err)
	}

	return teams, nil
}
