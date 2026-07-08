package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func InitDb() {
	//Creating a connection with database
	connStr := os.Getenv("DB_URL")

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	log.Println("Database connected ....")

	//Creation of channel table
	commandTag, err := dbpool.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS channel 
		(channel_id varchar(10) PRIMARY KEY NOT NULL, 
		name varchar(25) NOT NULL, 
		Bio varchar(500),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP);`,
	)
	if err != nil {
		log.Printf("Failed to create table: %v\n", err)
		return
	}
	fmt.Printf("Table status: %s\n", commandTag.String())

	commandTag, err = dbpool.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS video 
		(video_id varchar(10) PRIMARY KEY NOT NULL, 
		title varchar(100) NOT NULL, 
		description varchar(5000),
		duration INT(4) NOT NULL, 
		channel_id varchar(10) NOT NULL,
		FOREIGN KEY (channel_id) REFERENCES channel(channel_id));`,
	)
	if err != nil {
		log.Printf("Failed to create table: %v\n", err)
		return
	}
	fmt.Printf("Table status: %s\n", commandTag.String())
}
