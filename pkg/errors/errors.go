package errors

import (
	"errors"

	"github.com/rotisserie/eris"
)

var (
	NotImplementedError = eris.New("function not implemented")
	InvalidSchemaError  = eris.New("invalid schema")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}
