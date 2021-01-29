package validation

import (
	"regexp"
	"unicode"
)

var regEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Phone проверяет переданный номер телефона n на соответствие требованиям валидации
func Phone(n *int64, allowNil bool) bool {
	if allowNil {
		return n == nil || (*n >= 79000000000 && *n <= 79999999999)
	}
	return n != nil && *n >= 79000000000 && *n <= 79999999999
}

// Email проверяет переданный электронный адрес s на соответствие требованиям валидации
func Email(s string, allowNil bool) bool {
	if allowNil {
		return s == "" || regEmail.MatchString(s)
	}
	return s != "" && regEmail.MatchString(s)
}

// Password проверяет переданный пароль s на соответствие требованиям валидации
func Password(s *string) bool {
	if s == nil || len(*s) < 6 {
		return false
	}

	var hasLetter, hasDigit bool
	for _, c := range *s {
		switch {
		case !hasDigit && unicode.IsDigit(c):
			hasDigit = true
		case !hasLetter && unicode.IsLetter(c):
			hasLetter = true
		}
		if hasLetter && hasDigit {
			return true
		}
	}

	return hasLetter && hasDigit
}
