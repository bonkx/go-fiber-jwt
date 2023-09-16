package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

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

func getJSONField(data interface{}, snp string, fieldname string) error {
	if !strings.Contains(snp, fieldname) {
		return nil
	}

	fieldArr := strings.Split(snp, ".")
	// fmt.Println(fieldArr)

	val := reflect.ValueOf(data)

	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)

		fieldName := t.Name
		// fmt.Println(fieldArr[1])
		if fieldArr[1] == fieldName {
			if jsonTag := t.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				// check for possible comma as in "...,omitempty"
				var commaIdx int
				if commaIdx = strings.Index(jsonTag, ","); commaIdx < 0 {
					commaIdx = len(jsonTag)
				}
				fieldName = jsonTag[:commaIdx]
			}
			return errors.New(fieldName)
		}

		return nil
	}

	return nil
}

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
			element.Field = err.StructNamespace()
			// log.Println(err.StructNamespace())
			element.Tag = err.Tag()
			// element.Value = err.Param()
			element.Message = translatedErr.Error()

			errJson := getJSONField(data, err.StructNamespace(), err.Field())
			if errJson != nil {
				element.Field = errJson.Error()
			}

			if err.Tag() == "e164" {
				element.Message = fmt.Sprintf("%s %s", translatedErr.Error(), "i.e. +6281234567890 or +628 123 4567 890")
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
			element.Field = err.StructNamespace()
			// log.Println(err.StructNamespace())
			element.Tag = err.Tag()
			// element.Value = err.Param()
			element.Message = translatedErr.Error()
			if err.Tag() == "e164" {
				element.Message = fmt.Sprintf("%s %s", translatedErr.Error(), "i.e. +6281234567890 or +628 123 4567 890")
			}
			errors = append(errors, &element)
		}
	}

	return errors
}

func ValidateMyVal(fl validator.FieldLevel) bool {
	return fl.Field().String() == "e164"
}
