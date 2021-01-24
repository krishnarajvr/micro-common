package common

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//Log with context data
func Log(c *gin.Context, message string) {
	log := c.MustGet("log").(*logrus.Logger)

	reqMethod := c.Request.Method
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()
	tenantID, _ := c.Get("tenantId")

	//Log format
	log.WithFields(logrus.Fields{
		"statusCode": statusCode,
		"clientIp":   clientIP,
		"reqMethod":  reqMethod,
		"reqUri":     c.Request.RequestURI,
		"tenantId":   tenantID,
	}).Info(message)
}
