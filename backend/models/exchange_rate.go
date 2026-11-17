package models

import "time"

type ExchangeRate struct {
	ID           uint      `gorm:"primaryKey" json:"_id"`
	FromCurrency string    `gorm:"not null" json:"from_currency" binding:"required"`
	ToCurrency   string    `gorm:"not null" json:"to_currency" binding:"required"`
	Rate         float64   `gorm:"not null" json:"rate" binding:"required"`
	Date         time.Time `gorm:"not null" json:"date"`
}
