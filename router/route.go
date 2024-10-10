package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qyvx/controller"
	"qyvx/logger"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.POST("/", controller.HookHandler)

	//注册业务路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("/invite", controller.InviteHandler)
		v1.GET("/update_users", controller.UpdateUsersHandler)
	}

	//v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	//{
	//
	//}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
