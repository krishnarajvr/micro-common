package middleware

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/krishnarajvr/micro-common"
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

	// //Set rotatelogs
	// logWriter, err := rotatelogs.New(
	// 	//Split file name
	// 	logFile+".%Y%m%d.log",
	// 	//Generate soft chain, point to the latest log file
	// 	rotatelogs.WithLinkName(logFile),
	// 	//Set maximum save time (7 days)
	// 	rotatelogs.WithMaxAge(7*24*time.Hour),
	// 	//Set log cutting interval (1 day)
	// 	rotatelogs.WithRotationTime(24*time.Hour),
	// )

	// writeMap := lfshook.WriterMap{
	// 	logrus.InfoLevel:  logWriter,
	// 	logrus.FatalLevel: logWriter,
	// 	logrus.DebugLevel: logWriter,
	// 	logrus.WarnLevel:  logWriter,
	// 	logrus.ErrorLevel: logWriter,
	// 	logrus.PanicLevel: logWriter,
	// }

	// lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// })

	// //New hook
	// logger.AddHook(lfHook)
	/**
	Sample Json Structure

		{
	    "timestamp"     : "2004-02-12T15:19:21.000+00:00",
	    "domain"        : "health",
	    "tenantId"        : "1",
	    "application"   : "app1",
	    "module"        : "catalog-service",
	    "component"     : "product",
	    "level"         : "info",
	    "message"       : "Updated product",
	    "sessionId"     : "b222ba40-1aec-47b1-a917-70262a487dbe",
	    "ip"            : "192.168.1.1",
	    "userAgent"     : "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36  Chrome/51.0.2704.103 Safari/537.36",
	    "resource"      : "product",
	    "request"       : "Request object",
	    "service"       : "productService",
	    "appVersion"    : "v1",
	    "traceId"       : "db176f6b-b4c6-4a59-ac62-4d471eedf354",
	    "spanId"        : "db176f6b-b4c6-4a59-ac62-4d471eedf354",
	    "server"        : "i0e5b5c88",
	    "serverLocation": "us-east-1",
	    "protocol"      : "https",
	    "method"        : "PUT",
	    "status"        : "200",
	    "source"        : "/tmp/app.log",
	    "duration"      : 500,
	    "stack_trace"   : "trace object"
	}
	*/

	return func(c *gin.Context) {

		startTime := time.Now()

		reqMethod := c.Request.Method
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		tenantID := c.Request.Header.Get("Tenantid")
		domain := c.Request.Header.Get("Domain")
		userAgent := c.Request.Header.Get("User-Agent")
		traceID := c.Request.Header.Get("X-B3-Traceid")
		spanID := c.Request.Header.Get("X-B3-Spanid")

		fields := logrus.Fields{
			"domain":    domain,
			"tenantId":  tenantID,
			"status":    statusCode,
			"userAgent": userAgent,
			"traceId":   traceID,
			"spanId":    spanID,
			"ip":        clientIP,
			"method":    reqMethod,
			"uri":       c.Request.RequestURI,
		}

		microLog := common.MicroLog{
			C:      c,
			Fields: fields,
			Log:    logger,
		}

		c.Set("log", &microLog)

		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		fields["latency"] = latencyTime
		//Log format
		logger.WithFields(fields).Info()
	}
}
