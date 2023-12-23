package handler

import (
	"fmt"
	"net/http"
	api_model "spider-go/api/model"
	"spider-go/asset"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"spider-go/usecase"
	"spider-go/utils/uuid"
	"spider-go/utils/validator"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	AuthoritiesUsecase domain.Authorities
	log                *logger.Logger
}

func NewLoginHandler(authoritailUsecase domain.Authorities) *LoginHandler {
	return &LoginHandler{
		AuthoritiesUsecase: authoritailUsecase,
		log:                logger.L().Named("LoginHandler"),
	}
}

func (h *LoginHandler) Login(ctx *gin.Context) {
	log := h.log.WithContext(ctx)

	var req api_model.LoginRequest
	var resp api_model.LoginResponse

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	log.Infof("Login handler start with req: %+v", req)

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	accountInfo, token, err := h.AuthoritiesUsecase.Login(ctx, req.Data.Username, req.Data.Password)
	if err != nil {
		log.Errorf("usecase request failed: %+v", err)
		assetError := h.mapErrorLogin(err)
		resp.Header.ErrorCode = assetError.ErrorCode
		resp.Header.Message = assetError.ErrorMessageEN
		ctx.JSON(assetError.StatusCode, resp)
		return
	}

	resp.Data = h.prepareAccountInfoResponse(*accountInfo, token)

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = SUCCESS_MESSAGE

	log.Infof("response: %+v", resp)

	ctx.JSON(http.StatusOK, resp)
}

func (h *LoginHandler) prepareAccountInfoResponse(accountInfo model.Account, token string) (respData api_model.LoginInfo) {
	respData = api_model.LoginInfo{
		User: api_model.User{
			ID:        uuid.GernerateUUID32(),
			Username:  accountInfo.Username,
			Name:      fmt.Sprintf("%s.%s %s", accountInfo.Title, accountInfo.FirstName, accountInfo.LastName),
			Role:      accountInfo.Role,
			CreatedAt: time.Now(),
		},
		BackendToken: api_model.BackendToken{
			Token: token,
		},
	}

	return respData
}

func (h *LoginHandler) mapErrorLogin(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorAuthoritiesInvalidPassword, usecase.ErrorAuthoritiesFindAccountNotFound:
		return &asset.E().InvalidLoginAccount
	case usecase.ErrorAuthoritiesMongoConnection:
		return &asset.E().ErrorSpiderDB
	default:
		return &asset.E().GeneralSystemError
	}
}

func (h *LoginHandler) VerifyLogin(ctx *gin.Context) {

	resp := api_model.VerifyLoginResponser{
		Header: api_model.ResponseHeader{
			ErrorCode: "00",
			Message:   "user is login",
		},
	}

	ctx.JSON(http.StatusOK, resp)
}
