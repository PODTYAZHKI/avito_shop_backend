package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"avito-shop-test/config"
	"avito-shop-test/internal/adapter"
	"avito-shop-test/internal/handler"
	"avito-shop-test/internal/middleware"
	"avito-shop-test/internal/repository"
	"avito-shop-test/internal/token"
	"avito-shop-test/internal/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	dbCongig := config.DBConfig()

	db, err := gorm.Open(postgres.Open(dbCongig), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	log.Println("Подключение к базе данных успешно!!!")

	serverAddress := os.Getenv("SERVER_ADDRESS")
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(200)

	err = runMigration(db, "db/migration.sql")
	if err != nil {
		log.Fatalf("failed to run migration: %v", err)
	}

	storeRepo := repository.NewStoreRepository(db)

	userRepo := repository.NewUSerRepository(db)

	transactionRepo := repository.NewCoinTransactionRepository(db)
	transactionUC := usecase.NewCoinTransactionUseCase(transactionRepo, userRepo)

	purchaseRepo := repository.NewPurchaseRepository(db)
	purchaseUC := usecase.NewPurchaseUseCase(purchaseRepo, userRepo, storeRepo)

	userUC := usecase.NewUserUsecase(userRepo, purchaseRepo, transactionRepo, token.NewGenerator(jwtSecret))

	router := gin.Default()
	apiGroup := router.Group("/api")
	ginRouter := adapter.NewGinRouter(apiGroup)

	handler.NewUserHandler(ginRouter, userUC, middleware.AuthMiddleware(jwtSecret))
	handler.NewCoinTransactionHandler(ginRouter, transactionUC, middleware.AuthMiddleware(jwtSecret))
	handler.NewPurchaseHandler(ginRouter, purchaseUC, middleware.AuthMiddleware(jwtSecret))

	srv := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Ошибка запуска сервера: %s", err)
		}
	}()
	log.Println("Сервер запущен на", serverAddress)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Завершение работы сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка аварийного завершения: %s", err)
	}
	log.Println("Сервер остановлен")
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
