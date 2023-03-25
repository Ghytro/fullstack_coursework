package validation

import (
	"errors"
	"github.com/Ghytro/galleryapp/internal/common"
	"github.com/Ghytro/galleryapp/internal/entity"
	"net/url"
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

func validateUserName(username string) error {
	for _, c := range username {
		if !validateUserNameChar(byte(c)) {
			return errors.New("имя пользователя содержит некорректные символы")
		}
	}
	return nil
}

func validateAvatarUrl(u string) error {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return errors.New("некорректный url аватара")
	}
	return nil
}

func validateUserFirstName(firstName string) error {
	for _, c := range firstName {
		if !validateFirstNameChar(byte(c)) {
			return errors.New("в имени есть некорректные символы")
		}
	}
	return nil
}

func validateUserLastName(lastName string) error {
	for _, c := range lastName {
		if !validateLastNameChar(byte(c)) {
			return errors.New("в фамилии есть некорректные символы")
		}
	}
	return nil
}

func validateUserFirstLastName(firstName *string, lastName *string) error {
	if lastName != nil && firstName == nil {
		return errors.New("указана фамилия, но не указано имя")
	}
	if firstName != nil {
		if err := validateUserFirstName(*firstName); err != nil {
			return err
		}
	}
	if err := validateUserLastName(*lastName); err != nil {
		return err
	}
	return nil
}

func validateUserCountry(countryCode string) error {
	if c := common.GetCountryByAlpha2(countryCode); c == nil {
		return errors.New("указан некорректный код страны")
	}
	return nil
}

func validateUserPassword(password string) error {
	if len(password) < 6 {
		return errors.New("пароль должен быть не менее 6 символов")
	}
	return nil
}

func ValidateUser(user *entity.User) error {
	if err := validateUserName(user.Username); err != nil {
		return err
	}
	if err := validateUserFirstLastName(user.FirstName, user.LastName); err != nil {
		return err
	}

	if user.AvatarUrl != nil {
		if err := validateAvatarUrl(*user.AvatarUrl); err != nil {
			return err
		}
	}
	if user.Country != nil {
		if err := validateUserCountry(*user.Country); err != nil {
			return err
		}
	}
	if err := validateUserPassword(user.Password); err != nil {
		return err
	}
	return nil
}
