package validations

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func DateValidation(fl validator.FieldLevel) bool {

	var match bool = false
	date := fl.Field().String()

	dateRegex := `^(19|20)\d{2}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])$`
	var err error
	match, err = regexp.MatchString(dateRegex, date)
	if err != nil {
		return false
	}
	return match

}
