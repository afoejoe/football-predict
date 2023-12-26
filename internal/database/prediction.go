package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Prediction struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Keywords       string    `json:"keywords"`
	Body           string    `json:"body"`
	Odds           float64   `json:"odds"`
	PredictionType string    `json:"prediction_type"`
	ScheduledAt    time.Time `json:"scheduled_at"`
	IsFeatured     bool      `json:"is_featured"`
	IsArchived     bool      `json:"is_archived"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

var ErrPredictionNotFound = errors.New("prediction not found")

func (db *DB) GetPredictions() ([]*Prediction, error) {
	stmt := `
	SELECT id, title, slug, created_at, scheduled_at, odds, prediction_type
	FROM prediction
	WHERE is_archived = false
	ORDER BY scheduled_at, created_at;`

	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	predictions := []*Prediction{}

	for rows.Next() {
		p := &Prediction{}

		err = rows.Scan(&p.ID, &p.Title, &p.Slug, &p.CreatedAt, &p.ScheduledAt, &p.Odds, &p.PredictionType)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return predictions, nil
}

func (db *DB) GetFeatured() ([]*Prediction, error) {
	stmt := `
	SELECT id, title, slug, created_at, scheduled_at, odds, prediction_type
	FROM prediction
	WHERE is_featured = true
	AND is_archived = false
	ORDER BY scheduled_at, created_at;`

	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	predictions := []*Prediction{}

	for rows.Next() {
		p := &Prediction{}

		err = rows.Scan(&p.ID, &p.Title, &p.Slug, &p.CreatedAt, &p.ScheduledAt, &p.Odds, &p.PredictionType)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return predictions, nil
}

func (db *DB) GetPredictionBySlug(slug string) (*Prediction, error) {
	stmt := `
	SELECT id, title, slug, body, created_at, scheduled_at, odds, prediction_type
	FROM prediction
	WHERE slug = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, stmt, slug)

	p := &Prediction{}

	err := row.Scan(&p.ID, &p.Title, &p.Slug, &p.Body, &p.CreatedAt, &p.ScheduledAt, &p.Odds, &p.PredictionType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPredictionNotFound
		}
		return nil, err
	}

	return p, nil
}