package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBConnection() *pgxpool.Pool {
	dbHost := os.Getenv("DBHOST")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbPort := os.Getenv("DBPORT")
	dbName := os.Getenv("DBNAME")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	if dbURL == "" {
		log.Fatal("Database URL is not set correctly")
	}

	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 5; i++ {
        pool, err = pgxpool.New(context.Background(), dbURL)
        if err == nil {
            if err = pool.Ping(context.Background()); err == nil {
                break
            }
        }
        log.Printf("Retrying database connection... (%d/5)\n", i)
        time.Sleep(2 * time.Second)
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
