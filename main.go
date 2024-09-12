package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "recomment"
)

type Recommendation struct {
	MusicID int
	Score   float64
}

func main() {

	// targetUserID := 1 // Example user ID
	startTime := time.Now()
	// Set up PostgreSQL connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// ///////////////////////////////////////////////////////////////////////////////
	// musicID, err := fetchNextMusic(db, targetUserID)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(musicID)
	// ///////////////////////////////////////////////////////////////////////////////
	UpdateAllUserSimilars(db)

	fmt.Printf("Time taken: %s\n", time.Since(startTime))
}
