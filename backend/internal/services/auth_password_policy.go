package services

import "unicode/utf8"

func validatePasswordStrength(password string) error {
	if utf8.RuneCountInString(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	if len(password) > maxPasswordBytes {
		return ErrPasswordTooLong
	}

	return nil
}
