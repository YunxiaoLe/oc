package middleware

import (
	"github.com/OverClockTeam/OC-WordsInSongs-BE/cache"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/consts"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/services/logit"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

/*
该中间件需要登陆，用以限制登录用户对相应api的请求次数
`api_name` 表示在redis中记录被限制api存储的部分key，与实际限制使用的api无关
`rateLimitPerMinute`表示每分钟可以使用多少次该api
*/

func RateLimiter(apiName string, rateLimitPerMinute int) gin.HandlerFunc {
	return func(context *gin.Context) {
		emailAes, _ := util.GetEmailFromAuthorization(context)
		key := apiName + ":" + emailAes
		isInBlackList := cache.IsInvolvedInBlackList(emailAes)

		//To recognize which api in blacklist
		blacklist_key := consts.BLACK_LIST_PREFIX + apiName + ":" + emailAes

		if isInBlackList {
			re, _ := cache.RedisClient.Get(blacklist_key).Int64()
			time_left := strconv.Itoa(int(30 - (util.GetTimeStamp()-re)/60))
			context.Abort()
			logit.LogWithWarn("用户请求过于频繁被限制", map[string]interface{}{
				"Email":     emailAes,
				"Time left": time_left,
			})
			context.JSON(http.StatusBadRequest, gin.H{
				"msg": "请求过于频繁，请" + time_left + "分钟后再发言:)",
			})
			return
		} else if re, _ := cache.RedisClient.Get(key).Int(); re >= rateLimitPerMinute {
			context.Abort()
			cache.AddUserToBlackList(key, consts.Default_MUTE_TIME)
			logit.LogWithWarn("用户请求过于频繁,现在已被禁止", map[string]interface{}{
				"Email": emailAes,
			})
			context.JSON(http.StatusBadRequest, gin.H{
				"msg": "过于频繁请求!",
			})
			return
		} else {
			if re, _ := cache.RedisClient.Exists(key).Result(); re == 1 {
				cache.RedisClient.Incr(key)
			} else {
				cache.RedisClient.Set(key, 1, time.Minute)
			}
			context.Next()
		}
	}
}

/*
该中间件针对未登录用户对api的访问，通过请求的ip地址进行访问
*/
func IpRateLimiter(apiName string, rateLimitPerMinute int) gin.HandlerFunc {
	return func(context *gin.Context) {
		//context.Next()
		ipAddress := context.ClientIP()
		key := apiName + ":" + ipAddress
		isInBlackList := cache.IsInvolvedInBlackList(ipAddress)
		//
		////To recognize which api in blacklist
		blacklist_key := consts.BLACK_LIST_PREFIX + apiName + ":" + ipAddress
		if isInBlackList {
			re, _ := cache.RedisClient.Get(blacklist_key).Int64()
			time_left := strconv.Itoa(int(30 - (util.GetTimeStamp()-re)/60))
			context.Abort()
			logit.LogWithWarn("ip地址请求过于频繁被限制", map[string]interface{}{
				"Ip Address": ipAddress,
				"Time left":  time_left,
			})
			context.JSON(http.StatusTooManyRequests, gin.H{
				"msg": "请求过于频繁，请" + time_left + "分钟后再来:)",
			})
			return
		} else if re, _ := cache.RedisClient.Get(key).Int(); re >= rateLimitPerMinute {
			context.Abort()
			_ = util.AddToIptablesBlacklist(ipAddress)
			cache.AddUserToBlackList(key, consts.Default_MUTE_TIME)
			logit.LogWithWarn("ip地址请求过于频繁,现在已被禁止", map[string]interface{}{
				"Ip Address": ipAddress,
			})
			context.JSON(http.StatusTooManyRequests, gin.H{
				"msg": "过于频繁请求!",
			})
			return
		} else {
			if re, _ := cache.RedisClient.Exists(key).Result(); re == 1 {
				cache.RedisClient.Incr(key)
			} else {
				cache.RedisClient.Set(key, 1, time.Minute)
			}
			context.Next()
		}
	}
}
