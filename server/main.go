package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/rs/cors"

	"os"

	"strconv"
)

type Data struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	CoverLink   string   `json:"cover_link"`
	TypeID      string   `json:"type_id"`
	PublisherID string   `json:"publisher_id"`
	Mal_Id      int      `json:"mal_id"`
	Score       float64  `json:"score"`
	Popularity  int      `json:"popularity"`
	StatusID    string   `json:"status_id"`
	Volumes     []Volume `json:"volumes"`
}

// Release represents a release entry in the database
type Volume struct {
	ID           int    `json:"id"`
	Title_ID     string `json:"title_id"`
	VolumeNumber int    `json:"volume_number"`
	ReleaseDate  string `json:"release_date"`
	CoverLink    string `json:"cover_link"`
	Title        string `json:"title"`
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
	// Creating a new Gorilla Mux router
	r := mux.NewRouter()

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		rows, err := db.Query(`
		SELECT m.id, m.title, v.cover_link, m.mal_id, m.score, m.popularity, p.publishers, t.types, s.status
		FROM manga m
		LEFT JOIN types t ON m.type_id = t.id
		LEFT JOIN status s ON m.status_id = s.id
		LEFT JOIN publishers p ON m.publisher_id = p.id
		LEFT JOIN volumes v ON m.id = v.title_id WHERE v.volume_number = '1'
`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []Data

		for rows.Next() {
			var o Data
			var coverLink sql.NullString // Use sql.NullString to handle NULL values

			if err := rows.Scan(&o.ID, &o.Title, &coverLink, &o.Mal_Id, &o.Score, &o.Popularity, &o.PublisherID, &o.TypeID, &o.StatusID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Check if coverImage is valid before using its value
			if coverLink.Valid {
				o.CoverLink = coverLink.String
			} else {
				o.CoverLink = "" // Set a default value if coverImage is NULL
			}

			data = append(data, o)
		}

		// Encoding the data as JSON and sending it in the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}).Methods("GET")

	r.HandleFunc("/releases-data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		rows, err := db.Query(`
		SELECT v.id, v.title_id, m.title, v.volume_number, v.release_date, v.cover_link
		FROM volumes v
		LEFT JOIN manga m ON v.title_id = m.id
`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []Volume

		for rows.Next() {
			var o Volume

			if err := rows.Scan(&o.ID, &o.Title_ID, &o.Title, &o.VolumeNumber, &o.ReleaseDate, &o.CoverLink); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data = append(data, o)
		}

		// Encoding the data as JSON and sending it in the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}).Methods("GET")

	// Handler for /other_query/{id} route
	r.HandleFunc("/other_query/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Extracting id from the request URL
		vars := mux.Vars(r)
		ID := vars["id"]

		// Query the database to fetch manga details
		mangaQuery := `
        SELECT m.id, m.title, m.mal_id, m.score, m.popularity, p.publishers, t.types, s.status
		FROM manga m
		LEFT JOIN types t ON m.type_id = t.id
		LEFT JOIN status s ON m.status_id = s.id
		LEFT JOIN publishers p ON m.publisher_id = p.id WHERE m.id = $1
    `
		mangaRow := db.QueryRow(mangaQuery, ID)
		var mangaData Data
		err := mangaRow.Scan(&mangaData.ID, &mangaData.Title, &mangaData.Mal_Id, &mangaData.Score, &mangaData.Popularity, &mangaData.PublisherID, &mangaData.TypeID, &mangaData.StatusID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Query the database to fetch associated volume cover links
		volumesQuery := `
        SELECT cover_link, volume_number
        FROM volumes
        WHERE title_id = $1
    `
		volumesRows, err := db.Query(volumesQuery, ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer volumesRows.Close()

		// Slice to hold volume data
		var volumes []Volume

		// Iterating through the rows and appending data to the slice
		for volumesRows.Next() {
			var volume Volume
			err := volumesRows.Scan(&volume.CoverLink, &volume.VolumeNumber)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			volumes = append(volumes, volume)
		}

		// Assign volumes to mangaData
		mangaData.Volumes = volumes

		// Encoding the data as JSON and sending it in the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mangaData)
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

		// Convert the array of strings to a slice of cover images
		coverImages := make([]string, 0, len(data.CoverLink))

		// Insert data into the PostgreSQL database
		_, err := db.Exec("INSERT INTO manga (title, cover_image, type_id, publisher_id, mal_id, score, popularity, status_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			data.Title, pq.Array(coverImages), data.TypeID, data.PublisherID, data.Mal_Id, data.Score, data.Popularity, data.StatusID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	r.HandleFunc("/add-datas", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		decoder := json.NewDecoder(r.Body)
		var data Volume
		if err := decoder.Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Parse and format the release date
		releaseDate, err := time.Parse("2006-01", data.ReleaseDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		formattedReleaseDate := releaseDate.Format("2006-01-02") // Assuming you want to use the first day of the month

		// Insert data into the PostgreSQL database
		_, err = db.Exec("INSERT INTO volumes (title_id, volume_number, release_date, cover_link) VALUES ($1, $2, $3, $4)",
			data.Title_ID, data.VolumeNumber, formattedReleaseDate, data.CoverLink)

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
		rows, err := db.Query("SELECT id, " + tableName + " FROM " + tableName)
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
		row := db.QueryRow("SELECT id, "+tableName+" FROM "+tableName+" WHERE id = $1", rowID)

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

	r.HandleFunc("/fetch-manga-titles", func(w http.ResponseWriter, r *http.Request) {

		// Query the database to fetch manga titles
		rows, err := db.Query("SELECT id, title FROM manga")
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

	r.HandleFunc("/upload-cover-image", func(w http.ResponseWriter, r *http.Request) {
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

		// Extract other form data
		title := r.FormValue("title")
		volumeNumber := r.FormValue("volumeNumber")

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

		// Construct the URL string for the cover image
		coverImageURL := fmt.Sprintf("http://localhost:8080/covers/%s", handler.Filename)

		// Update the database
		_, err = db.Exec("UPDATE manga SET cover_image[$1] = $2 WHERE title = $3", volumeNumber, coverImageURL, title)
		if err != nil {
			http.Error(w, "Unable to update the database", http.StatusInternalServerError)
			return
		}

		// Respond with success message
		w.Write([]byte("File uploaded successfully"))
	})

	// Applying CORS middleware to the router
	handlerr := c.Handler(r)

	// Handle upload requests
	uploadHandler := c.Handler(http.HandlerFunc(uploadHandler))

	http.Handle("/upload", uploadHandler)

	http.Handle("/", handlerr)

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

// convertVolumeNumberToInt converts the volume number string to an integer index.
func convertVolumeNumberToInt(volumeNumber string) int {
	// Parse the volume number string to an integer
	volumeIndex, err := strconv.Atoi(volumeNumber)
	if err != nil {
		// If there's an error parsing the string to an integer, return -1
		// You can handle this error differently based on your application's requirements
		return -1
	}

	// Subtract 1 from the volume number to convert it to a zero-based index
	volumeIndex-- // Assuming volume numbers are 1-indexed in the database

	return volumeIndex
}
