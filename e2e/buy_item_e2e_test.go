package e2e

import (
	"bytes"
	"encoding/json"
	"log"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"avito-shop-test/internal/adapter"
	"avito-shop-test/internal/middleware"
	"avito-shop-test/internal/models"

	"avito-shop-test/internal/handler"
	"avito-shop-test/internal/repository"
	"avito-shop-test/internal/usecase"

	"avito-shop-test/internal/token"
)

func TestBuyItemE2ESuccesful(t *testing.T) {

	db := setupTestDB()
	err := runMigration(db, "../db/migration.sql")
	if err != nil {
		log.Fatalf("failed to run migration: %v", err)
	}

	jwtSecret := []byte("1234")

	userRepo := repository.NewUSerRepository(db)
	coinTransactionRepo := repository.NewCoinTransactionRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	purchaseUC := usecase.NewPurchaseUseCase(purchaseRepo, userRepo, storeRepo)

	userUc := usecase.NewUserUsecase(userRepo, purchaseRepo, coinTransactionRepo, token.NewGenerator(jwtSecret))

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	apiGroup := router.Group("/api")
	ginRouter := adapter.NewGinRouter(apiGroup)
	handler.NewUserHandler(ginRouter, userUc, middleware.AuthMiddleware(jwtSecret))
	handler.NewPurchaseHandler(ginRouter, purchaseUC, middleware.AuthMiddleware(jwtSecret))

	user := models.User{Username: "user", Password: "password"}
	authRequestBody := map[string]string{"username": user.Username, "password": user.Password}
	body, _ := json.Marshal(authRequestBody)

	req1, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req1)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var authResponse struct {
		Token string `json:"token"`
	}
	json.Unmarshal(recorder.Body.Bytes(), &authResponse)

	reqTransfer, _ := http.NewRequest(http.MethodGet, "/api/buy/book", nil)
	reqTransfer.Header.Set("Authorization", "Bearer "+authResponse.Token)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, reqTransfer)
	assert.Equal(t, http.StatusOK, recorder.Code, "Ошибка покупки: %v, Response Body: %s", recorder.Code, recorder.Body.String())

	reqInfo, _ := http.NewRequest(http.MethodGet, "/api/info", nil)
	reqInfo.Header.Set("Authorization", "Bearer "+authResponse.Token)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, reqInfo)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var userInfo struct {
		Balance   int                    `json:"coins"`
		Inventory []models.PurchasedItem `json:"inventory"`
	}
	json.Unmarshal(recorder.Body.Bytes(), &userInfo)
	assert.Equal(t, 950, userInfo.Balance, "Баланс пользователя должен быть 950")
	assert.Equal(t, "book", userInfo.Inventory[0].ItemName, "Покупка пользователя должна быть книгой")
	assert.Equal(t, 1, userInfo.Inventory[0].Quantity, "Количество купленных предметов должно быть 1")

	db.Exec("TRUNCATE users, transactions RESTART IDENTITY CASCADE")
}
