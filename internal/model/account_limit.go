package model

type AccountLimit struct {
	AccountID            uint  `gorm:"primaryKey;column:account_id" json:"account_id"`
	AvailableCreditLimit int64 `gorm:"column:available_credit_limit" json:"available_credit_limit"`
}
