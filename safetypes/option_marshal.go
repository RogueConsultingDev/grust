package st

import "encoding/json"

// MarshalJSON implements json.Marshaler.
func (o *Option[T]) MarshalJSON() ([]byte, error) {
	if o == nil || !o.ok {
		return []byte("null"), nil
	}

	return json.Marshal(o.val) //nolint:wrapcheck // The source error is fine here
}

// UnmarshalJSON implements json.Unmarshaler.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.ok = false

		return nil
	}

	err := json.Unmarshal(data, &o.val)
	if err != nil {
		return err //nolint:wrapcheck // The source error is fine here
	}

	o.ok = true

	return nil
}
