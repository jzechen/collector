/**
 * @Time: 2023/10/20 16:15
 * @Author: jzechen
 * @File: router.go
 * @Software: GoLand collector
 */

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jzechen/toresa/pkg/manager/dto"
	"github.com/jzechen/toresa/pkg/manager/server/engine/gin/response"
	"k8s.io/klog/v2"
	"net/http"
)

func (g *Engine) CreateHandler() http.Handler {
	sina := g.Group("/v1/collector/manager/sina")
	sina.POST("/login", g.login)
	// TODO: add router handle here

	return g
}

func (g *Engine) hello(ctx *gin.Context) {
	klog.V(4).InfoS("hello")
	rsp, err := g.services.Hello(ctx.Request.Context(), &dto.NullRsp{})
	response.JsonResponse(ctx, "hello", rsp, err)
}

func (g *Engine) login(ctx *gin.Context) {
	req := dto.LoginReq{}
	err := ctx.ShouldBind(&req)
	response.CheckError(ctx, http.StatusInternalServerError, "sina/login", "validate", err)

	klog.V(4).InfoS("sina/login", "req", req)
	rsp, err := g.services.Login(ctx.Request.Context(), &req)
	response.JsonResponse(ctx, "sina/login", rsp, err)
}
