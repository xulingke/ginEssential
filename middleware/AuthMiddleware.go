package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"xlk/ginessential/common"
	"xlk/ginessential/model"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authoerization header
		tokenString := ctx.GetHeader("Authorization")

		//validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			//if tokenstring为空，或者不是Bearer为开头，则返回权限不足，就是格式不对的处理
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足0",
			})
			ctx.Abort() //抛弃掉这一次请求
			return
		}
		//如果格式正确
		tokenString = tokenString[7:] //前面设置了头bearer为七位，所以从七位开始截取有效信息

		//解析token
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid { //解析失败或者解析权限不足，返回权限不足
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足1",
			})
			ctx.Abort()
			return
		}
		//否则通过验证，则获取claim中的userId

		DB := common.GetDB()
		var user model.User
		DB.First(&user, claims.UserId)

		//用户不存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足2",
			})
			ctx.Abort() //抛弃掉这一次请求
			return
		}
		//用户存在将user信息写进上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
