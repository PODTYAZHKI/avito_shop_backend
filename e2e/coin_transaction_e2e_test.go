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



func TestTransferCoinsE2ESuccesful(t *testing.T) {

	db := setupTestDB()
	

	jwtSecret := []byte("1234")

	userRepo := repository.NewUSerRepository(db)
	transactionRepo := repository.NewCoinTransactionRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)

	transactionUC := usecase.NewCoinTransactionUseCase(transactionRepo, userRepo)

	userUc := usecase.NewUserUsecase(userRepo, purchaseRepo, transactionRepo, token.NewGenerator(jwtSecret))

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	apiGroup := router.Group("/api")
	ginRouter := adapter.NewGinRouter(apiGroup)
	handler.NewUserHandler(ginRouter, userUc, middleware.AuthMiddleware(jwtSecret))

	handler.NewCoinTransactionHandler(ginRouter, transactionUC, middleware.AuthMiddleware(jwtSecret))

	user1 := models.User{Username: "user1", Password: "password1"}
	user2 := models.User{Username: "user2", Password: "password2"}

	authRequestBody := map[string]string{"username": user1.Username, "password": user1.Password}
	body, _ := json.Marshal(authRequestBody)

	req1, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req1)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var authResponse1 struct {
		Token string `json:"token"`
	}
	json.Unmarshal(recorder.Body.Bytes(), &authResponse1)


	authRequestBody = map[string]string{"username": user2.Username, "password": user2.Password}
	body, _ = json.Marshal(authRequestBody)

	req2, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req2)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var authResponse2 struct {
		Token string `json:"token"`
	}
	json.Unmarshal(recorder.Body.Bytes(), &authResponse2)

	log.Println("authResponse2", authResponse2.Token)

	requestBody := models.SendCoinRequest{ToUser: user2.Username, Amount: 30}
	body, _ = json.Marshal(requestBody)
	log.Println("requestBody", string(body))

	reqTransfer, _ := http.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBuffer(body))
	reqTransfer.Header.Set("Authorization", "Bearer "+authResponse1.Token)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, reqTransfer)
	if recorder.Code != http.StatusOK {
		t.Logf("Ошибка отправки монет: %d, Response Body: %s", recorder.Code, recorder.Body.String())
	}

	assert.Equal(t, http.StatusOK, recorder.Code)

	reqInfo, _ := http.NewRequest(http.MethodGet, "/api/info", nil)
	reqInfo.Header.Set("Authorization", "Bearer "+authResponse1.Token)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, reqInfo)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var userInfo struct {
		Balance int `json:"coins"`
	}
	json.Unmarshal(recorder.Body.Bytes(), &userInfo)
	log.Println("userInfo", userInfo)

	assert.Equal(t, 970, userInfo.Balance, "Баланс пользователя 1 должен быть 970")

	reqInfo, _ = http.NewRequest(http.MethodGet, "/api/info", nil)
	reqInfo.Header.Set("Authorization", "Bearer "+authResponse2.Token)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, reqInfo)
	assert.Equal(t, http.StatusOK, recorder.Code)

	json.Unmarshal(recorder.Body.Bytes(), &userInfo)

	assert.Equal(t, 1030, userInfo.Balance, "Баланс пользователя 2 должен быть 1030")

	db.Exec("TRUNCATE users, transactions RESTART IDENTITY CASCADE")
}
