package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/krishnarajvr/micro-common"
	"github.com/sirupsen/logrus"
)

//LoggerToFile Middleware
func LoggerToFile(filePath string, fileName string) gin.HandlerFunc {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

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
		tenantID := c.Request.Header.Get("X-Tenant-Id")
		userID := c.Request.Header.Get("X-User-Id")
		domain := c.Request.Header.Get("X-Domain")
		userAgent := c.Request.Header.Get("User-Agent")
		traceID := c.Request.Header.Get("X-B3-Traceid")
		spanID := c.Request.Header.Get("X-B3-Spanid")
		serviceName := os.Getenv("SERVICE_NAME")

		if len(serviceName) == 0 {
			serviceName = "micro"
		}

		fields := logrus.Fields{
			"app":       "micro",
			"domain":    domain,
			"tenantId":  tenantID,
			"userId":    userID,
			"status":    statusCode,
			"userAgent": userAgent,
			"traceId":   traceID,
			"spanId":    spanID,
			"ip":        clientIP,
			"method":    reqMethod,
			"uri":       c.Request.RequestURI,
			"service":   serviceName,
		}

		microLog := common.MicroLog{
			C:      c,
			Fields: fields,
			Log:    logger,
		}

		c.Set("log", &microLog)
		//Process Request
		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		fields["latency"] = latencyTime
		//Log format
		logger.WithFields(fields).Info()
	}
}
