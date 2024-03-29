package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"network-disk/tool"
	"strings"
	"time"
)

var jwtKey = []byte("niu_rou_ban_mian")

type MyClaims struct {
	Username string `json:"username"`
	Identity string `json:"identity"`
	jwt.StandardClaims
}

//生成token

func SetToken(username, identity string) (string, bool) {
	SetClaims := MyClaims{
		username,
		identity,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,
			ExpiresAt: time.Now().Unix() + 60*60*2,
			Issuer:    "Twentyue",
			Subject:   "my_network_disk",
		},
	}

	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	token, err := reqClaim.SignedString(jwtKey)
	if err != nil {
		return "", false
	}
	return token, true
}

//jwt中间件

func JwtToken(c *gin.Context) {
	var code string
	tokenHeader := c.Request.Header.Get("Authorization")
	if tokenHeader == "" {
		code = "token 不存在"
		c.JSON(200, gin.H{
			"code": code,
			"msg":  "请登陆后再操作",
		})
		c.Abort()
		return
	}
	checkToken := strings.SplitN(tokenHeader, "", 2)
	if len(checkToken) != 2 && checkToken[0] != "Bearer" {
		code = "token格式错误"
		c.JSON(200, gin.H{
			"msg": code,
		})
		c.Abort()
		return
	}

	//解析token
	token, err := jwt.ParseWithClaims(tokenHeader, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	//获取token中的字段
	username := token.Claims.(*MyClaims).Username
	Time := token.Claims.(*MyClaims).ExpiresAt
	administrator := token.Claims.(*MyClaims).Identity
	if Time < time.Now().Unix() {
		code = "token已过期"
		c.JSON(200, gin.H{
			"msg": code,
		})
		c.Abort()
		return
	}

	if err != nil {
		fmt.Println("check token failed,err:", err)
		tool.RespInternetError(c)
		return
	}

	if token.Valid == false {
		code = "token不正确"
		c.JSON(200, gin.H{
			"msg": code,
		})
		c.Abort()
		return
	}
	if administrator != "管理员" {
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "非管理员，无权限操作",
		})
		return
	}

	c.Set("username", username)
}
