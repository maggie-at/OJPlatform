package router

import (
	_ "OJPlatform/docs"
	"OJPlatform/middlewares"
	"OJPlatform/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	ginServer := gin.Default()

	// Swagger配置
	ginServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 测试gin
	ginServer.GET("/ping", service.Ping)

	// 公开方法
	ginServer.GET("/problem/list", service.GetProblemList)
	ginServer.GET("/problem/detail", service.GetProblemDetail)

	ginServer.POST("/user/login", service.UserLogin)
	ginServer.POST("/user/register", service.UserRegister)
	ginServer.POST("/user/send_code", service.SendCodeMail)
	ginServer.GET("/user/detail", service.GetUserDetail)
	ginServer.GET("/user/rank", service.GetRankList)

	ginServer.GET("/submission/list", service.GetSubmissionList)

	// 管理员方法
	adminGroup := ginServer.Group("/admin", middlewares.AuthAdmin())
	{
		adminGroup.POST("/problem/create", service.ProblemCreate)
		adminGroup.PUT("/problem/modify", service.ProblemModify)
		adminGroup.GET("/category/list", service.GetCategoryList)
		adminGroup.POST("/category/create", service.CategoryCreate)
		adminGroup.PUT("/category/modify", service.CategoryModify)
		adminGroup.DELETE("/category/delete", service.CategoryDelete)
	}

	// 用户私有方法
	userAuthGroup := ginServer.Group("/user", middlewares.AuthUser())
	{
		userAuthGroup.POST("/submit", service.SubmitAnswer)
	}
	return ginServer
}
