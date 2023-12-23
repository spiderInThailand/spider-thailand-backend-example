package validator

import "github.com/go-playground/validator/v10"

var vldt *validator.Validate

func V() *validator.Validate {
	return vldt
}

func init() {
	vldt = validator.New()
}

func Struct(data interface{}) error {
	return vldt.Struct(data)
}
