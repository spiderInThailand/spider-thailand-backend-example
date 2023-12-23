package domain

import (
	"context"
	"spider-go/model"
)

//go:generate mockgen -source=account_domain.go -destination=./mock/account_domain.go
type AccountRepository interface {
	CreateAccout(ctx context.Context, acc model.Account) (err error)
	FindAccountByUsername(ctx context.Context, username string) (AccountInfo *model.Account, err error)
}
