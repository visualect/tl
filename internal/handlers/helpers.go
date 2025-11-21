package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.Validator.Struct(i); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			errorMsg := GetValidationError(validationErrors)
			return echo.NewHTTPError(http.StatusBadRequest, errorMsg)
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func GetValidationError(vErrors validator.ValidationErrors) (errMsg string) {
	for _, fe := range vErrors {
		fieldName := strings.ToLower(fe.Field())
		tag := fe.Tag()
		param := fe.Param()

		switch tag {
		case "required":
			errMsg += fmt.Sprintf("field \"%s\" is required\n", fieldName)
		case "max":
			errMsg += fmt.Sprintf("field \"%s\" should not contain more than %s symbols\n", fieldName, param)
		default:
			errMsg += fmt.Sprintf("validation error of rule %s in the \"%s\" field", tag, fieldName)
		}
	}

	return
}
