package uuid

import (
	"github.com/google/uuid"
)

func GernerateUUID32() (key string) {
	uuid_key := uuid.New()

	return uuid_key.String()
}
