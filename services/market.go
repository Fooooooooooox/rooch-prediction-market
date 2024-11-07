package services

import (
	"github.com/rooch-prediction-market/backend/models"
	"gorm.io/gorm"
)

type MarketService struct {
	Db                    *gorm.DB
	MarketRepo            *BaseService[models.Market]
	TradeRepo             *BaseService[models.Trade]
	VoteRepo              *BaseService[models.Vote]
	UserMarketBalanceRepo *BaseService[models.UserMarketBalance]
	UserBalanceRepo       *BaseService[models.UserBalance]
}

func NewMarketService(db *gorm.DB) *MarketService {
	return &MarketService{
		MarketRepo:            NewBaseService[models.Market](db),
		TradeRepo:             NewBaseService[models.Trade](db),
		VoteRepo:              NewBaseService[models.Vote](db),
		UserMarketBalanceRepo: NewBaseService[models.UserMarketBalance](db),
		UserBalanceRepo:       NewBaseService[models.UserBalance](db),
		Db:                    db,
	}
}

// UpdatePrices updates the prices of yes and no tokens based on the current pool amounts
func (ms *MarketService) UpdatePrices(market *models.Market) error {
	totalAmount := market.YesAmount + market.NoAmount
	if totalAmount == 0 {
		market.Prob = 0.5
		market.ProbNo = 0.5
	} else {
		// Here, we calculate the price as the ratio of the amounts
		k := float64(market.YesAmount) * float64(market.NoAmount)
		market.Prob = k / float64(market.NoAmount)
		market.ProbNo = k / float64(market.YesAmount)
	}

	if err := ms.Db.Save(market).Error; err != nil {
		return err
	}
	return nil
}

func (ms *MarketService) BetOnYes(userAddress string, marketID uint, amount uint) error {
	market, err := ms.MarketRepo.GetById(marketID)
	if err != nil {
		return err
	}
	market.YesAmount += amount
	if err := ms.UpdatePrices(market); err != nil {
		return err
	}

	// Check if userMarketBalance exists
	var userMarketBalance models.UserMarketBalance
	err = ms.UserMarketBalanceRepo.First(&userMarketBalance, "user_address = ? AND market_id = ? AND tick = ?", userAddress, marketID, "yes")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new UserMarketBalance if not found
			userMarketBalance = models.UserMarketBalance{
				UserAddress: userAddress,
				MarketID:    marketID,
				Tick:        "yes",
				Balance:     amount,
			}
			if err := ms.UserMarketBalanceRepo.Create(&userMarketBalance); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		userMarketBalance.Balance += amount
		if err := ms.Db.Save(&userMarketBalance).Error; err != nil {
			return err
		}
	}
	return nil
}

func (ms *MarketService) BetOnNo(userAddress string, marketID uint, amount uint) error {
	market, err := ms.MarketRepo.GetById(marketID)
	if err != nil {
		return err
	}
	market.NoAmount += amount
	if err := ms.UpdatePrices(market); err != nil {
		return err
	}

	// Check if userMarketBalance exists
	var userMarketBalance models.UserMarketBalance
	err = ms.UserMarketBalanceRepo.First(&userMarketBalance, "user_address = ? AND market_id = ? AND tick = ?", userAddress, marketID, "no")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new UserMarketBalance if not found
			userMarketBalance = models.UserMarketBalance{
				UserAddress: userAddress,
				MarketID:    marketID,
				Tick:        "no",
				Balance:     amount,
			}
			if err := ms.UserMarketBalanceRepo.Create(&userMarketBalance); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		userMarketBalance.Balance += amount
		if err := ms.Db.Save(&userMarketBalance).Error; err != nil {
			return err
		}
	}

	return nil
}

func (ms *MarketService) CalculateClaimableAmount(address string, market *models.Market) (uint, error) {
	userMarketBalanceYesList, err := ms.UserMarketBalanceRepo.Find(map[string]interface{}{"user_address": address, "market_id": market.ID, "tick": "yes"})
	if err != nil {
		return 0, err
	}
	var userMarketBalanceYes models.UserMarketBalance
	if len(userMarketBalanceYesList) > 0 {
		userMarketBalanceYes = userMarketBalanceYesList[0]
	}

	userMarketBalanceNoList, err := ms.UserMarketBalanceRepo.Find(map[string]interface{}{"user_address": address, "market_id": market.ID, "tick": "no"})
	if err != nil {
		return 0, err
	}
	var userMarketBalanceNo models.UserMarketBalance
	if len(userMarketBalanceNoList) > 0 {
		userMarketBalanceNo = userMarketBalanceNoList[0]
	}

	var claimableAmount uint
	if market.Result { // yes wins
		claimableAmount += market.NoAmount * userMarketBalanceYes.Balance / market.YesAmount
	} else {
		claimableAmount += market.YesAmount * userMarketBalanceNo.Balance / market.NoAmount
	}

	return claimableAmount, nil
}
