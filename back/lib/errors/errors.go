// This package helps to manipulate errors.
package errors

import (
	"fmt"
	"strings"
)

// Aggregate multiple errors into one.
func Aggregate(err []error) error {
	var errStrings []string
	for _, e := range err {
		errStrings = append(errStrings, e.Error())
	}
	return fmt.Errorf(strings.Join(errStrings, "\n"))
}
