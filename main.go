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

func fetchNextMusic(db *sql.DB, userID int) (int, error) {
	var musicID int
	err := db.QueryRow(`
		select 
			music_id
		from ratings
		where user_id != $1 and 
			music_id not in (select unnest(music_ids) from l_musics where user_id = $1) and
			user_id in (select unnest(similar_user_ids) from l_musics where user_id = $1)
		order by rating desc
		limit $1;
		`, userID).Scan(&musicID)
	return musicID, err
}

func main() {

	targetUserID := 1 // Example user ID
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
	musicID, err := fetchNextMusic(db, targetUserID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(musicID)
	// ///////////////////////////////////////////////////////////////////////////////

	fmt.Printf("Time taken: %s\n", time.Since(startTime))
}
