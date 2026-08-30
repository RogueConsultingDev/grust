package st

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// Scan implements the Scanner interface.
func (o *Option[T]) Scan(value any) error {
	if value == nil {
		*o = None[T]()

		return nil
	}

	sValue := reflect.ValueOf(value)
	sType := sValue.Type()

	tValue := reflect.ValueOf(o.val)
	if tValue.Kind() != reflect.Pointer {
		tValue = reflect.ValueOf(&o.val)
	}

	if tValue.IsNil() {
		return fmt.Errorf("can't assign value to nil pointer of type %T", o.val)
	}

	tValue = tValue.Elem()
	tType := tValue.Type()

	if sType.AssignableTo(tType) {
		tValue.Set(sValue)

		o.ok = true

		return nil
	}

	if sType.ConvertibleTo(tType) {
		tValue.Set(sValue.Convert(tType))

		o.ok = true

		return nil
	}

	return fmt.Errorf("invalid type for %T: %T", o, value)
}

// Value implements the driver Valuer interface.
func (o Option[T]) Value() (driver.Value, error) {
	if o.ok {
		return o.val, nil
	}

	return nil, nil //nolint:nilnil // There's no error here, we return a valid nil value
}
