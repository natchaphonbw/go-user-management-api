package validator

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
	once     sync.Once
)

type InvalidField struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func Init() {
	once.Do(func() {
		validate = validator.New()

		// Translate messages
		english := en.New()
		uni := ut.New(english, english)
		trans, _ = uni.GetTranslator("en")
		_ = enTranslations.RegisterDefaultTranslations(validate, trans)

		// Use JSON tag name
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	})
}


func ValidateStruct(i interface{}) []InvalidField {
	err := validate.Struct(i)
	if err == nil {
		return nil
	}

	var errs []InvalidField
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			errs = append(errs, InvalidField{
				Field:  e.Field(),
				Reason: e.Translate(trans),
			})
		}
	}

	return errs
}
