package validator

import (
	"errors"
	"strings"
)

type Errors []error

func (errs Errors) Is(err error) bool {
	for _, e := range errs {
		if errors.Is(e, err) {
			return true
		}
	}
	return false
}

func (errs Errors) Error() string {
	sb := strings.Builder{}
	for i, v := range errs {
		sb.WriteString(v.Error())
		// don't write newline after last element
		if i != len(errs)-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}
