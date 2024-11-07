package handlers

import (
	"net/http"
	"strconv"

	"github.com/rooch-prediction-market/backend/dtos"
	"github.com/rooch-prediction-market/backend/models"
	"github.com/rooch-prediction-market/backend/services"
	"gopkg.in/macaron.v1"
)

func CreateMarket(ctx *macaron.Context, req dtos.Market, marketService *services.MarketService) {
	market := models.Market{
		Title:         req.Title,
		Description:   req.Description,
		YesAmount:     100,
		NoAmount:      100,
		Prob:          1,
		ProbNo:        1,
		VoteYesAmount: 0,
		VoteNoAmount:  0,
		Status:        models.MarketStatusOpen,
	}

	err := marketService.MarketRepo.Create(&market)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, market)
}

func GetMarkets(ctx *macaron.Context, marketService *services.MarketService) {
	markets, err := marketService.MarketRepo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, markets)
}

func GetMarket(ctx *macaron.Context, marketService *services.MarketService) {
	marketID := ctx.Params("marketId")
	marketIDInt, err := strconv.Atoi(marketID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid market ID")
		return
	}
	market, err := marketService.MarketRepo.GetById(uint(marketIDInt))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, market)
}

func UpdateMarket(ctx *macaron.Context, req dtos.UpdateMarket, marketService *services.MarketService) {
	market, err := marketService.MarketRepo.GetById(req.MarketID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	market.Title = req.Title
	market.Description = req.Description
	market.Status = models.MarketStatus(req.Status)
	market.Prob = req.Price
	market.ProbNo = req.PriceNo
	market.YesAmount = req.YesAmount
	market.NoAmount = req.NoAmount
	market.VoteYesAmount = req.VoteYesAmount
	market.VoteNoAmount = req.VoteNoAmount

	if err := marketService.Db.Save(&market).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, market)
}

func SettleMarket(ctx *macaron.Context, req dtos.SettleMarket, marketService *services.MarketService) {
	market, err := marketService.MarketRepo.GetById(req.MarketID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if market.VoteYesAmount > market.VoteNoAmount {
		market.Result = true
	} else {
		market.Result = false
	}

	market.Status = models.MarketStatusClosed

	err = marketService.Db.Save(&market).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, market)
}

func CreateTrade(ctx *macaron.Context, req dtos.Trade, marketService *services.MarketService) {
	// check if market exists
	market, err := marketService.MarketRepo.GetById(req.MarketID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// if market.Status != models.MarketStatusOpen {
	// 	ctx.JSON(http.StatusBadRequest, "Market is not open")
	// 	return
	// }

	trade := models.Trade{
		MarketID: market.ID,
		Address:  req.Address,
		Side:     req.Side,
		Tick:     req.Tick,
		Amount:   req.Amount,
	}

	err = marketService.TradeRepo.Create(&trade)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Update market amounts based on trade
	if req.Side == "buy" {
		if req.Tick == "yes" {
			err = marketService.BetOnYes(req.Address, req.MarketID, req.Amount)
		} else {
			err = marketService.BetOnNo(req.Address, req.MarketID, req.Amount)
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, trade)
}

func GetTrades(ctx *macaron.Context, marketService *services.MarketService) {
	marketID, err := strconv.Atoi(ctx.Params("marketId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid market ID")
		return
	}

	trades, err := marketService.TradeRepo.Find(map[string]interface{}{"marketId": marketID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, trades)
}

func CreateVote(ctx *macaron.Context, req dtos.Vote, marketService *services.MarketService) {
	market, err := marketService.MarketRepo.GetById(req.MarketID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// if market.Status != models.MarketStatusVoting {
	// 	ctx.JSON(http.StatusBadRequest, "Market is not voting")
	// 	return
	// }

	vote := models.Vote{
		MarketID: market.ID,
		Address:  req.Address,
		Tick:     req.Tick,
		Amount:   req.Amount,
	}

	err = marketService.VoteRepo.Create(&vote)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if req.Tick == "yes" {
		market.VoteYesAmount += req.Amount
	} else {
		market.VoteNoAmount += req.Amount
	}

	err = marketService.Db.Save(&market).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, vote)
}

func GetVotes(ctx *macaron.Context, marketService *services.MarketService) {
	marketID, err := strconv.Atoi(ctx.Params("marketId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid market ID")
		return
	}

	votes, err := marketService.VoteRepo.Find(map[string]interface{}{"marketId": marketID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, votes)
}

func GetClaimableAmount(ctx *macaron.Context, marketService *services.MarketService) {
	address := ctx.Params("address")
	marketID, err := strconv.Atoi(ctx.Params("marketId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid market ID")
		return
	}

	market, err := marketService.MarketRepo.GetById(uint(marketID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	// if market.Status != models.MarketStatusClosed {
	// 	ctx.JSON(http.StatusBadRequest, "Market is not setteled")
	// 	return
	// }

	claimableAmount, err := marketService.CalculateClaimableAmount(address, market)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, claimableAmount)
}

func ClaimReward(ctx *macaron.Context, req dtos.Claim, marketService *services.MarketService) {
	market, err := marketService.MarketRepo.GetById(req.MarketID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	claimableAmount, err := marketService.CalculateClaimableAmount(req.Address, market)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// increase user balance
	userBalances, err := marketService.UserBalanceRepo.Find(map[string]interface{}{"address": req.Address})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if len(userBalances) == 0 {
		userBalance := models.UserBalance{
			BaseModel: models.BaseModel{
				ID: userBalances[0].ID,
			},
			Address: req.Address,
			Balance: userBalances[0].Balance + claimableAmount,
		}

		err = marketService.Db.Save(&userBalance).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, claimableAmount)
}
