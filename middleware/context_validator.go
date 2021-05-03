package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	common "github.com/krishnarajvr/micro-common"
)

//TenantValidator Middleware
func TenantValidator(excludeList map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		public := false

		if _, ok := excludeList[c.Request.URL.Path]; ok {
			public = true
		}

		if !public && strings.Contains(c.Request.URL.Path, "/swagger/") {
			public = true
		}

		if !public && strings.Contains(c.Request.URL.Path, "/thirdpartySwagger/") {
			public = true
		}

		if !public {
			tenantID := c.Request.Header.Get("X-Tenant-Id")

			if len(tenantID) == 0 {
				common.AccessDenied(c, "")
				c.Abort()
				return
			}

			c.Set("tenantId", tenantID)
		}

		c.Next()
	}
}

//VendorValidator Middleware
func VendorValidator(excludeList map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		public := false

		if _, ok := excludeList[c.Request.URL.Path]; ok {
			public = true
		}

		if !public && strings.Contains(c.Request.URL.Path, "/swagger/") {
			public = true
		}

		if !public && strings.Contains(c.Request.URL.Path, "/thirdpartySwagger/") {
			public = true
		}

		if !public {
			vendorID := c.Request.Header.Get("X-Reference-Id")

			if len(vendorID) == 0 {
				common.AccessDenied(c, "")
				c.Abort()
				return
			}

			c.Set("vendorId", vendorID)
		}

		c.Next()
	}
}
