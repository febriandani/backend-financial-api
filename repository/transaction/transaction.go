package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	mt "github.com/febriandani/backend-financial-api/domain/model/transaction"
	"github.com/febriandani/backend-financial-api/infra"
	"github.com/sirupsen/logrus"
)

type TransactionConfig struct {
	db  *infra.DatabaseList
	log *logrus.Logger
}

func newDatabaseTransaction(db *infra.DatabaseList, logger *logrus.Logger) TransactionConfig {
	return TransactionConfig{
		db:  db,
		log: logger,
	}
}

type Transaction interface {
	CreateTransaction(ctx context.Context, tx *sql.Tx, data mt.TransactionRequest) (int64, error)
	GetSummaryHome(ctx context.Context, userID int) (*mt.SummaryHomeResponse, error)
	GetTransactions(ctx context.Context, filter mt.Filter) ([]mt.TransactionResponseDetail, error)
	GetCurrentBalanceTransactions(ctx context.Context, filter mt.Filter) (int64, error)
	UpdateTransaction(ctx context.Context, tx *sql.Tx, data mt.TransactionRequestUpdate) (int64, error)
	DeleteTransaction(ctx context.Context, tx *sql.Tx, data mt.TransactionRequestUpdate) error
}

const (
	tQCreateTransaction = `INSERT INTO public.transactions
	(  
	user_id, 
	category_id,
	category_type, 
	amount, 
	description,
	created_at, 
	created_by)
	VALUES(?,?,?,?,?,?,?) returning id;
	`

	tQGetSummaryHome = `SELECT 
	t.user_id,
    -- Menampilkan Total Pemasukan (category_type = 'IN')
    SUM(CASE WHEN t.category_type = 'IN' THEN cast(t.amount as numeric) ELSE 0 END) AS total_income,
    -- Menampilkan Total Pengeluaran (category_type = 'OUT')
    SUM(CASE WHEN t.category_type = 'OUT' THEN cast(t.amount as numeric) ELSE 0 END) AS total_spending,
    -- Menampilkan Saldo saat ini (Pemasukan - Pengeluaran)
    SUM(CASE WHEN t.category_type = 'IN' THEN cast(t.amount as numeric) ELSE 0 END) - SUM(CASE WHEN t.category_type = 'OUT' THEN cast(t.amount as numeric) ELSE 0 END) AS current_balance
FROM 
    public.transactions t
WHERE 
    t.user_id = ?
    group by t.user_id;`
)

func (tc TransactionConfig) CreateTransaction(ctx context.Context, tx *sql.Tx, data mt.TransactionRequest) (int64, error) {
	param := make([]interface{}, 0)

	param = append(param, data.UserID)
	param = append(param, data.CategoryID)
	param = append(param, data.CategoryType)
	param = append(param, data.Amount)
	param = append(param, data.Description)
	param = append(param, data.CreatedAt)
	param = append(param, data.CreatedBy)

	query, args, err := tc.db.Backend.Write.In(tQCreateTransaction, param...)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	query = tc.db.Backend.Write.Rebind(query)

	var res *sql.Row
	if tx == nil {
		res = tc.db.Backend.Write.QueryRow(ctx, query, args...)
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

func (tc TransactionConfig) GetSummaryHome(ctx context.Context, userID int) (*mt.SummaryHomeResponse, error) {
	var result mt.SummaryHomeResponse

	query, args, err := tc.db.Backend.Read.In(tQGetSummaryHome, userID)
	if err != nil {
		return nil, err
	}

	query = tc.db.Backend.Read.Rebind(query)
	err = tc.db.Backend.Read.GetContext(ctx, &result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &result, nil
}

func (tc TransactionConfig) GetTransactions(ctx context.Context, filter mt.Filter) ([]mt.TransactionResponseDetail, error) {
	var result []mt.TransactionResponseDetail

	uQGetTransactions := `SELECT 
	t.id as transaction_id,
	t.category_id ,
	t.category_type ,
	c.category_name ,
	t.description ,
	cast(t.amount as numeric) as amount,
	to_char(t.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
FROM
	public.transactions t 
LEFT JOIN 
	public.categories c ON c.id = t.category_id 
`

	queryStatement, args2 := BuildQueryGetTransactions(uQGetTransactions, filter)

	query, args, err := tc.db.Backend.Read.In(queryStatement, args2...)
	if err != nil {
		return result, err
	}

	query = tc.db.Backend.Read.Rebind(query)
	err = tc.db.Backend.Read.Select(&result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}

func (tc TransactionConfig) GetCurrentBalanceTransactions(ctx context.Context, filter mt.Filter) (int64, error) {
	var result int64

	uQGetTransactions := `SELECT 
    -- Menampilkan Saldo saat ini (Pemasukan - Pengeluaran)
    SUM(CASE WHEN t.category_type = 'IN' THEN cast(t.amount as numeric) ELSE 0 END) - SUM(CASE WHEN t.category_type = 'OUT' THEN cast(t.amount as numeric) ELSE 0 END) AS current_balance
FROM 
    public.transactions t
`

	queryStatement, args2 := BuildQueryGetCurrentBalanceTransactions(uQGetTransactions, filter)

	query, args, err := tc.db.Backend.Read.In(queryStatement, args2...)
	if err != nil {
		return result, err
	}

	query = tc.db.Backend.Read.Rebind(query)
	err = tc.db.Backend.Read.Get(&result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return result, nil
}

func (tc TransactionConfig) UpdateTransaction(ctx context.Context, tx *sql.Tx, data mt.TransactionRequestUpdate) (int64, error) {
	q := `UPDATE public.transactions
	SET amount=?, description=?, updated_at=?, updated_by=? where id = ? and user_id = ?
	returning id;`

	param := make([]interface{}, 0)

	param = append(param, data.Amount)
	param = append(param, data.Description)
	param = append(param, data.UpdatedAt)
	param = append(param, data.UpdatedBy)
	param = append(param, data.ID)
	param = append(param, data.UserID)

	query, args, err := tc.db.Backend.Write.In(q, param...)
	if err != nil {
		return 0, err
	}

	query = tc.db.Backend.Write.Rebind(query)

	query = tc.db.Backend.Write.Rebind(query)
	_, err = tc.db.Backend.Write.ExecContext(ctx, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return 1, nil
}

func (tc TransactionConfig) DeleteTransaction(ctx context.Context, tx *sql.Tx, data mt.TransactionRequestUpdate) error {

	q := `delete from transactions where id = ? and user_id = ?;`

	query, args, err := tc.db.Backend.Read.In(q, data.ID, data.UserID)
	if err != nil {
		return err
	}

	query = tc.db.Backend.Read.Rebind(query)

	var res sql.Result
	if tx == nil {
		res, err = tc.db.Backend.Write.ExecContext(ctx, query, args...)
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
