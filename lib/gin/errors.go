package gin

import (
	"github.com/AndryHardiyanto/dealltest/lib/errors"
	"github.com/AndryHardiyanto/dealltest/lib/response"
	"github.com/gin-gonic/gin"
)

func RegisterErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// call c.ext () to execute the middleware
		// Start from here after all middleware and router processing is complete
		// Check for errors in c. elors
		for _, e := range c.Errors {
			// Return code and MSG if there is a custom error
			if err, ok := e.Err.(*errors.Error); ok {
				response.New(c).Error(err)
			}
		}
	}
}
