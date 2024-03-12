package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/LukasDeco/crypto-bot-data-aggregator/golang-crypto-data-agg/discord"
	_ "github.com/mattn/go-sqlite3"
)

// CryptoToken represents a cryptocurrency with a name and some metadata
type CryptoToken struct {
	Name             string
	DiscordInviteURL string // This should be just the server ID part if using the previously defined DiscordAdapter
	TwitterHandle    string // This should be just the server ID part if using the previously defined DiscordAdapter
}

func setupDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "crypto_data.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS discord (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        token_name TEXT NOT NULL,
        guild_id TEXT,
        guild_name TEXT,
        member_count INTEGER,
        presence_count INTEGER,
        created_at TEXT
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func fetchAndStoreInvites(db *sql.DB, tokens []CryptoToken, adapter *discord.DiscordAdapter) {
	for _, token := range tokens {
		time.Sleep(time.Second * 3)
		apiInviteUrl, err := convertDiscordURLToAPIURL(token.DiscordInviteURL)
		if err != nil {
			log.Printf("Error parsing URL %s: %v\n", token.Name, err)
			continue
		}
		inviteInfo, err := adapter.GetInviteInfo(apiInviteUrl)
		if err != nil {
			log.Printf("Error fetching invite info for %s: %v\n", token.Name, err)
			continue
		}

		_, err = db.Exec(`INSERT INTO discord 
		(token_name, guild_id, guild_name, member_count, presence_count, created_at) 
		VALUES (?, ?, ?, ?, ?)`,
			token.Name,
			inviteInfo.Guild.ID,
			inviteInfo.Guild.Name,
			inviteInfo.ApproximateMemberCount,
			inviteInfo.ApproximatePresenceCount,
			inviteInfo.Guild.Name,
			time.Now(),
		)
		if err != nil {
			log.Printf("Error inserting invite info for %s into DB: %v\n",
				token.Name,
				err)
		}
	}
}

func main() {
	// Initialize the DiscordAdapter
	adapter := discord.NewDiscordAdapter()

	// Setup the SQLite DB
	db := setupDatabase()
	defer db.Close()

	// Open the CSV file
	file, err := os.Open("discord/discord_urls.csv")
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	// Define your list of crypto tokens and their Discord server IDs
	tokens := []CryptoToken{}

	// Skip the header row by starting with i = 1
	for _, record := range records[1:] {
		tokens = append(tokens, CryptoToken{Name: record[0], DiscordInviteURL: record[1]})
	}

	specialTokens := []CryptoToken{
		{"MATIC", "https://discord.com/invite/XvpHAxZ", ""},
		{"VET", "https://discord.com/invite/vechain", ""},
		{"DIMO", "https://discord.com/invite/dimonetwork", ""},
		{"MXC", "https://discord.com/invite/mxcfoundation", ""},
		// Add more tokens as needed...
	}
	tokens = append(tokens, specialTokens...)

	// Fetch and store the invite info for each token
	fetchAndStoreInvites(db, tokens, adapter)
}

func convertDiscordURLToAPIURL(inviteURL string) (string, error) {
	// Parse the input URL
	parsedURL, err := url.Parse(inviteURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %v", err)
	}

	// Ensure the URL is a Discord invite link
	if !strings.Contains(parsedURL.Host, "discord.com") || !strings.HasPrefix(parsedURL.Path, "/invite/") {
		return "", fmt.Errorf("invalid Discord invite URL")
	}

	// Replace the path to point to the API endpoint and add query parameter
	parsedURL.Path = strings.Replace(parsedURL.Path, "/invite/", "/api/invite/", 1)
	parsedURL.RawQuery = "with_counts=true"

	// Return the modified URL as a string
	return parsedURL.String(), nil
}
