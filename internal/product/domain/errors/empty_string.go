package errors

import "fmt"

type ErrEmptyString struct {
	label string
}

func (e ErrEmptyString) Error() string {
	return fmt.Sprintf("%s cannot be empty", e.label)
}

func NewNonEmptyString(label, value string) error {
	if value == "" {
		return ErrEmptyString{label: label}
	}

	return nil
}
