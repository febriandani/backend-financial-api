package transaction

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

func (r *TransactionRequest) Validate() map[string]string {
	if r.UserID == 0 {
		return map[string]string{
			"en": "user id cannot be empty",
			"id": "id user tidak boleh kosong",
		}
	}

	if r.CategoryID == 0 {
		return map[string]string{
			"en": "Category id cannot be empty",
			"id": "id kategori tidak boleh kosong",
		}
	}

	if r.CategoryType == "" {
		return map[string]string{
			"en": "Category type cannot be empty",
			"id": "Tipe kategori tidak boleh kosong",
		}
	}

	if r.Amount == "" {
		return map[string]string{
			"en": "Amount cannot be empty",
			"id": "Amount tidak boleh kosong",
		}
	}

	if r.Description == "" {
		return map[string]string{
			"en": "Description cannot be empty",
			"id": "Deskripsi tidak boleh kosong",
		}
	}

	return nil
}

type TransactionRequest struct {
	UserID       int64     `json:"user_id"`
	CategoryID   int64     `json:"category_id"`
	CategoryType string    `json:"category_type"`
	Amount       string    `json:"amount"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
}

type Filter struct {
	UserID       int
	CategoryType null.String `json:"category_type" db:"category_type"`
	StartDate    null.String `json:"start_date"`
	EndDate      null.String `json:"end_date"`
	Offset       null.Int    `json:"offset"`
	Limit        null.Int    `json:"limit"`
}

type SummaryHomeResponse struct {
	UserID         int64  `json:"user_id" db:"user_id"`
	Name           string `json:"name" db:"name"`
	CurrentBalance int64  `json:"current_balance" db:"current_balance"`
	TotalSpending  int64  `json:"total_spending" db:"total_spending"`
	TotalIncome    int64  `json:"total_income" db:"total_income"`
}

type TransactionResponseDetail struct {
	TransactionID int64       `json:"transaction_id" db:"transaction_id"`
	CategoryID    null.Int    `json:"category_id" db:"category_id"`
	CategoryType  string      `json:"category_type" db:"category_type"`
	CategoryName  null.String `json:"category_name" db:"category_name"`
	Description   string      `json:"description" db:"description"`
	Amount        int64       `json:"amount" db:"amount"`
	CreatedAt     string      `json:"created_at" db:"created_at"`
}

type TransactionResponse struct {
	CurrentBalance    int64                       `json:"current_balance"`
	TransactionDetail []TransactionResponseDetail `json:"transaction_detail"`
}

type TransactionRequestUpdate struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id"`
	Amount      string    `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}
