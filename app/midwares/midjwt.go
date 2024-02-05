package midwares

import (
	"SelectionSystem-Back/config/config"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)



type Claims struct{
	UserID int
	jwt.StandardClaims
}

func ParseToken(tokenStr string)(*Claims,error){
	secret :=config.Config.GetString("jwt.pass")
	var Secret = []byte(secret)
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
        tokenStr = tokenStr[7:]
    }
	token,err:=jwt.ParseWithClaims(tokenStr,&Claims{},func (token *jwt.Token)(interface{},error){
		return Secret,nil
	} )
	if err !=nil{
		return nil,err
	}
	//检验token
	if claims,ok:=token.Claims.(*Claims);ok&&token.Valid{
		return claims,nil
	}
	return nil,errors.New("invalid token")
}
func JWTAuthMiddleware()func(c *gin.Context){
	return func(c *gin.Context) {
		tokenStr:=c.Request.Header.Get("Authorization")
		if tokenStr ==""{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
				"code":200,
				"msg":"auth is null",

			})
			c.Abort()
			return
		}
		
		parts:=strings.Split(tokenStr,".")
		if len(parts)!=3{
			c.JSON(http.StatusUnauthorized,gin.H{
				"code":200,
				"msg":"auth is error",
			})
			c.Abort()
			return
		}
		mc,err:=ParseToken(tokenStr)
		if err !=nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"code":200,
				"msg":"Token is not vaild",
			})
			c.Abort()
			return
		}else if time.Now().Unix()>mc.ExpiresAt{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
				"code":200,
				"msg": "Token is overdue",
			})
			c.Abort()
			return
		}
		c.Set("UserID",mc.UserID)
		c.Next()
	}
}