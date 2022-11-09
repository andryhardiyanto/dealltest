package gin

import (
	"github.com/AndryHardiyanto/dealltest/lib/errors"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, err any) {
		c.Error(errors.NewError(err.(error).Error()).SetType(errors.TypePanic))
		c.Abort()
	})
}
