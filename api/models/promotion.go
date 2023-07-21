package api

type CreatePromotionRequest struct {
	PromotionId    string  `json:"promotion_id"`
	Price          float64 `json:"price"`
	ExpirationDate string  `json:"expiration_date"`
}

type CreatePromotionResponse struct {
	Message           string               `json:"message"`
	PromotionResponse GetPromotionResponse `json:"response"`
}

type GetPromotionResponse struct {
	ID             int64   `json:"id"`
	PromotionId    string  `json:"promotion_id"`
	Price          float64 `json:"price"`
	ExpirationDate string  `json:"expiration_date"`
}
