package dto

// ApplyPromoCodeRequest represents the request body for applying a promo code
type ApplyPromoCodeRequest struct {
	Promocode string `json:"promocode"`
	RequestID string `json:"requestID"`
}

// ApplyPromoCodeResponse represents the response when promo code is successfully applied
type ApplyPromoCodeResponse struct {
	Promocode     string `json:"promocode"`
	OriginalPrice int    `json:"original_price"`
	FinalPrice    int    `json:"final_price"`
	Currency      string `json:"currency"`
}

// PromoCodeErrorResponse represents an error response for promo code operations
type PromoCodeErrorResponse struct {
	Error string `json:"error"`
}
