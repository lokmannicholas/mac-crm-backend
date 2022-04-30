package models

import (
	"database/sql"
	"time"
)

type StoragesReportResult struct {
	BranchID     string        `json:"branch_id"`
	ID           string        `json:"id"`
	Code         int           `json:"code"`
	Status       string        `json:"status"`
	MonthlyPrice float64       `json:"monthly_price"`
	YearlyPrice  float64       `json:"yearly_price"`
	Schedule     string        `json:"schedule"`
	ExpiredAt    *time.Time    `json:"expired_at"`
	OrderNo      *string       `json:"order_no"`
	OrderID      *string       `json:"order_id"`
	FinalPrice   *float64      `json:"final_price"`
	CustomerID   string        `json:"customer_id"`
	Categories   *sql.RawBytes `json:"categories"`
	Features     *sql.RawBytes `json:"features"`
	Customer     []byte        `json:"customer"`
}

type RentalOrdersReportResult struct {
	BranchID        string
	ID              string
	OrderNo         string
	RentType        string
	OrderDate       time.Time
	ActualStartDate time.Time
	ActualEndDate   time.Time
	Status          string
	DiscountAmount  float64
	FinalPrice      float64
	PaymentDate     *time.Time
	PaymentMethod   *string
	Categories      *sql.RawBytes
}
