package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type League struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrLeagueNotFound = errors.New("league not found")

func (db *DB) GetLeagues() ([]*League, error) {
	stmt := `
	SELECT
		id, title
	FROM league
	LIMIT 20;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	leagues := []*League{}

	for rows.Next() {
		p := &League{}

		err = rows.Scan(&p.ID, &p.Title)
		if err != nil {
			return nil, err
		}
		leagues = append(leagues, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return leagues, nil
}

func (db *DB) GetLeague(id int64) (*League, error) {
	stmt := `
	SELECT
		id, title
	FROM league
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, stmt, id)

	l := &League{}
	err := row.Scan(&l.ID, &l.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLeagueNotFound
		}
		return nil, err
	}

	return l, nil
}

func (db *DB) InsertLeague(l *League) error {
	stmt := `
	INSERT INTO league (
		title
	)
	VALUES ($1)
	RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	arg := []interface{}{l.Title}
	err := db.QueryRowContext(ctx, stmt, arg...).Scan(&l.ID)

	return err
}

func (db *DB) UpdateLeague(l *League) error {
	stmt := `
	UPDATE league SET
		title = $1
	WHERE id = $2
	RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	arg := []interface{}{l.Title, l.ID}

	err := db.QueryRowContext(ctx, stmt, arg...).Scan(&l.ID)

	return err
}

func (db *DB) DeleteLeague(id int64) error {
	stmt := `
	DELETE FROM league
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, stmt, id)

	return err
}
