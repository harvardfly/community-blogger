package alipay

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

func MustLength(str string, length int) (string, error) {
	if utf8.RuneCountInString(str) <= length {
		return str, nil
	}
	return "", errors.New(fmt.Sprintf("%s must less %v", str, length))
}

func TruncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}

	count := utf8.RuneCountInString(str)
	if count <= length {
		return str
	}

	trunc := []rune(str)[:length]
	return string(trunc)
}
