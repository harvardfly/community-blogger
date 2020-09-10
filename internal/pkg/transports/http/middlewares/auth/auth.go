package auth

/*
token认证中间件
*/

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"community-blogger/internal/pkg/baseerror"
	"community-blogger/internal/pkg/baseresponse"
	"community-blogger/internal/pkg/utils/middlewareutil"
)

var (
	// DefaultField 默认字段
	DefaultField = "Authorization"
	// AccessTokenValidErr Token验证失败
	AccessTokenValidErr = baseerror.NewBaseError("AccessToken 验证失败")
	// AccessTokenValidationErrorExpiredErr  AccessToken过期
	AccessTokenValidationErrorExpiredErr = baseerror.NewBaseError("AccessToken过期")
	// AccessTokenValidationErrorMalformedErr  AccessToken格式错误
	AccessTokenValidationErrorMalformedErr = baseerror.NewBaseError("AccessToken格式错误")
)

// Next 允许跨域
func Next() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Authorization,User-Agent, Keep-Alive, Content-Type, X-Requested-With,X-CSRF-Token,AccessToken,Token")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusAccepted)
		}
		c.Next()
	}
}

// ValidAccessToken 验证JWT TOKEN
func ValidAccessToken(context *gin.Context) {
	authorization := context.GetHeader(DefaultField)
	user, err := middlewareutil.ParseToken(authorization)
	if err != nil {
		if err, ok := err.(*jwt.ValidationError); ok {
			if err.Errors&jwt.ValidationErrorMalformed != 0 {
				baseresponse.HTTPResponse(context, nil, AccessTokenValidationErrorMalformedErr)
				context.Abort()
				return
			}
			if err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				baseresponse.HTTPResponse(context, nil, AccessTokenValidationErrorExpiredErr)
				context.Abort()
				return
			}
		}
		baseresponse.HTTPResponse(context, nil, AccessTokenValidErr)
		context.Abort()
		return
	}
	if user != nil {
		context.Set("username", user.Username)
		context.Next()
		return
	}
	baseresponse.HTTPResponse(context, nil, AccessTokenValidErr)
	context.Abort()
	return

}
