package usecase

import (
	"context"
	"fmt"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"spider-go/utils/cryptography"
)

var (
	ErrorAuthoritiesConfirmPasswordNotMatch  = fmt.Errorf("[Authorities Usecase]: password and confirm password not match")
	ErrorAuthoritiesHashPasswordFail         = fmt.Errorf("[Authorities Usecase]: hashing password failed")
	ErrorAuthoritiesInsertAccountToMongoFail = fmt.Errorf("[Authorities Usecase]: insert account to mongo failed")
	ErrorAuthoritiesFindAccountNotFound      = fmt.Errorf("[Authorities Usecase]: find data in mongo not found")
	ErrorAuthoritiesMongoConnection          = fmt.Errorf("[Authorities Usecase]: mongo error")
	ErrorAuthoritiesInvalidPassword          = fmt.Errorf("[Authorities Usecase]: invalid password")
	ErrorAuthoritiesTempDataConnection       = fmt.Errorf("[Authorities Usecase]: redis error")
	ErrorAuthoritiesGenerateTokenFail        = fmt.Errorf("[Authorities Usecase]: generate jwt token failed")
)

type Authorities struct {
	accRepo    domain.AccountRepository
	JWTService domain.JWTService
	log        *logger.Logger
}

func NewAuthoritiesUsecase(accRepo domain.AccountRepository, JWTService domain.JWTService) domain.Authorities {
	return &Authorities{
		accRepo:    accRepo,
		JWTService: JWTService,
		log:        logger.L().Named("CreateAccountUsecase"),
	}
}

func (u *Authorities) CreateAccout(ctx context.Context, data model.Account, password, confirmPassowrd string) (err error) {
	log := u.log.WithContext(ctx)

	// ==========================================================
	// validate strong password and match password with confirm password
	// ==========================================================
	if password != confirmPassowrd {
		return ErrorAuthoritiesConfirmPasswordNotMatch
	}

	// ==========================================================
	// hashing password
	// ==========================================================
	cryptoFunc := cryptography.NewCrypto()

	hashPass, err := cryptoFunc.HashPassword(ctx, password)
	if err != nil {
		log.Errorf("hashing password error: %+v", err)
		return ErrorAuthoritiesHashPasswordFail
	}

	data.HashPassword = hashPass

	// ==========================================================
	// save data to mongo
	// ==========================================================
	if err := u.accRepo.CreateAccout(ctx, data); err != nil {
		log.Errorf("insert data to mongo error: %+v", err)
		return ErrorAuthoritiesInsertAccountToMongoFail
	}

	return nil
}

func (u *Authorities) Login(ctx context.Context, username, password string) (accountInfo *model.Account, token string, err error) {
	log := u.log.WithContext(ctx)

	accountInfo, err = u.accRepo.FindAccountByUsername(ctx, username)
	if err != nil {
		if err.Error() == MONGO_NOT_FOUND {
			log.Errorf("find account by username %v, error not found", username)
			return nil, "", ErrorAuthoritiesFindAccountNotFound
		}
		log.Errorf("find account by username %v, but error: %+v", username, err)
		return nil, "", ErrorAuthoritiesMongoConnection
	}

	crypto := cryptography.NewCrypto()

	if !crypto.CheckPasswordHash(password, accountInfo.HashPassword) {
		return nil, "", ErrorAuthoritiesInvalidPassword
	}

	token, err = u.JWTService.GenerateNewToken(accountInfo.Username, accountInfo.Role)
	if err != nil {
		return nil, "", ErrorAuthoritiesGenerateTokenFail
	}

	return accountInfo, token, nil
}
