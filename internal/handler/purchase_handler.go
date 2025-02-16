package handler

import (
	"net/http"
)

type PurchaseUseCase interface {
	BuyItem(username string, itemName string) error
}

type PurchaseDelivery struct {
	PurchaseUC PurchaseUseCase
}

func (d *PurchaseDelivery) BuyItem(c Context) {
	item := c.Param("item")

	if item == "" {
		c.JSON(http.StatusBadRequest, map[string]string{"Errors": "Неверный запрос"})
		return
	}

	username := c.MustGet("username").(string)

	if err := d.PurchaseUC.BuyItem(username, item); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"Errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"Message": "Товар куплен успешно"})
}

func NewPurchaseHandler(api Router, purchaseUC PurchaseUseCase, middleware Middleware) {
	handler := &PurchaseDelivery{
		PurchaseUC: purchaseUC,
	}

	protected := api.Group("/")
	protected.Use(middleware)

	protected.GET("/buy/:item", handler.BuyItem)
}
