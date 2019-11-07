package metadata

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

func (p Payload) Validate() (validator.ValidationErrorsTranslations, error) {
	validate := validator.New()

	en := en.New()
	uni := ut.New(en, en)
	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, fmt.Errorf("could not find 'en' translation")
	}
	err := en_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, fmt.Errorf("could not setup 'en' translation: %s", err)
	}

	err = validate.Struct(p)
	if errs, ok := err.(validator.ValidationErrors); ok {
		return errs.Translate(trans), nil
	}
	return nil, err
}
