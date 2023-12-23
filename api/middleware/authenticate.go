package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"spider-go/asset"
	"spider-go/domain"
	"spider-go/logger"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authenticate(jwtService domain.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.L().Named("Authenticate").WithContext(ctx)

		var req MiddlewareRequest
		var resp MiddlewareResponse

		// get data in body
		bodyReq, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Errorf("[authenticate] read all request body failed, error: %+v", err)
			resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
			resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN

			ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
			return
		}

		defer ctx.Request.Body.Close()

		// recover data to body
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyReq))

		// read body request
		if err := json.Unmarshal(bodyReq, &req); err != nil {
			log.Errorf("[authenticate] unmarshal request body failed, error: %+v", err)
			resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
			resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN

			ctx.AbortWithStatusJSON(asset.E().GeneralSystemError.StatusCode, resp)
			return
		}

		// =========================================================
		// validate user login
		// =========================================================

		token, err := jwtService.ValidateToken(fmt.Sprint(req.Header["token"]))
		if err != nil {
			log.Errorf("[authenticate] find user in redis failed, error: %+v", err)
			resp.Header.ErrorCode = asset.E().UserNotLogin.ErrorCode
			resp.Header.Message = asset.E().UserNotLogin.ErrorMessageEN
			ctx.AbortWithStatusJSON(asset.E().UserNotLogin.StatusCode, resp)
		}
		// **********************************************************

		// var userInfo interface{}

		// json.Unmarshal(redisResult, &userInfo)
		// log.Debugf("[authenticate] user login info: %+v", userInfo)
		// ctx.Set("user_info", userInfo)
		log.Infof("[authenticate] token %+v", token.Claims.(jwt.MapClaims))

		// =========================================================
		//  update data expire
		// =========================================================

		// jwtService.RefreshToken(token)

		// **********************************************************

		ctx.Next()

	}
}
