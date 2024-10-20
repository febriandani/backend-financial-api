package main

import (
	"fmt"
	"net/http"

	"github.com/febriandani/backend-financial-api/cmd/api/routes"
	mg "github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/handler/api"
	"github.com/febriandani/backend-financial-api/infra"
	"github.com/febriandani/backend-financial-api/repository"
	"github.com/febriandani/backend-financial-api/service"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	conf, err := getConfigKey()
	if err != nil {
		panic(err)
	}

	handler, log, err := newRepoContext(conf)
	if err != nil {
		logrus.WithError(err).Errorf("error")
		panic(err)
	}

	headers := handlers.AllowedHeaders(conf.Route.Headers)
	methods := handlers.AllowedMethods(conf.Route.Methods)
	origins := handlers.AllowedOrigins(conf.Route.Origins)
	credentials := handlers.AllowCredentials()

	router := routes.GetCoreEndpoint(conf, handler, log)

	port := fmt.Sprintf(":%s", conf.App.Port)
	log.Info("server listen to port ", port)
	log.Fatal(http.ListenAndServe(port, handlers.CORS(headers, methods, origins, credentials)(router)))
}

func getConfigKey() (*mg.AppService, error) {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	conf := mg.AppService{
		App: mg.AppAccount{
			Name:         viper.GetString("APP.NAME"),
			Environtment: viper.GetString("APP.ENV"),
			URL:          viper.GetString("APP.URL"),
			Port:         viper.GetString("APP.PORT"),
		},
		Route: mg.RouteAccount{
			Methods: viper.GetStringSlice("ROUTE.METHODS"),
			Headers: viper.GetStringSlice("ROUTE.HEADERS"),
			Origins: viper.GetStringSlice("ROUTE.ORIGIN"),
		},
		Database: mg.Database{
			Read: mg.DBDetail{
				Username:     viper.GetString("DATABASE.READ.USERNAME"),
				Password:     viper.GetString("DATABASE.READ.PASSWORD"),
				URL:          viper.GetString("DATABASE.READ.URL"),
				Port:         viper.GetString("DATABASE.READ.PORT"),
				DBName:       viper.GetString("DATABASE.READ.DB_NAME"),
				MaxIdleConns: viper.GetInt("DATABASE.READ.MAXIDLECONNS"),
				MaxOpenConns: viper.GetInt("DATABASE.READ.MAXOPENCONNS"),
				MaxLifeTime:  viper.GetInt("DATABASE.READ.MAXLIFETIME"),
				Timeout:      viper.GetString("DATABASE.READ.TIMEOUT"),
				SSLMode:      viper.GetString("DATABASE.READ.SSL_MODE"),
			},
			Write: mg.DBDetail{
				Username:     viper.GetString("DATABASE.WRITE.USERNAME"),
				Password:     viper.GetString("DATABASE.WRITE.PASSWORD"),
				URL:          viper.GetString("DATABASE.WRITE.URL"),
				Port:         viper.GetString("DATABASE.WRITE.PORT"),
				DBName:       viper.GetString("DATABASE.WRITE.DB_NAME"),
				MaxIdleConns: viper.GetInt("DATABASE.WRITE.MAXIDLECONNS"),
				MaxOpenConns: viper.GetInt("DATABASE.WRITE.MAXOPENCONNS"),
				MaxLifeTime:  viper.GetInt("DATABASE.WRITE.MAXLIFETIME"),
				Timeout:      viper.GetString("DATABASE.WRITE.TIMEOUT"),
				SSLMode:      viper.GetString("DATABASE.WRITE.SSL_MODE"),
			},
		},
		Authorization: mg.AuthAccount{
			JWT: mg.JWTCredential{
				IsActive:             viper.GetBool("AUTHORIZATION.JWT.IS_ACTIVE"),
				AccessTokenSecretKey: viper.GetString("AUTHORIZATION.JWT.ACCESS_TOKEN_SECRET_KEY"),
				AccessTokenDuration:  viper.GetInt("AUTHORIZATION.JWT.ACCESS_TOKEN_DURATION"),
			},
		},
		KeyData: mg.KeyAccount{
			User: viper.GetString("KEY.USER"),
		},
	}

	return &conf, nil
}

func newRepoContext(conf *mg.AppService) (api.Handler, *logrus.Logger, error) {
	var handler api.Handler

	// Init Log
	logger := infra.NewLogger(conf)

	// Init DB Read Connection.
	dbRead := infra.NewDB(logger)
	dbRead.ConnectDB(&conf.Database.Read)
	if dbRead.Err != nil {
		return handler, logger, dbRead.Err
	}

	// Init DB Write Connection.
	dbWrite := infra.NewDB(logger)
	dbWrite.ConnectDB(&conf.Database.Write)
	if dbWrite.Err != nil {
		return handler, logger, dbWrite.Err
	}

	dbList := &infra.DatabaseList{
		Backend: infra.DatabaseType{
			Read:  &dbRead,
			Write: &dbWrite,
		},
	}

	repo := repository.NewRepo(dbList, logger)
	usecase := service.NewService(repo, *conf, dbList, logger)
	handler = api.NewHandler(usecase, *conf, logger)

	return handler, logger, nil
}
