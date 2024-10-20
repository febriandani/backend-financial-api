package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/febriandani/backend-financial-api/domain/model/user"
	"github.com/febriandani/backend-financial-api/infra"
	"github.com/sirupsen/logrus"
)

type UserConfig struct {
	db  *infra.DatabaseList
	log *logrus.Logger
}

func newDatabaseUser(db *infra.DatabaseList, logger *logrus.Logger) UserConfig {
	return UserConfig{
		db:  db,
		log: logger,
	}
}

type User interface {
	Registration(ctx context.Context, tx *sql.Tx, data user.RegistrationUser) (int64, error)
	IsExistUserByEmail(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*user.ResponseDataUser, error)
	UpdatePassword(ctx context.Context, email, password string) error
}

const (
	uQCreateUser = `INSERT INTO public.users
	(  
	name, 
	username,
	email, 
	"password", 
	created_at, 
	created_by, 
	phone_number)
	VALUES(?,?,?,?,?,?,?) returning id;
	`

	uqGetDataUserByEmail = `SELECT id, 
	email, password, name
FROM public.users u where u.email = ? or u.username = ?;
	`
)

func (uc UserConfig) Registration(ctx context.Context, tx *sql.Tx, data user.RegistrationUser) (int64, error) {
	param := make([]interface{}, 0)

	param = append(param, data.Name)
	param = append(param, data.Username)
	param = append(param, data.Email)
	param = append(param, data.Password)
	param = append(param, data.CreatedAt)
	param = append(param, data.CreatedBy)
	param = append(param, data.PhoneNumber)

	query, args, err := uc.db.Backend.Write.In(uQCreateUser, param...)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	query = uc.db.Backend.Write.Rebind(query)

	var res *sql.Row
	if tx == nil {
		res = uc.db.Backend.Write.QueryRow(ctx, query, args...)
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

func (uc UserConfig) IsExistUserByEmail(ctx context.Context, email string) (bool, error) {
	var result bool

	uQIsExistUserByEmailActive := ` select exists(select u.id from users u where ( u.email = ? or u.username = ? ) and is_active = true)`

	query, args, err := uc.db.Backend.Read.In(uQIsExistUserByEmailActive, email, email)
	if err != nil {
		return result, err
	}

	query = uc.db.Backend.Read.Rebind(query)
	err = uc.db.Backend.Read.GetContext(ctx, &result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return result, err
	}

	return result, nil
}

func (uc UserConfig) GetUserByEmail(ctx context.Context, email string) (*user.ResponseDataUser, error) {
	var result user.ResponseDataUser

	query, args, err := uc.db.Backend.Read.In(uqGetDataUserByEmail, email, email)
	if err != nil {
		return nil, err
	}

	query = uc.db.Backend.Read.Rebind(query)
	err = uc.db.Backend.Read.GetContext(ctx, &result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &result, nil
}

func (uc UserConfig) UpdatePassword(ctx context.Context, email, password string) error {
	q := `	UPDATE public.users
	SET "password"=?, updated_at=?, updated_by=?
	WHERE email = ?;`

	query, args, err := uc.db.Backend.Read.In(q, password, time.Now().UTC(), "system", email)
	if err != nil {
		return err
	}

	query = uc.db.Backend.Read.Rebind(query)
	res, err := uc.db.Backend.Write.ExecContext(ctx, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		err = fmt.Errorf("no rows inserted")
		return err
	}

	return nil
}
