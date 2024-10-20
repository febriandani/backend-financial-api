package transaction

import "time"

func (r *CategoryRequest) Validate() map[string]string {

	if r.CategoryType == "" {
		return map[string]string{
			"en": "Category type cannot be empty",
			"id": "Tipe kategori tidak boleh kosong",
		}
	}

	if r.CategoryName == "" {
		return map[string]string{
			"en": "Category name cannot be empty",
			"id": "Nama kategory tidak boleh kosong",
		}
	}

	return nil
}

type CategoryRequest struct {
	UserID              int       `json:"user_id"`
	CategoryType        string    `json:"category_type"`
	CategoryName        string    `json:"category_name"`
	CategoryDescription string    `json:"category_description"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
}

type CategoryResponse struct {
	UserID              int    `json:"user_id" db:"user_id"`
	ID                  int    `json:"category_id" db:"id"`
	CategoryType        string `json:"category_type" db:"category_type"`
	CategoryName        string `json:"category_name" db:"category_name"`
	CategoryDescription string `json:"category_description" db:"category_description"`
}

type CategoryRequestUpdate struct {
	ID                  int       `json:"id" db:"id"`
	UserID              int       `json:"user_id"`
	CategoryName        string    `json:"category_name" db:"category_name"`
	CategoryDescription string    `json:"category_description" db:"category_description"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
}
