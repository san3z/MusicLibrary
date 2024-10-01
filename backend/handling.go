package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func Handling() {
	db := ConnDB()
	defer db.Close()
	err := RunMigrations(db)
	if err != nil {
		fmt.Println("Problem with migrations: ", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/insert-music", insertMusicHandler)
	r.HandleFunc("/get-music", getMusicHandler)
	r.HandleFunc("/json-test", testJsonHandler)
	r.HandleFunc("/del-music", delMusicHandler)
	r.HandleFunc("/upd-music", updateMusicHandler)
	fmt.Println("Server started")
	http.ListenAndServe(":9090", r)
}

func renderHTML(w http.ResponseWriter, r *http.Request, filename string) {
	tmpl, err := template.ParseFiles("frontend/pages/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, r)
}

func home(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, r, "home.html")
}

func getMusicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db := ConnDB()
	defer db.Close()

	var music []Music
	rows, err := db.Query(`SELECT "band", "song" FROM public."musicLibrary"`)
	fmt.Println(rows)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Error Get Music: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var m Music
		err = rows.Scan(&m.Band, &m.Song)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		music = append(music, m)
	}
	fmt.Printf("Music: %+v\n", music)
	json.NewEncoder(w).Encode(music)
}

func insertMusicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db := ConnDB()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	music := Music{
		Band: r.FormValue("music-band"),
		Song: r.FormValue("music-song"),
	}
	fmt.Println("Music", music.Band, music.Song)
	_, err = db.Exec((`INSERT INTO public."musicLibrary"("band", "song") VALUES($1, $2)`), music.Band, music.Song)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Music inserted successfully")
}

func testJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Println("Received request at /json-test")
	fmt.Println("Request Method:", r.Method)
	fmt.Println("Request Headers:", r.Header)
	if r.Method == "GET" {
		fmt.Fprint(w, "GET request received")
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db := ConnDB()
	defer db.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var jsonInput map[string]interface{}
	err = json.Unmarshal(body, &jsonInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Полученный JSON: %+v\n", jsonInput)

	music := Music{
		Band: jsonInput["band"].(string),
		Song: jsonInput["song"].(string),
	}

	fmt.Println("INSERT to DB")
	_, err = db.Exec((`INSERT INTO public."musicLibrary"("band", "song") VALUES($1, $2)`), music.Band, music.Song)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "JSON успешно обработан и записан в базу данных")
}

func delMusicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delMusicHandler called")
	if r.Method != "POST" {
		fmt.Fprintf(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	db := ConnDB()
	defer db.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var jsonInput map[string]interface{}
	err = json.Unmarshal(body, &jsonInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Given JSON:", jsonInput)
	music := Music{
		Band: jsonInput["band"].(string),
		Song: jsonInput["song"].(string),
	}
	fmt.Println("JSON Parsed:", music.Band, music.Song)
	cmd, err := db.Exec(`DELETE FROM public."musicLibrary" WHERE "band" = $1 AND "song" = $2`, music.Band, music.Song)
	fmt.Println(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Music deleted successfully")
}

func updateMusicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updateMusicHandler called")
	if r.Method != "POST" {
		fmt.Fprintf(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	db := ConnDB()
	defer db.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var jsonInput map[string]interface{}
	err = json.Unmarshal(body, &jsonInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Given JSON:", jsonInput)

	keys := []string{"bandNew", "songNew", "band", "song"}
	musicNew := MusicNew{}

	for _, key := range keys {
		value, ok := jsonInput[key].(string)
		if !ok || value == "" {
			http.Error(w, fmt.Sprintf("%s is missing or not a string", key), http.StatusBadRequest)
			return
		}

		switch key {
		case "bandNew":
			musicNew.BandNew = value
		case "songNew":
			musicNew.SongNew = value
		case "band":
			musicNew.Band = value
		case "song":
			musicNew.Song = value
		}
	}
	fmt.Println("JSON Parsed:", musicNew.Band, musicNew.Song)
	_, err = db.Exec(`UPDATE public."musicLibrary" SET "band" = $1, "song" = $2 WHERE "band" = $3 AND "song" = $4; `, musicNew.BandNew, musicNew.SongNew, musicNew.Band, musicNew.Song)
	if err != nil {
		http.Error(w, "Problem with update", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Music Data changed")
}
