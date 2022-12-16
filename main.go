package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Score string `json:"score"`
	Date  string `json:"date"`
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func appleFaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "apple-touch-icon.png")
}

func formBody(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusCreated)

	//Populate template
	var body bytes.Buffer
	t, _ := template.ParseFiles("add.html")
	err := t.Execute(&body, nil)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(body.Bytes())
	if err != nil {
		log.Fatal(err)
	}

}

func movieBody(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusCreated)

	db, err := sql.Open("sqlite3", "reviews.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row, err := db.Query("SELECT * FROM movies ORDER BY ID DESC")
	if err != nil {
		log.Fatal(err)
	}
	var movies = new([]Entry)
	for row.Next() {
		var id int
		var title string
		var score string
		var date string

		err := row.Scan(&id, &title, &score, &date)
		if err != nil {
			log.Fatal(err)
		}
		*movies = append(*movies, Entry{
			Title: title,
			Score: score,
			Date:  date,
		})
	}

	jsonMovies, err := json.Marshal(&movies)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(jsonMovies)
	if err != nil {
		log.Fatal(err)
	}
}

func showBody(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusCreated)

	db, err := sql.Open("sqlite3", "reviews.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row, err := db.Query("SELECT * FROM shows ORDER BY ID DESC")
	if err != nil {
		log.Fatal(err)
	}
	var shows = new([]Entry)

	for row.Next() {
		var id int
		var title string
		var score string
		var date string

		err := row.Scan(&id, &title, &score, &date)
		if err != nil {
			log.Fatal(err)
		}
		*shows = append(*shows, Entry{
			Title: title,
			Score: score,
			Date:  date,
		})
	}
	jsonShows, err := json.Marshal(&shows)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(jsonShows)
	if err != nil {
		log.Fatal(err)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusCreated)

	//Populate template
	var body bytes.Buffer
	t, _ := template.ParseFiles("template.html")
	err := t.Execute(&body, nil)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(body.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

// Handle REST post
// Test with curl --header "Content-Type: application/json" --request POST --data '{"firstName":"Tom","lastName":"Rat","email":"tom@email.com","password":"fake"}' http://localhost:9091/
func post(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	// This works for standard (non-json) post
	//w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	//w.WriteHeader(http.StatusCreated)
	//r.ParseForm()
	//fmt.Println(r)
	//for key, value := range r.Form {
	//	fmt.Printf("%s = %s\n", key, value)
	//}
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var entry Entry
	err := decoder.Decode(&entry)
	if err != nil {
		log.Fatal(err)
	}
	postBody, _ := json.Marshal(entry)
	// print decoded json (raw []bytes if no bytes.NewBuffer)
	pb := bytes.NewBuffer(postBody)
	fmt.Println(pb)

	db, err := sql.Open("sqlite3", "reviews.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var insertSQL string

	if entry.Type == "1" {
		insertSQL = `INSERT INTO movies(Title, Score, Entry) VALUES (?, ?, ?)`
	} else {
		insertSQL = `INSERT INTO shows(Title, Score, Entry) VALUES (?, ?, ?)`
	}

	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = statement.Exec(entry.Title, entry.Score, entry.Date)
	if err != nil {
		log.Fatal(err)
	}

	//respCode, _ := json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	//_, err = w.Write([]byte(respCode))
	_, err = w.Write([]byte(`{"success": true}`))
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/tyreviews").Subrouter()
	api.HandleFunc("/favicon.ico", faviconHandler).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/apple-touch-icon.png", appleFaviconHandler).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/", get).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/formBody", formBody).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/movieBody", movieBody).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/showBody", showBody).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/post", post).Methods(http.MethodPost, http.MethodOptions)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Origin", "Accept", "token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(":9091", handlers.LoggingHandler(os.Stdout,
		handlers.CORS(originsOk, headersOk, methodsOk)(r))))
}
