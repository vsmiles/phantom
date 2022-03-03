package api

import "github.com/go-playground/validator/v10"

var validLogin validator.Func = func(fl validator.FieldLevel) bool {
	if loginParams, ok := fl.Field().Interface().(string); ok {
		return loginParams == "login"
	}
	return false
}
