package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Prediction struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Keywords       string    `json:"keywords"`
	Body           string    `json:"body"`
	Odds           float64   `json:"odds"`
	PredictionType string    `json:"prediction_type"`
	ScheduledAt    time.Time `json:"scheduled_at"`
	IsFeatured     bool      `json:"is_featured"`
	IsArchived     bool      `json:"is_archived"`
	Campaigned     bool      `json:"campaigned"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LeagueID       int64     `json:"league_id"`
	League         League    `json:"league,omitempty"`
}

var ErrPredictionNotFound = errors.New("prediction not found")

func (db *DB) GetPredictions(showArchived bool) ([]*Prediction, error) {
	stmt := `
	SELECT
		p.id, p.title, slug, p.created_at, scheduled_at, odds, prediction_type, campaigned,
		league_id, l.title
	FROM prediction p
	LEFT JOIN league l ON l.id = p.league_id
	WHERE ($1 = true OR is_archived = false)
	AND DATE(scheduled_at) >= DATE(now())
	ORDER BY l.title, scheduled_at, p.created_at
	LIMIT 30;`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, stmt, showArchived)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	predictions := []*Prediction{}

	for rows.Next() {
		p := &Prediction{}

		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.CreatedAt,
			&p.ScheduledAt,
			&p.Odds,
			&p.PredictionType,
			&p.Campaigned,
			&p.LeagueID,
			&p.League.Title,
		)

		p.League.ID = p.LeagueID
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println(predictions[0])

	return predictions, nil
}

func (db *DB) GetFeatured() ([]*Prediction, error) {
	stmt := `
	SELECT id, title, slug, created_at, scheduled_at, odds, prediction_type
	FROM prediction
	WHERE is_featured = true
	AND is_archived = false
	AND DATE(scheduled_at) >= DATE(now())
	ORDER BY scheduled_at, created_at;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, stmt)
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

func (db *DB) GetPrediction(id int64) (*Prediction, error) {
	fmt.Println("GetPrediction", id)
	stmt := `
	SELECT
		p.id, p.title, slug, body, p.created_at, scheduled_at, odds, prediction_type, is_featured, is_archived, keywords,
		l.id, l.title

	FROM prediction p
	LEFT JOIN league l ON l.id = p.league_id

	WHERE p.id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	fmt.Println(stmt, id)

	row := db.QueryRowContext(ctx, stmt, id)

	p := &Prediction{}
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Body,
		&p.CreatedAt,
		&p.ScheduledAt,
		&p.Odds,
		&p.PredictionType,
		&p.IsFeatured,
		&p.IsArchived,
		&p.Keywords,
		&p.LeagueID,
		&p.League.Title,
	)
	p.League.ID = p.LeagueID

	fmt.Println(err, p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPredictionNotFound
		}
		return nil, err
	}

	return p, nil
}

func (db *DB) GetPredictionBySlug(slug string) (*Prediction, error) {
	stmt := `
	SELECT id, title, slug, body, created_at, scheduled_at, odds, prediction_type, keywords
	FROM prediction
	WHERE slug = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, stmt, slug)

	p := &Prediction{}

	err := row.Scan(&p.ID, &p.Title, &p.Slug, &p.Body, &p.CreatedAt, &p.ScheduledAt, &p.Odds, &p.PredictionType, &p.Keywords)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPredictionNotFound
		}
		return nil, err
	}

	return p, nil
}

func (db *DB) InsertPrediction(p *Prediction) error {
	stmt := `
	INSERT INTO prediction (
		title, slug, body, created_at, scheduled_at, odds, prediction_type, is_archived, is_featured, keywords, league_id
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	arg := []interface{}{
		p.Title,
		p.Slug,
		p.Body,
		p.CreatedAt,
		p.ScheduledAt,
		p.Odds,
		p.PredictionType,
		p.IsArchived,
		p.IsFeatured,
		p.Keywords,
		p.LeagueID,
	}
	err := db.QueryRowContext(ctx, stmt, arg...).Scan(&p.ID)

	return err
}

func (db *DB) UpdatePrediction(p *Prediction) error {
	stmt := `
	UPDATE prediction SET
		title = $1,
		slug = $2,
		body = $3,
		created_at = $4,
		scheduled_at = $5,
		odds = $6,
		prediction_type = $7,
		is_archived = $8,
		is_featured = $9,
		keywords = $10,
		campaigned = $11,
		league_id = $12
	WHERE id = $13
	RETURNING id
	`

	fmt.Println(stmt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	arg := []interface{}{
		p.Title,
		p.Slug,
		p.Body,
		p.CreatedAt,
		p.ScheduledAt,
		p.Odds,
		p.PredictionType,
		p.IsArchived,
		p.IsFeatured,
		p.Keywords,
		p.Campaigned,
		p.LeagueID,
		p.ID,
	}

	err := db.QueryRowContext(ctx, stmt, arg...).Scan(&p.ID)
	fmt.Println(err)
	return err
}

func (db *DB) DeletePrediction(id int64) error {
	stmt := `
	DELETE FROM prediction
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, stmt, id)

	return err
}
