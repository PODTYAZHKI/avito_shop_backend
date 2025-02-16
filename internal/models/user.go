package models

type User struct {
	ID       string `gorm:"column:id;type:uuid;default:uuid+generate_v4()"`
	Username string `json:"username,omitempty" gorm:"column:username"`
	Password string `json:"password,omitempty" gorm:"column:password"`
	Balance  int    `gorm:"column:balance"`
}

type UserInfo struct {
	Coins       int             `json:"coins"`
	Inventory   []PurchasedItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}
