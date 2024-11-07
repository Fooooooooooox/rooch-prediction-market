package testutils

import (
	"log"

	"github.com/rooch-prediction-market/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DropTablesAndReCreate(db *models.DB) error {
	tables := []string{
		"machine_infos",
	}

	for _, table := range tables {
		err := db.Pg.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error // Ensure CASCADE is used

		if err != nil {
			log.Printf("Failed to drop table %s: %v", table, err)
			return err
		} else {
			log.Printf("Table %s dropped successfully", table)
		}
	}

	db.Pg.AutoMigrate(&models.Trade{}, &models.Vote{}, &models.UserMarketBalance{}, &models.Market{})

	log.Printf("auto migrate success\n")
	return nil
}

func DeleteTables(db *models.DB) error {
	// List all the tables you want to delete
	tables := []interface{}{
		&models.Trade{},
		&models.Vote{},
	}

	for _, model := range tables {
		// Using Delete method with a model type and no where clause
		// Passing a slice of model ensures that Delete affects all records
		err := db.Pg.Unscoped().Where("1 = 1").Delete(model).Error

		if err != nil {
			log.Printf("Failed to delete contents of table %T: %v", model, err)
			return err
		} else {
			log.Printf("Contents of table %T deleted successfully", model)
		}
	}

	return nil
}

func InitUserLevelData(db *models.DB) error {
	return nil
}

func SetupTestDB(dsn string) (*models.DB, error) {
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	testDB := &models.DB{Pg: db}
	if err := DropTablesAndReCreate(testDB); err != nil {
		return nil, err
	}

	return testDB, nil
}
