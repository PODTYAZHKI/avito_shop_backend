package handler

import (
	"net/http"

	"avito-shop-test/internal/models"
)

type CoinTransactionUseCase interface {
	SendCoins(fromUser, toUser string, amount int) error
}

type coinTransactionDelivery struct {
	coinTransactionUC CoinTransactionUseCase
}

func (d *coinTransactionDelivery) SendCoin(c Context) {
	var requestBody models.SendCoinRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"Errors": "Неверный запрос"})
		return
	}
	if requestBody.Amount < 0 {
		c.JSON(http.StatusBadRequest, map[string]string{"Errors": "Сумма не может быть отрицательной"})
		return
	}

	username := c.MustGet("username").(string)

	err := d.coinTransactionUC.SendCoins(username, requestBody.ToUser, requestBody.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"Errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"Message": "Монеты отправлены успешно"})
}

func NewCoinTransactionHandler(api Router, coinTransactionUC CoinTransactionUseCase, middleware Middleware) {
	handler := &coinTransactionDelivery{
		coinTransactionUC: coinTransactionUC,
	}

	protected := api.Group("/")
	protected.Use(middleware)

	protected.POST("/sendCoin", handler.SendCoin)

}
