package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	_ "math/rand"
	"net/http"
	_ "time"
	"xlk/ginessential/common"
	"xlk/ginessential/dto"
	"xlk/ginessential/model"
	"xlk/ginessential/response"
	"xlk/ginessential/util"
)

func Register(ctx *gin.Context) {
	//获取参数
	DB := common.GetDB()
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//数据验证
	if len(telephone) != 11 { //手机号不是11位
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须是11位")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"手机号必须是11位"})
		return
	} //这里的H其实是type H map[string]interface{} ，是个别名,他是个map，KEY是字符串，值为接口，可以改为map[string]interface{}{"code":422,"msg":"手机号必须是11位"}
	//这里的422常量：响应状态码表示服务器理解请求实体的内容类型，请求实体的语法正确，但无法处理包含的指令。

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能小于6位")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"密码不能小于6位"})
		return
	}
	//名称没有传，给一个10位的随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Println(name, telephone, password)
	//判断手机号是否存
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"用户已经存在"})
		return

	}

	//创建用户
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		//ctx.JSON(http.StatusInternalServerError,gin.H{"code":500,"msg":"加密错误"})//内部错误500
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashPassword),
	}
	DB.Create(&newUser)

	//返回结果
	response.Success(ctx, nil, "注册成功")
	/*ctx.JSON(200, gin.H{
		"msg": "注册成功",
	})*/
}
func Login(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数

	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(telephone) != 11 { //手机号不是11位
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须是11位")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"手机号必须是11位"})
		return
	} //这里的H其实是type H map[string]interface{} ，是个别名,他是个map，KEY是字符串，值为接口，可以改为map[string]interface{}{"code":422,"msg":"手机号必须是11位"}
	//这里的422常量：响应状态码表示服务器理解请求实体的内容类型，请求实体的语法正确，但无法处理包含的指令。

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能小于6位")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"密码不能小于6位"})
		return
	}

	//判断手机号是否存在
	var user model.User //前面定义了结构体User
	db.Where("telephone =?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"用户不存在"})
		return
	}
	//判断密码是否正确
	//这里要加密，所以创建用户的时候要加密码属性
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		//这个hash比较函数，第一个参数是原始密码hash后的密码，第二个参数是要对比的密码
		//这个hashcomp函数会返回密码，如果加密成功的话。反而，加密不成功err就是空
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		//ctx.JSON(http.StatusBadRequest,gin.H{"code":400,"msg":"密码错误"})
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		/*ctx.JSON(http.StatusInternalServerError,gin.H{
			"code":500,
			"msg":"系统异常",
		})*/
		log.Printf("token generate error:%v", err) //记录日志
	}

	//返回结果
	/*ctx.JSON(http.StatusOK,gin.H{
		"code": 200,
		"data": gin.H{"token":token},
		"msg":"登录成功",
	})*/
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User //前面定义了结构体User
	db.Where("telephone =?", telephone).First(&user)
	if user.ID != 0 {
		return true
	} else {
		return false
	}
}
