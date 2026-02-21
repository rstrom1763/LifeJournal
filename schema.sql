CREATE TABLE video_games (
    id INT,
    title VARCHAR(255),
    notes TEXT,
    multiplayer BOOLEAN,
    PRIMARY KEY (id)
);
CREATE TABLE users (
    uuid CHAR(36),
    email VARCHAR(255),
    password_hash VARCHAR(255),
    salt VARCHAR(255),
    created TIMESTAMP,
    PRIMARY KEY (uuid)
);
CREATE TABLE theater_movies (
    id INT,
    title VARCHAR(255),
    date DATE,
    people_went_with TEXT,
    notes TEXT,
    PRIMARY KEY (id)
);
CREATE TABLE watched_movies (
    title VARCHAR(255),
    rating VARCHAR(50),
    tier VARCHAR(50),
    notes TEXT,
    PRIMARY KEY (title)
);
CREATE TABLE travel (
    title VARCHAR(255),
    places TEXT,
    people_went_with TEXT,
    notes TEXT,
    dates DATE,
    id INT,
    PRIMARY KEY (id)
);
CREATE TABLE tv_shows (
    title VARCHAR(255),
    date DATE,
    notes TEXT,
    seasons_watched TEXT,
    childhood_show BOOLEAN,
    PRIMARY KEY (title)
);
CREATE TABLE books (
    title VARCHAR(255),
    date_finished DATE,
    author VARCHAR(255),
    rating FLOAT,
    series VARCHAR(255),
    pages INT,
    series_sequence INT,
    finished BOOLEAN,
    PRIMARY KEY (title)
);
CREATE TABLE food_places (
    name VARCHAR(255),
    type VARCHAR(255),
    location VARCHAR(255),
    notes TEXT,
    category VARCHAR(50),
    PRIMARY KEY (name)
);
CREATE TABLE life_events (
    id INT,
    title VARCHAR(255),
    month INT,
    day INT,
    year INT,
    notes TEXT,
    PRIMARY KEY (id)
);
CREATE TABLE concerts (
    date DATE,
    artists TEXT,
    notes TEXT,
    people_went_with TEXT,
    PRIMARY KEY (date)
);
CREATE TABLE people (
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
CREATE TABLE random_memories (
    id INT,
    date DATE,
    notes TEXT,
    involved_people TEXT,
    PRIMARY KEY (id)
)
CREATE TABLE journal_entries (

id INTEGER PRIMARY KEY AUTOINCREMENT,
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
entry text NOT NULL,
title varchar(255)

);


