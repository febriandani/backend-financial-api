package user

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/febriandani/backend-financial-api/domain/model/general"
	mu "github.com/febriandani/backend-financial-api/domain/model/user"
	"github.com/febriandani/backend-financial-api/domain/utils"
	"github.com/febriandani/backend-financial-api/infra"
	ru "github.com/febriandani/backend-financial-api/repository/user"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	db     ru.DatabaseUser
	log    *logrus.Logger
	conf   general.AppService
	dbConn *infra.DatabaseList
}

func newUserService(database ru.DatabaseUser, logger *logrus.Logger, dbConn *infra.DatabaseList, conf general.AppService) UserService {
	return UserService{
		db:     database,
		log:    logger,
		conf:   conf,
		dbConn: dbConn,
	}
}

type User interface {
	CreateRegistrationUser(ctx context.Context, data mu.RegistrationUser) (map[string]string, error)
	ForgotPassword(ctx context.Context, data mu.ForgotPasswordRequest) (map[string]string, error)
	Login(ctx context.Context, data mu.LoginRequest) (*mu.LoginResponse, map[string]string, error)
}

func (us UserService) CreateRegistrationUser(ctx context.Context, data mu.RegistrationUser) (map[string]string, error) {
	messages := data.Validate()
	if messages != nil {
		return messages, errors.New("data not valid")
	}

	us.log.WithField("request", utils.StructToString(data.Email)).WithError(nil).Info("Registration | start registration")

	getDataUser, err := us.db.User.GetUserByEmail(ctx, strings.ToLower(data.Email))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Registration | fail to get user by email")
		return map[string]string{
			"en": "An error occurred during registration, please try again",
			"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
		}, err
	}

	if getDataUser.Email == strings.ToLower(data.Email) {
		return map[string]string{
			"en": "Email has been registered, please try with another email",
			"id": "Email telah terdaftar, silakan coba dengan email lain",
		}, errors.New("error")
	} else {
		tx, err := us.dbConn.Backend.Write.Begin()
		if err != nil {
			us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Registration | fail to begin transaction")
			return map[string]string{
				"en": "An error occurred during registration, please try again",
				"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
			}, err
		}

		password, err := utils.GeneratePassword(data.Password)
		if err != nil {
			tx.Rollback()
			us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Registration | fail to generate password")
			return map[string]string{
				"en": "An error occurred during registration, please try again",
				"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
			}, err
		}

		_, err = us.db.User.Registration(ctx, tx, mu.RegistrationUser{
			Name:        data.Name,
			Email:       strings.ToLower(strings.ToLower(data.Email)),
			Username:    data.Username,
			CreatedAt:   time.Now().UTC(),
			CreatedBy:   "system",
			Password:    password,
			PhoneNumber: utils.FormatPhoneNumber(data.PhoneNumber),
		})
		if err != nil {
			tx.Rollback()
			us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Registration | fail to create registration")
			return map[string]string{
				"en": "An error occurred during registration, please try again",
				"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
			}, err
		}

		err = tx.Commit()
		if err != nil {
			us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Registration | fail to commit transaction")
			tx.Rollback()
			return map[string]string{
				"en": "An error occurred during registration, please try again",
				"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
			}, err
		}
	}

	us.log.WithField("request", utils.StructToString(data.Email)).WithError(nil).Errorf("Registration | finish registration")

	return map[string]string{
		"en": "User is successfully registered",
		"id": "Pengguna berhasil terdaftar",
	}, nil
}

