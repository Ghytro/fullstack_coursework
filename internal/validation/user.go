package validation

import (
	"errors"
)

func validateUserNameChar(c byte) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c == '_'
}

func validateFirstNameChar(c byte) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z'
}

func validateLastNameChar(c byte) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z'
}

func ValidateUserName(username string) error {
	for _, c := range username {
		if !validateUserNameChar(byte(c)) {
			return errors.New("имя пользователя содержит некорректные символы")
		}
	}
	return nil
}

func ValidateUserPassword(password string) error {
	if len(password) < 6 {
		return errors.New("пароль должен быть не менее 6 символов")
	}
	return nil
}
