package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate(data any) error {
	err := validate.Struct(data)
	if err != nil {
		strErr := ""

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, vErr := range validationErrors {
				strErr += fmt.Sprintf("%s", vErr.Error())
			}
		} else {
			strErr = err.Error()
		}

		return errors.New(strErr)
	}

	return nil
}