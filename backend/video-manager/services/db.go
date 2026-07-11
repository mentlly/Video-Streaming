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

	var err error
	dbpool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Database connected ....")

	// //Creation of channel table
	// commandTag, err := dbpool.Exec(
	// 	ctx,
	// 	`CREATE TABLE IF NOT EXISTS channels (
	// 	account_id serial NOT NULL,
	// 	channel_id varchar(10) PRIMARY KEY,
	// 	name varchar(25) NOT NULL,
	// 	Bio varchar(500),
	// 	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	// 	FOREIGN KEY (account_id) REFERENCES users(account_id));`,
	// )
	// if err != nil {
	// 	log.Printf("Failed to create table: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Table status: %s\n", commandTag.String())

	commandTag, err := dbpool.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS videos 
		(video_id varchar(10) PRIMARY KEY, 
		title varchar(100) NOT NULL, 
		description varchar(5000),
		duration INT NOT NULL, 
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		account_id SERIAL NOT NULL,
		FOREIGN KEY (account_id) REFERENCES users(account_id));`,
	)
	if err != nil {
		log.Printf("Failed to create table: %v\n", err)
		return
	}
	fmt.Printf("Table status: %s\n", commandTag.String())
}

func UploadVideoDb(account_id int, video_id string, title string, description string, duration int) {
	ctx := context.Background()

	// //Checking if that channel and account exists
	// var exists bool
	// err := dbpool.QueryRow(
	// 	ctx,
	// 	`SELECT * FROM channel
	// 	WHERE channel_id = ${1}
	// 	AND account_id = ${2}`,
	// 	channel_id,
	// 	account_id,
	// ).Scan(&exists)
	// if err != nil && !exists {
	// 	log.Printf("%v\n", err)
	// 	return
	// }

	_, err := dbpool.Exec(
		ctx,
		`INSERT INTO VIDEOS 
		(video_id, title, description, duration, account_id)
		VALUES
		($1, $2, $3, $4, $5)`,
		video_id,
		title,
		description,
		duration,
		account_id,
	)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}
