package common

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MicroLog struct {
	C      *gin.Context
	Fields map[string]interface{}
	Log    *logrus.Logger
	Level  string
}

//Log with context data
func (ml MicroLog) Message(message interface{}) {
	ml.Fields["sample"] = ml.C.Request.Header.Get("User-Agent")
	//Log format
	ml.Log.WithFields(ml.Fields).Info(message)
}

//Log with context data
func Log(c *gin.Context, message string) {
	log := c.MustGet("log").(*logrus.Logger)

	reqMethod := c.Request.Method
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()
	tenantID := c.Request.Header.Get("Tenantid")
	domain := c.Request.Header.Get("Domain")
	userAgent := c.Request.Header.Get("User-Agent")
	traceID := c.Request.Header.Get("X-B3-Traceid")
	spanID := c.Request.Header.Get("X-B3-Spanid")

	//Log format
	log.WithFields(logrus.Fields{
		"domain":    domain,
		"tenantId":  tenantID,
		"status":    statusCode,
		"userAgent": userAgent,
		"traceId":   traceID,
		"spanId":    spanID,
		"ip":        clientIP,
		"method":    reqMethod,
		"uri":       c.Request.RequestURI,
	}).Info(message)
}
