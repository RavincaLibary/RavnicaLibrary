package cards

// ScryfallCard represents the card data structure returned by Scryfall
type ScryfallCard struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	ImageURIs  *ImageURI  `json:"image_uris,omitempty"`
	ManaCost   string     `json:"mana_cost,omitempty"`
	TypeLine   string     `json:"type_line,omitempty"`
	OracleText string     `json:"oracle_text,omitempty"`
	Power      string     `json:"power,omitempty"`
	Toughness  string     `json:"toughness,omitempty"`
	Rarity     string     `json:"rarity,omitempty"`
	Set        string     `json:"set,omitempty"`
	SetName    string     `json:"set_name,omitempty"`
	Prices     *Prices    `json:"prices,omitempty"`
	CardFaces  []CardFace `json:"card_faces,omitempty"` // For double-faced cards
}

// ImageURI contains URLs to card images of different sizes
type ImageURI struct {
	Small      string `json:"small,omitempty"`
	Normal     string `json:"normal,omitempty"`
	Large      string `json:"large,omitempty"`
	PNG        string `json:"png,omitempty"`
	ArtCrop    string `json:"art_crop,omitempty"`
	BorderCrop string `json:"border_crop,omitempty"`
}

// CardFace represents one face of a double-faced card
type CardFace struct {
	Name       string    `json:"name,omitempty"`
	TypeLine   string    `json:"type_line,omitempty"`
	OracleText string    `json:"oracle_text,omitempty"`
	ManaCost   string    `json:"mana_cost,omitempty"`
	ImageURIs  *ImageURI `json:"image_uris,omitempty"`
	Power      string    `json:"power,omitempty"`
	Toughness  string    `json:"toughness,omitempty"`
}

// Prices contains pricing information for a card
type Prices struct {
	USD     string `json:"usd,omitempty"`
	USDFoil string `json:"usd_foil,omitempty"`
	EUR     string `json:"eur,omitempty"`
	EURFoil string `json:"eur_foil,omitempty"`
	TIX     string `json:"tix,omitempty"`
}

// ScryfallSearchResponse represents the response from a search query
type ScryfallSearchResponse struct {
	Object     string         `json:"object"`
	TotalCards int            `json:"total_cards"`
	HasMore    bool           `json:"has_more"`
	NextPage   string         `json:"next_page,omitempty"`
	Data       []ScryfallCard `json:"data"`
}

// ScryfallErrorResponse represents an error response from Scryfall
type ScryfallErrorResponse struct {
	Object  string `json:"object"`
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Details string `json:"details"`
}

// SearchResult represents our API's search response
type SearchResult struct {
	Cards      []ScryfallCard `json:"cards"`
	TotalCards int            `json:"total_cards"`
	HasMore    bool           `json:"has_more"`
}
