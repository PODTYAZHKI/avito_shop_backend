package e2e

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {

	dbName := "test_db"
	dbUser := "test"
	dbPassword := "test"
	dbHost := "localhost"
	dbPort := "5433"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к тестовой базе данных: %v", err)
	}

	return db
}
func runMigration(db *gorm.DB, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := db.Exec(string(data)).Error; err != nil {

		if errors.Is(err, gorm.ErrDuplicatedKey) {

			return nil
		}
		return err
	}

	return nil
}
