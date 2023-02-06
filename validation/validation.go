package validation

import "regexp"

func AlphanumString(str string) bool {
	regex, err := regexp.MatchString("^[A-Za-z0-9_-]*$", str)
	if err != nil {
		return false
	}

	return regex
}

func AlphanumObjKey(data map[string]string) bool {
	regex, err := regexp.Compile("^[A-Za-z0-9_-]*$")
	if err != nil {
		return false
	}

	for v := range data {
		isTrue := regex.MatchString(v)
		if !isTrue {
			return isTrue
		}
	}

	return true
}
