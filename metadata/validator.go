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
	validate.RegisterValidation("property-exists", func(fl validator.FieldLevel) bool {
		_, found := p.FindPropertyBlueprintFromPropertyInput(fl.Parent().Interface().(PropertyInput))
		return found
	})

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

	validate.RegisterTranslation("property-exists", trans, func(ut ut.Translator) error {
		return ut.Add("property-exists", "References a property blueprint ('{0}') that does not exist", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("property-exists", fe.Value().(string))
		return t
	})

	err = validate.Struct(p)
	if errs, ok := err.(validator.ValidationErrors); ok {
		return errs.Translate(trans), nil
	}
	return nil, err
}

func (p Payload) FindPropertyBlueprintFromPropertyInput(pi PropertyInput) (PropertyBlueprint, bool) {
	for _, pb := range p.PropertyBlueprints {
		if fmt.Sprintf(".properties.%s", pb.Name) == pi.Reference {
			return pb, true
		}
	}
	return PropertyBlueprint{}, false
}