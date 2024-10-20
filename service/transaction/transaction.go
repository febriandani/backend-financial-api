package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/febriandani/backend-financial-api/domain/model/general"
	mt "github.com/febriandani/backend-financial-api/domain/model/transaction"
	"github.com/febriandani/backend-financial-api/domain/utils"
	"github.com/febriandani/backend-financial-api/infra"
	rc "github.com/febriandani/backend-financial-api/repository/category"
	rt "github.com/febriandani/backend-financial-api/repository/transaction"
	"github.com/sirupsen/logrus"
)

type TransactionService struct {
	db         rt.DatabaseTransaction
	dbCategory rc.DatabaseCategory
	log        *logrus.Logger
	conf       general.AppService
	dbConn     *infra.DatabaseList
}

func newTransactionService(database rt.DatabaseTransaction, databaseCategory rc.DatabaseCategory, logger *logrus.Logger, dbConn *infra.DatabaseList, conf general.AppService) TransactionService {
	return TransactionService{
		db:         database,
		dbCategory: databaseCategory,
		log:        logger,
		conf:       conf,
		dbConn:     dbConn,
	}
}

type Transaction interface {
	CreateTransaction(ctx context.Context, data mt.TransactionRequest) (map[string]string, error)
	GetHomeSummaryByUserId(ctx context.Context, userId int) (*mt.SummaryHomeResponse, map[string]string, error)
	GetTransactions(ctx context.Context, filter mt.Filter) (*mt.TransactionResponse, map[string]string, error)
	UpdateTransaction(ctx context.Context, data mt.TransactionRequestUpdate) (map[string]string, error)
	DeleteTransaction(ctx context.Context, data mt.TransactionRequestUpdate) (map[string]string, error)
}

func (ts TransactionService) CreateTransaction(ctx context.Context, data mt.TransactionRequest) (map[string]string, error) {
	messages := data.Validate()
	if messages != nil {
		return messages, errors.New("data not valid")
	}

	ts.log.WithField("request", utils.StructToString(data)).WithError(nil).Infoln("Transaction | start Transaction")

	tx, err := ts.dbConn.Backend.Write.Begin()
	if err != nil {
		ts.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Transaction | fail to begin transaction")
		return map[string]string{
			"en": "An error occurred during Transaction, please try again",
			"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
		}, err
	}

	if data.CategoryType != "IN" && data.CategoryType != "OUT" {
		return map[string]string{
			"en": "Category type not supported, please input IN or OUT",
			"id": "Tipe kategori tidak didukung, tolong input IN atau OUT",
		}, err
	}

	isValid, err := ts.dbCategory.Category.ValidateCategory(ctx, int(data.CategoryID), data.CategoryType)
	if err != nil {
		ts.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Transaction | fail to validation category")
		return map[string]string{
			"en": "An error occurred during Transaction, please try again",
			"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
		}, err
	}

	if !isValid {
		return map[string]string{
			"en": "Category and category type not match",
			"id": "Kategori dan Tipe kategori tidak didukung",
		}, err
	}

	_, err = ts.db.Transaction.CreateTransaction(ctx, tx, mt.TransactionRequest{
		UserID:       data.UserID,
		CategoryID:   data.CategoryID,
		CategoryType: data.CategoryType,
		Amount:       data.Amount,
		Description:  data.Description,
		CreatedAt:    time.Now().UTC(),
		CreatedBy:    data.CreatedBy,
	})
	if err != nil {
		tx.Rollback()
		ts.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Transaction | fail to create Transaction")
		return map[string]string{
			"en": "An error occurred during Transaction, please try again",
			"id": "Terjadi kesalahan saat transaksi, silakan coba lagi",
		}, err
	}

	err = tx.Commit()
	if err != nil {
		ts.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Transaction | fail to commit transaction")
		tx.Rollback()
		return map[string]string{
			"en": "An error occurred during Transaction, please try again",
			"id": "Terjadi kesalahan saat transaksi, silakan coba lagi",
		}, err
	}

	ts.log.WithField("request", utils.StructToString(data)).WithError(nil).Errorf("Transaction | finish Transaction")

	return map[string]string{
		"en": "Successfully",
		"id": "Berhasil",
	}, nil
}

func (ts TransactionService) GetHomeSummaryByUserId(ctx context.Context, userId int) (*mt.SummaryHomeResponse, map[string]string, error) {

	ts.log.WithField("request", utils.StructToString(userId)).WithError(nil).Errorf("Transaction | start Transaction")

	data, err := ts.db.Transaction.GetSummaryHome(ctx, userId)
	if err != nil {
		ts.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Transaction | fail to get summary transaction")
		return nil, map[string]string{
			"en": "An error occurred during get summary Transaction, please try again",
			"id": "Terjadi kesalahan saat menampilkan ringkasan transaksi, silakan coba lagi",
		}, err
	}

	ts.log.WithField("request", utils.StructToString(data)).WithError(nil).Infoln("Transaction | finish Transaction")

	return data, map[string]string{
		"en": "Successfully",
		"id": "Berhasil",
	}, nil
}

