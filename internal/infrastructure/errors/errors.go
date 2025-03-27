package errors

import "fmt"

type ErrNotFound struct {
	Structure string
}

func NewErrNotFound(structure string) error {
	return ErrNotFound{Structure: structure}
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.Structure)
}

type ErrDatabase struct{}

func NewErrDatabase() error {
	return ErrDatabase{}
}

func (e ErrDatabase) Error() string {
	return "database error"
}

type ErrInvalidInput struct {
	Field string
}

func NewErrInvalidInput(field string) error {
	return ErrInvalidInput{Field: field}
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input in field %s", e.Field)
}

type ErrExternal struct {
	Err error
}

func NewErrExternal(err error) error {
	return ErrExternal{Err: err}
}

func (e ErrExternal) Error() string {
	return fmt.Sprintf("external error: %s", e.Err)
}
