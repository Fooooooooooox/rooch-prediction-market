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
		market.Price = 1
		market.PriceNo = 1
	} else {
		// Here, we calculate the price as the ratio of the amounts
		k := float64(market.YesAmount) * float64(market.NoAmount)
		market.Price = k / float64(market.NoAmount)
		market.PriceNo = k / float64(market.YesAmount)
	}

	if err := ms.Db.Save(market).Error; err != nil {
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
