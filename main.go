package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kawabatas/m-bank/infra/database"
)

func main() {
	_, err := setupDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		log.Fatalf("setup DB error: %v", err)
	}

	server, err := newServer()
	if err != nil {
		log.Fatalf("new Server error: %v", err)
	}
	defer func() {
		_ = server.Shutdown()
	}()

	server.Port = 3000
	if err := server.Serve(); err != nil {
		log.Fatalf("serve Server error: %v", err)
	}
}

func setupDB(dbHost, dbName, dbUser, dbPassword string) (*sql.DB, error) {
	dsn := database.DSN(dbHost, dbUser, dbPassword, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
