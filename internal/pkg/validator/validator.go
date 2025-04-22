package validator

import (
	"errors"
	"log"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator interface {
	Validate(data any) error
}

type Validation struct {
	v *validator.Validate
	t ut.Translator
}

func New() *Validation {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	t, ok := uni.GetTranslator("en")
	if !ok {
		log.Fatalln("translator not found")
	}

	if err := enTranslations.
		RegisterDefaultTranslations(validate, t); err != nil {
		log.Fatalln(err)
	}

	return &Validation{
		validate,
		t,
	}
}

// Validate validates the data (struct)
// returning an error if the data is invalid.
func (v *Validation) Validate(
	data any,
) error {
	err := v.v.Struct(data)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	strErrs := make([]string, len(validationErrs))
	for i, validationErr := range validationErrs {
		strErrs[i] = validationErr.Translate(v.t)
	}

	errMsg := strings.Join(
		strErrs,
		", ",
	)

	return errors.New(errMsg)
}

var _ Validator = (*Validation)(nil)
