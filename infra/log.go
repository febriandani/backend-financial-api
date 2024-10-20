package infra

import (
	"os"

	constants "github.com/febriandani/backend-financial-api/domain/constants/general"
	"github.com/febriandani/backend-financial-api/domain/model/general"
	"github.com/febriandani/backend-financial-api/domain/utils"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

var logger *logrus.Logger

// func NewLogger(conf *general.AppService) *logrus.Logger {
// 	if logger == nil {
// 		path := "log/"

// 		isExist, err := utils.DirExists(path)
// 		if err != nil {
// 			panic(err)
// 		}

// 		if !isExist {
// 			err = os.MkdirAll(path, os.ModePerm)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}

// 		writer, err := rotatelogs.New(
// 			path+conf.App.Name+"-"+"%Y%m%d.log",
// 			rotatelogs.WithMaxAge(-1),
// 			rotatelogs.WithRotationCount(constants.MaxRotationFile),
// 			rotatelogs.WithRotationTime(constants.LogRotationTime),
// 		)
// 		if err != nil {
// 			panic(err)
// 		}

// 		logger = logrus.New()

// 		logger.Hooks.Add(lfshook.NewHook(
// 			writer,
// 			&logrus.TextFormatter{
// 				DisableColors:   false,
// 				FullTimestamp:   true,
// 				TimestampFormat: constants.FullTimeFormat,
// 			},
// 		))

// 		// Set formatter for os.Stdout
// 		logger.SetFormatter(&logrus.TextFormatter{
// 			DisableColors:   false,
// 			FullTimestamp:   true,
// 			TimestampFormat: constants.FullTimeFormat,
// 		})

// 		logger.SetFormatter(&ecslogrus.Formatter{})

// 		return logger
// 	}

// 	return logger
// }

func NewLogger(conf *general.AppService) *logrus.Logger {
	if logger == nil {
		path := "log/"

		isExist, err := utils.DirExists(path)
		if err != nil {
			panic(err)
		}

		if !isExist {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		// Create a rotating log file writer
		writer, err := rotatelogs.New(
			path+conf.App.Name+"-"+"%Y%m%d.log",
			rotatelogs.WithMaxAge(-1),
			rotatelogs.WithRotationCount(constants.MaxRotationFile),
			rotatelogs.WithRotationTime(constants.LogRotationTime),
		)
		if err != nil {
			panic(err)
		}

		// Initialize the logger
		logger = logrus.New()

		// File logging hook (JSON format)
		logger.Hooks.Add(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.InfoLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.WarnLevel:  writer,
				logrus.DebugLevel: writer,
			},
			&logrus.JSONFormatter{
				TimestampFormat: constants.FullTimeFormat,
				PrettyPrint:     true, // Beautify JSON output
			},
		))
		// Set formatter for stdout (colorized, plain-text for development)
		if conf.App.Environtment == "beta" {
			logger.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: constants.FullTimeFormat,
				ForceColors:     true,  // Enable colors in terminal
				DisableColors:   false, // Set false if you want colored output in stdout
			})
		} else {
			// Use JSON formatter for production logs to stdout as well
			logger.SetFormatter(&ecslogrus.Formatter{
				PrettyPrint: true, // Enables nicely formatted JSON
			})
		}

		// Set log level based on the environment
		if conf.App.Environtment == "development" {
			logger.SetLevel(logrus.InfoLevel)
		} else {
			logger.SetLevel(logrus.DebugLevel)
		}

		return logger
	}

	return logger
}
