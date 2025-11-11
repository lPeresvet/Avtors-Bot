package repository

import (
	"avtor.ru/bot/analyse_service/internal/model"
	"avtor.ru/bot/server"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var createTabeCommand = `CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,
    zone_id VARCHAR NOT NULL,
    user_id BIGINT NOT NULL
);`

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotExists     = errors.New("not exists")
)

func initDB(config *model.DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open sql connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL!")

	_, err = db.Exec(createTabeCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}

type Repository struct {
	db *sql.DB
}

func NewRepository(config *model.DBConfig) (*Repository, error) {
	db, err := initDB(config)
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) InsertLike(like model.Like) error {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE zone_id = $1 AND user_id = $2)", like.ZoneID, like.UserID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrAlreadyExists
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO likes (zone_id, user_id) 
        VALUES ($1, $2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(like.ZoneID, like.UserID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) DeleteLike(like model.Like) error {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE zone_id = $1 AND user_id = $2)", like.ZoneID, like.UserID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNotExists
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        DELETE FROM likes
        WHERE zone_id = $1 AND user_id = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(like.ZoneID, like.UserID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetLikes(userID string) (*server.Zones, error) {
	rows, err := r.db.Query("SELECT zone_id FROM likes WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get likes: %w", err)
	}

	defer rows.Close()

	zones := make(server.Zones, 0)

	for rows.Next() {
		var zone server.Zone
		err := rows.Scan(
			&zone.Id,
		)
		if err != nil {
			return nil, err
		}
		zones = append(zones, zone)
	}

	return &zones, nil
}
