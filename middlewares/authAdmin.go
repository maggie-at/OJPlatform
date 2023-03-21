package middlewares

import (
	"OJPlatform/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthAdmin 验证管理员身份的中间件
func AuthAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 从 header 中获得 token 并解析
		auth := context.GetHeader("Authorization")
		userClaims, err := helper.AnalyseToken(auth)
		// 拦截1: 解析失败
		if err != nil {
			// 阻止后面HandlerFunc的执行, 不影响当前HandlerFunc
			context.Abort()
			context.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "Unauthorized",
			})
			return
		}
		// 拦截2: 没有 admin 权限
		if userClaims == nil || userClaims.IsAdmin != 1 {
			context.Abort()
			context.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "Unauthorized",
			})
			return
		}
		context.Next()
	}
}