func (ts TransactionService) GetTransactions(ctx context.Context, filter mt.Filter) (*mt.TransactionResponse, map[string]string, error) {

	ts.log.WithField("request", utils.StructToString(filter)).WithError(nil).Infoln("Transaction | start GetTransaction")

	if filter.StartDate.String != "" && filter.EndDate.String != "" {
		_, err := time.Parse("2006-01-02", filter.StartDate.String)
		if err != nil {
			return nil, map[string]string{
				"en": "An error occurred formatting start date filter",
				"id": "Terjadi kesalahan dalam memformat filter tanggal mulai",
			}, err
		}
		_, err = time.Parse("2006-01-02", filter.EndDate.String)
		if err != nil {
			return nil, map[string]string{
				"en": "An error occurred formatting end date filter",
				"id": "Terjadi kesalahan dalam memformat filter tanggal akhir",
			}, err
		}
	}

	if filter.Offset.Int64 == 0 && filter.Limit.Int64 == 0 {
		filter.Offset.Valid = false
		filter.Limit.Valid = false
	}

	dataCurrentBalance, err := ts.db.Transaction.GetCurrentBalanceTransactions(ctx, filter)
	if err != nil {
		ts.log.WithField("request", utils.StructToString(filter)).WithError(err).Errorf("Transaction | fail to get current balance transaction")
		return nil, map[string]string{
			"en": "An error occurred during get Transaction, please try again",
			"id": "Terjadi kesalahan saat menampilkan transaksi, silakan coba lagi",
		}, err
	}

	data, err := ts.db.Transaction.GetTransactions(ctx, filter)
	if err != nil {
		ts.log.WithField("request", utils.StructToString(filter)).WithError(err).Errorf("Transaction | fail to get transaction")
		return nil, map[string]string{
			"en": "An error occurred during get Transaction, please try again",
			"id": "Terjadi kesalahan saat menampilkan transaksi, silakan coba lagi",
		}, err
	}

	ts.log.WithField("request", utils.StructToString(filter)).WithError(nil).Infoln("Transaction | finish GetTransaction")

	return &mt.TransactionResponse{
			CurrentBalance:    dataCurrentBalance,
			TransactionDetail: data,
		}, map[string]string{
			"en": "Successfully",
			"id": "Berhasil",
		}, nil
}

func (ts TransactionService) UpdateTransaction(ctx context.Context, data mt.TransactionRequestUpdate) (map[string]string, error) {

	ts.log.WithField("request", utils.StructToString(data)).WithError(nil).Infoln("Transaction | start UpdatedTransaction")

	_, err := ts.db.Transaction.UpdateTransaction(ctx, nil, mt.TransactionRequestUpdate{
		ID:          data.ID,
		UserID:      data.UserID,
		Amount:      data.Amount,
		Description: data.Description,
		UpdatedAt:   time.Now().UTC(),
		UpdatedBy:   data.UpdatedBy,
	})
	if err != nil {
		ts.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Transaction | fail to update transaction")
		return map[string]string{
			"en": "An error occurred during update Transaction, please try again",
			"id": "Terjadi kesalahan saat memperbarui transaksi, silakan coba lagi",
		}, err
	}

	ts.log.WithField("request", utils.StructToString(data)).WithError(nil).Infoln("Transaction | finish UpdatedTransaction")

	return map[string]string{
		"en": "Successfully",
		"id": "Berhasil",
	}, nil
}

func (ts TransactionService) DeleteTransaction(ctx context.Context, data mt.TransactionRequestUpdate) (map[string]string, error) {

	ts.log.WithField("request", utils.StructToString(data)).WithError(nil).Infoln("Transaction | start DeleteTransaction")

	err := ts.db.Transaction.DeleteTransaction(ctx, nil, mt.TransactionRequestUpdate{
		ID:     data.ID,
		UserID: data.UserID,
	})
	if err != nil {
		ts.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Transaction | fail to delete transaction")
		return map[string]string{
			"en": "An error occurred during delete Transaction, please try again",
			"id": "Terjadi kesalahan saat hapus transaksi, silakan coba lagi",
		}, err
	}

	ts.log.WithField("request", utils.StructToString(data)).WithError(nil).Infoln("Transaction | finish DeleteTransaction")

	return map[string]string{
		"en": "Successfully",
		"id": "Berhasil",
	}, nil
}
