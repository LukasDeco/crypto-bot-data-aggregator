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
	Chat []string `json:"chat"`
}

func TestCoinMarketCapGetDiscords(t *testing.T) {
	apiKey := ""

	// Define the first URL
	listingsURL := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"

	// Open a new CSV file for writing
	file, err := os.Create("discord_urls.csv")
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	if err := writer.Write([]string{"Token Symbol", "Discord URL"}); err != nil {
		log.Fatalf("Failed to write CSV header: %v", err)
	}

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
		for _, infos := range info.Data {
			for _, ci := range infos {
				for _, chatURL := range ci.Urls.Chat {
					if strings.Contains(chatURL, "discord.com") {
						if err := writer.Write([]string{token.Symbol, chatURL}); err != nil {
							log.Fatalf("Failed to write to CSV: %v", err)
						}
						fmt.Println("Wrote to CSV:", token.Symbol, chatURL)
					}
				}
			}
		}

	}
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
