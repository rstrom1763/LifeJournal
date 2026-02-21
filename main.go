package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	utils "github.com/rstrom1763/goUtils"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

// Create the DB connection and create the tables if they don't already exist
func initDB(path string) *sql.DB {

	createTablesQuery := `CREATE TABLE IF NOT EXISTS video_games (
    id INT,
    title VARCHAR(255),
    notes TEXT,
    multiplayer BOOLEAN,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS users (
    uuid CHAR(36),
    email VARCHAR(255),
    password_hash VARCHAR(255),
    salt VARCHAR(255),
    created TIMESTAMP,
    PRIMARY KEY (uuid)
);
CREATE TABLE IF NOT EXISTS theater_movies (
    id INT,
    title VARCHAR(255),
    date DATE,
    people_went_with TEXT,
    notes TEXT,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS watched_movies (
    title VARCHAR(255),
    rating VARCHAR(50),
    tier VARCHAR(50),
    notes TEXT,
    PRIMARY KEY (title)
);
CREATE TABLE IF NOT EXISTS travel (
    title VARCHAR(255),
    places TEXT,
    people_went_with TEXT,
    notes TEXT,
    dates DATE,
    id INT,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS tv_shows (
    title VARCHAR(255),
    date DATE,
    notes TEXT,
    seasons_watched TEXT,
    childhood_show BOOLEAN,
    PRIMARY KEY (title)
);
CREATE TABLE IF NOT EXISTS books (
    title VARCHAR(255),
    date_finished DATE,
    author VARCHAR(255),
    rating FLOAT,
    series VARCHAR(255),
    owned BOOLEAN,
    pages INT,
    series_sequence INT,
    finished BOOLEAN,
    PRIMARY KEY (title)
);
CREATE TABLE IF NOT EXISTS food_places (
    name VARCHAR(255),
    type VARCHAR(255),
    location VARCHAR(255),
    notes TEXT,
    category VARCHAR(50),
    PRIMARY KEY (name)
);
CREATE TABLE IF NOT EXISTS life_events (
    id INT,
    title VARCHAR(255),
    month INT,
    day INT,
    year INT,
    notes TEXT,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS concerts (
    date DATE,
    artists TEXT,
    notes TEXT,
    people_went_with TEXT,
    PRIMARY KEY (date)
);
CREATE TABLE IF NOT EXISTS people (
    id INT,
    first VARCHAR(255),
    middle VARCHAR(255),
    last VARCHAR(255),
    address VARCHAR(255),
    birth_day INT,
    birth_month INT,
    birth_year INT,
    gift_ideas TEXT[],
    email VARCHAR(255),
    category VARCHAR(50),
    notes TEXT,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS random_memories (
    id INT,
    date DATE,
    notes TEXT,
    involved_people TEXT,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS journal_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    entry TEXT NOT NULL,
    title VARCHAR(255),
    tags TEXT,
    photos TEXT
);
CREATE TABLE IF NOT EXISTS files (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	bytes BLOB NOT NULL,
	file_name VARCHAR(255) NOT NULL,
	created DATETIME DEFAULT CURRENT_TIMESTAMP
);`

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not open DB: %s", err))
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not ping DB: %s", err))
	}

	_, err = db.Exec(createTablesQuery)
	if err != nil {
		log.Fatalf("Could not create tables: %s", err)
	}

	return db

}

type Concert struct {
	Date    string
	Artists string
	Notes   string
	People  string
}

type Movie struct {
	Title string
	Tier  string
}

type Book struct {
	Title    string
	Rating   float64
	Pages    int
	Author   string
	Series   string
	Finished bool
}

type FoodPlace struct {
	Name     string
	Location string
	Notes    string
	Type     string
	Category string
}

type TVShow struct {
	Title          string
	Notes          string
	SeasonsWatched string
}

type JournalEntry struct {
	Title  string
	Entry  string
	Tags   string
	Photos string
}

func env(key string) string {

	// load .env file
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func generateSSL() {

	// Generate a private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Error generating private key:", err)
		return
	}

	// Generate a self-signed certificate
	certTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal("Error creating certificate:", err)
		return
	}

	// Write the private key and certificate to files
	keyOut, err := os.Create("./private.key")
	if err != nil {
		log.Fatal("Error creating private key file:", err)
		return
	}
	defer keyOut.Close()

	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err != nil {
		log.Fatal("Error creating certificate file: ", err)
		return
	}

	certOut, err := os.Create("./cert.pem")
	if err != nil {
		log.Fatal("Error creating certificate file: ", err)
		return
	}
	defer certOut.Close()

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		log.Fatal("Error creating certificate file: ", err)
		return
	}

	fmt.Println("TLS certificate and private key generated successfully.")
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true // File exists
	}
	if os.IsNotExist(err) {
		return false // File does not exist
	}
	return false // Error occurred (e.g., permission denied)
}

