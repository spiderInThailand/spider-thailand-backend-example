package domain

import "context"

//go:generate mockgen -source=get_rsa_key_domain.go -destination=./mock/get_rsa_key_domain.go

type GetRsaKeyUsecase interface {
	GenerateRsaKey(ctx context.Context) (publicKey, searchKey string, err error)
}
