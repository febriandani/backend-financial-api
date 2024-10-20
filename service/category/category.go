package category

import (
	"context"
	"errors"
	"time"

	mc "github.com/febriandani/backend-financial-api/domain/model/category"
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/domain/utils"
	"github.com/febriandani/backend-financial-api/infra"
	rc "github.com/febriandani/backend-financial-api/repository/category"
	"github.com/sirupsen/logrus"
)

type CategoryService struct {
	db     rc.DatabaseCategory
	log    *logrus.Logger
	conf   general.AppService
	dbConn *infra.DatabaseList
}

func newCategoryService(database rc.DatabaseCategory, logger *logrus.Logger, dbConn *infra.DatabaseList, conf general.AppService) CategoryService {
	return CategoryService{
		db:     database,
		log:    logger,
		conf:   conf,
		dbConn: dbConn,
	}
}

type Category interface {
	CreateCategory(ctx context.Context, data mc.CategoryRequest) (map[string]string, error)
	GetCategoryByUserId(ctx context.Context, userId int, categoryType string) ([]mc.CategoryResponse, map[string]string, error)
	UpdateCategory(ctx context.Context, data mc.CategoryRequestUpdate) (map[string]string, error)
	DeleteCategory(ctx context.Context, data mc.CategoryRequestUpdate) (map[string]string, error)
}

func (cs CategoryService) CreateCategory(ctx context.Context, data mc.CategoryRequest) (map[string]string, error) {
	messages := data.Validate()
	if messages != nil {
		return messages, errors.New("data not valid")
	}

	cs.log.WithField("request", utils.StructToString(data)).WithError(nil).Info("Category | start Category")

	tx, err := cs.dbConn.Backend.Write.Begin()
	if err != nil {
		cs.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Category | fail to begin Category")
		return map[string]string{
			"en": "An error occurred during Category, please try again",
			"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
		}, err
	}

	if data.CategoryType != "IN" && data.CategoryType != "OUT" {
		return map[string]string{
			"en": "Category type not supported, please input IN or OUT",
			"id": "Tipe kategori tidak didukung, tolong input IN atau OUT",
		}, err
	}

	_, err = cs.db.Category.CreateCategory(ctx, tx, mc.CategoryRequest{
		UserID:              data.UserID,
		CategoryType:        data.CategoryType,
		CategoryName:        data.CategoryName,
		CategoryDescription: data.CategoryDescription,
		CreatedAt:           time.Now().UTC(),
		CreatedBy:           data.CreatedBy,
	})
	if err != nil {
		tx.Rollback()
		cs.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Category | fail to create Category")
		return map[string]string{
			"en": "An error occurred during Category, please try again",
			"id": "Terjadi kesalahan saat menambahkan kategori, silakan coba lagi",
		}, err
	}

	err = tx.Commit()
	if err != nil {
		cs.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Category | fail to commit Category")
		tx.Rollback()
		return map[string]string{
			"en": "An error occurred during Category, please try again",
			"id": "Terjadi kesalahan saat menambahkan kategori, silakan coba lagi",
		}, err
	}

	cs.log.WithField("request", utils.StructToString(data)).WithError(nil).Errorf("Category | finish Category")

	return map[string]string{
		"en": "Successfully",
		"id": "Berhasil",
	}, nil
}

func (cs CategoryService) GetCategoryByUserId(ctx context.Context, userId int, categoryType string) ([]mc.CategoryResponse, map[string]string, error) {

	cs.log.WithFields(logrus.Fields{
		categoryType: categoryType,
		"userID":     userId,
	}).WithError(nil).Infoln("Category | start GetCategoryByUserId")

	if categoryType != "IN" && categoryType != "OUT" {
		cs.log.WithFields(logrus.Fields{
			categoryType: categoryType,
			"userID":     userId,
		}).WithError(nil).Infoln("Category | category type not match ")
		return nil, map[string]string{
			"en": "Failed to get category",
			"id": "Gagal menampilkan kategory",
		}, errors.New("error")
	}

	data, err := cs.db.Category.GetCategoryByUserId(ctx, userId, categoryType)
	if err != nil {
		cs.log.WithFields(logrus.Fields{
			categoryType: categoryType,
			"userID":     userId,
		}).WithError(nil).Infoln("Category | failed get GetCategoryByUserId")
		return nil, map[string]string{
			"en": "Failed to get category",
			"id": "Gagal menampilkan kategory",
		}, err
	}

	cs.log.WithFields(logrus.Fields{
		categoryType: categoryType,
		"userID":     userId,
	}).WithError(nil).Infof("Category | finish GetCategoryByUserId : %v ", data)

	return data, map[string]string{
		"en": "Successfully",
		"id": "Berhasil",
	}, nil
}

func (cs CategoryService) UpdateCategory(ctx context.Context, data mc.CategoryRequestUpdate) (map[string]string, error) {

	cs.log.WithField("request", utils.StructToString(data)).WithError(nil).Infoln("Category | start UpdateCategory")

	_, err := cs.db.Category.UpdateCategory(ctx, nil, mc.CategoryRequestUpdate{
		ID:                  data.ID,
		UserID:              data.UserID,
		CategoryName:        data.CategoryName,
		CategoryDescription: data.CategoryDescription,
		UpdatedAt:           time.Now().UTC(),
		UpdatedBy:           data.UpdatedBy,
	})
	if err != nil {
		cs.log.WithField("request", utils.StructToString(data)).WithError(nil).Errorf("Category | failed update Category")
		return map[string]string{
			"en": "Failed to update category",
			"id": "Gagal update kategory",
		}, err
	}

	cs.log.WithField("request", utils.StructToString(data)).WithError(nil).Infoln("Category | start UpdateCategory")

	return map[string]string{
		"en": "Successfully",
		"id": "Berhasil",
	}, nil
}

func (cs CategoryService) DeleteCategory(ctx context.Context, data mc.CategoryRequestUpdate) (map[string]string, error) {

	cs.log.WithField("request", utils.StructToString(data)).WithError(nil).Infoln("Category | start DeleteCategory")

	err := cs.db.Category.DeleteCategory(ctx, nil, mc.CategoryRequestUpdate{
		ID:     data.ID,
		UserID: data.UserID,
	})
	if err != nil {
		cs.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Category | failed delete Category")
		return map[string]string{
			"en": "Failed to delete category",
			"id": "Gagal menghapus kategory",
		}, err
	}

	cs.log.WithField("request", utils.StructToString(data)).WithError(nil).Infoln("Category | start DeleteCategory")

	return map[string]string{
		"en": "Successfully",
		"id": "Berhasil",
	}, nil
}
