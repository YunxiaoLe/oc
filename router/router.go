package router

import (
	"github.com/OverClockTeam/OC-WordsInSongs-BE/consts"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/middleware"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/router/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UseMyRouter(r *gin.Engine) {
	//示范写法 具体函数不得在路由写 只是这里足够简单罢了
	hello := r.Group("/hello")
	{
		hello.GET("/go", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "Hello Go!",
			})
			return
		})
	}

	songApi := r.Group("/api")
	user := songApi.Group("/auth")
	{
		user.POST("/register", middleware.IpRateLimiter("post/auth/register", consts.EXTREMLY_LOW_RATE),
			api.RegisterWithoutVerifyCode)
		user.POST("/isSessionExpired", middleware.IpRateLimiter("post/auth/isSessionExpired", consts.LOW_RATE),
			api.IsTokenExpired)
		user.POST("/sendVerifyCode", middleware.IpRateLimiter("post/auth/sendVerifyCode", consts.EXTREMLY_LOW_RATE),
			api.SendVerifyCodeHandler)
		user.POST("/login", middleware.IpRateLimiter("post/auth/login", consts.EXTREMLY_LOW_RATE),
			api.Login)

	}
}
