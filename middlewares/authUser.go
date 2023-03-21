package middlewares

import (
	"OJPlatform/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthUser 验证用户身份的中间件
func AuthUser() gin.HandlerFunc {
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
		// 拦截2: userClaim为空
		if userClaims == nil {
			context.Abort()
			context.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "Unauthorized",
			})
			return
		}
		context.Set("user", userClaims)
		context.Next()
	}
}
