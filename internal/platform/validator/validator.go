package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

var validate = func() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Use json tag names in error field keys instead of Go struct field names.
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return field.Name
		}

		return name
	})

	return v
}()

func Validate(v any) *apierrors.APIError {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return apierrors.UnprocessableEntity("validation failed", nil)
	}

	details := make(apierrors.ErrorDetails, len(ve))
	for _, fe := range ve {
		field := fe.Field()
		details[field] = append(details[field], errorKey(fe))
	}

	return apierrors.UnprocessableEntity("validation failed", details)
}

// errorKey returns an i18n-friendly key for the frontend.
// For parameterised tags the param is appended after ":" (e.g. "min:8").
func errorKey(fe validator.FieldError) string {
	switch fe.Tag() {
	case "min", "max", "len", "gt", "gte", "lt", "lte", "oneof":
		return fe.Tag() + ":" + fe.Param()
	default:
		return fe.Tag()
	}
}
