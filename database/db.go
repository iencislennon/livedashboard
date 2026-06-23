package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	DB = db
	migrate()
}

func migrate() {
	schema := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS quick_links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		label TEXT NOT NULL,
		url TEXT NOT NULL,
		page TEXT NOT NULL DEFAULT 'studies'
	);

	CREATE TABLE IF NOT EXISTS classes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL DEFAULT '',
		material_emoji TEXT NOT NULL DEFAULT '📘',
		exam_dates TEXT NOT NULL DEFAULT '',
		hw_deadline TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS dev_projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL DEFAULT '',
		repo_url TEXT NOT NULL DEFAULT '',
		claude_url TEXT NOT NULL DEFAULT '',
		other_url TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS knowledge (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL DEFAULT '',
		topic TEXT NOT NULL DEFAULT '',
		link TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS curiosities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL DEFAULT '',
		link TEXT NOT NULL DEFAULT '',
		description TEXT NOT NULL DEFAULT '',
		notes TEXT NOT NULL DEFAULT '',
		icon_emoji TEXT NOT NULL DEFAULT '✨'
	);

	CREATE TABLE IF NOT EXISTS job_apps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		company TEXT NOT NULL DEFAULT '',
		title TEXT NOT NULL DEFAULT '',
		app_url TEXT NOT NULL DEFAULT '',
		cv_url TEXT NOT NULL DEFAULT '',
		priority TEXT NOT NULL DEFAULT 'medium'
	);

	CREATE TABLE IF NOT EXISTS future_ideas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL DEFAULT '',
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS diet_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		content TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS gym_notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		day_name TEXT NOT NULL,
		exercises TEXT NOT NULL DEFAULT '',
		active BOOLEAN NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS skincare_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		cleanser BOOLEAN NOT NULL DEFAULT 0,
		moisturizer BOOLEAN NOT NULL DEFAULT 0,
		spf BOOLEAN NOT NULL DEFAULT 0,
		serum BOOLEAN NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS hair_notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL DEFAULT '',
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS weekly_scores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		week_label TEXT NOT NULL,
		score INTEGER NOT NULL DEFAULT 5
	);

	CREATE TABLE IF NOT EXISTS reminder_text (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS stats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		label TEXT NOT NULL,
		value INTEGER NOT NULL DEFAULT 5
	);

	CREATE TABLE IF NOT EXISTS other_notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL DEFAULT '',
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		url TEXT NOT NULL
	);
	`
	if _, err := DB.Exec(schema); err != nil {
		log.Fatal(err)
	}

	seedDefaults()
}

func seedDefaults() {
	var count int

	DB.QueryRow("SELECT COUNT(*) FROM reminder_text").Scan(&count)
	if count == 0 {
		DB.Exec("INSERT INTO reminder_text (content) VALUES ('Drink more water today')")
	}

	DB.QueryRow("SELECT COUNT(*) FROM future_ideas").Scan(&count)
	if count == 0 {
		DB.Exec("INSERT INTO future_ideas (content) VALUES ('')")
	}

	DB.QueryRow("SELECT COUNT(*) FROM hair_notes").Scan(&count)
	if count == 0 {
		DB.Exec("INSERT INTO hair_notes (content) VALUES ('')")
	}

	DB.QueryRow("SELECT COUNT(*) FROM other_notes").Scan(&count)
	if count == 0 {
		DB.Exec("INSERT INTO other_notes (content) VALUES ('')")
	}

	DB.QueryRow("SELECT COUNT(*) FROM stats").Scan(&count)
	if count == 0 {
		labels := []string{"Studies", "Work", "Health", "Fitness", "Mood"}
		for _, l := range labels {
			DB.Exec("INSERT INTO stats (label, value) VALUES (?, 5)", l)
		}
	}

	DB.QueryRow("SELECT COUNT(*) FROM gym_notes").Scan(&count)
	if count == 0 {
		days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		for _, d := range days {
			DB.Exec("INSERT INTO gym_notes (day_name, exercises, active) VALUES (?, '', 0)", d)
		}
	}
}
