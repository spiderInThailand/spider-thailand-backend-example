package domain

import (
	"context"
	"spider-go/model"
)

//go:generate mockgen -source=auth_domain.go -destination=./mock/auth_domain.go

type Authorities interface {
	CreateAccout(ctx context.Context, data model.Account, password, confirmPassowrd string) (err error)
	Login(ctx context.Context, username, password string) (accountInfo *model.Account, key string, err error)
}
