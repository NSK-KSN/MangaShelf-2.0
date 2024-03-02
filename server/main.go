package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type Data struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	CoverImage string  `json:"cover_image"`
	Type       string  `json:"type"`
	Publisher  string  `json:"publisher"`
	Mal_Id     int     `json:"mal_id"`
	Score      float64 `json:"score"`
	Popularity int     `json:"popularity"`
}

func main() {
	// PostgreSQL connection string
	connStr := "user=username dbname=dbname password=pass sslmode=disable"

	// Connecting to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Setting the connection encoding explicitly
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.Exec("SET NAMES 'utf8mb4'") // Specifying the appropriate encoding

	// Testing the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the PostgreSQL database")

	// Creating a new CORS middleware with default options
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},  // Allow only localhost:3000
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Allow only specified methods
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Creating a new HTTP handler with the CORS middleware
	handler := c.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Fetching data from the database
		rows, err := db.Query("SELECT * FROM manga")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Creating a slice to hold the data
		var data []Data

		// Iterating through the rows and appending data to the slice
		for rows.Next() {
			var d Data
			if err := rows.Scan(&d.ID, &d.Title, &d.CoverImage, &d.Type, &d.Publisher, &d.Mal_Id, &d.Score, &d.Popularity); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data = append(data, d)
		}

		fmt.Println(data)

		// Encoding the data as JSON and sending it in the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}))

	// Creating a new Gorilla Mux router
	r := mux.NewRouter()

	// Handler for /other_query/{mal_id} route
	r.HandleFunc("/other_query/{mal_id}", func(w http.ResponseWriter, r *http.Request) {
		// Extracting mal_id from the request URL
		vars := mux.Vars(r)
		malID := vars["mal_id"]

		// Query for the database using mal_id
		rows, err := db.Query("SELECT * FROM manga WHERE mal_id = $1", malID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Creating a slice to hold the data for the other query
		var otherData []Data

		// Iterating through the rows and appending data to the slice
		for rows.Next() {
			var o Data
			if err := rows.Scan(&o.ID, &o.Title, &o.CoverImage, &o.Type, &o.Publisher, &o.Mal_Id, &o.Score, &o.Popularity); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			otherData = append(otherData, o)
		}

		// Encoding the data as JSON and sending it in the response for the other query
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(otherData)
	}).Methods("GET")

	// Applying CORS middleware to the router
	handlerr := c.Handler(r)

	http.Handle("/", handlerr)

	// CORS-enabled handler with the default ServeMux
	http.Handle("/data", handler)

	// Serve images
	http.Handle("/covers/", http.StripPrefix("/covers/", http.FileServer(http.Dir("path/to/the/folder"))))

	// Start the HTTP server
	http.ListenAndServe(":8080", nil)
}
