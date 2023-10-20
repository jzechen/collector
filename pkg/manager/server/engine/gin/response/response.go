/**
 * @Time: 2023/10/20 16:18
 * @Author: jzechen
 * @File: response.go
 * @Software: GoLand collector
 */

package response

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"net/http"
)

func JsonResponse(ctx *gin.Context, module string, data interface{}, err error) {
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "Success", "data": data, "module": module})
}

func CheckError(ctx *gin.Context, code int, module, msg string, err error) {
	if err != nil {
		klog.ErrorS(err, msg, "code", code, "module", module)
		ctx.JSON(code, gin.H{"code": code, "msg": msg, "module": module})
		panic(err)
	}
}

func FailWithMsg(ctx *gin.Context, code int, module, msg string) {
	klog.Error(msg)
	ctx.Abort()
	ctx.JSON(code, gin.H{"code": code, "msg": msg, "module": module})
}
