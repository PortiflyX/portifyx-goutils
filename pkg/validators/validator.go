package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/portifyx/portifyx-goutils/pkg/apperror"
)

var symbolRegex = regexp.MustCompile(`^[a-zA-Z0-9]{1,20}(_[a-zA-Z0-9]{1,20})?$`)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() (*Validator, error) {
	v := validator.New()

	// Register custom tag: symbol_format
	_ = v.RegisterValidation("symbol_format", func(fl validator.FieldLevel) bool {
		return symbolRegex.MatchString(fl.Field().String())
	})

	return &Validator{validate: v}, nil
}

func (v *Validator) ValidateStruct(s any) error {
	if err := v.validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				tag := fieldErr.Tag()
				field := fieldErr.Field()

				if tag == "required" {
					return apperror.NewMissingRequiredField(field)
				}
				return apperror.NewValidationError(field, tag)
			}
		}
		return apperror.NewValidationError("struct", err.Error())
	}
	return nil
}
