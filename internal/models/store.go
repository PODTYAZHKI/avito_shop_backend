package models

type Product struct {
	Name  string `json:"name" gorm:"column:name"`
	Price int    `json:"price" gorm:"column:price"`
}

func (Product) TableName() string {
	return "items"
}

type Inventory struct {
	ID       string `gorm:"type:uuid;default:uuid+generate_v4()"`
	UserID   string `gorm:"column:user_id;type:uuid"`
	ItemType string `gorm:"column:item_type"`
	Quantity int    `gorm:"column:quantity;"`
}

func (Inventory) TableName() string {
	return "inventory"
}

type PurchasedItem struct {
	ItemName string `json:"type"`
	Quantity int    `json:"quantity"`
}
