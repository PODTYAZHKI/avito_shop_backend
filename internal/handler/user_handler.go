package handler

import (
	"net/http"

	"avito-shop-test/internal/models"
)

type UserUseCase interface {
	Authenticate(username, password string) (string, error)
	GetUserInfo(username string) (*models.UserInfo, error)
	GetCoinHistory(userID string) (models.CoinHistory, error)
	GetPurchasedItems(userID string) ([]models.PurchasedItem, error)
}

type UserDelivery struct {
	UserUC UserUseCase
}

func (d *UserDelivery) Authenticate(c Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"Errors": "Неверный запрос."})
		return
	}

	token, err := d.UserUC.Authenticate(req.Username, req.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"Errors": err.Error()})
		return
	}
	response := map[string]string{
		"token": token,
	}

	c.JSON(http.StatusOK, response)
}

func (d *UserDelivery) GetUserInfo(c Context) {
	username := c.MustGet("username").(string)

	userInfo, err := d.UserUC.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"Errors": err.Error()})

		return
	}

	response := models.UserInfo{
		Coins:       userInfo.Coins,
		Inventory:   userInfo.Inventory,
		CoinHistory: userInfo.CoinHistory,
	}

	c.JSON(http.StatusOK, response)

}

func NewUserHandler(api Router, userUC UserUseCase, middleware Middleware) {
	handler := &UserDelivery{
		UserUC: userUC,
	}

	api.POST("/auth", handler.Authenticate)

	protected := api.Group("/")
	protected.Use(middleware)
	protected.GET("/info", handler.GetUserInfo)
}
