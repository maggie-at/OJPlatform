package service

import (
	"OJPlatform/define"
	"OJPlatform/helper"
	"OJPlatform/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "user_identity"
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /user/detail [get]
func GetUserDetail(context *gin.Context) {
	identity := context.Query("identity")
	if identity == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户唯一标识不能为空",
		})
		return
	}
	u, err := models.GetUserDetail(identity)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "未找到相关用户",
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetUserDetail - " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": u,
	})
}

// GetRankList
// @Tags 公共方法
// @Summary 用户排行榜
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /user/rank [get]
func GetRankList(context *gin.Context) {
	// 接收QueryString参数, 用于分页参数
	size, err := strconv.Atoi(context.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("GetProblemList - Size Parse Error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetProblemList - Page Parse Error: " + err.Error(),
		})
		return
	}
	page, err := strconv.Atoi(context.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetProblemList - Page Parse Error:", err)
		return
	}
	page = page - 1

	var count int64

	tx := models.GeuUserRank()
	// 分页
	var uList []models.User
	err = tx.Count(&count).Offset(page * size).Limit(size).Find(&uList).Error
	if err != nil {
		log.Println("GetRankList - Pagination Error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetRankList - Pagination Error: " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"count": count,
			"data":  uList,
		},
	})
}

// UserLogin
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /user/login [post]
func UserLogin(context *gin.Context) {
	username := context.PostForm("username")
	password := context.PostForm("password")
	if username == "" || password == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息缺失",
		})
		return
	}
	passwordMD5 := helper.GetMd5(password)
	userInfo, err := models.UserLogin(username, passwordMD5)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名或密码错误",
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "UserLogin - " + err.Error(),
		})
		return
	}
	// 调用Helper生成token
	token, err := helper.GenerateToken(userInfo.Identity, userInfo.Name, userInfo.IsAdmin)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GenerateToken - " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"token":     token,
			"user_info": userInfo,
		},
	})
}

// SendCodeMail
// @Tags 公共方法
// @Summary 发送验证码邮件
// @Param email formData string false "email"
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /user/send_code [post]
func SendCodeMail(context *gin.Context) {
	emailAddress := context.PostForm("email")
	if emailAddress == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱为空",
		})
		return
	}
	// 生成随机验证码 -> 保存到redis & 发送邮件
	code := helper.GenerateRandomCode()
	models.RDB.Set(context, emailAddress, code, time.Second*300)
	err := helper.SendCodeMail(emailAddress, code)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码发送失败 - " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功, 请注意接收",
	})
}

// UserRegister
// @Tags 公共方法
// @Summary 用户注册
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string false "phone"
// @Param mail formData string true "mail"
// @Param code formData string true "code"
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /user/register [post]
func UserRegister(context *gin.Context) {
	name := context.PostForm("name")
	password := context.PostForm("password")
	phone := context.PostForm("phone")
	mail := context.PostForm("mail")
	inputCode := context.PostForm("code")
	if name == "" || password == "" || mail == "" || inputCode == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	// 比较用户输入的验证码和Redis中的验证码
	rdbCode, err := models.RDB.Get(context, mail).Result()
	if err != nil {
		log.Printf("Get Code Error - %v", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确, 请重新获取验证码",
		})
		return
	}
	if inputCode != rdbCode {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码错误",
		})
		return
	}
	valid, err := models.CheckEmailExist(mail)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "CheckEmailExist - " + err.Error(),
		})
		return
	}
	if valid > 0 {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该邮箱已被注册",
		})
		return
	}
	// 将密码转为md5
	pwdMd5 := helper.GetMd5(password)
	// 生成uuid作为identity
	identity := helper.GetUUID()
	user := models.User{
		Identity: identity,
		Name:     name,
		Password: pwdMd5,
		Email:    mail,
		Phone:    phone,
	}
	err = models.UserRegister(user)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "UserRegister - Insert Error - " + err.Error(),
		})
		return
	}
	// 生成token
	token, err := helper.GenerateToken(identity, name, 0)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "UserRegister - GenerateToken Error - " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "用户注册成功",
		"data": map[string]interface{}{
			"token": token,
		},
	})
}
