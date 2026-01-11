package utils

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

var ErrorRecordNotFound = errors.New("record not found")

func WrapError(err error) error {
	if err == nil {
		return nil
	}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return errors.WithStack(err)
	}

	fn := runtime.FuncForPC(pc)
	return errors.Wrap(err, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
}
