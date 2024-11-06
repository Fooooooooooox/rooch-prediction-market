package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/rooch-prediction-market/backend/models"
	"github.com/rooch-prediction-market/backend/services"
	"gopkg.in/macaron.v1"
)

type Claims struct {
	UserId        uint   `json:"userid"`
	UserTwitterId string `json:"usertwitterid"`
	jwt.StandardClaims
}

func GetPaginationParams(ctx *macaron.Context) services.PaginationParams {
	page := ctx.QueryInt("page")
	if page < 1 {
		page = 1
	}

	pageSize := ctx.QueryInt("page_size")
	if pageSize < 1 {
		pageSize = 10 // default page size
	}

	return services.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

func HealthHandlerDb(ctx *macaron.Context, logger *log.Logger, db *models.DB) {
	err := db.Ping(context.Background())
	if err != nil {
		// If pinging the database fails, log the error
		log.Println("Failed to connect to the database:", err)
	} else {
		// If pinging is successful, log a success message
		log.Println("Successfully connected to the database.")
		ctx.JSON(http.StatusOK, "success")
	}
}

func HealthHandler(ctx *macaron.Context, logger *log.Logger, db *models.DB) {
	ctx.JSON(http.StatusOK, "success")
}
