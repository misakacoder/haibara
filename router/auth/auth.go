package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"haibara/errs"
	"haibara/response"
	"haibara/ro"
	"time"
)

const (
	salt       = "Haibara Ai"
	expireTime = 7200
)

type Token struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}

type TokenRO struct {
	Token string `form:"token" binding:"required" msg:"token不能为空"`
}

type JWTClaims struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	jwt.StandardClaims
}

func LoadRouter(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.GET("token", getToken)
		auth.GET("check_token", checkToken)
		auth.GET("refresh_token", refreshToken)
	}
}

func ParseToken(context *gin.Context) (*JWTClaims, error) {
	tokenRO := ro.ValidateRO(context, &TokenRO{})
	claims, err := ParseTokenString(tokenRO.Token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func ParseTokenString(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(salt), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errs.ParseTokenError
	}
	if err = token.Claims.Valid(); err != nil {
		return nil, err
	}
	_, err = ro.ValidateUserRO(&ro.UserRO{Username: claims.Username, Password: claims.Password})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func getToken(context *gin.Context) {
	userRO := ro.ValidateRO(context, &ro.UserRO{})
	_, err := ro.ValidateUserRO(userRO)
	if err != nil {
		response.ErrorWithMessageJSON(context, err.Error())
		return
	}
	claims := &JWTClaims{
		Username: userRO.Username,
		Password: userRO.Password,
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Duration(expireTime) * time.Second).Unix()
	token, err := generateToken(claims)
	if err != nil {
		response.ErrorWithMessageJSON(context, err.Error())
		return
	}
	response.OkWithDataJSON(context, Token{Token: token, ExpiresIn: expireTime})
}

func checkToken(context *gin.Context) {
	claims, err := ParseToken(context)
	if err != nil {
		response.ErrorWithMessageJSON(context, err.Error())
		return
	}
	claims.Password = ""
	response.OkWithDataJSON(context, claims)
}

func refreshToken(context *gin.Context) {
	claims, err := ParseToken(context)
	if err != nil {
		response.ErrorWithMessageJSON(context, err.Error())
		return
	}
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	token, err := generateToken(claims)
	if err != nil {
		response.ErrorWithMessageJSON(context, err.Error())
		return
	}
	response.OkWithDataJSON(context, Token{Token: token, ExpiresIn: expireTime})
}

func generateToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(salt))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
