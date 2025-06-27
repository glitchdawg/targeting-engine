package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	if connStr == "" {
		log.Fatal("DB connection error")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS campaigns (
        id TEXT PRIMARY KEY,
        name TEXT,
        img TEXT,
        cta TEXT,
        status TEXT
    );
    CREATE TABLE IF NOT EXISTS targeting_rules (
        campaign_id TEXT PRIMARY KEY,
        include_country TEXT,
        exclude_country TEXT,
        include_os TEXT,
        exclude_os TEXT,
        include_app TEXT,
        exclude_app TEXT,
        include_state TEXT,
        exclude_state TEXT
    );
    CREATE TABLE IF NOT EXISTS country_states (
        country TEXT PRIMARY KEY,
        states TEXT
    );
    DELETE FROM campaigns;
    DELETE FROM targeting_rules;
    DELETE FROM country_states;
    INSERT INTO campaigns (id, name, img, cta, status) VALUES
        ('spotify', 'Spotify', 'https://somelink', 'Download', 'ACTIVE'),
        ('duolingo', 'Duolingo', 'https://somelink2', 'Install', 'ACTIVE'),
        ('subwaysurfer', 'Subway Surfer', 'https://somelink3', 'Play', 'ACTIVE');
    INSERT INTO targeting_rules (campaign_id, include_country, exclude_country, include_os, exclude_os, include_app, exclude_app, include_state, exclude_state) VALUES
        ('spotify', 'US,Canada', NULL, 'Android,iOS', NULL, NULL, NULL, NULL, NULL),
        ('duolingo', NULL, 'US', 'Android,iOS', NULL, NULL, NULL, NULL, NULL),
        ('subwaysurfer', NULL, NULL, 'Android', NULL, 'com.gametion.ludokinggame', NULL, NULL, NULL);
    INSERT INTO country_states (country, states) VALUES
        ('US', 'California,Texas,New York,Florida'),
        ('Canada', 'Ontario,Quebec,British Columbia,Alberta');
    `)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database seeded successfully.")
}
