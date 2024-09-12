package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"sort"
)

// Compute cosine similarity between two users
func cosineSimilarity(userRatings1, userRatings2 map[int]float64) float64 {
	var dotProduct, norm1, norm2 float64

	// Calculate dot product and norms
	for musicID, rating1 := range userRatings1 {
		norm1 += rating1 * rating1 // norm1 for all userRatings1

		// If the second user has rated the same music, calculate dot product
		if rating2, exists := userRatings2[musicID]; exists {
			dotProduct += rating1 * rating2
		}
	}

	// Calculate norm2 for all userRatings2
	for _, rating2 := range userRatings2 {
		norm2 += rating2 * rating2
	}

	// If either norm is 0, it means at least one of the users has no ratings, so return 0 similarity
	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	// Return cosine similarity
	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

func UpdateAllUserSimilars(db *sql.DB) {
	// var userID = 2
	ratings, _ := fetchRatings(db)

	for userID := range ratings {
		fmt.Println("start for this user: ", userID)
		users := generateRecommendations(ratings, userID)
		fmt.Println("generated recom users for this user: ", userID)
		insertSimilarUsers(db, userID, users)

	}

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
		limit 1;
		`, userID).Scan(&musicID)
	return musicID, err
}

// Generate recommendations for a user based on similar users
func generateRecommendations(ratings map[int]map[int]float64, targetUserID int) []int {
	similarities := make(map[int]float64)
	targetRatings := ratings[targetUserID]

	for userID, userRatings := range ratings {

		if userID == targetUserID {
			continue
		}
		similarity := cosineSimilarity(targetRatings, userRatings)
		similarities[userID] = similarity

	}

	return getTopNSimilarUsers(similarities, 100)
}

func insertSimilarUsers(db *sql.DB, targetUser int, users []int) {

	for _, userID := range users {
		_, err := db.Exec(
			`
			
    INSERT INTO l_musics (user_id, similar_user_ids)
        VALUES ($1, ARRAY[$2]::int[])  -- Insert a new row with the user_id and music_id
        ON CONFLICT (user_id)   -- If the user_id already exists, update the similar_user_ids array
    DO UPDATE
        SET similar_user_ids = CASE
            WHEN array_position(l_musics.similar_user_ids, $2) IS NULL THEN array_append(l_musics.similar_user_ids, $2)
            ELSE l_musics.similar_user_ids
        END;
			`,
			targetUser, userID,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getTopNSimilarUsers(similarities map[int]float64, n int) []int {
	// Create a slice to hold user IDs and their similarity scores
	type userSimilarity struct {
		userID     int
		similarity float64
	}

	var userSimList []userSimilarity

	// Fill the slice with data from the map
	for userID, similarity := range similarities {
		userSimList = append(userSimList, userSimilarity{
			userID:     userID,
			similarity: similarity,
		})
	}

	// Sort the slice by similarity in descending order
	sort.Slice(userSimList, func(i, j int) bool {
		return userSimList[i].similarity > userSimList[j].similarity
	})

	// Get the top N users
	var topNUsers []int

	for i := 0; i < n && i < len(userSimList); i++ {

		topNUsers = append(topNUsers, userSimList[i].userID)

	}

	return topNUsers
}

// Fetch ratings from the database
func fetchRatings(db *sql.DB) (map[int]map[int]float64, error) {
	rows, err := db.Query(`
		
		SELECT user_id, music_id, rating 
		FROM ratings
	
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ratings := make(map[int]map[int]float64)
	for rows.Next() {
		var userID, musicID int
		var rating float64
		if err := rows.Scan(&userID, &musicID, &rating); err != nil {
			return nil, err
		}
		if _, exists := ratings[userID]; !exists {
			ratings[userID] = make(map[int]float64)
		}
		ratings[userID][musicID] = rating
	}
	return ratings, nil
}

// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"

// 	_ "github.com/lib/pq"
// )

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "12345"
// 	dbname   = "recomment"
// )

// // Generate random username
// func randomUsername(n int) string {
// 	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	rand.Seed(time.Now().UnixNano())

// 	username := make([]byte, n)
// 	for i := range username {
// 		username[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(username)
// }

// func main() {
// 	// Set up PostgreSQL connection string
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)

// 	// Open database connection
// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Test the connection
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("Successfully connected to PostgreSQL")

// 	// Insert 100,000 random users
// 	batchSize := 1000
// 	userCount := 100000

// 	// Start the insertion process
// 	for i := 0; i < userCount; i += batchSize {
// 		// Prepare the bulk insert statement
// 		query := "INSERT INTO musics (title) VALUES "

// 		// Generate batch of users
// 		values := []interface{}{}
// 		for j := 0; j < batchSize && i+j < userCount; j++ {
// 			if j > 0 {
// 				query += ","
// 			}
// 			query += fmt.Sprintf("($%d)", j+1)
// 			values = append(values, randomUsername(10)) // Generating 10-character usernames
// 		}

// 		// Execute the bulk insert
// 		_, err := db.Exec(query, values...)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		log.Printf("Inserted %d users", i+batchSize)
// 	}

// 	log.Println("Finished inserting users")
// }

// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"

// 	_ "github.com/lib/pq"
// )

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "12345"
// 	dbname   = "recomment"
// )

// func randomRating() float64 {
// 	return float64(rand.Intn(41)+10) / 10.0 // Generates a number between 1.0 and 5.0
// }

// func randomMusicID() int {
// 	return rand.Intn(500000) + 1 // Generates a random music ID between 1 and 500,000
// }

// func randomRatingsCount() int {
// 	return rand.Intn(100) + 1 // Generates a random number of ratings between 1 and 100
// }

// func main() {
// 	// Set up PostgreSQL connection string
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)

// 	// Open database connection
// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Test the connection
// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("Successfully connected to PostgreSQL")

// 	rand.Seed(time.Now().UnixNano()) // Seed random number generator

// 	// Insert ratings for 100,000 users
// 	batchSize := 1000
// 	userCount := 100000

// 	// Start inserting ratings
// 	for userID := 1; userID <= userCount; userID++ {
// 		ratingCount := randomRatingsCount()

// 		for i := 0; i < ratingCount; i += batchSize {
// 			// Prepare the bulk insert statement
// 			query := "INSERT INTO ratings (music_id, user_id, rating) VALUES "

// 			values := []interface{}{}
// 			for j := 0; j < batchSize && i+j < ratingCount; j++ {
// 				musicID := randomMusicID()
// 				rating := randomRating()

// 				if j > 0 {
// 					query += ","
// 				}
// 				query += fmt.Sprintf("($%d, $%d, $%d)", j*3+1, j*3+2, j*3+3)

// 				values = append(values, musicID, userID, rating)
// 			}

// 			// Execute the bulk insert
// 			_, err := db.Exec(query, values...)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		}

// 		// log.Printf("Inserted ratings for user %d", userID)
// 	}

// 	log.Println("Finished inserting ratings")
// }
