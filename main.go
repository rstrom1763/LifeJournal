package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"memories/daos"
	. "memories/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	utils "github.com/rstrom1763/goUtils"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func env(key string) string {

	// load .env file
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func envOrDefault(key, fallback string) string {
	value := env(key)
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func main() {

	port := env("PORT") // Port to listen on
	protocol := strings.ToLower(env("PROTOCOL"))
	daoName := strings.ToLower(env("DAO"))

	// Initialize Gin
	gin.SetMode(gin.ReleaseMode) // Turn off debugging mode
	r := gin.Default()           // Initialize Gin

	//Ensure valid protocol env entry
	if protocol != "http" && protocol != "https" {
		log.Fatal("Invalid protocol. Must be HTTP or HTTPS")
	}

	//Generate TLS keys if they do not already exist
	if !(utils.FileExists("./cert.pem") && utils.FileExists("./private.key")) && protocol == "https" {
		utils.GenerateSSL()
	}

	var dao LifeJournalDAO
	var closeDB func()

	switch daoName {
	case "sqlite":
		dbPath := envOrDefault("SQLITE_PATH", "./life_journal.sqlite")
		db := daos.InitSQLiteDB(dbPath)
		closeDB = func() {
			if err := db.Close(); err != nil {
				log.Println("Error closing DB: ", err)
			}
		}
		fmt.Println("Connected to DB (sqlite)")
		dao = daos.NewSQLiteDAO(db)
	case "postgres":
		dsn := env("POSTGRES_DSN")
		db := daos.InitPostgresDB(dsn)
		closeDB = func() {
			if err := db.Close(); err != nil {
				log.Println("Error closing DB: ", err)
			}
		}
		fmt.Println("Connected to DB (postgres)")
		dao = daos.NewPostgresDAO(db)
	default:
		log.Fatalf("Invalid DAO")
	}

	if closeDB != nil {
		defer closeDB()
	}

	// Home page
	r.GET("/", func(c *gin.Context) {
		html, _ := os.ReadFile("./assets/html/home.html")
		c.Data(http.StatusOK, "text/html", html)
	})

	// Get all concerts (HTML page)
	r.GET("/concerts", func(c *gin.Context) {
		html, _ := os.ReadFile("./assets/html/concerts.html")
		c.Data(http.StatusOK, "text/html", html)
	})

	// Get all concerts (JSON API)
	r.GET("/api/concerts", func(c *gin.Context) {
		concerts, err := dao.GetAllConcerts()
		if err != nil {
			log.Printf("Could not get concerts: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get concerts"})
			return
		}

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
		movies, err := dao.GetAllMovies()
		if err != nil {
			log.Printf("Could not get movies: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get movies"})
			return
		}

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

		movies, err := dao.GetMoviesByTier(tier)
		if err != nil {
			log.Printf("Could not get movies: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get movies"})
			return
		}

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
		books, err := dao.GetAllBooks()
		if err != nil {
			log.Printf("Could not get books: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get books"})
			return
		}

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
		foodPlaces, err := dao.GetAllFoodPlaces()
		if err != nil {
			log.Printf("Could not get food: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get food"})
			return
		}

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

	// Get all people (HTML page)
	r.GET("/people", func(c *gin.Context) {
		html, _ := os.ReadFile("./assets/html/people.html")
		c.Data(http.StatusOK, "text/html", html)
	})

	// Get all people (JSON API)
	r.GET("/api/people", func(c *gin.Context) {
		people, err := dao.GetAllPeople()
		if err != nil {
			log.Printf("Could not get people: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get people"})
			return
		}

		jsonData, err := json.Marshal(people)
		if err != nil {
			log.Printf("Could not marshal people: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encode data"})
			return
		}

		gzipData := utils.GzipData(jsonData)

		c.Header("Content-Encoding", "gzip")
		c.Data(http.StatusOK, "application/json", gzipData)
	})

	// Get food places by location (JSON API)
	r.GET("/api/food/:location", func(c *gin.Context) {
		location := c.Param("location")
		location = strings.ToLower(location)

		foodPlaces, err := dao.GetFoodPlacesByLocation(location)
		if err != nil {
			log.Printf("Could not get food: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get food"})
			return
		}

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
		tvShows, err := dao.GetAllTVShows()
		if err != nil {
			log.Printf("Could not get TV shows: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get TV shows"})
			return
		}

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

			fileHandle, err := file.Open()
			if err != nil {
				log.Println("Failed to open file:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
				return
			}
			data, err := io.ReadAll(fileHandle)
			if err != nil {
				log.Println("Failed to read file:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
				return
			}

			// Use DAO to create photo and get ID
			id, err := dao.CreatePhoto(file.Filename, data)
			if err != nil {
				log.Println("Failed to create photo:", err)
				c.String(http.StatusInternalServerError, "Failed to insert data")
				return
			}

			fileIds = append(fileIds, id)
			fileIdsJson, _ = json.Marshal(fileIds)

		}

		fmt.Println(fileIdsJson)
		c.Data(http.StatusOK, "text/plain", fileIdsJson)
	})

	r.POST("/journal/upload", func(c *gin.Context) {

		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Struct for parsing JSON input (without ID and Created fields)
		var e struct {
			Title  string `json:"title"`
			Entry  string `json:"entry"`
			Tags   string `json:"tags"`
			Photos string `json:"photos"`
		}
		err = json.Unmarshal(data, &e)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing body"})
			return
		}

		// Use DAO to create journal entry
		err = dao.CreateJournalEntry(e.Title, e.Entry, e.Tags, e.Photos)
		if err != nil {
			log.Println("Failed to create journal entry:", err)
			c.String(http.StatusInternalServerError, "Failed to insert data")
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Success", "entries": e})

	})

	r.GET("/api/journal", func(c *gin.Context) {
		entries, err := dao.GetAllJournalEntries()
		if err != nil {
			log.Printf("Could not get journal entries: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get journal entries"})
			return
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
		idStr := c.Param("id")
		var id int
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
			return
		}

		photo, err := dao.GetPhotoByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			return
		}

		contentType := "image/jpeg"
		if strings.HasSuffix(strings.ToLower(photo.FileName), ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(strings.ToLower(photo.FileName), ".gif") {
			contentType = "image/gif"
		}

		c.Data(http.StatusOK, contentType, photo.Bytes)
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
