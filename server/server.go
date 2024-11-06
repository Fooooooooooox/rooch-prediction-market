package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"github.com/go-macaron/toolbox"
	"github.com/rooch-prediction-market/backend/config"
	"github.com/rooch-prediction-market/backend/dtos"
	"github.com/rooch-prediction-market/backend/handlers"
	"github.com/rooch-prediction-market/backend/models"
	"github.com/rooch-prediction-market/backend/services"
	"gopkg.in/macaron.v1"
)

var srv *http.Server

func Initialize(addr string, handler http.Handler) {
	srv = &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}

func StartServer() {
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server start error: %v", err)
		}
	}()
}

// StopServer stops the server with a timeout
func StopServer() {
	fmt.Print("🎃🎃 StopServer: ")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	} else {
		log.Println("Server gracefully shut down")
	}
}

func CreateServer(modes ...string) *macaron.Macaron {
	mode := "" // Default mode
	if len(modes) > 0 {
		mode = modes[0] // Use the first mode if provided
	}

	c := config.New(mode, "")
	dbInstance := models.New(c)
	// create services
	marketService := services.NewMarketService(dbInstance.Pg)

	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	m.Use(toolbox.Toolboxer(m))
	m.Use(func(ctx *macaron.Context) {
		ctx.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, ngrok-skip-browser-warning")
		ctx.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Req.Method == "OPTIONS" {
			ctx.Status(http.StatusOK)
			return
		}

		ctx.Next()
	})

	m.Options("/*", func(ctx *macaron.Context) {
		ctx.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		ctx.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		ctx.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Status(http.StatusOK)
	})

	m.Use(DatabaseMiddleware(dbInstance))
	m.Use(ConfigMiddleware(c))
	m.Use(MarketServiceMiddleware(marketService))

	m.Get("api/v1/health", handlers.HealthHandler)
	m.Get("api/v1/healthDb", handlers.HealthHandlerDb)

	m.Post("api/v1/markets", binding.Bind(dtos.Market{}), handlers.CreateMarket)
	m.Get("api/v1/markets", handlers.GetMarkets)
	m.Get("api/v1/markets/:id", handlers.GetMarket)
	m.Put("api/v1/markets/:id", binding.Bind(dtos.UpdateMarket{}), handlers.UpdateMarket)

	m.Get("api/v1/trades/:id", handlers.GetTrades)
	m.Post("api/v1/trade", binding.Bind(dtos.Trade{}), handlers.CreateTrade)

	m.Get("api/v1/votes/:id", handlers.GetVotes)
	m.Post("api/v1/vote", binding.Bind(dtos.Vote{}), handlers.CreateVote)

	return m
}

func DatabaseMiddleware(dBInstance *models.DB) macaron.Handler {
	return func(ctx *macaron.Context) {
		ctx.Map(dBInstance)
	}
}

func ConfigMiddleware(cfg *config.Config) macaron.Handler {
	return func(ctx *macaron.Context) {
		ctx.Map(cfg)
	}
}

func MarketServiceMiddleware(marketService *services.MarketService) macaron.Handler {
	return func(ctx *macaron.Context) {
		ctx.Map(marketService)
	}
}
