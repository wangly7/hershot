package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	List(ctx context.Context) ([]Team, error)

	GetById(ctx context.Context, id string) (*Team, error)

	Create(ctx context.Context, team *Team) error

	Update(ctx context.Context, team *Team) error

	Delete(ctx context.Context, id string) error
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
	const query = `
	SELECT
		id::text,
		name,
		city,
		abbreviation
	FROM teams
	ORDER BY city, name
	`

	row, err := r.db.Query(ctx, query)
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

func (r *PostgresRepository) GetById(ctx context.Context, id string) (*Team, error) {
	const query = `
		SELECT
			id::text,
			name,
			city,
			abbreviation
		FROM teams
		WHERE id = $1
	`

	var team Team

	err := r.db.QueryRow(ctx, query, id).Scan(
		&team.ID,
		&team.Name,
		&team.City,
		&team.Abbreviation,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTeamNotFound
		}
		return nil, fmt.Errorf("query team by id: %w", err)
	}
	return &team, nil
}

func (r *PostgresRepository) Create(ctx context.Context, team *Team) error {
	query := `
	INSERT INTO teams (
		id,
		name,
		city,
		abbreviation
	)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		team.ID,
		team.Name,
		team.City,
		team.Abbreviation,
	)
	if err != nil {
		return fmt.Errorf("insert team: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, team *Team) error {
	const query = `
	UPDATE teams
	SET
		name = $2,
		city = $3,
		abbreviation = $4
	WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		team.ID,
		team.Name,
		team.City,
		team.Abbreviation,
	)
	if err != nil {
		return fmt.Errorf("update team: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrTeamNotFound
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	const query = `
	Delete FROM teams
	WHERE id = $1	
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete team: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrTeamNotFound
	}
	return nil
}
