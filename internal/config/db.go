package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBConnection() *pgxpool.Pool {
	dburl := os.Getenv("DBURL")
	if dburl == "" {
		log.Fatal("DBURL is not set in environment variables")
	}

	pool, err := pgxpool.New(context.Background(), dburl)
	if err != nil {
		log.Fatal("failed to create connection pool: ", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("failed to ping database: ", err)
	}

	tableQuery, err := os.ReadFile("create_table.sql")
	if err != nil {
		log.Fatal("failed to open schema file: ", err)
	}

	_, err = pool.Exec(context.Background(), string(tableQuery))
	if err != nil {
		log.Fatal("failed to execute schema query: ", err)
	}

	log.Println("database connected successfully")
	return pool
}
