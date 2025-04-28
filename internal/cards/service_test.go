package cards

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScryfallClient is our mocked client
type MockScryfallClient struct {
	mock.Mock
}

// Implement ScryfallClientInterface methods
func (m *MockScryfallClient) SearchCards(query string) (*ScryfallSearchResponse, error) {
	args := m.Called(query)

	// Handle the case where the first return value might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ScryfallSearchResponse), args.Error(1)
}

func (m *MockScryfallClient) GetCardByID(id string) (*ScryfallCard, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ScryfallCard), args.Error(1)
}

func TestCardService_SearchCards(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		query         string
		setupMock     func(*MockScryfallClient)
		expectedError bool
		expectedCards int
	}{
		{
			name:  "successful search",
			query: "lightning bolt",
			setupMock: func(m *MockScryfallClient) {
				m.On("SearchCards", "lightning bolt").Return(
					&ScryfallSearchResponse{
						Object:     "list",
						TotalCards: 1,
						HasMore:    false,
						Data: []ScryfallCard{
							{ID: "1", Name: "Lightning Bolt"},
						},
					},
					nil,
				)
			},
			expectedError: false,
			expectedCards: 1,
		},
		{
			name:  "empty query",
			query: "",
			setupMock: func(m *MockScryfallClient) {
				// No mock setup needed for empty query, it should fail before API call
			},
			expectedError: true,
			expectedCards: 0,
		},
		{
			name:  "API error",
			query: "error",
			setupMock: func(m *MockScryfallClient) {
				m.On("SearchCards", "error").Return(
					nil,
					errors.New("API error"),
				)
			},
			expectedError: true,
			expectedCards: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock
			mockClient := new(MockScryfallClient)

			// Setup mock expectations
			if tc.setupMock != nil {
				tc.setupMock(mockClient)
			}

			// Create service with mock
			service := NewCardService(mockClient)

			// Call the function
			result, err := service.SearchCards(tc.query)

			// Check error expectations
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Check result expectations
			if !tc.expectedError {
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedCards, len(result.Cards))

				// For successful "lightning bolt" search, check card name
				if tc.query == "lightning bolt" {
					assert.Equal(t, "Lightning Bolt", result.Cards[0].Name)
				}
			}

			// Verify all expectations were met
			mockClient.AssertExpectations(t)
		})
	}
}
