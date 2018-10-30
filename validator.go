package main

import(
	"gopkg.in/go-playground/validator.v9"		
)

var validate *validator.Validate

var errormessage string

func StartValidator() {
	
	validate = validator.New()

	validate.RegisterValidation("iluminati-ursal", ValidateIlur)

}

func ValidateIlur(gender validator.FieldLevel) bool {

	g := gender.Field().String()

	if (g == "iluminati" || g == "ursal") {
		return false
	}

	return true
}