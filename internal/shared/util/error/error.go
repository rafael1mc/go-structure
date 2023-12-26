package errorutil

import "errors"

func UnwrapRecursive(err error) error {
	var originalErr = err

	for originalErr != nil {
		var internalErr = errors.Unwrap(originalErr)
		if internalErr == nil {
			break
		}
		originalErr = internalErr
	}

	return originalErr
}
