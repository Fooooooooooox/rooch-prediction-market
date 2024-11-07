package handlers

import (
	"net/http"

	"github.com/rooch-prediction-market/backend/dtos"
	"github.com/rooch-prediction-market/backend/models"
	"gopkg.in/macaron.v1"
)

func GetUserBalance(ctx *macaron.Context, db *models.DB) {
	address := ctx.Params("address")

	var userBalance models.UserBalance
	if err := db.Pg.Where("address = ?", address).First(&userBalance).Error; err != nil {
		ctx.JSON(http.StatusNotFound, "User not found")
		return
	}
	ctx.JSON(http.StatusOK, userBalance)
}

func AddUserBalance(ctx *macaron.Context, db *models.DB, req dtos.UpdateBalanceRequest) {
	var userBalance models.UserBalance
	if err := db.Pg.Where("address = ?", req.Address).First(&userBalance).Error; err != nil {
		// Create new user if not found
		userBalance = models.UserBalance{
			Address: req.Address,
			Balance: 0,
		}
		if err := db.Pg.Create(&userBalance).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, "Failed to create user")
			return
		}
	}

	userBalance.Balance += req.Amount
	if err := db.Pg.Save(&userBalance).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, "Failed to update balance")
		return
	}

	ctx.JSON(http.StatusOK, userBalance)
}

func DecreaseUserBalance(ctx *macaron.Context, db *models.DB, req dtos.UpdateBalanceRequest) {
	var userBalance models.UserBalance
	if err := db.Pg.Where("address = ?", req.Address).First(&userBalance).Error; err != nil {
		// Create new user if not found
		userBalance = models.UserBalance{
			Address: req.Address,
			Balance: 0,
		}
		if err := db.Pg.Create(&userBalance).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, "Failed to create user")
			return
		}
	}

	if userBalance.Balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, "Insufficient balance")
		return
	}

	userBalance.Balance -= req.Amount
	if err := db.Pg.Save(&userBalance).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, "Failed to update balance")
		return
	}

	ctx.JSON(http.StatusOK, userBalance)
}
