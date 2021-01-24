package middleware

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

//LoggerToFile Middleware
func LoggerToFile(filePath string, fileName string) gin.HandlerFunc {

	logFile := path.Join(filePath, fileName)

	// Create log folder
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, 0700)
	}

	src, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("err", err)
	}

	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	//Set rotatelogs
	logWriter, err := rotatelogs.New(
		//Split file name
		fileName+".%Y%m%d.log",
		//Generate soft chain, point to the latest log file
		rotatelogs.WithLinkName(fileName),
		//Set maximum save time (7 days)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		//Set log cutting interval (1 day)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	//New hook
	logger.AddHook(lfHook)

	return func(c *gin.Context) {
		startTime := time.Now()

		c.Set("log", logger)
		c.Next()

		endTime := time.Now()
		reqMethod := c.Request.Method
		latencyTime := endTime.Sub(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		//Log format
		logger.WithFields(logrus.Fields{
			"statusCode":  statusCode,
			"latencyTime": latencyTime,
			"clientIp":    clientIP,
			"reqMethod":   reqMethod,
			"reqUri":      c.Request.RequestURI,
		}).Info()

	}
}
