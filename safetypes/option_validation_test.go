package st

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestValidations(t *testing.T) {
	v := validator.New()

	type S struct {
		Opt Option[any] `validate:"required,max=5"`
	}

	tests := []struct {
		name                     string
		value                    Option[any]
		expectedFailedValidation string
	}{
		{
			"valid string",
			*Some[any]("valid"),
			"",
		},
		{
			"valid empty string",
			*Some[any](""),
			"",
		},
		{
			"invalid string",
			*Some[any]("too long"),
			"max",
		},
		{
			"valid int value",
			*Some[any](3),
			"",
		},
		{
			"valid zero int",
			*Some[any](0),
			"",
		},
		{
			"invalid int",
			*Some[any](6),
			"max",
		},
		{
			"none",
			*None[any](),
			"required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := S{Opt: tt.value}

			err := v.Struct(s)

			if tt.expectedFailedValidation == "" {
				require.NoError(t, err)

				return
			}

			var vErr validator.ValidationErrors
			require.ErrorAs(t, err, &vErr)
			require.Len(t, vErr, 1)
			require.Equal(t, "S.Opt", vErr[0].Namespace())
			require.Equal(t, tt.expectedFailedValidation, vErr[0].Tag())
		})
	}
}
