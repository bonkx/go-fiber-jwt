package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	en_translations "github.com/go-playground/validator/v10/translations/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

// use a single instance , it caches struct info
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    *ut.Translator
)

// func getJSONField(data interface{}, snp string, fieldname string) error {
// 	if !strings.Contains(snp, fieldname) {
// 		return nil
// 	}

// 	fieldArr := strings.Split(snp, ".")
// 	// fmt.Println(fieldArr)

// 	val := reflect.ValueOf(data)

// 	for i := 0; i < val.Type().NumField(); i++ {
// 		t := val.Type().Field(i)

// 		fieldName := t.Name
// 		// fmt.Println(fieldName)
// 		// fmt.Println(fieldArr[1])
// 		if fieldArr[1] == fieldName {
// 			if jsonTag := t.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
// 				// check for possible comma as in "...,omitempty"
// 				var commaIdx int
// 				if commaIdx = strings.Index(jsonTag, ","); commaIdx < 0 {
// 					commaIdx = len(jsonTag)
// 				}
// 				fieldName = jsonTag[:commaIdx]
// 			}
// 			return errors.New(fieldName)
// 		}

// 		return nil
// 	}

// 	return nil
// }

func getJSONFieldName(s interface{}, snp string) (fieldname string) {
	// fmt.Println("tag: ", tag)
	// fmt.Println("snp: ", snp)
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}

	fieldArr := strings.Split(snp, ".")

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		fieldName := f.Name
		if fieldArr[1] == fieldName {
			if jsonTag := f.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				// check for possible comma as in "...,omitempty"
				var commaIdx int
				if commaIdx = strings.Index(jsonTag, ","); commaIdx < 0 {
					commaIdx = len(jsonTag)
				}
				fieldName = jsonTag[:commaIdx]
			}
			return fieldName
		}
	}
	return ""
}

func ValidateStruct(data interface{}) ResponseHTTP {
	validate = validator.New()

	// register function to get tag name from json tags.
	// validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
	// 	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	// 	if name == "-" {
	// 		return ""
	// 	}
	// 	return name
	// })

	en := en.New()
	uni = ut.New(en, en)
	// id := id.New()
	// uni = ut.New(id, id)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")
	// trans, _ := uni.GetTranslator("id")

	en_translations.RegisterDefaultTranslations(validate, trans)
	// id_translations.RegisterDefaultTranslations(validate, trans)

	var errors []*ErrorDetailsResponse
	err := validate.Struct(data)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorDetailsResponse
			translatedErr := fmt.Errorf(err.Translate(trans))
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Message = translatedErr.Error()

			element.Field = getJSONFieldName(data, err.StructNamespace())

			if err.Tag() == "e164" {
				element.Message = fmt.Sprintf("%s %s", translatedErr.Error(), "i.e. +6281234567890 or +628 123 4567 890")
			}
			errors = append(errors, &element)
		}
	}

	return ResponseHTTP{
		Code:    fiber.StatusUnprocessableEntity,
		Message: fiber.ErrUnprocessableEntity.Message,
		Errors:  errors,
	}
}

func ValidatePhoneNumber(phone_number string) []*ErrorDetailsResponse {
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

	var errors []*ErrorDetailsResponse
	// validate.RegisterValidation("e164", ValidateMyVal)

	s := Phone{Phone: phone_number}
	err := validate.Struct(s)

	if err != nil {
		for idx, err := range err.(validator.ValidationErrors) {
			print(idx)
			var element ErrorDetailsResponse
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
