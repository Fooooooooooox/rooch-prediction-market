package handlers

import (
	"fmt"
	"os"
	"testing"

	"github.com/rooch-prediction-market/backend/config"
	"github.com/rooch-prediction-market/backend/models"
	"github.com/rooch-prediction-market/backend/services"
	"github.com/rooch-prediction-market/backend/testutils"
	"gopkg.in/macaron.v1"
)

var (
	dsn           string
	db            *models.DB
	marketService *services.MarketService
	cfg           *config.Config
)

func setup() {
	dsn = testutils.GetTestDSN()
	var err error
	db, err = testutils.SetupTestDB(dsn)
	testutils.AssertNil(err)

	marketService = services.NewMarketService(db.Pg)

	cfg = config.New("", "")

	fmt.Println("👀👀👀👀 this is cfg.TestMode: ", cfg.TestMode)
}

func TestMain(m *testing.M) {
	// Define a custom test runner
	runTests := func() int {
		setup() // Call setup before running tests
		return m.Run()
	}

	// Execute the custom test runner
	code := runTests()
	os.Exit(code)
}

func TestGetUserInfo(t *testing.T) {
}

func MarketServiceMiddleware(marketService *services.MarketService) macaron.Handler {
	return func(ctx *macaron.Context) {
		ctx.Map(marketService)
	}
}

func MockClerkMiddleware(ctx *macaron.Context) {
}
