package models

type SendCoinRequest struct {
	ToUser string `json:"toUser,omitempty" binding:"required"`
	Amount int    `json:"amount,omitempty" binding:"required"`
}

type CoinTransaction struct {
	ID       string `gorm:"column:id;type:uuid;default:uuid+generate_v4()"`
	FromUser string `gorm:"column:from_user_id;type:uuid"`
	ToUser   string `gorm:"column:to_user_id;type:uuid"`
	Amount   int    `gorm:"column:amount"`
}

type CoinHistory struct {
	Received []CoinTransactionInfo `json:"received"`
	Sent     []CoinTransactionInfo `json:"sent"`
}

type CoinTransactionInfo struct {
	Amount   int    `json:"amount"`
	Username string `json:"username"`
}

func (CoinTransaction) TableName() string {
	return "transactions"
}
