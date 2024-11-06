package services

import (
	"os"
	"testing"

	"github.com/rooch-prediction-market/backend/models"
	"github.com/rooch-prediction-market/backend/testutils"
)

var (
	dsn           string
	db            *models.DB
	marketService *MarketService
)

func setup() {
	dsn = testutils.GetTestDSN()
	var err error
	db, err = testutils.SetupTestDB(dsn)
	testutils.AssertNil(err)

	err = testutils.InitUserLevelData(db)
	testutils.AssertNil(err)

	marketService = NewMarketService(db.Pg)
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
