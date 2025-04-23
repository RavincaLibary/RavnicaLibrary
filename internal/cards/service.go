package cards

import (
	"fmt"
)

// CardService handles card-related business logic
type CardService struct {
	scryfallClient *ScryfallClient
}

// NewCardService creates a new card service
func NewCardService(client ScryfallClientInterface) *CardService {
	if client == nil {
		client = NewScryfallClient()
	}
	return &CardService{
		scryfallClient: NewScryfallClient(),
	}
}

// SearchCards searches for cards based on the provided query
func (s *CardService) SearchCards(query string) (*SearchResult, error) {
	// Validate query
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Call Scryfall API
	resp, err := s.scryfallClient.SearchCards(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search cards: %w", err)
	}

	// Transform to our response format
	result := &SearchResult{
		Cards:      resp.Data,
		TotalCards: resp.TotalCards,
		HasMore:    resp.HasMore,
	}

	return result, nil
}

// GetCardByID retrieves a specific card by its Scryfall ID
func (s *CardService) GetCardByID(id string) (*ScryfallCard, error) {
	if id == "" {
		return nil, fmt.Errorf("card ID cannot be empty")
	}

	return s.scryfallClient.GetCardByID(id)
}
