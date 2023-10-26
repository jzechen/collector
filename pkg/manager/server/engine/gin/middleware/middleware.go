/**
 * @Time: 2023/10/20 16:19
 * @Author: jzechen
 * @File: middleware.go
 * @Software: GoLand collector
 */

package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jzechen/toresa/pkg/manager/server/rate"
	"github.com/jzechen/toresa/pkg/manager/utils/errcode"
	"k8s.io/klog/v2"
	"net/http"
	"time"
)

var Middlewares = []gin.HandlerFunc{gin.Recovery(), RateLimiter(), ErrorHandler()}

func LoggerHandler() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		str := fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
		klog.V(4).Infof(str)
		return str
	})
}

func RateLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limiter := rate.GetRateLimiter()
		if limiter.Allow() {
			ctx.Next()
		} else {
			ctx.Abort()
		}
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// check gin return error
		for _, e := range ctx.Errors {
			err := e.Err
			// our define error
			// return it to client
			if ec, ok := err.(errcode.ErrCode); ok {
				if ec == errcode.NotFound {
					ctx.JSON(http.StatusNotFound, gin.H{"code": ec, "msg": ec.Error()})
					return
				}
				ctx.JSON(http.StatusBadRequest, gin.H{
					"code": ec,
					"msg":  ec.Error(),
				})
			} else {
				// not our define error
				// print in server log and not tell client error detail
				klog.Error(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "server error",
				})
			}

			return
		}
	}
}
