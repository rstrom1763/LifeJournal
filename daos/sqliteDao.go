package daos

import (
	"database/sql"
	"fmt"
	"log"

	. "memories/model"
)

// SQLiteDAO implementation
type SQLiteDAO struct {
	db *sql.DB
}

// NewSQLiteDAO creates a new SQLite DAO
func NewSQLiteDAO(db *sql.DB) *SQLiteDAO {
	return &SQLiteDAO{db: db}
}

// Concert methods
func (dao *SQLiteDAO) GetAllConcerts() ([]Concert, error) {
	rows, err := dao.db.Query("SELECT artists, people_went_with, notes FROM concerts")
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
	rows, err := dao.db.Query("SELECT title, tier FROM watched_movies")
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
	rows, err := dao.db.Query("SELECT title, tier FROM watched_movies WHERE tier = ?", tier)
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
	rows, err := dao.db.Query("SELECT title, rating, pages, author, series, finished FROM books")
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
	rows, err := dao.db.Query("SELECT name, location, notes, type, category FROM food_places ORDER BY name")
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
	rows, err := dao.db.Query("SELECT name, location, notes, type, category FROM food_places WHERE location = ? ORDER BY name", location)
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

// TV methods
func (dao *SQLiteDAO) GetAllTVShows() ([]TVShow, error) {
	rows, err := dao.db.Query("SELECT title, notes, seasons_watched FROM tv_shows")
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
	rows, err := dao.db.Query("SELECT id, created, title, entry, tags, photos FROM journal_entries ORDER BY created DESC")
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
	err := dao.db.QueryRow("SELECT id, file_name, bytes, created FROM files WHERE id = ?", id).
		Scan(&photo.ID, &photo.FileName, &photo.Bytes, &photo.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("photo not found")
		}
		return nil, fmt.Errorf("failed to query photo: %w", err)
	}

	return &photo, nil
}
