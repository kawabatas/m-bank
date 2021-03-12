package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kawabatas/m-bank/infra/database"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	databases := []string{dbname, database.TestDBName(dbname)}
	for _, dbname := range databases {
		if err := createDB(host, user, password, dbname); err != nil {
			log.Fatal(err)
		}

		dsn := database.DSN(host, user, password, dbname)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}

		migrations := &migrate.FileMigrationSource{
			Dir: "db/migrate",
		}
		migrate.SetTable("migrations")
		n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("applied %d migrations to %v\n", n, dbname)
	}
	os.Exit(0)
}

func createDB(host, user, password, dbname string) error {
	dsn := database.DSN(host, user, password, "mysql")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname)); err != nil {
		return err
	}

	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", dbname)); err != nil {
		return err
	}
	return nil
}
