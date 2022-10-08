package business

import (
	"github.com/OverClockTeam/OC-WordsInSongs-BE/error_code"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/model"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/services/http_param"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/services/logit"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/util"
	"strings"
)

func (svc Service) GetUserByEmailAes(emailAes string) (user model.AuthUser, err error) {
	user, err = svc.dao.GetUser(emailAes)
	return
}

func (svc Service) Register(user *model.AuthUser) (errorCode error_code.ErrorCode) {
	err := svc.dao.Register(user)
	if err != nil {
		errorCode.Err = err
		errorCode.Status = 400
		errorCode.Description = "Register failed"
	} else {
		errorCode.Status = 200
	}
	return
}

func (svc Service) Login(param http_param.LoginArgument) (errorCode error_code.ErrorCode, token string) {
	email := util.CheckHustEmailSuffixOrAddIt(strings.ToLower(param.Email))

	cipherEmail, err := util.EncryptWithAes(email)
	if err != nil {
		errorCode = error_code.ErrorCode{Err: err, Description: err.Error(), Status: 400}
		return
	}
	// TODO 虽然前端加密是无效行为 但是还是应该有BASE64加密后传入后端的行为 这里预留好了decode部分为了方便测试所以暂时注掉 实际上还是应该有解密过程的
	//password, err := util.RsaDecryption(param.Password)
	password := param.Password
	if err != nil {
		errorCode = error_code.ErrorCode{Err: err, Description: err.Error(), Status: 400}
		return
	}
	err, role := svc.dao.GetRoleWhileLogin(cipherEmail)
	if err != nil {
		errorCode = error_code.ErrorCode{Err: err, Description: err.Error(), Status: 400}
		return
	}
	logit.LogWithInfo("正在登陆", map[string]interface{}{
		"password":    password,
		"email":       email,
		"cipherEmail": cipherEmail,
	})
	e, isRegistered := svc.dao.IsEmailRegistered(cipherEmail)
	if !isRegistered {
		errorCode = error_code.ErrorCode{Err: e, Description: "该账户未注册，请首先注册", Status: 400}
		return
	}

	if isMatch := svc.dao.IsEmailAndPasswordMatch(cipherEmail, password); !isMatch {
		errorCode = error_code.ErrorCode{Err: e, Description: "用户名密码不匹配", Status: 400}
		return
	}
	token, _ = util.GenerateTokenByJwt(cipherEmail, role)
	errorCode = error_code.ErrorCode{Err: nil, Description: "登录成功", Status: 200}
	return
}
