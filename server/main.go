package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"os"
)

type Data struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	CoverImage  string  `json:"cover_image"`
	TypeID      int     `json:"type_id"`
	PublisherID int     `json:"publisher_id"`
	Mal_Id      int     `json:"mal_id"`
	Score       float64 `json:"score"`
	Popularity  int     `json:"popularity"`
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
			if err := rows.Scan(&d.ID, &d.Title, &d.CoverImage, &d.Mal_Id, &d.Score, &d.Popularity, &d.PublisherID, &d.TypeID); err != nil {
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

	// Handler for /other_query/{id} route
	r.HandleFunc("/other_query/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Extracting id from the request URL
		vars := mux.Vars(r)
		ID := vars["id"]

		// Query for the database using id
		rows, err := db.Query("SELECT * FROM manga WHERE id = $1", ID)
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
			if err := rows.Scan(&o.ID, &o.Title, &o.CoverImage, &o.Mal_Id, &o.Score, &o.Popularity, &o.PublisherID, &o.TypeID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			otherData = append(otherData, o)
		}

		// Encoding the data as JSON and sending it in the response for the other query
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(otherData)
	}).Methods("GET")

	r.HandleFunc("/add-data", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		decoder := json.NewDecoder(r.Body)
		var data Data
		if err := decoder.Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Insert data into the PostgreSQL database
		_, err := db.Exec("INSERT INTO manga (title, cover_image, type_id, publisher_id, mal_id, score, popularity) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			data.Title, data.CoverImage, data.TypeID, data.PublisherID, data.Mal_Id, data.Score, data.Popularity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	r.HandleFunc("/dbquery/{tableName}", func(w http.ResponseWriter, r *http.Request) {
		// Extract the table name from the request URL
		vars := mux.Vars(r)
		tableName := vars["tableName"]

		// Query the database to fetch data from the specified table
		rows, err := db.Query("SELECT id, name FROM " + tableName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Create a slice to hold the data
		var data []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		// Iterate through the rows and append data to the slice
		for rows.Next() {
			var item struct {
				ID   int
				Name string
			}
			if err := rows.Scan(&item.ID, &item.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data = append(data, struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{ID: item.ID, Name: item.Name})
		}

		// Encode the data as JSON and send it in the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	r.HandleFunc("/dbquery/{tableName}/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Extracting id from the request URL
		vars := mux.Vars(r)
		rowID := vars["id"]
		tableName := vars["tableName"]

		// Query the database using items ID
		row := db.QueryRow("SELECT id, name FROM "+tableName+" WHERE id = $1", rowID)

		// Creating a struct to hold the item data
		var item struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		// Scan the row into the struct
		if err := row.Scan(&item.ID, &item.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encoding the data as JSON and sending it in the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	}).Methods("GET")

	r.HandleFunc("/fetch-mal-data/{malId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		malId := vars["malId"]

		// Construct the URL for the MyAnimeList API
		apiUrl := fmt.Sprintf("https://api.myanimelist.net/v2/manga/%s?fields=num_list_users,mean,media_type", malId)

		// Create a new request to the MyAnimeList API
		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set headers for the MyAnimeList API request
		req.Header.Set("X-MAL-CLIENT-ID", "client_id") // Replace with your actual Client ID

		// Create a new HTTP client
		client := &http.Client{}

		// Send the request to the MyAnimeList API
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		fmt.Println(resp)
		// Copy the response from the MyAnimeList API to the response writer
		io.Copy(w, resp.Body)
	})

	// Applying CORS middleware to the router
	handlerr := c.Handler(r)

	// Handle upload requests
	uploadHandler := c.Handler(http.HandlerFunc(uploadHandler))

	http.Handle("/upload", uploadHandler)

	http.Handle("/", handlerr)

	// CORS-enabled handler with the default ServeMux
	http.Handle("/data", handler)

	// Serve images
	http.Handle("/covers/", http.StripPrefix("/covers/", http.FileServer(http.Dir("path/to/the/folder"))))

	// Start the HTTP server
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB maximum
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the file on the server
	f, err := os.Create("./covers/" + handler.Filename)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the file contents to the server file
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Write([]byte("File uploaded successfully"))
}
