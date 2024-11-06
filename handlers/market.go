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
		Price:         0.5,
		PriceNo:       0.5,
		VoteYesAmount: 0,
		VoteNoAmount:  0,
		Status:        models.MarketStatusVoting,
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

func CreateTrade(ctx *macaron.Context, req dtos.Trade, marketService *services.MarketService) {
	// check if market exists
	market, err := marketService.MarketRepo.GetById(req.MarketID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if market.Status != models.MarketStatusOpen {
		ctx.JSON(http.StatusBadRequest, "Market is not open")
		return
	}

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
			err = marketService.BuyYesToken(req.MarketID, req.Amount)
		} else {
			err = marketService.BuyNoToken(req.MarketID, req.Amount)
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, trade)
}

func GetTrades(ctx *macaron.Context, marketService *services.MarketService) {
	marketID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid market ID")
		return
	}

	trades, err := marketService.TradeRepo.Find(map[string]interface{}{"market_id": marketID})
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

	if market.Status != models.MarketStatusVoting {
		ctx.JSON(http.StatusBadRequest, "Market is not voting")
		return
	}

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

	err = marketService.MarketRepo.Update(func() models.Market { return *market }, "vote_yes_amount", "vote_no_amount")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, vote)
}

func GetVotes(ctx *macaron.Context, marketService *services.MarketService) {
	marketID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid market ID")
		return
	}

	votes, err := marketService.VoteRepo.Find(map[string]interface{}{"market_id": marketID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, votes)
}