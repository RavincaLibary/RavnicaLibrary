package cards

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// ScryfallClientInterface defines the methods the Scryfall client should implement
type ScryfallClientInterface interface {
	SearchCards(query string) (*ScryfallSearchResponse, error)
	GetCardByID(id string) (*ScryfallCard, error)
}

// ScryfallClient handles communication with the Scryfall API
type ScryfallClient struct {
	httpClient *http.Client
	baseURL    string
}

var _ ScryfallClientInterface = (*ScryfallClient)(nil)

// NewScryfallClient creates a new Scryfall API client
func NewScryfallClient() *ScryfallClient {
	return &ScryfallClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.scryfall.com",
	}
}

// SearchCards searches for cards matching the given query
func (c *ScryfallClient) SearchCards(query string) (*ScryfallSearchResponse, error) {
	// URL encode the query
	escapedQuery := url.QueryEscape(query)

	// Build the Scryfall API URL
	requestURL := fmt.Sprintf("%s/cards/search?q=%s", c.baseURL, escapedQuery)

	// Make the HTTP request
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error making request to Scryfall API: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		// Scryfall returns 404 when no cards match the search
		if resp.StatusCode == http.StatusNotFound {
			return &ScryfallSearchResponse{
				Object:     "list",
				TotalCards: 0,
				HasMore:    false,
				Data:       []ScryfallCard{},
			}, nil
		}

		// Try to parse error message
		var errorResp ScryfallErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("scryfall API error: %s", errorResp.Details)
		}

		return nil, fmt.Errorf("received non-200 response from Scryfall API: %d", resp.StatusCode)
	}

	// Parse the JSON response
	var searchResp ScryfallSearchResponse
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	return &searchResp, nil
}

// GetCardByID fetches a specific card by its Scryfall ID
func (c *ScryfallClient) GetCardByID(id string) (*ScryfallCard, error) {
	requestURL := fmt.Sprintf("%s/cards/%s", c.baseURL, id)

	// Make the HTTP request
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error making request to Scryfall API: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response from Scryfall API: %d", resp.StatusCode)
	}

	// Parse the JSON response
	var card ScryfallCard
	err = json.Unmarshal(body, &card)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	return &card, nil
}
