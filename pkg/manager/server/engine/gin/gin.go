/**
 * @Time: 2023/10/20 16:17
 * @Author: jzechen
 * @File: gin.go
 * @Software: GoLand collector
 */

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jzechen/toresa/pkg/manager/server/engine"
	"github.com/jzechen/toresa/pkg/manager/server/engine/gin/middleware"
	"github.com/jzechen/toresa/pkg/manager/services"
	"io"
	"k8s.io/klog/v2"
)

type Engine struct {
	*gin.Engine
	services services.Interface
}

func NewGinEngine(svc services.Interface) engine.Interface {
	e := &Engine{
		Engine:   createEngine(),
		services: svc,
	}

	return e
}

func createEngine() *gin.Engine {
	if !klog.V(4).Enabled() {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	e := gin.New()
	if klog.V(4).Enabled() {
		e.Use(middleware.LoggerHandler())
	}

	// register middleware
	e.Use(middleware.Middlewares...)

	return e
}
