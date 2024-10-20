package category

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	mc "github.com/febriandani/backend-financial-api/domain/model/category"
	"github.com/febriandani/backend-financial-api/infra"
	"github.com/sirupsen/logrus"
)

type CategoryConfig struct {
	db  *infra.DatabaseList
	log *logrus.Logger
}

func newDatabaseCategory(db *infra.DatabaseList, logger *logrus.Logger) CategoryConfig {
	return CategoryConfig{
		db:  db,
		log: logger,
	}
}

type Category interface {
	CreateCategory(ctx context.Context, tx *sql.Tx, data mc.CategoryRequest) (int64, error)
	GetCategoryByUserId(ctx context.Context, userID int, categoryType string) ([]mc.CategoryResponse, error)
	ValidateCategory(ctx context.Context, categoryId int, categoryType string) (bool, error)
	UpdateCategory(ctx context.Context, tx *sql.Tx, data mc.CategoryRequestUpdate) (int64, error)
	DeleteCategory(ctx context.Context, tx *sql.Tx, data mc.CategoryRequestUpdate) error
}

const (
	cQCreateCategory = `INSERT INTO public.categories
	(  
	user_id,
	category_type, 
	category_name, 
	category_description,
	created_at, 
	created_by)
	VALUES(?,?,?,?,?,?) returning id;
	`
)

func (cc CategoryConfig) CreateCategory(ctx context.Context, tx *sql.Tx, data mc.CategoryRequest) (int64, error) {
	param := make([]interface{}, 0)

	param = append(param, data.UserID)
	param = append(param, data.CategoryType)
	param = append(param, data.CategoryName)
	param = append(param, data.CategoryDescription)
	param = append(param, data.CreatedAt)
	param = append(param, data.CreatedBy)

	query, args, err := cc.db.Backend.Write.In(cQCreateCategory, param...)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	query = cc.db.Backend.Write.Rebind(query)

	var res *sql.Row
	if tx == nil {
		res = cc.db.Backend.Write.QueryRow(ctx, query, args...)
	} else {
		res = tx.QueryRowContext(ctx, query, args...)
	}

	if err != nil {
		return 0, err
	}

	err = res.Err()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	var id int64
	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (cc CategoryConfig) GetCategoryByUserId(ctx context.Context, userID int, categoryType string) ([]mc.CategoryResponse, error) {
	var result []mc.CategoryResponse

	query, args, err := cc.db.Backend.Read.In(`SELECT id, 
	user_id, 
	category_type,
	category_name, 
	category_description
	FROM public.categories c where c.user_id = ? and c.category_type = '`+categoryType+`'`, userID)
	if err != nil {
		return nil, err
	}

	query = cc.db.Backend.Read.Rebind(query)
	err = cc.db.Backend.Read.Select(&result, query, args...)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}

func (cc CategoryConfig) ValidateCategory(ctx context.Context, categoryId int, categoryType string) (bool, error) {
	var result bool

	query, args, err := cc.db.Backend.Read.In(`select exists(select id from public.categories c where c.id = ? and c.category_type = '`+categoryType+`');`, categoryId)
	if err != nil {
		return result, err
	}

	query = cc.db.Backend.Read.Rebind(query)
	err = cc.db.Backend.Read.Get(&result, query, args...)

	if err != nil && err != sql.ErrNoRows {
		return result, err
	}

	return result, nil
}

func (cc CategoryConfig) UpdateCategory(ctx context.Context, tx *sql.Tx, data mc.CategoryRequestUpdate) (int64, error) {
	q := `UPDATE public.categories
	SET category_name=?, category_description=?, updated_at=?, updated_by=? where id = ? and user_id = ?
	returning id;`

	param := make([]interface{}, 0)

	param = append(param, data.CategoryName)
	param = append(param, data.CategoryDescription)
	param = append(param, data.UpdatedAt)
	param = append(param, data.UpdatedBy)
	param = append(param, data.ID)
	param = append(param, data.UserID)

	query, args, err := cc.db.Backend.Write.In(q, param...)
	if err != nil {
		return 0, err
	}

	query = cc.db.Backend.Write.Rebind(query)

	query = cc.db.Backend.Write.Rebind(query)
	_, err = cc.db.Backend.Write.ExecContext(ctx, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return 1, nil
}

func (cc CategoryConfig) DeleteCategory(ctx context.Context, tx *sql.Tx, data mc.CategoryRequestUpdate) error {

	q := `delete from categories where id = ? and user_id = ?;`

	query, args, err := cc.db.Backend.Read.In(q, data.ID, data.UserID)
	if err != nil {
		return err
	}

	query = cc.db.Backend.Read.Rebind(query)

	var res sql.Result
	if tx == nil {
		res, err = cc.db.Backend.Write.ExecContext(ctx, query, args...)
	} else {
		res, err = tx.ExecContext(ctx, query, args...)
	}

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		err = fmt.Errorf("no rows updated")
		return err
	}

	return nil
}
