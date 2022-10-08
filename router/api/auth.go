package api

import (
	"fmt"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/business"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/client/emailclt"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/client/redisclt"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/consts"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/model"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/services/http_param"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/services/logit"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"strings"
)

// 不需要邮件验证码即可注册
func RegisterWithoutVerifyCode(context *gin.Context) {
	param := http_param.RegisterWithoutVerifyCodeArgument{}
	svc := business.NewService(context)
	err := context.ShouldBind(&param)
	if err != nil {
		logit.LogWithError("参数不全", map[string]interface{}{
			"Error": err,
		})
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数不全",
		})
		return
	}
	userEmail := util.CheckHustEmailSuffixOrAddIt(strings.ToLower(param.Email))
	fmt.Println("email: ", userEmail)
	//password in request is encrypted by rsa,so I should decrypt it first
	//password, err := util.RsaDecryption(param.Password)
	//if err != nil {
	//	context.JSON(http.StatusBadRequest, gin.H{
	//		"msg": "注册失败1",
	//	})
	//	return
	//}
	/*if !util.IsEmailAllowedRegister(email) {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": "你必须用@hust.edu.cn后缀的邮箱注册",
		})
		return
	}*/
	aesEmail, err := util.EncryptWithAes(userEmail)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": "注册失败",
		})
		return
	}

	if u, _ := svc.GetUserByEmailAes(aesEmail); u.IsEmailActivated {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": "该用户已经注册",
		})
		return
	}
	passwordHash := util.HashWithSalt(param.Password)
	user := model.AuthUser{
		Email:             aesEmail,
		Password:          passwordHash,
		RegisterTimestamp: util.GetTimeStamp(),
		IsEmailActivated:  false,
		Role:              consts.USER}

	if e := svc.Register(&user); e.Status == 400 {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": "该用户已经注册",
		})
		return
	} else if e.Status == 200 {
		token, _ := util.GenerateTokenByJwt(aesEmail, consts.USER)
		context.SetCookie(consts.COOKIE_NAME, token, consts.EXPIRE_TIME_TOKEN, "/", "localhost", false, true)
		context.JSON(http.StatusOK, gin.H{
			"email": user.Email,
			"msg":   "注册成功",
		})
	}
}
func IsTokenExpired(context *gin.Context) {
	tokenString, err := util.GetTokenFromAuthorization(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": "没有携带token信息",
		})
		return
	}
	claim, err := util.GetClaimFromToken(tokenString)
	if err != nil {
		context.Abort()
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": "token无效",
		})
		return
	}
	tokenTimeStamp := claim.(jwt.MapClaims)["timeStamp"].(float64)
	time := util.GetTimeStamp() - int64(tokenTimeStamp)
	var isExpired bool
	if time > consts.EXPIRE_TIME_TOKEN {
		isExpired = true
	}
	context.JSON(http.StatusOK, gin.H{
		"is_expired": isExpired,
	})
}
func SendVerifyCodeHandler(context *gin.Context) {
	var email string
	param := http_param.SendVerifyCodeArgument{}
	svc := business.NewService(context)
	err := context.ShouldBind(&param)
	if err != nil {
		logit.LogWithError("参数不全", map[string]interface{}{
			"Error": err,
		})
		context.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	isResetPassword, _ := strconv.ParseBool(param.IsResetPassword)
	email_aes, err := util.GetEmailFromAuthorization(context)
	if err != nil {
		emailAddr := strings.ToLower(param.Email)
		email = util.CheckHustEmailSuffixOrAddIt(emailAddr)
	} else {
		email, err = util.DecryptWithAes(email_aes)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "邮箱有误",
			})
			return
		}
	}
	if u, _ := svc.GetUserByEmailAes(email_aes); !isResetPassword && u.IsEmailActivated {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": "该邮箱已验证，请勿重复验证。",
		})
		return
	}
	verifyCode := util.GenerateVerifyCode(consts.VERIFYCODE_LENGTH)
	redisclt.StoreEmailAndVerifyCodeInRedis(verifyCode, email)
	go func() {
		var emailType int = 0
		if isResetPassword {
			emailType = 1
		}
		subject, content := emailclt.ParseEmailTemplate(emailType, verifyCode)
		// TODO 交给HY和LYX
		emailclt.RequestSendEmail(email, subject, content)
	}()
	context.JSON(http.StatusOK, gin.H{
		"msg": "成功发送验证码，请注意查收。如果没有收到，请稍后重试",
	})
}
func Login(context *gin.Context) {
	param := http_param.LoginArgument{}
	err := context.ShouldBind(&param)
	if err != nil {
		logit.LogWithError("参数不全", map[string]interface{}{
			"Error": err,
		})
		context.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	svc := business.NewService(context)
	e, token := svc.Login(param)
	if !e.CheckValid() {
		fmt.Println(e.Err)
	}
	if token == "" {
		context.JSON(e.Status, gin.H{"msg": e.Description})
		return
	}
	context.JSON(e.Status, gin.H{
		"msg":   e.Description,
		"token": token,
	})
}
