package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/consts"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CheckError returns false if err is not nil and output a info as log
func CheckError(err error) bool {
	if err != nil {
		logrus.Error(err)
		return true
	}
	return false
}

// CheckHustEmailSuffixOrAddIt 为邮箱添加hust.edu.cn的后缀如果已有则不变，同时起检查邮箱的作用
func CheckHustEmailSuffixOrAddIt(EmailOrUid string) (email string) {
	emailSuffix := "@hust.edu.cn"
	if !strings.Contains(EmailOrUid, emailSuffix) {
		email = EmailOrUid + emailSuffix
	} else {
		email = EmailOrUid
	}
	return email
}

// RsaDecryption chiper is a string with base64 encoding,and return a plaintext
func RsaDecryption(chiper string) (plaintext string, err error) {
	des, err := base64.StdEncoding.DecodeString(chiper)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	privateKeyStr, err := ioutil.ReadFile("private.pem")
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	priblock, _ := pem.Decode(privateKeyStr)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	priKey, err := x509.ParsePKCS1PrivateKey(priblock.Bytes)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	plain, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, des)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	if err != nil {
		logrus.Error(err)
		log.Panic(err)
	}
	plaintext = string(plain)
	return plaintext, nil
}

// HashWithSalt A Hash function using salt with bcrypt libriry to hash password
// 将纯字符串hash
func HashWithSalt(plainText string) (HashText string) {

	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.MinCost)
	CheckError(err)
	HashText = string(hash)
	return
}

/*
AES加密函数，将明文stringToEncrypt加密成为密文encryptedString
AES的密钥从配置文件中读入
*/
func EncryptWithAes(stringToEncrypt string) (encryptedString string, err error) {
	key, err := hex.DecodeString(aesKey)
	if err != nil {

		return "", err
	}
	plaintext := []byte(stringToEncrypt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return fmt.Sprintf("%x", ciphertext), nil
}

func GetTimeStamp() (t int64) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t = time.Now().In(loc).Unix()
	return
}

// GenerateTokenByJwt 使用Jwt生成token作身份认证使用
// 其中存储了登录的时间，登陆的email以及用户的角色信息
func GenerateTokenByJwt(email string, role string) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"email":     email,
		"timeStamp": GetTimeStamp(),
		"role":      role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(consts.TOKEN_SCRECT_KEY))
	return
}
func GetTokenFromAuthorization(c *gin.Context) (token string, err error) {
	defer func() {
		if e := recover(); e != any(nil) {
			err = errors.New("您未登录，请登陆后查看")
		}
	}()
	authHeader := c.Request.Header.Get("Authorization")
	token = strings.Fields(authHeader)[1]
	if err != nil {
		err = errors.New("您未登录，请登陆后查看")
		return
	}
	return
}
func GetClaimFromToken(tokenString string) (claims jwt.Claims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return []byte(consts.TOKEN_SCRECT_KEY), err
	})
	if err != nil {
		return nil, err
	} else {
		claims = token.Claims.(jwt.MapClaims)
		return claims, nil
	}
}

// GetEmailFromAuthorization Get encrypted user email
// err is not nil when authorization token is not valid
func GetEmailFromAuthorization(c *gin.Context) (email string, err error) {
	defer func() {
		if e := recover(); e != any(nil) {
			err = errors.New("您未登录，请登陆后查看")
		}
	}()
	authHeader := c.Request.Header.Get("Authorization")
	token := strings.Fields(authHeader)[1]
	if err != nil {
		err = errors.New("您未登录，请登陆后查看")
		return
	}
	claim, _ := GetClaimFromToken(token)
	email = claim.(jwt.MapClaims)["email"].(string)
	return
}

func DecryptWithAes(CipherText string) (decryptedString string, err error) {
	key, err := hex.DecodeString(aesKey)
	if err != nil {
		return "", err
	}
	ciphertext, _ := hex.DecodeString(CipherText)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return fmt.Sprintf("%s", ciphertext), nil
}

// GenerateVerifyCode 随机生成生成length长度的激活码
func GenerateVerifyCode(length int) (verifyString string) {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(any(err))
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	verifyString = string(b)
	return
}
func AddToIptablesBlacklist(ip string) error {
	err := os.Setenv("black", ip)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	} else {
		fmt.Println("set " + ip + "into blacklist OK!")
	}
	command := `./black.sh .`
	cmd := exec.Command("/bin/bash", "-c", command)

	output, err1 := cmd.Output()
	if err1 != nil {
		fmt.Printf("Execute Shell:%s failed with error:%s", command, err1.Error())
		return err1
	}
	fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(output))
	return nil
}
