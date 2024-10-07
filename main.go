package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	apiPathConst = "/apis/v1/books"
)

type book struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Isbn string `json:"isbn"`
}
type library struct {
	dbHost string
	dbPass string
	dbName string
	dbPort string
	dbUser string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "your_username"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "your_password"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = apiPathConst
	}

	l := library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
		dbPort: dbPort,
		dbUser: dbUser,
	}

	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodGet)
	r.HandleFunc(apiPath, l.postBooks).Methods(http.MethodPost)
	err := http.ListenAndServe(":8082", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}

func (l library) getBooks(w http.ResponseWriter, r *http.Request) {
	//open db connection
	db := l.openConnection()

	//query db
	rows, err := db.Query("SELECT * FROM book")
	if err != nil {
		log.Fatalf("Error querying db: %v", err)
	}
	var books []book
	for rows.Next() {
		var b book
		err := rows.Scan(&b.ID, &b.Name, &b.Isbn)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		books = append(books, b)
	}
	json.NewEncoder(w).Encode(books)
	//close db connection
	l.closeConnection(db)
}

func (l library) postBooks(w http.ResponseWriter, r *http.Request) {
	//open db connection
	db := l.openConnection()

	var b book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Fatalf("Error decoding request body: %v", err)
	}
	_, err = db.Exec("INSERT INTO book (id, name, isbn) VALUES ($1, $2,$3)", b.ID, b.Name, b.Isbn)
	if err != nil {
		log.Fatalf("Error inserting book: %v", err)
	}
	w.WriteHeader(http.StatusCreated)
	//close db connection
	l.closeConnection(db)
}

func (l library) openConnection() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		l.dbHost, l.dbPort, l.dbUser, l.dbPass, l.dbName)

	fmt.Printf("Connecting to db with connection string: %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	return db
}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("Error closing db connection: %v", err)
	}

}
