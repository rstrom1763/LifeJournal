package model

// DAO interface
type LifeJournalDAO interface {
	// Concert methods
	GetAllConcerts() ([]Concert, error)

	// Movie methods
	GetAllMovies() ([]Movie, error)
	GetMoviesByTier(tier string) ([]Movie, error)

	// Book methods
	GetAllBooks() ([]Book, error)

	// Food methods
	GetAllFoodPlaces() ([]FoodPlace, error)
	GetFoodPlacesByLocation(location string) ([]FoodPlace, error)

	// TV methods
	GetAllTVShows() ([]TVShow, error)

	// Journal methods
	GetAllJournalEntries() ([]JournalEntry, error)
	CreateJournalEntry(title, entry, tags, photos string) error

	// Photo methods
	CreatePhoto(fileName string, bytes []byte) (int, error)
	GetPhotoByID(id int) (*Photo, error)
}

// Concert represents a concert entry
type Concert struct {
	Date    string `json:"Date"`
	Artists string `json:"Artists"`
	Notes   string `json:"Notes"`
	People  string `json:"People"`
}

// Movie represents a movie entry
type Movie struct {
	Title string `json:"Title"`
	Tier  string `json:"Tier"`
}

// Book represents a book entry
type Book struct {
	Title    string  `json:"Title"`
	Rating   float64 `json:"Rating"`
	Pages    int     `json:"Pages"`
	Author   string  `json:"Author"`
	Series   string  `json:"Series"`
	Finished bool    `json:"Finished"`
}

// FoodPlace represents a food place entry
type FoodPlace struct {
	Name     string `json:"Name"`
	Location string `json:"Location"`
	Notes    string `json:"Notes"`
	Type     string `json:"Type"`
	Category string `json:"Category"`
}

// TVShow represents a TV show entry
type TVShow struct {
	Title          string `json:"Title"`
	Notes          string `json:"Notes"`
	SeasonsWatched string `json:"SeasonsWatched"`
}

// JournalEntry represents a journal entry
type JournalEntry struct {
	ID      int    `json:"id"`
	Created string `json:"created"`
	Title   string `json:"title"`
	Entry   string `json:"entry"`
	Tags    string `json:"tags"`
	Photos  string `json:"photos"`
}

// Photo represents a photo/file entry
type Photo struct {
	ID       int    `json:"id"`
	FileName string `json:"fileName"`
	Bytes    []byte `json:"-"`
	Created  string `json:"created"`
}
