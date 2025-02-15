package hw02unpackstring

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result []rune
	runes := []rune(s)
	length := len(runes)

	i := 0
	for i < length {
		if unicode.IsDigit(runes[i]) {
			return "", ErrInvalidString
		}

		if runes[i] == '\\' {
			if i < length-1 && (runes[i+1] == '\\' || unicode.IsDigit(runes[i+1])) {
				i++
			} else {
				return "", ErrInvalidString
			}
		}

		if i < length-1 && unicode.IsDigit(runes[i+1]) {
			for j := 0; j < int(runes[i+1]-'0'); j++ {
				result = append(result, runes[i])
			}

			i++
		} else {
			result = append(result, runes[i])
		}

		i++
	}

	return string(result), nil
}
