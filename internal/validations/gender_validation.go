package validations

import (
	"github.com/go-playground/validator/v10"
	"github.com/juseph-q/SchoolPr/internal/database/models"
)

func ValidateGender(gender validator.FieldLevel) bool {

	if gender != nil {
		gender := gender.Field().String()

		switch gender {
		case string(models.GenderFemale), string(models.GenderMale):
			return true
		default:
			return false
		}
	}
	return true

}
