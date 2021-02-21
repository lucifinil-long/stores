package utils

import "errors"

// AddErrorPrefix adds prefix to error description and return the new error
// @param err is original error
// @param prefix is prefix of description
// @return new error
func AddErrorPrefix(err error, prefix string) error {
	return errors.New(prefix + err.Error())
}

// AddErrorSuffix adds suffix to error description and return the new error
// @param err is original error
// @param suffix is suffix of description
// @return new error
func AddErrorSuffix(err error, suffix string) error {
	return errors.New(err.Error() + suffix)
}
