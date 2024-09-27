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
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/insert-music", insertMusicHandler)
	r.HandleFunc("/get-music", getMusicHandler)
	r.HandleFunc("/json-test", testJsonHandler)
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

	db := connDB()
	defer db.Close()

	var music []Music
	rows, err := db.Query(`SELECT "Group", "Song" FROM public."musicLibrary"`)
	fmt.Println(rows)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var m Music
		err = rows.Scan(&m.Group, &m.Song)
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

	db := connDB()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	music := Music{
		Group: r.FormValue("music-group"),
		Song:  r.FormValue("music-song"),
	}

	_, err = db.Exec((`INSERT INTO public."musicLibrary"("Group", "Song") VALUES($1, $2)`), music.Group, music.Song)
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

	db := connDB()
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
		Group: jsonInput["group"].(string),
		Song:  jsonInput["song"].(string),
	}

	fmt.Println("INSERT to DB")
	_, err = db.Exec((`INSERT INTO public."musicLibrary"("Group", "Song") VALUES($1, $2)`), music.Group, music.Song)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "JSON успешно обработан и записан в базу данных")
}
