package dataagg

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type ListingsResponse struct {
	Data []Token `json:"data"`
}

type Token struct {
	Symbol string `json:"symbol"`
}

type InfoResponse struct {
	Data map[string][]CryptoInfo `json:"data"`
}

type CryptoInfo struct {
	Urls Urls `json:"urls"`
}

type Urls struct {
	Chat    []string `json:"chat"`
	Twitter []string `json:"twitter"`
	Reddit  []string `json:"reddit"`
	Github  []string `json:"source_code"`
}

func TestCoinMarketCapGetDiscords(t *testing.T) {
	apiKey := ""
	limit := 3623
	start := 1

	// Open a new CSV file for writing
	fileName := fmt.Sprintf("token_metadata-%s+.csv", time.Now().Format(time.RFC3339))
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	if err := writer.Write([]string{"Token Symbol", "Discord URL", "Twitter", "Reddit", "Github"}); err != nil {
		log.Fatalf("Failed to write CSV header: %v", err)
	}

	for {
		// Update listingsURL with pagination parameters
		listingsURL := fmt.Sprintf("https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?start=%d&limit=%d", start, limit)
		// Call the first endpoint
		listingsResp, err := callAPI(listingsURL, apiKey)
		if err != nil {
			log.Fatalf("Failed to call listings API: %v", err)
		}

		// Parse the response
		var listings ListingsResponse
		err = json.Unmarshal(listingsResp, &listings)
		if err != nil {
			log.Fatalf("Failed to unmarshal listings response: %v", err)
		}

		// Iterate through tokens
		time.Sleep(time.Second * 2)
		for _, token := range listings.Data {
			infoURL := fmt.Sprintf("https://pro-api.coinmarketcap.com/v2/cryptocurrency/info?symbol=%s", token.Symbol)

			time.Sleep(time.Second * 2)
			// Call the second endpoint for each symbol
			infoResp, err := callAPI(infoURL, apiKey)
			if err != nil {
				log.Printf("Failed to call info API for symbol %s: %v", token.Symbol, err)
				continue
			}

			// Parse the response
			var info InfoResponse
			err = json.Unmarshal(infoResp, &info)
			if err != nil {
				log.Printf("Failed to unmarshal info response for symbol %s: %v", token.Symbol, err)
				continue
			}

			// Check for Discord URLs in the chat array
			row := getCsvRowString(token.Symbol, info)
			if len(row) > 0 {
				writer.Write(row)
			}

		}
		// Break out of the loop if we've received fewer than 'limit' results
		if len(listings.Data) < limit {
			break
		}

		// Prepare for the next batch
		start += limit
		time.Sleep(1)
	}
}

func getCsvRowString(symbol string, info InfoResponse) []string {
	// Check for URLs in the various arrays
	for _, infos := range info.Data {
		for _, ci := range infos {
			// Initialize variables to hold URLs for each platform
			discordURL := ""
			redditURL := ""
			twitterURL := ""
			githubURL := ""

			// Extract Discord URL if available
			for _, url := range ci.Urls.Chat {
				if strings.Contains(url, "discord.com") {
					discordURL = url
					break // Assuming only one Discord URL is needed
				}
			}

			// Extract Reddit URL if available
			if len(ci.Urls.Reddit) > 0 {
				redditURL = ci.Urls.Reddit[0] // Assuming taking the first Reddit URL if multiple
			}

			// Extract Twitter URL if available
			if len(ci.Urls.Twitter) > 0 {
				twitterURL = ci.Urls.Twitter[0] // Assuming taking the first Twitter URL if multiple
			}

			// Extract Github URL if available
			if len(ci.Urls.Github) > 0 {
				for _, url := range ci.Urls.Github {
					if strings.Contains(url, "github.com") {
						githubURL = url
						break // Assuming only one Github URL is needed
					}
				}
			}

			// Write to CSV if at least one URL is present
			if discordURL != "" || redditURL != "" || twitterURL != "" || githubURL != "" {
				return []string{symbol, discordURL, twitterURL, redditURL, githubURL}
			}
		}
	}

	return []string{}
}

// Helper function to call the API
func callAPI(url string, apiKey string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
