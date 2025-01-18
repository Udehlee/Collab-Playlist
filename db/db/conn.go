package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

type PgConn struct {
	Conn *sql.DB
}

func ConnectDB(config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		config.Username, config.Password, config.Host, config.Port, config.DbName)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error creating database connection: %w", err)
	}

	if err := runMigrations(dbConn); err != nil {
		log.Fatal("migration unsuccessful")
	}

	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(10)

	if err := dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("connection not alive: %w", err)
	}

	fmt.Println("connected successfully")
	return dbConn, nil
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}

func InitPG(cfg Config) (*PgConn, error) {
	pg := PgConn{}

	db, err := ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	pg.Conn = db
	return &pg, nil
}
