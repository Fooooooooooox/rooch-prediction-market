package models

import "time"

type MarketStatus string

const (
	MarketStatusOpen   MarketStatus = "open"
	MarketStatusClosed MarketStatus = "closed"
	MarketStatusVoting MarketStatus = "voting"
)

type Market struct {
	BaseModel
	Title            string       `json:"title"`
	Description      string       `json:"description"`
	YesAmount        uint         `json:"yes_amount"`
	NoAmount         uint         `json:"no_amount"`
	Prob             float64      `json:"price"`    // price of yes
	ProbNo           float64      `json:"price_no"` // price of no
	VoteYesAmount    uint         `json:"vote_yes_amount"`
	VoteNoAmount     uint         `json:"vote_no_amount"`
	Status           MarketStatus `json:"status"`
	Result           bool         `json:"result"` // true if yes, false if no
	JudgementStartAt time.Time    `json:"judgement_start_at"`
	JudgementEndAt   time.Time    `json:"judgement_end_at"`
}

type Trade struct {
	BaseModel
	MarketID uint   `json:"market_id"`
	Address  string `json:"address"`
	Side     string `json:"side"` // buy or sell
	Tick     string `json:"tick"` // yes or no
	Amount   uint   `json:"amount"`
}

type Vote struct {
	BaseModel
	MarketID uint   `json:"market_id"`
	Address  string `json:"address"`
	Tick     string `json:"tick"` // yes or no
	Amount   uint   `json:"amount"`
}

type UserBalance struct {
	BaseModel
	Address string `json:"address"`
	Balance uint   `json:"balance"`
}

type UserMarketBalance struct {
	BaseModel
	UserAddress string `json:"user_address"`
	MarketID    uint   `json:"market_id"`
	Tick        string `json:"tick"` // yes or no
	Balance     uint   `json:"balance"`
}
