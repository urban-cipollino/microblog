package validator

import (
	"errors"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"micro-blog/internal/exception"
)

// TODO Replace with Wire overtime
var (
	Validate   = validator.New()
	Translator ut.Translator
)

func Check(s any) error {
	if err := Validate.Struct(s); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			messages := make(map[string]string)
			for _, e := range ve {
				messages[e.Field()] = e.Translate(Translator)
			}
			return &exception.UnprocessableEntityError{Msg: messages}
		}
		return err
	}

	return nil
}
