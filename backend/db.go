package backend

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	//"encoding/json"
)

func handle() {
	// http.HandleFunc("insert-music", insertMusicHandler)
	// http.HandleFunc("/get-music", getMusicHandler)
	//	http.HandleFunc("/json-test", testJsonHandler)
}

type Music struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

func connDB() *sql.DB {
	fmt.Println("Loading .env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file -------------------------------------------")
	}
	fmt.Println("DBHOST:", os.Getenv("DBHOST"))
	fmt.Println("DBPORT:", os.Getenv("DBPORT"))
	fmt.Println("DBNAME:", os.Getenv("DBNAME"))
	fmt.Println("DBLOGIN:", os.Getenv("DBLOGIN"))
	fmt.Println("DBPASS:", os.Getenv("DBPASS"))
	dbConf := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME"),
		os.Getenv("DBLOGIN"), os.Getenv("DBPASS"))
	fmt.Println("dbConf: ", dbConf)
	db, err := sql.Open("postgres", dbConf)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("Database pinging succesfully")
	return db
}

func insertMusic(db *sql.DB, s *Music) error {
	fmt.Println("Workin insertMusic func")
	insert := `INSERT INTO musicLibrary (group, song) VALUES ($1 ,$2, $3)`
	_, err := db.Exec(insert, s.Group, s.Song)
	return err
}

func GetMusic(db *sql.DB) ([]Music, error) {
	rows, err := db.Query(`SELECT "ID", "Group", "Song" FROM public."musicLibrary";`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var music []Music
	for rows.Next() {
		var s Music
		err := rows.Scan(&s.Group, &s.Song)
		if err != nil {
			return nil, err
		}
		music = append(music, s)
	}
	return music, nil
}
