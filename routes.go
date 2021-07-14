package main

import (
	"github.com/gin-gonic/gin"
	"xlk/ginessential/controller"
	"xlk/ginessential/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)                    //注册
	r.POST("/api/auth/login", controller.Login)                          //登录
	r.GET("api/auth/info", middleware.AuthMiddleware(), controller.Info) //用户信息,用中间件保护
	return r
}
