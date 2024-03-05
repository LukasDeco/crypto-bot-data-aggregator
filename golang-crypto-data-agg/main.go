package main

import (
	"database/sql"
	"log"

	"github.com/LukasDeco/crypto-bot-data-aggregator/golang-crypto-data-agg/discord"
	_ "github.com/mattn/go-sqlite3"
)

// CryptoToken represents a cryptocurrency with a name and some metadata
type CryptoToken struct {
	Name          string
	DiscordID     string // This should be just the server ID part if using the previously defined DiscordAdapter
	TwitterHandle string // This should be just the server ID part if using the previously defined DiscordAdapter
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
        presence_count INTEGER
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func fetchAndStoreInvites(db *sql.DB, tokens []CryptoToken, adapter *discord.DiscordAdapter) {
	for _, token := range tokens {
		inviteInfo, err := adapter.GetInviteInfo(token.DiscordID)
		if err != nil {
			log.Printf("Error fetching invite info for %s: %v\n", token.Name, err)
			continue
		}

		_, err = db.Exec(`INSERT INTO discord 
		(token_name, guild_id, guild_name, member_count, presence_count) 
		VALUES (?, ?, ?, ?, ?)`,
			token.Name,
			inviteInfo.Guild.ID,
			inviteInfo.Guild.Name,
			inviteInfo.ApproximateMemberCount,
			inviteInfo.ApproximatePresenceCount,
			inviteInfo.Guild.Name)
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

	// Define your list of crypto tokens and their Discord server IDs
	tokens := []CryptoToken{
		{"MATIC", "XvpHAxZ", ""},
		{"VET", "vechain", ""},
		{"ETH", "ethereum", ""},
		{"DIMO", "dimonetwork", ""},
		{"MXC", "mxcfoundation", ""},
		// Add more tokens as needed...
	}

	// Fetch and store the invite info for each token
	fetchAndStoreInvites(db, tokens, adapter)
}
