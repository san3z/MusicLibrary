package main

import (
	"fmt"

	b "musiclibrary/backend"
)

func main() {
	b.Handling()
	db := b.ConnDB()
	defer db.Close()
	b.RunMigrations(db)
	fmt.Println("Booting migrations")
}
