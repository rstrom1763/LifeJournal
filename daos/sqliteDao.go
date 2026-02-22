package daos

import (
	"database/sql"
	"fmt"
	"log"

	. "memories/model"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDAO implementation
type SQLiteDAO struct {
	db *sql.DB
}

// Create the SQLite DB connection and create the tables if they don't already exist
func InitSQLiteDB(path string) *sql.DB {
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

// NewSQLiteDAO creates a new SQLite DAO
func NewSQLiteDAO(db *sql.DB) *SQLiteDAO {
	return &SQLiteDAO{db: db}
}

// Concert methods
func (dao *SQLiteDAO) GetAllConcerts() ([]Concert, error) {
	rows, err := dao.db.Query("SELECT COALESCE(artists, ''), COALESCE(people_went_with, ''), COALESCE(notes, '') FROM concerts")
	if err != nil {
		return nil, fmt.Errorf("failed to query concerts: %w", err)
	}
	defer rows.Close()

	var concerts []Concert
	for rows.Next() {
		var concert Concert
		err = rows.Scan(&concert.Artists, &concert.People, &concert.Notes)
		if err != nil {
			log.Printf("Failed to scan concert row: %v", err)
			continue
		}
		concerts = append(concerts, concert)
	}

	return concerts, nil
}

// Movie methods
func (dao *SQLiteDAO) GetAllMovies() ([]Movie, error) {
	rows, err := dao.db.Query("SELECT COALESCE(title, ''), COALESCE(tier, '') FROM watched_movies")
	if err != nil {
		return nil, fmt.Errorf("failed to query movies: %w", err)
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.Title, &movie.Tier)
		if err != nil {
			log.Printf("Failed to scan movie row: %v", err)
			continue
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (dao *SQLiteDAO) GetMoviesByTier(tier string) ([]Movie, error) {
	rows, err := dao.db.Query("SELECT COALESCE(title, ''), COALESCE(tier, '') FROM watched_movies WHERE tier = ?", tier)
	if err != nil {
		return nil, fmt.Errorf("failed to query movies by tier: %w", err)
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.Title, &movie.Tier)
		if err != nil {
			log.Printf("Failed to scan movie row: %v", err)
			continue
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// Book methods
func (dao *SQLiteDAO) GetAllBooks() ([]Book, error) {
	rows, err := dao.db.Query("SELECT COALESCE(title, ''), COALESCE(rating, 0), COALESCE(pages, 0), COALESCE(author, ''), COALESCE(series, ''), COALESCE(finished, 0) FROM books")
	if err != nil {
		return nil, fmt.Errorf("failed to query books: %w", err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.Title, &book.Rating, &book.Pages, &book.Author, &book.Series, &book.Finished)
		if err != nil {
			log.Printf("Failed to scan book row: %v", err)
			continue
		}
		books = append(books, book)
	}

	return books, nil
}

// Food methods
func (dao *SQLiteDAO) GetAllFoodPlaces() ([]FoodPlace, error) {
	rows, err := dao.db.Query("SELECT COALESCE(name, ''), COALESCE(location, ''), COALESCE(notes, ''), COALESCE(type, ''), COALESCE(category, '') FROM food_places ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to query food places: %w", err)
	}
	defer rows.Close()

	var foodPlaces []FoodPlace
	for rows.Next() {
		var place FoodPlace
		err = rows.Scan(&place.Name, &place.Location, &place.Notes, &place.Type, &place.Category)
		if err != nil {
			log.Printf("Failed to scan food place row: %v", err)
			continue
		}
		foodPlaces = append(foodPlaces, place)
	}

	return foodPlaces, nil
}

func (dao *SQLiteDAO) GetFoodPlacesByLocation(location string) ([]FoodPlace, error) {
	rows, err := dao.db.Query("SELECT COALESCE(name, ''), COALESCE(location, ''), COALESCE(notes, ''), COALESCE(type, ''), COALESCE(category, '') FROM food_places WHERE location = ? ORDER BY name", location)
	if err != nil {
		return nil, fmt.Errorf("failed to query food places by location: %w", err)
	}
	defer rows.Close()

	var foodPlaces []FoodPlace
	for rows.Next() {
		var place FoodPlace
		err = rows.Scan(&place.Name, &place.Location, &place.Notes, &place.Type, &place.Category)
		if err != nil {
			log.Printf("Failed to scan food place row: %v", err)
			continue
		}
		foodPlaces = append(foodPlaces, place)
	}

	return foodPlaces, nil
}

// People methods
func (dao *SQLiteDAO) GetAllPeople() ([]Person, error) {
	rows, err := dao.db.Query("SELECT id, COALESCE(first, ''), COALESCE(middle, ''), COALESCE(last, ''), COALESCE(address, ''), COALESCE(birth_day, 0), COALESCE(birth_month, 0), COALESCE(birth_year, 0), COALESCE(gift_ideas, ''), COALESCE(email, ''), COALESCE(category, ''), COALESCE(notes, '') FROM people ORDER BY last, first")
	if err != nil {
		return nil, fmt.Errorf("failed to query people: %w", err)
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var person Person
		err = rows.Scan(
			&person.ID,
			&person.First,
			&person.Middle,
			&person.Last,
			&person.Address,
			&person.BirthDay,
			&person.BirthMonth,
			&person.BirthYear,
			&person.GiftIdeas,
			&person.Email,
			&person.Category,
			&person.Notes,
		)
		if err != nil {
			log.Printf("Failed to scan person row: %v", err)
			continue
		}
		people = append(people, person)
	}

	return people, nil
}

// TV methods
func (dao *SQLiteDAO) GetAllTVShows() ([]TVShow, error) {
	rows, err := dao.db.Query("SELECT COALESCE(title, ''), COALESCE(notes, ''), COALESCE(seasons_watched, '') FROM tv_shows")
	if err != nil {
		return nil, fmt.Errorf("failed to query TV shows: %w", err)
	}
	defer rows.Close()

	var tvShows []TVShow
	for rows.Next() {
		var show TVShow
		err = rows.Scan(&show.Title, &show.Notes, &show.SeasonsWatched)
		if err != nil {
			log.Printf("Failed to scan TV show row: %v", err)
			continue
		}
		tvShows = append(tvShows, show)
	}

	return tvShows, nil
}

// Journal methods
func (dao *SQLiteDAO) GetAllJournalEntries() ([]JournalEntry, error) {
	rows, err := dao.db.Query("SELECT id, COALESCE(created, ''), COALESCE(title, ''), COALESCE(entry, ''), COALESCE(tags, ''), COALESCE(photos, '') FROM journal_entries ORDER BY created DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query journal entries: %w", err)
	}
	defer rows.Close()

	var entries []JournalEntry
	for rows.Next() {
		var entry JournalEntry
		err = rows.Scan(&entry.ID, &entry.Created, &entry.Title, &entry.Entry, &entry.Tags, &entry.Photos)
		if err != nil {
			log.Printf("Failed to scan journal entry row: %v", err)
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (dao *SQLiteDAO) CreateJournalEntry(title, entry, tags, photos string) error {
	insertQuery := `INSERT INTO journal_entries (title, entry, tags, photos) VALUES (?, ?, ?, ?)`
	_, err := dao.db.Exec(insertQuery, title, entry, tags, photos)
	if err != nil {
		return fmt.Errorf("failed to insert journal entry: %w", err)
	}
	return nil
}

// Photo methods
func (dao *SQLiteDAO) CreatePhoto(fileName string, bytes []byte) (int, error) {
	insertQuery := `INSERT INTO files (file_name, bytes) VALUES (?, ?)`
	result, err := dao.db.Exec(insertQuery, fileName, bytes)
	if err != nil {
		return 0, fmt.Errorf("failed to insert photo: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

func (dao *SQLiteDAO) GetPhotoByID(id int) (*Photo, error) {
	var photo Photo
	err := dao.db.QueryRow("SELECT id, COALESCE(file_name, ''), bytes, COALESCE(created, '') FROM files WHERE id = ?", id).
		Scan(&photo.ID, &photo.FileName, &photo.Bytes, &photo.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("photo not found")
		}
		return nil, fmt.Errorf("failed to query photo: %w", err)
	}

	return &photo, nil
}
