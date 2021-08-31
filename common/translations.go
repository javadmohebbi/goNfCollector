package common

import "strings"

// this function will check the actual field
// in flow recordset and compare it to trans.yml
// configuration file
// 'want' argument is a comma separated string
// and will be converted to array and
// after that it will be check against 'actual'
// string. If one of the array objects are the same as
// 'actual', it will return the 'actual' value,
// otherwise it will return '---> UNKNOWN FIELDS'
func CheckTranslationField(actual, want string) string {

	// split and trim 'want'
	// and compare it to actual
	arrWant := strings.Split(want, ",")
	for i := range arrWant {
		// trim spaces
		arrWant[i] = strings.TrimSpace(arrWant[i])

		// comapre them together
		if strings.ToLower(actual) == strings.ToLower(arrWant[i]) {
			// return True because 'actual' is equal to
			// of of 'want' values
			return actual
		}
	}

	return "UNKNOWN FIELDS"
}
