package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	"github.com/vasapolrittideah/moneylog-api/shared/contract"
)

var (
	val   = validator.New()
	trans = registerTranslation()
)

func ValidateStruct(input any) []contract.APIValidationError {
	var errs []contract.APIValidationError
	if err := val.Struct(input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errs = translateErrorMessage(validationErrors)
		}
	}

	return errs
}

func translateErrorMessage(validationErrors validator.ValidationErrors) []contract.APIValidationError {
	var errs []contract.APIValidationError
	var invalidField contract.APIValidationError
	for _, err := range validationErrors {
		invalidField = contract.APIValidationError{
			Field:   err.Field(),
			Message: err.Translate(trans),
			Value:   err.Value(),
		}

		errs = append(errs, invalidField)
	}

	return errs
}

func registerTranslation() ut.Translator {
	english := en.New()
	universalTranslator := ut.New(english, english)
	trans, _ := universalTranslator.GetTranslator("en")
	_ = enTrans.RegisterDefaultTranslations(val, trans)

	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		const jsonTagParts = 2
		name := strings.SplitN(fld.Tag.Get("json"), ",", jsonTagParts)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return trans
}
