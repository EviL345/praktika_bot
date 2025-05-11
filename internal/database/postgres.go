package database

import (
	"database/sql"
	"fmt"
	"github.com/EviL345/praktika_bot/internal/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
)

type Database struct {
	Db *sql.DB
}

func New(cfg *config.Db) *Database {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if err = goose.Up(db, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return &Database{Db: db}
}

func (d *Database) Close() {
	if err := d.Db.Close(); err != nil {
		log.Fatalf("failed to close database: %v", err)
	}
}

func (d *Database) GetTopicId(userId int64) int {
	log.SetPrefix("GetTopicId")
	query := "SELECT topic_id FROM users WHERE id = $1"
	row := d.Db.QueryRow(query, userId)
	var topicId int
	if err := row.Scan(&topicId); err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		log.Printf("failed to scan topic ID: %v", err)

		return 0
	}

	return topicId
}

func (d *Database) CreateTopic(userId int64, topicId int) {
	log.SetPrefix("SetTopicId")
	query := "INSERT INTO users (id, topic_id) VALUES ($1, $2)"
	_, err := d.Db.Exec(query, userId, topicId)
	if err != nil {
		log.Printf("failed to update topic ID: %v", err)
	}
}
