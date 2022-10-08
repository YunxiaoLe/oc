package cache

import (
	"github.com/OverClockTeam/OC-WordsInSongs-BE/consts"
	"time"
)

func AddUserToBlackList(email_aes string, mute_time float64) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(loc).Unix()
	Set(consts.BLACK_LIST_PREFIX+email_aes, t, time.Duration(mute_time)*time.Hour)
}
func RemoveFromBlackList(email_aes string) {
	Del(consts.BLACK_LIST_PREFIX + email_aes)
}
func IsInvolvedInBlackList(email_aes string) (isInvolved bool) {
	re, _ := Exists(consts.BLACK_LIST_PREFIX + email_aes).Result()
	if re == 1 {
		isInvolved = true
	} else {
		isInvolved = false
	}
	return
}
