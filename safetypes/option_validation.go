package st

// ValidatorValue implements the validator.Valuer interface.
//
// It returns the actual value of the Option, so the validation can be done on
// it instead of on the Option itself. If the Option is a None, nil is returned
// instead.
//
// In order for an Option that contains a zero-value to pass the `required`
// validation, we return a pointer to the value instead of the value itself.
//
// Note, this requires validate@v10.30.2.
func (o Option[T]) ValidatorValue() any {
	if !o.ok {
		return nil
	}

	return &o.val
}
