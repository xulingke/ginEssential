package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
	_"github.com/go-sql-driver/mysql"
)

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(110);not null unique"`
	Password string `gorm:"size:255;not null"`
}

func main(){
	db:=InitDB()
	defer db.Close()
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		//获取参数
			name:=ctx.PostForm("name")
			telephone:=ctx.PostForm("telephone")
			password:=ctx.PostForm("password")
		//数据验证
			if len(telephone)!=11{    //手机号不是11位
				ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"手机号必须是11位"})
				return
			}//这里的H其实是type H map[string]interface{} ，是个别名,他是个map，KEY是字符串，值为接口，可以改为map[string]interface{}{"code":422,"msg":"手机号必须是11位"}
               //这里的422常量：响应状态码表示服务器理解请求实体的内容类型，请求实体的语法正确，但无法处理包含的指令。

            if len(password)<6{
            	ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"密码不能小于6位"})
				return
			}
			//名称没有传，给一个10位的随机字符串
			if len(name)==0{
				name=RandomString(10)
			}
			log.Println(name,telephone,password)
               //判断手机号是否存
 			if isTelephoneExist(db,telephone){
				ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"用户已经存在"})
				return

		}

		//创建用户
		newUser:=User{
			Name:	name,
			Telephone:  telephone,
			Password:	password,
		}
		db.Create(&newUser)


		//返回结果

		ctx.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	panic(r.Run())// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func RandomString(n int)string{//随机n位字符串
	var letters=[]byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result:=make([]byte,n)//创建一个byte数组，长度为n，名为result

	rand.Seed(time.Now().Unix())//根据时间戳获得随机数种子
	for i:=range result{
		result[i]=letters[rand.Intn(len(letters))]//rand是一个库，调用Intn方法，这个方法随机一个数
	}
	return string(result)
}

func InitDB()*gorm.DB{
	driverName:="mysql"
	host:="localhost"
	port:="3306"
	database:="ginessential"
	username:="root"
	password:="deskjei3538"
	charset:="utf8mb4"
	args:=fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db,err:=gorm.Open(driverName,args)
	if err!=nil{
		panic("failed to connect database,err: "+err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}

func isTelephoneExist(db *gorm.DB,telephone string)bool {
	var user User //前面定义了结构体User
	db.Where("telephone =?", telephone).First(&user)
	if user.ID != 0 {
		return true
	} else {
		return false
	}
}