// internal/cards/handler_test.go
package cards

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_SearchCards(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		query          string
		setupMock      func(*MockScryfallClient)
		expectedStatus int
		checkBody      bool
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
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "empty query",
			query:          "",
			setupMock:      nil, // No mock needed, fails at validation
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
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

			// Create handler with service
			handler := &Handler{
				service: service,
			}

			// Create a test router
			router := mux.NewRouter()
			handler.RegisterRoutes(router)

			// Create a test server
			server := httptest.NewServer(router)
			defer server.Close()

			// Create the request URL
			url := server.URL + "/api/cards/search"
			if tc.query != "" {
				url += "?q=" + tc.query
			}

			// Make the request
			resp, err := http.Get(url)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Check status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// Check response body for successful requests
			if tc.checkBody && resp.StatusCode == http.StatusOK {
				var result SearchResult
				err := json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)

				// Verify the cards
				assert.Equal(t, 1, len(result.Cards))
				if len(result.Cards) > 0 {
					assert.Equal(t, "Lightning Bolt", result.Cards[0].Name)
				}
			}

			// Verify all expectations were met
			mockClient.AssertExpectations(t)
		})
	}
}
