package services

import (
	"github.com/rooch-prediction-market/backend/models"
	"gorm.io/gorm"
)

type MarketService struct {
	Db         *gorm.DB
	MarketRepo *BaseService[models.Market]
	TradeRepo  *BaseService[models.Trade]
	VoteRepo   *BaseService[models.Vote]
}

func NewMarketService(db *gorm.DB) *MarketService {
	return &MarketService{
		MarketRepo: NewBaseService[models.Market](db),
		TradeRepo:  NewBaseService[models.Trade](db),
		VoteRepo:   NewBaseService[models.Vote](db),
		Db:         db,
	}
}

// UpdatePrices updates the prices of yes and no tokens based on the current pool amounts
func (ms *MarketService) UpdatePrices(market *models.Market) error {
	totalAmount := market.YesAmount + market.NoAmount
	if totalAmount == 0 {
		market.Price = 0.5
		market.PriceNo = 0.5
	} else {
		market.Price = float64(market.YesAmount) / float64(totalAmount)
		market.PriceNo = float64(market.NoAmount) / float64(totalAmount)
	}

	err := ms.MarketRepo.Update(func() models.Market { return *market }, "price", "price_no")
	if err != nil {
		return err
	}
	return nil
}

// BuyYesToken simulates buying yes tokens and updates the pool
func (ms *MarketService) BuyYesToken(marketID uint, amount uint) error {
	market, err := ms.MarketRepo.GetById(marketID)
	if err != nil {
		return err
	}
	market.YesAmount += amount
	return ms.UpdatePrices(market)
}

// BuyNoToken simulates buying no tokens and updates the pool
func (ms *MarketService) BuyNoToken(marketID uint, amount uint) error {
	market, err := ms.MarketRepo.GetById(marketID)
	if err != nil {
		return err
	}
	market.NoAmount += amount
	return ms.UpdatePrices(market)
}
