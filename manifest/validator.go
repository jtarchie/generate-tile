package manifest

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

func (p *Payload) Validate() (validator.ValidationErrorsTranslations, error) {
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, fmt.Errorf("could not find 'english' translation")
	}
	err := en_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, fmt.Errorf("could not setup 'english' translation: %s", err)
	}

	err = validate.Struct(p)
	if errs, ok := err.(validator.ValidationErrors); ok {
		return errs.Translate(trans), nil
	}
	return nil, err
}
