package redisclt

import (
	"fmt"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/consts"
	"time"
)

func StoreEmailAndVerifyCodeInRedis(verifyCode string, email string) {
	Set(consts.REDIS_VERIFY_CODE_SUFFIX+email, verifyCode, consts.VERIFYCODE_VALID_TIME*time.Second)
}

// IsVerifyCodeMatchToRegisterAccount Check if verify code and user specified by email inputs are match
func IsVerifyCodeMatchToRegisterAccount(verifyCode string, email string) (IsMatch bool) {
	re := Get(consts.REDIS_VERIFY_CODE_SUFFIX + email)
	fmt.Println("\n\n", re.Val())
	if re.Val() == verifyCode && re.Val() != "" {
		IsMatch = true
	} else {
		IsMatch = false
	}
	fmt.Println("IsMatch: ", IsMatch)
	return
}
func DeleteVerifyFromRedis(email string) {
	Del(consts.REDIS_VERIFY_CODE_SUFFIX + email)
}
