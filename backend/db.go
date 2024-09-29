package backend

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Music struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type MusicNew struct {
	GroupNew string `json:"groupNew"`
	SongNew  string `json:"songNew"`
	Group    string `json:"group"`
	Song     string `json:"song"`
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
	// Получаем абсолютный путь к директории с миграциями
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(currentDir)
	migrationsDir := filepath.Join(currentDir, "/", "migrations")

	// Проверяем, что директория существует
	_, err = os.Stat(migrationsDir)
	if err != nil {
		return fmt.Errorf("migration directory not found: %s", migrationsDir)
	}
	fmt.Println("makin migration source")
	// Создаем источник миграций из файловой системы
	source, err := iofs.New(os.DirFS("musicLibrary"), migrationsDir)
	if err != nil {
		return err
	}
	fmt.Println("creating migrator")
	// Создаем мигратор
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}
	fmt.Println("STARTING migrations")
	// Выполняем миграции
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Println("Migrations applied successfully")

	return nil
}

/*func insertMusic(db *sql.DB, s *Music) error {
	fmt.Println("Workin insertMusic func")
	insert := `INSERT INTO musicLibrary (group, song) VALUES ($1 ,$2, $3)`
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
