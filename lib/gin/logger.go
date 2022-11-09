package gin

import (
	"fmt"
	"time"

	"github.com/AndryHardiyanto/dealltest/lib/errors"
	"github.com/AndryHardiyanto/dealltest/lib/log"
	"github.com/gin-gonic/gin"
)

type GinLogger struct {
	Path     string
	Latency  time.Duration
	Method   string
	ClientIP string
	Request  string
}

func RegisterGinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now().In(time.FixedZone("Asia/Jakarta", 7*60*60))
		// before request
		path := c.Request.URL.Path
		request := NewRequest(c).GetRequest()

		c.Next()

		if len(c.Errors.Errors()) > 0 {
			for _, e := range c.Errors {
				// 	// Return code and MSG if there is a custom error
				if err, ok := e.Err.(*errors.Error); ok {
					logging(request, c, err, path, fmt.Sprintf("%v", time.Since(t)), err.Type)
				}
			}
		} else {
			logging(request, c, nil, path, fmt.Sprintf("%v", time.Since(t)), "success")
		}
	}
}

func logging(request interface{}, c *gin.Context, err *errors.Error, path, latency, errType string) {
	switch errType {
	case "success":
		log.Info().
			Str("path", path).
			Interface("request", request).
			Str("latency", latency).Msg("")
	default:
		log.Error().
			Str("path", path).
			Interface("request", request).
			Str("latency", latency).
			Str("error_type", errType).Msg(err.Error())
	}
}
