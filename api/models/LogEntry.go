package models

import "time"

type LogEntry struct {
	NumberOrder string `gorm:"number_order,omitempty"`
	IdSession   string `gorm:"id_session,omitempty"`
	Status      string `gorm:"status"`
	Amount      int    `gorm:"amount"`
	//BuyOrder          string     `gorm:"buy_order"`
	//SessionID         string     `gorm:"session_id"`
	AccountingDate    string     `gorm:"accounting_date"`
	TransactionDate   time.Time  `gorm:"transaction_date"`
	PaymentTypeCode   string     `gorm:"payment_type_code"`
	CardDetail        CardDetail `gorm:"-"`
	AuthorizationCode string     `gorm:"authorization_code"`
}

type CardDetail struct {
	CardNumber string `gorm:"card_number"`
}
