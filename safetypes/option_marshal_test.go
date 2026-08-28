package st

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{
			"nil option",
			nil,
			"null",
		},
		{
			"none",
			None[any](),
			"null",
		},
		{
			"some int",
			Some[any](123),
			"123",
		},
		{
			"some string",
			Some[any]("one-two-three"),
			`"one-two-three"`,
		},
		{
			"some struct",
			Some[any](struct {
				Name string
				Age  int
			}{"Firstname", 21}),
			`{"Name": "Firstname", "Age": 21}`,
		},
		{
			"none inside a struct",
			struct {
				Name string
				Job  *Option[string]
			}{"Poor Schmuck", None[string]()},
			`{"Name": "Poor Schmuck", "Job": null}`,
		},
		{
			"some inside a struct",
			struct {
				Name string
				Job  *Option[string]
			}{"Mega Chad", Some("Big Boss")},
			`{"Name": "Mega Chad", "Job": "Big Boss"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := json.Marshal(tt.value)
			require.NoError(t, err)

			assert.JSONEq(t, tt.expected, string(res))
		})
	}
}

func TestUnmarshalText(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var o Option[int]

		err := json.Unmarshal([]byte("null"), &o)
		require.NoError(t, err)

		expected := Option[int]{ok: false, val: 0}
		assert.Equal(t, expected, o)
	})

	t.Run("some int", func(t *testing.T) {
		var o Option[int]

		err := json.Unmarshal([]byte("123"), &o)
		require.NoError(t, err)

		expected := Option[int]{ok: true, val: 123}
		assert.Equal(t, expected, o)
	})

	t.Run("some string", func(t *testing.T) {
		var o Option[string]

		err := json.Unmarshal([]byte(`"one-two-three"`), &o)
		require.NoError(t, err)

		expected := Option[string]{ok: true, val: "one-two-three"}
		assert.Equal(t, expected, o)
	})

	t.Run("some struct", func(t *testing.T) {
		type S struct {
			Name string
			Age  int
		}

		var o Option[S]
		err := json.Unmarshal([]byte(`{"Name": "Firstname", "Age": 21}`), &o)
		require.NoError(t, err)

		expected := Option[S]{ok: true, val: S{Name: "Firstname", Age: 21}}
		assert.Equal(t, expected, o)
	})

	t.Run("none inside a struct", func(t *testing.T) {
		type S struct {
			Name string
			Job  Option[string]
		}

		var o Option[S]
		err := json.Unmarshal([]byte(`{"Name": "Poor Schmuck", "Job": null}`), &o)
		require.NoError(t, err)

		expected := Option[S]{ok: true, val: S{"Poor Schmuck", *None[string]()}}
		assert.Equal(t, expected, o)
	})

	t.Run("some string inside a struct", func(t *testing.T) {
		type S struct {
			Name string
			Job  Option[string]
		}

		var o Option[S]
		err := json.Unmarshal([]byte(`{"Name": "Mega Chad", "Job": "Big Boss"}`), &o)
		require.NoError(t, err)

		expected := Option[S]{ok: true, val: S{"Mega Chad", *Some("Big Boss")}}
		assert.Equal(t, expected, o)
	})

	t.Run("invalid int", func(t *testing.T) {
		var o Option[int]

		err := json.Unmarshal([]byte(`"actually-a-string"`), &o)
		require.ErrorContains(t, err, "cannot unmarshal")
	})
}
