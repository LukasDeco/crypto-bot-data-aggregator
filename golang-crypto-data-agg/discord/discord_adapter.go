package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DiscordAdapter struct to interact with Discord API
type DiscordAdapter struct {
	BaseURL string // Base URL for Discord API
}

// NewDiscordAdapter creates a new instance of DiscordAdapter with default settings
func NewDiscordAdapter() *DiscordAdapter {
	return &DiscordAdapter{
		BaseURL: "https://discord.com/api/",
	}
}

// GetInviteInfo makes a GET request to Discord API to retrieve invite information
func (adapter *DiscordAdapter) GetInviteInfo(serverID string) (*DiscordInviteResponse, error) {
	// Construct the request URL using the provided server ID
	requestURL := fmt.Sprintf("%sinvite/%s?with_counts=true", adapter.BaseURL, serverID)

	// Create the HTTP GET request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add any required headers here. For example, Authorization if needed.
	// req.Header.Add("Authorization", "Bot YOUR_BOT_TOKEN_HERE")

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to Discord API: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON response into the DiscordInviteResponse struct
	var inviteInfo DiscordInviteResponse
	err = json.Unmarshal(body, &inviteInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Check the HTTP response status code
	if resp.StatusCode != http.StatusOK {
		return &inviteInfo, fmt.Errorf("API request error: %s", resp.Status)
	}

	return &inviteInfo, nil
}

// DiscordInviteResponse represents the JSON response structure for a Discord invite request
type DiscordInviteResponse struct {
	Type                     int         `json:"type"`
	Code                     any         `json:"code"`
	ExpiresAt                *time.Time  `json:"expires_at"` // Using *time.Time to handle null values
	Guild                    GuildInfo   `json:"guild"`
	GuildID                  string      `json:"guild_id"`
	Channel                  ChannelInfo `json:"channel"`
	ApproximateMemberCount   int         `json:"approximate_member_count"`
	ApproximatePresenceCount int         `json:"approximate_presence_count"`
}

// GuildInfo represents detailed information about the guild (server)
type GuildInfo struct {
	ID                       string   `json:"id"`
	Name                     string   `json:"name"`
	Splash                   string   `json:"splash"`
	Banner                   string   `json:"banner"`
	Description              *string  `json:"description"` // Using *string to handle null values
	Icon                     string   `json:"icon"`
	Features                 []string `json:"features"`
	VerificationLevel        int      `json:"verification_level"`
	VanityURLCode            string   `json:"vanity_url_code"`
	NsfwLevel                int      `json:"nsfw_level"`
	Nsfw                     bool     `json:"nsfw"`
	PremiumSubscriptionCount int      `json:"premium_subscription_count"`
}

// ChannelInfo represents basic information about a channel within the guild
type ChannelInfo struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
}
