package handler

import (
	"net/http"
	api_model "spider-go/api/model"
	"spider-go/asset"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"spider-go/usecase"
	"spider-go/utils/validator"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS_CODE    = "00"
	SUCCESS_MESSAGE = "success"
)

type CreateAccoutHandler struct {
	authUseCase domain.Authorities
	log         *logger.Logger
}

func NewCreateAccountHandler(authUsecase domain.Authorities) *CreateAccoutHandler {
	return &CreateAccoutHandler{
		authUseCase: authUsecase,
		log:         logger.L().Named("CreateAccoutHandler"),
	}
}

func (h *CreateAccoutHandler) CreateAccout(ctx *gin.Context) {

	log := h.log.WithContext(ctx)

	var req api_model.CreateAccoutReq
	var resp api_model.CreateAccoutResp

	if err := ctx.ShouldBind(&req); err != nil {
		log.Errorf("should bind request failed: %+v", err)
		resp.Header.ErrorCode = asset.E().GeneralSystemError.ErrorCode
		resp.Header.Message = asset.E().GeneralSystemError.ErrorMessageEN
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	if err := validator.Struct(req); err != nil {
		log.Errorf("validate request data fail, error: %+v", err)
		resp.Header.ErrorCode = asset.E().RequestDataFail.ErrorCode
		resp.Header.Message = asset.E().RequestDataFail.ErrorMessageEN
		ctx.JSON(asset.E().RequestDataFail.StatusCode, resp)
		return
	}

	preAccount := h.prepareAccount(req)

	if err := h.authUseCase.CreateAccout(ctx, preAccount, req.Data.Password, req.Data.ConfirmPassword); err != nil {
		log.Errorf("usecase request failed: %+v", err)
		assetError := h.mapErrorCreateAccount(err)
		resp.Header.ErrorCode = assetError.ErrorCode
		resp.Header.Message = assetError.ErrorMessageEN
		ctx.JSON(assetError.StatusCode, resp)
		return
	}

	resp.Header.ErrorCode = SUCCESS_CODE
	resp.Header.Message = SUCCESS_MESSAGE

	ctx.JSON(http.StatusOK, resp)

}

func (h *CreateAccoutHandler) prepareAccount(reqAccount api_model.CreateAccoutReq) (account model.Account) {
	account = model.Account{
		Username:  reqAccount.Data.Username,
		Title:     reqAccount.Data.Title,
		FirstName: reqAccount.Data.FirstName,
		LastName:  reqAccount.Data.LastName,
		Age:       reqAccount.Data.Age,
		MobileNO:  reqAccount.Data.MobileNO,
		Role:      reqAccount.Data.Role,
	}
	return account
}

func (h *CreateAccoutHandler) mapErrorCreateAccount(err error) *asset.ErrorCode {
	switch err {
	case usecase.ErrorAuthoritiesConfirmPasswordNotMatch:
		return &asset.E().PasswordMatchingError
	case usecase.ErrorAuthoritiesInsertAccountToMongoFail:
		return &asset.E().ErrorSpiderDB
	case usecase.ErrorAuthoritiesHashPasswordFail:
		return &asset.E().HashingError
	default:
		return &asset.E().GeneralSystemError
	}
}
