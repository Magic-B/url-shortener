package apperr

import "fmt"

func ErrWrapper(place string, err error, descr ...string) error {
	if len(descr) > 0 && descr[0] != "" {
			return fmt.Errorf("%s: %s: %w", place, descr[0], err)
	}
	return fmt.Errorf("%s: %w", place, err)
}
