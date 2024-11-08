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

type Claim struct {
	Address  string `json:"address" binding:"required"`
	MarketID uint   `json:"market_id" binding:"required"`
}

type Market struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateMarket struct {
	MarketID      uint    `json:"market_id" binding:"required"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Status        string  `json:"status"`
	Price         float64 `json:"price"`
	PriceNo       float64 `json:"price_no"`
	YesAmount     uint    `json:"yes_amount"`
	NoAmount      uint    `json:"no_amount"`
	VoteYesAmount uint    `json:"vote_yes_amount"`
	VoteNoAmount  uint    `json:"vote_no_amount"`
}

type SettleMarket struct {
	MarketID uint `json:"market_id" binding:"required"`
}

type UpdateBalanceRequest struct {
	Address string `json:"address" binding:"required"`
	Amount  uint   `json:"amount" binding:"required"`
}
