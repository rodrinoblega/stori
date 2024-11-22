package entities

type Account struct {
	AccountID int    `gorm:"primaryKey;column:account_id"`
	Name      string `gorm:"column:name"`
	Mail      string `gorm:"column:mail"`
}
