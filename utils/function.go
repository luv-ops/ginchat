package utils

import "regexp"

func IsPhone(num string) bool {
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phoneRegex.MatchString(num) {

		return false
	}
	return true
}
