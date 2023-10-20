/**
 * @Time: 2023/10/20 16:15
 * @Author: jzechen
 * @File: router.go
 * @Software: GoLand collector
 */

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jzechen/collector/pkg/manager/dto"
	"github.com/jzechen/collector/pkg/manager/server/engine/gin/response"
	"k8s.io/klog/v2"
	"net/http"
)

func (g *Engine) CreateHandler() http.Handler {
	g.Group("/v1/collector/manager")

	// TODO: add router handle here

	return g
}

func (g *Engine) hello(ctx *gin.Context) {
	klog.V(4).InfoS("hello")
	rsp, err := g.services.Hello(ctx.Request.Context(), &dto.NullRsp{})
	response.JsonResponse(ctx, "hello", rsp, err)
}
