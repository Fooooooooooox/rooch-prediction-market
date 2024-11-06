package dtos

type Trade struct {
	Address  string `json:"address" binding:"required"`
	MarketID uint   `json:"market_id" binding:"required"`
	Side     string `json:"side" binding:"required"` // buy or sell
	Tick     string `json:"tick" binding:"required"` // yes or no
	Amount   uint   `json:"amount" binding:"required"`
}

type Vote struct {
	Address  string `json:"address" binding:"required"`
	MarketID uint   `json:"market_id" binding:"required"`
	Tick     string `json:"tick" binding:"required"` // yes or no
	Sig      string `json:"sig" binding:"required"`
	Amount   uint   `json:"amount" binding:"required"`
}

type Market struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}
