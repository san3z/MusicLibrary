package backend

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type Music struct {
	Band string `json:"band"`
	Song string `json:"song"`
}

type MusicNew struct {
	BandNew string `json:"bandNew"`
	SongNew string `json:"songNew"`
	Band    string `json:"band"`
	Song    string `json:"song"`
}

func ConnDB() *sql.DB {
	fmt.Println("Loading .env")
	err := godotenv.Load()
	dbConf := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("DBHOST"), os.Getenv("DBPORT"),
		os.Getenv("DBNAME"), os.Getenv("DBLOGIN"),
		os.Getenv("DBPASS"))
	fmt.Println(dbConf)
	if err != nil {
		log.Fatal("Error Loading .env file -------------------------------------------")
	}
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

func RunMigrations(db *sql.DB) error {
	fmt.Println("Running migrations")
	err := goose.Up(db, "migrations/")
	if err != nil {
		fmt.Printf("Error running migrations: %s\n ", err)
	}
	fmt.Println("Migrations completed successfully")
	return nil
}

/*func insertMusic(db *sql.DB, s *Music) error {
	fmt.Println("Workin insertMusic func")
	insert := `INSERT INTO musicLibrary (band, song) VALUES ($1 ,$2, $3)`
	_, err := db.Exec(insert, s.Group, s.Song)
	return err
}*/

/*func GetMusic(db *sql.DB) ([]Music, error) {
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
}*/
