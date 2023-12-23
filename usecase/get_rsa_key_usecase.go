package usecase

import (
	"context"
	"fmt"
	"spider-go/config"
	"spider-go/domain"
	"spider-go/logger"
	"spider-go/model"
	"spider-go/utils/cryptography"
	"spider-go/utils/random"
)

type GetRsaKey struct {
	redisRepo domain.RedisRepository
	config    *config.Root
	log       *logger.Logger
}

var (
	ErrorGenerateRSAKey  = fmt.Errorf("[Get RSA Key Usecase]: generate rsa key error")
	ErrorSaveDataToRedis = fmt.Errorf("[Get RSA Key Usecase]: save data to redis failed")
)

func NewGetRsaKey(redisRepo domain.RedisRepository, conf *config.Root) domain.GetRsaKeyUsecase {
	return &GetRsaKey{
		redisRepo: redisRepo,
		config:    conf,
		log:       logger.L().Named("GetRsaKey"),
	}
}

func (u *GetRsaKey) GenerateRsaKey(ctx context.Context) (publicKey, searchKey string, err error) {

	log := u.log.WithContext(ctx)

	//  generate random string for make redis key
	r := random.NewRandom()

	searchKey = r.RandomString(u.config.RSAOption.RandomKeySize)

	redisKey := fmt.Sprintf(u.config.RedisOption.RSA.KeyFormat, searchKey)

	log.Infof("redis key: %v, ttl: %v", redisKey, u.config.RedisOption.RSA.TTL)

	// create rsa key
	cryto := cryptography.NewCrypto()

	privateKey, publicKey, err := cryto.GenerateRSAKey(u.config.RSAOption.RSASize)
	if err != nil {
		log.Errorf("generate RSA key error: %+v", err)
		return "", "", ErrorGenerateRSAKey
	}

	RsaKey := model.RedisRsaKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	// save to redis
	if err := u.redisRepo.SetDataToRedisWithTTL(ctx, redisKey, RsaKey, u.config.RedisOption.RSA.TTL); err != nil {
		log.Errorf("save data to redis error: %+v", err)
		return "", "", ErrorSaveDataToRedis
	}

	return publicKey, searchKey, nil

}
