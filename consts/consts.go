package consts

const (
	//Set-Cookie name prosperity
	COOKIE_NAME = "token"

	//Token is invalid in one month
	EXPIRE_TIME_TOKEN = 30 * 24 * 60 * 60

	TOKEN_SCRECT_KEY = "tokenemail"

	//Verify code is vaild in 24h
	VERIFYCODE_VALID_TIME = 24 * 60 * 60

	VERIFYCODE_LENGTH = 4

	RATA_NUMBER_LIMITER = 2

	//Redis黑名单prefix string
	BLACK_LIST_PREFIX = "blacklist:"

	CONFIG_FILE_NAME = "Config.json"

	USER  = "user"
	ADMIN = "admin"
	ROOT  = "root"

	REDIS_VERIFY_CODE_SUFFIX = "verify_code:"
	//默认的禁言时间，单位为小时
	Default_MUTE_TIME = 0.5

	TITLE_MUTE       = "用户封禁"
	PREFIX_INCOGNITO = "IN"

	DELETED_REPLY_CONTENT = "该评论已被删除"

	EXTREMLY_LOW_RATE  = 5
	LOW_RATE           = 30 //每分钟运行访问30次
	MIDDLE_RATE        = 60 //每分钟运行访问
	HIGH_RATE          = 120
	MAX_BASE_ALIAS_NUM = 500000
)