func listSlice(slice []string) string {
	var output string

	length := len(slice)
	for i, item := range slice {

		output += item
		if i != length-1 {
			output += ", "
		}
	}

	return output
}

func listIntSlice(slice []int) string {
	var output string

	length := len(slice)
	for i, item := range slice {

		output += fmt.Sprintf("%d", item)
		if i != length-1 {
			output += ", "
		}

	}

	return output
}

func main() {

	port := env("PORT") // Port to listen on
	protocol := strings.ToLower(env("PROTOCOL"))

	// Initialize Gin
	gin.SetMode(gin.ReleaseMode) // Turn off debugging mode
	r := gin.Default()           // Initialize Gin

	//Ensure valid protocol env entry
	if protocol != "http" && protocol != "https" {
		log.Fatal("Invalid protocol. Must be HTTP or HTTPS")
	}

	//Generate TLS keys if they do not already exist
	if !(fileExists("./cert.pem") && fileExists("./private.key")) && protocol == "https" {
		generateSSL()
	}

	db := initDB("./life_journal.sqlite")

	// Close session at the end
	defer func(db *sql.DB) {

		err := db.Close()
		if err != nil {
			log.Println("Error closing DB: ", err)
		}
	}(db)
	fmt.Println("Connected to DB")

	// Get all concerts (HTML page)
	r.GET("/concerts", func(c *gin.Context) {
		html, _ := os.ReadFile("./assets/html/concerts.html")
		c.Data(http.StatusOK, "text/html", html)
	})

	// Get all concerts (JSON API)
	r.GET("/api/concerts", func(c *gin.Context) {
		rows, err := db.Query("SELECT artists, people_went_with, notes FROM concerts")
		if err != nil {
			log.Printf("Could not get concerts: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get concerts"})
			return
		}

		var concerts []Concert

		for rows.Next() {
			var concert Concert
			err = rows.Scan(&concert.Artists, &concert.People, &concert.Notes)
			if err != nil {
				log.Printf("Could not scan rows: %v", err)
				continue
			}
			concerts = append(concerts, concert)
		}
		_ = rows.Close()

		jsonData, err := json.Marshal(concerts)
		if err != nil {
			log.Printf("Could not marshal concerts: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode data"})
			return
		}

		gzipData := utils.GzipData(jsonData)

		c.Header("Content-Encoding", "gzip")
		c.Data(http.StatusOK, "application/json", gzipData)
	})

	// Get all movies (HTML page)
	r.GET("/movies", func(c *gin.Context) {
		html, _ := os.ReadFile("./assets/html/movies.html")
		c.Data(http.StatusOK, "text/html", html)
	})

	// Get all movies (JSON API)
	r.GET("/api/movies", func(c *gin.Context) {
		getQuery := `SELECT title, tier FROM watched_movies`

		rows, err := db.Query(getQuery)
		if err != nil {
			log.Printf("Could not get movies: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get movies"})
			return
		}

		var movies []Movie

		for rows.Next() {
			var movie Movie
			err = rows.Scan(&movie.Title, &movie.Tier)
			if err != nil {
				log.Printf("Could not scan rows: %v", err)
				continue
			}
			movies = append(movies, movie)
		}
		_ = rows.Close()

		jsonData, err := json.Marshal(movies)
		if err != nil {
			log.Printf("Could not marshal movies: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode data"})
			return
		}

		gzipData := utils.GzipData(jsonData)

		c.Header("Content-Encoding", "gzip")
		c.Data(http.StatusOK, "application/json", gzipData)
	})

	// Get all movies of the provided tier (JSON API)
	r.GET("/api/movies/:tier", func(c *gin.Context) {
		tier := strings.ToUpper(c.Param("tier"))
		getQuery := `SELECT title, tier FROM watched_movies WHERE tier = ?`

		rows, err := db.Query(getQuery, tier)
		if err != nil {
			log.Printf("Could not get movies: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get movies"})
			return
		}

		var movies []Movie

		for rows.Next() {
			var movie Movie
			err := rows.Scan(&movie.Title, &movie.Tier)
			if err != nil {
				log.Printf("Could not scan rows: %v", err)
				continue
			}
			movies = append(movies, movie)
		}
		_ = rows.Close()

		jsonData, err := json.Marshal(movies)
		if err != nil {
			log.Printf("Could not marshal movies: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode data"})
			return
		}

		gzipData := utils.GzipData(jsonData)

		c.Header("Content-Encoding", "gzip")
		c.Data(http.StatusOK, "application/json", gzipData)
	})

	// Get all books (HTML page)
	r.GET("/books", func(c *gin.Context) {
		html, _ := os.ReadFile("./assets/html/books.html")
		c.Data(http.StatusOK, "text/html", html)
	})

	// Get all books (JSON API)
	r.GET("/api/books", func(c *gin.Context) {
		getQuery := `SELECT title, rating, pages, author, series, finished FROM books`

		rows, err := db.Query(getQuery)
		if err != nil {
			log.Printf("Could not get books: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get books"})
			return
		}

		var books []Book

		for rows.Next() {
			var book Book
			err = rows.Scan(&book.Title, &book.Rating, &book.Pages, &book.Author, &book.Series, &book.Finished)
			if err != nil {
				log.Printf("Could not scan rows: %v", err)
				continue
			}
			books = append(books, book)
		}
		_ = rows.Close()

		jsonData, err := json.Marshal(books)
		if err != nil {
			log.Printf("Could not marshal books: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode data"})
			return
		}

		gzipData := utils.GzipData(jsonData)

		c.Header("Content-Encoding", "gzip")
		c.Data(http.StatusOK, "application/json", gzipData)
	})

	// Get all food places (HTML page)
	r.GET("/food", func(c *gin.Context) {
		html, _ := os.ReadFile("./assets/html/food.html")
		c.Data(http.StatusOK, "text/html", html)
	})

	// Get all food places (JSON API)
	r.GET("/api/food", func(c *gin.Context) {
		getQuery := `SELECT name, location, notes, type, category FROM food_places ORDER BY name`

		rows, err := db.Query(getQuery)
		if err != nil {
			log.Printf("Could not get food: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get food"})
			return
		}

		var foodPlaces []FoodPlace

		for rows.Next() {
			var place FoodPlace
			_ = rows.Scan(&place.Name, &place.Location, &place.Notes, &place.Type, &place.Category)
			foodPlaces = append(foodPlaces, place)
		}
		_ = rows.Close()

		jsonData, err := json.Marshal(foodPlaces)
		if err != nil {
			log.Printf("Could not marshal food places: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode data"})
			return
		}

		gzipData := utils.GzipData(jsonData)

		c.Header("Content-Encoding", "gzip")
		c.Data(http.StatusOK, "application/json", gzipData)
	})

	// Get food places by location (JSON API)
	r.GET("/api/food/:location", func(c *gin.Context) {
		getQuery := `SELECT name, location, notes, type, category FROM food_places WHERE location = ? ORDER BY name`

		location := c.Param("location")
		location = strings.ToLower(location)

		rows, err := db.Query(getQuery, location)
		if err != nil {
			log.Printf("Could not get food: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get food"})
			return
		}

		var foodPlaces []FoodPlace

		for rows.Next() {
			var place FoodPlace
			_ = rows.Scan(&place.Name, &place.Location, &place.Notes, &place.Type, &place.Category)
			foodPlaces = append(foodPlaces, place)
		}
		_ = rows.Close()

		jsonData, err := json.Marshal(foodPlaces)
		if err != nil {
			log.Printf("Could not marshal food places: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode data"})
			return
		}

		gzipData := utils.GzipData(jsonData)

		c.Header("Content-Encoding", "gzip")
		c.Data(http.StatusOK, "application/json", gzipData)
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		img, _ := os.ReadFile("./assets/images/img.png")

		c.Data(http.StatusOK, "image/x-icon", img)
	})

	// Get all TV shows (HTML page)
	r.GET("/tv", func(c *gin.Context) {
		html, _ := os.ReadFile("./assets/html/tv.html")
		c.Data(http.StatusOK, "text/html", html)
	})

	// Get all TV shows (JSON API)
	r.GET("/api/tv", func(c *gin.Context) {
		getQuery := `SELECT title, notes, seasons_watched FROM tv_shows`

		rows, err := db.Query(getQuery)
		if err != nil {
			log.Printf("Could not get TV shows: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get TV shows"})
			return
		}

		var tvShows []TVShow

		for rows.Next() {
			var show TVShow
			_ = rows.Scan(&show.Title, &show.Notes, &show.SeasonsWatched)
			tvShows = append(tvShows, show)
		}
		_ = rows.Close()

		jsonData, err := json.Marshal(tvShows)
		if err != nil {
			log.Printf("Could not marshal TV shows: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode data"})
			return
		}

		gzipData := utils.GzipData(jsonData)

		c.Header("Content-Encoding", "gzip")
		c.Data(http.StatusOK, "application/json", gzipData)
	})

	// Endpoint for file upload
	r.POST("/journal/upload/photos", func(c *gin.Context) {

		var fileIds []int
		var fileIdsJson []byte

		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		files := form.File["images[]"] // Get multiple file uploads
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
			return
		}

		for _, file := range files {
			// Generate file path
			filePath := filepath.Join("./uploads", file.Filename)

			// Save file to disk
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}

			// Define the insert query
			insertQuery := `INSERT INTO files (file_name,bytes) VALUES (?,?)`

			fileHandle, _ := file.Open()
			if err != nil {
				return
			}
			data, err := io.ReadAll(fileHandle)

			// Execute the insert query
			_, err = db.Exec(insertQuery, file.Filename, data)
			if err != nil {
				log.Println("Failed to execute query:", err)
				c.String(http.StatusInternalServerError, "Failed to insert data")
			}

			// Define the get query
			getQuery := `SELECT id FROM files WHERE file_name = ? ORDER BY created DESC LIMIT 1`
			rows, err := db.Query(getQuery, file.Filename)
			if err != nil {
				fmt.Println("Failed to execute query:", err)
				continue
			}

			var id int
			if rows.Next() {
				_ = rows.Scan(&id)
				fileIds = append(fileIds, id)
			}
			rows.Close()

			fileIdsJson, _ = json.Marshal(fileIds)

		}

		fmt.Println(fileIdsJson)
		c.Data(http.StatusOK, "text/plain", fileIdsJson)
	})

	r.POST("/journal/upload", func(c *gin.Context) {

		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		}

		var e JournalEntry
		err = json.Unmarshal(data, &e)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing body"})
		}

		// Define the insert query
		insertQuery := `INSERT INTO journal_entries (title,entry,tags,photos) VALUES (?,?,?,?)`

		// Execute the insert query
		_, err = db.Exec(insertQuery, e.Title, e.Entry, e.Tags, e.Photos)
		if err != nil {
			log.Println("Failed to execute query:", err)
			c.String(http.StatusInternalServerError, "Failed to insert data")
		}

		c.JSON(http.StatusOK, gin.H{"message": "Success", "entries": e})

	})

	r.GET("/api/journal", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, created, title, entry, tags, photos FROM journal_entries ORDER BY created DESC")
		if err != nil {
			log.Printf("Could not get journal entries: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get journal entries"})
			return
		}
		defer rows.Close()

		var entries []map[string]interface{}
		for rows.Next() {
			var id int
			var createdAt, title, entry, tags, photos string
			err = rows.Scan(&id, &createdAt, &title, &entry, &tags, &photos)
			if err != nil {
				log.Printf("Could not scan journal entry: %v", err)
				continue
			}

			entries = append(entries, map[string]interface{}{
				"id":      id,
				"created": createdAt,
				"title":   title,
				"entry":   entry,
				"tags":    tags,
				"photos":  photos,
			})
		}

		jsonData, err := json.Marshal(entries)
		if err != nil {
			log.Printf("Could not marshal journal entries: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode data"})
			return
		}

		gzipData := utils.GzipData(jsonData)

		c.Header("Content-Encoding", "gzip")
		c.Data(http.StatusOK, "application/json", gzipData)
	})

	r.GET("/api/photos/:id", func(c *gin.Context) {
		id := c.Param("id")
		var fileName string
		var data []byte
		err := db.QueryRow("SELECT file_name, bytes FROM files WHERE id = ?", id).Scan(&fileName, &data)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			return
		}

		contentType := "image/jpeg"
		if strings.HasSuffix(strings.ToLower(fileName), ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(strings.ToLower(fileName), ".gif") {
			contentType = "image/gif"
		}

		c.Data(http.StatusOK, contentType, data)
	})

	r.GET("/journal", func(c *gin.Context) {
		html, _ := os.ReadFile("./assets/html/journal.html")
		c.Data(http.StatusOK, "text/html", html)
	})

	r.GET("/entries", func(c *gin.Context) {
		html, _ := os.ReadFile("./assets/html/journal_list.html")
		c.Data(http.StatusOK, "text/html", html)
	})

	r.GET("/style.css", func(c *gin.Context) {
		css, _ := os.ReadFile("./assets/css/style.css")
		c.Data(http.StatusOK, "text/css", css)
	})

	r.Static("/uploads", "./uploads")

	fmt.Printf("Listening for %v on port %v...\n", protocol, port) //Notifies that server is running on X port
	if protocol == "http" {                                        //Start running the Gin server
		err := r.Run(":" + port)
		if err != nil {
			fmt.Println(err)
		}
	} else if protocol == "https" {
		err := r.RunTLS(":"+port, "./cert.pem", "./private.key")
		if err != nil {
			fmt.Println(err)
		}
	} else {
		log.Fatal("Something went wrong starting the Gin server")
	}

}
