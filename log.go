package common

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MicroLog struct {
	C          *gin.Context
	Fields     map[string]interface{}
	Log        *logrus.Logger
	ContextLog *logrus.Entry
	Level      string
}

func New(level string, fields map[string]interface{}) *MicroLog {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logLevel, err := logrus.ParseLevel(level)

	if err != nil {
		logLevel, _ = logrus.ParseLevel("info")
	}

	logger.SetLevel(logLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return &MicroLog{
		Fields:     fields,
		Log:        logger,
		ContextLog: logger.WithFields(fields),
	}
}

func (ml MicroLog) SetField(key string, val interface{}) *MicroLog {
	ml.Fields[key] = val
	ml.ContextLog = ml.Log.WithFields(ml.Fields)
	return &ml
}

func (ml MicroLog) Logger() *logrus.Entry {
	return ml.Log.WithFields(ml.Fields)
}

//Message - log the message, deprecated
func (ml MicroLog) Message(message interface{}) {
	ml.Log.WithFields(ml.Fields).Info(message)
}

//Log with context - latest
func (ml MicroLog) Info(m interface{}, message ...interface{}) {
	ml.ContextLog.Info(m, message)
}

func (ml MicroLog) Debug(m interface{}, message ...interface{}) {
	ml.ContextLog.Debug(m, message)
}

func (ml MicroLog) Warn(m interface{}, message ...interface{}) {
	ml.ContextLog.Warn(m, message)
}

func (ml MicroLog) Error(m interface{}, message ...interface{}) {
	ml.ContextLog.Error(m, message)
}

func (ml MicroLog) Fatal(m interface{}, message ...interface{}) {
	ml.ContextLog.Fatal(m, message)
}

func (ml MicroLog) Panic(m interface{}, message ...interface{}) {
	ml.ContextLog.Panic(m, message)
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