func (us UserService) ForgotPassword(ctx context.Context, data mu.ForgotPasswordRequest) (map[string]string, error) {
	messages := data.Validate()
	if messages != nil {
		return messages, errors.New("data not valid")
	}

	us.log.WithField("request", utils.StructToString(data.Email)).WithError(nil).Errorf("ForgotPassword | start request")

	isExist, err := us.db.User.IsExistUserByEmail(ctx, strings.ToLower(data.Email))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data.Email)).WithError(err).Errorf("ForgotPassword | fail to get exist user")
		return map[string]string{
			"en": "There was an error in checking user data",
			"id": "Ada kesalahan dalam memeriksa data pengguna",
		}, err
	}

	if !isExist {
		us.log.WithField("request", utils.StructToString(data.Email)).Errorf("ForgotPassword | user not found")
		return map[string]string{
			"en": "Email not registered",
			"id": "Email tidak terdaftar",
		}, errors.New("user not exist")
	}

	userData, err := us.db.User.GetUserByEmail(ctx, strings.ToLower(data.Email))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data.Email)).WithError(err).Errorf("ForgotPassword | fail to get company user data")
		return map[string]string{
			"en": "There was an error in checking user data",
			"id": "Ada kesalahan dalam memeriksa data pengguna",
		}, err
	}

	isValid, err := utils.ComparePassword(userData.Password, data.OldPassword)
	if err != nil {
		us.log.WithField("request", utils.StructToString(data.Email)).WithError(err).Errorf("ForgotPassword | fail to compare password")
		return map[string]string{
			"en": "There was an error changing the password",
			"id": "Ada kesalahan dalam mengubah kata sandi",
		}, err
	}

	if !isValid {
		us.log.WithField("request", utils.StructToString(data.Email)).Errorf("ForgotPassword | password not valid")
		return map[string]string{
			"en": "The current password failed",
			"id": "Kata sandi saat ini salah",
		}, errors.New("PS-011")
	}

	password, err := utils.GeneratePassword(data.NewPassword)
	if err != nil {
		us.log.WithField("request", utils.StructToString(data.Email)).WithError(err).Errorf("ForgotPassword | fail to generate password")
		return map[string]string{
			"en": "There was an error changing the password",
			"id": "Ada kesalahan dalam mengubah kata sandi",
		}, err
	}

	err = us.db.User.UpdatePassword(ctx, strings.ToLower(data.Email), password)
	if err != nil {
		us.log.WithField("request", utils.StructToString(data.Email)).WithError(err).Errorf("ForgotPassword | fail to change password from repo")
		return map[string]string{
			"en": "There was an error changing the password",
			"id": "Ada kesalahan dalam mengubah kata sandi",
		}, err
	}

	us.log.WithField("request", utils.StructToString(data.Email)).WithError(nil).Errorf("ForgotPassword | finish change password")

	return map[string]string{
		"en": "success",
		"id": "password berhasil diganti",
	}, nil
}

func (us UserService) Login(ctx context.Context, data mu.LoginRequest) (*mu.LoginResponse, map[string]string, error) {
	messages := data.Validate()
	if messages != nil {
		return &mu.LoginResponse{
			Token: mu.JWTAccess{
				AccessToken:        "",
				AccessTokenExpired: "",
			},
			NamaLengkap: "",
		}, messages, errors.New("data not valid")
	}

	us.log.WithField("request", utils.StructToString(data.Email)).WithError(nil).Errorf("Login | start request login")

	isExist, err := us.db.User.IsExistUserByEmail(ctx, strings.ToLower(strings.ToLower(data.Email)))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data.Email)).WithError(err).Errorf("Login | fail to get exist user")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "There was an error in checking user data",
				"id": "Ada kesalahan dalam memeriksa data pengguna",
			}, errors.New("FailedServer")
	}

	if !isExist {
		us.log.WithField("request", utils.StructToString(data.Email)).Errorf("Login | user not found")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "Email not registered",
				"id": "Email tidak terdaftar",
			}, errors.New("UserNA")
	}

	userData, err := us.db.User.GetUserByEmail(ctx, strings.ToLower(strings.ToLower(data.Email)))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data.Email)).WithError(err).Errorf("Login | fail to get company user data")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "There was an error in checking user data",
				"id": "Ada kesalahan dalam memeriksa data pengguna",
			}, errors.New("FailedServer")
	}

	isValid, err := utils.ComparePassword(userData.Password, data.Password)
	if err != nil {
		us.log.WithField("request", utils.StructToString(data.Email)).WithError(err).Errorf("Login | fail to compare password")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "There was an error changing the password",
				"id": "Ada kesalahan dalam mengubah kata sandi",
			}, errors.New("FailedServer")
	}

	if !isValid {
		us.log.WithField("request", utils.StructToString(data.Email)).Errorf("Login | incorrect password")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "You entered an incorrect password",
				"id": "Kata sandi yang anda masukkan salah",
			}, errors.New("FailedPassword")
	}

	session, err := utils.GetEncrypt([]byte(us.conf.KeyData.User), utils.StructToString(mu.CredentialData{
		ID:       userData.ID,
		Fullname: userData.Name,
		Email:    userData.Email,
	}))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data.Email)).Errorf("Login | fail generate session jwt")
		log.Println(err.Error())
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "internal server error",
				"id": "terjadi kesalahan sistem, silahkan coba dilain waktu ",
			}, errors.New("FailedServer")
	}

	generateTime := time.Now().UTC()

	accessToken, err := utils.GenerateJWT(session)
	if err != nil {
		us.log.WithField("request", utils.StructToString(data.Email)).Errorf("Login | fail generate jwt token")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "internal server error",
				"id": "terjadi kesalahan sistem, silahkan coba dilain waktu ",
			}, errors.New("FailedServer")
	}

	resLogin := mu.LoginResponse{
		Token: mu.JWTAccess{
			AccessToken:        accessToken,
			AccessTokenExpired: generateTime.Add(time.Duration(us.conf.Authorization.JWT.AccessTokenDuration) * time.Minute).Format(time.RFC3339),
		},
		NamaLengkap: userData.Name,
		UserID:      userData.ID,
	}

	us.log.WithField("request", utils.StructToString(data.Email)).WithError(nil).Errorf("Login | finish request login")

	return &resLogin, nil, nil
}
