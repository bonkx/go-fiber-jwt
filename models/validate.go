package models

import (
	"fmt"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	en_translations "github.com/go-playground/validator/v10/translations/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

// use a single instance , it caches struct info
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    *ut.Translator
)

func ValidateStruct(data interface{}) []*ErrorResponse {

	validate = validator.New()

	en := en.New()
	id := id.New()
	uni = ut.New(en, en, id)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")
	// trans, _ := uni.GetTranslator("id")

	en_translations.RegisterDefaultTranslations(validate, trans)
	id_translations.RegisterDefaultTranslations(validate, trans)

	var errors []*ErrorResponse
	err := validate.Struct(data)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			translatedErr := fmt.Errorf(err.Translate(trans))
			element.FailedField = err.StructNamespace()
			// element.FailedField = strings.ToLower(err.StructField())
			// log.Println(err.StructNamespace())
			element.Tag = err.Tag()
			// element.Value = err.Param()
			element.Value = translatedErr.Error()
			if err.Tag() == "e164" {
				element.Value = fmt.Sprintf("%s %s", translatedErr.Error(), "i.e. +6281234567890 or +628 123 4567 890")
			}
			errors = append(errors, &element)
		}
	}

	return errors
}

func ValidatePhoneNumber(phone_number string) []*ErrorResponse {
	type Phone struct {
		Phone string `validate:"e164"`
	}
	validate = validator.New()

	en := en.New()
	id := id.New()
	uni = ut.New(en, en, id)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")
	// trans, _ := uni.GetTranslator("id")

	en_translations.RegisterDefaultTranslations(validate, trans)
	id_translations.RegisterDefaultTranslations(validate, trans)

	var errors []*ErrorResponse
	// validate.RegisterValidation("e164", ValidateMyVal)

	s := Phone{Phone: phone_number}
	err := validate.Struct(s)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			translatedErr := fmt.Errorf(err.Translate(trans))
			element.FailedField = err.StructNamespace()
			// element.FailedField = strings.ToLower(err.StructField())
			// log.Println(err.StructNamespace())
			element.Tag = err.Tag()
			// element.Value = err.Param()
			element.Value = translatedErr.Error()
			if err.Tag() == "e164" {
				element.Value = fmt.Sprintf("%s %s", translatedErr.Error(), "i.e. +6281234567890 or +628 123 4567 890")
			}
			errors = append(errors, &element)
		}
	}

	return errors
}

func ValidateMyVal(fl validator.FieldLevel) bool {
	return fl.Field().String() == "e164"
}
