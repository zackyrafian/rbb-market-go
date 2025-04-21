package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func Connect() {
	var err error
	dsn := "postgresql://postgres:root@localhost:5432/rbbmarket"
	DB, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println("Connected to database ✅")
}

func Close() {
	err := DB.Close(context.Background())
	if err != nil {
		log.Fatalf("Error closing connection: %v\n", err)
	}
}
