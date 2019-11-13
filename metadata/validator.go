package metadata

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

func (p Payload) Validate() (validator.ValidationErrorsTranslations, error) {
	validate := validator.New()
	err := validate.RegisterValidation("property-exists", func(fl validator.FieldLevel) bool {
		_, found := p.FindPropertyBlueprintFromPropertyInput(fl.Parent().Interface().(PropertyInput).Reference)
		return found
	})

	if err != nil {
		return nil, err
	}

	english := en.New()
	uni := ut.New(english, english)
	trans, found := uni.GetTranslator("en")

	if !found {
		return nil, fmt.Errorf("could not find 'english' translation")
	}

	err = en_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, fmt.Errorf("could not setup 'english' translation: %s", err)
	}

	err = validate.RegisterTranslation("property-exists", trans, func(ut ut.Translator) error {
		return ut.Add("property-exists", "References a property blueprint ('{0}') that does not exist", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("property-exists", fe.Value().(string))
		return t
	})
	if err != nil {
		return nil, err
	}

	err = validate.Struct(p)
	if errs, ok := err.(validator.ValidationErrors); ok {
		return errs.Translate(trans), nil
	}
	return nil, err
}

func (p Payload) FindPropertyBlueprintFromPropertyInput(reference string) (PropertyBlueprint, bool) {
	parts := strings.Split(reference, ".")
	if parts[1] == "properties" {
		return propertyBlueprint(".properties", reference, p.PropertyBlueprints)
	}

	for _, jobType := range p.JobTypes {
		if parts[1] == jobType.Name {
			jobPrefix := fmt.Sprintf(".%s", jobType.Name)
			if strings.HasPrefix(reference, jobPrefix) {
				return propertyBlueprint(jobPrefix, reference, jobType.PropertyBlueprints)
			}
		}
	}

	return PropertyBlueprint{}, false
}

func propertyBlueprint(prefix string, reference string, blueprints []PropertyBlueprint) (PropertyBlueprint, bool) {
	for _, pb := range blueprints {
		currentPrefix := fmt.Sprintf("%s.%s", prefix, pb.Name)
		if currentPrefix == reference {
			return pb, true
		}
		if strings.HasPrefix(reference, currentPrefix) {
			if len(pb.OptionTemplates) > 0 {
				for _, optionTemplate := range pb.OptionTemplates {
					optionTemplatePrefix := fmt.Sprintf("%s.%s.%s", prefix, pb.Name, optionTemplate.Name)
					pb, found := propertyBlueprint(optionTemplatePrefix, reference, optionTemplate.PropertyBlueprints)
					if found {
						return pb, found
					}
				}
			} else {
				return propertyBlueprint(currentPrefix, reference, pb.PropertyBlueprints)
			}
		}
	}

	return PropertyBlueprint{}, false
}
