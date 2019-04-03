package validation

import (
	// "fmt"
	"regexp"
)

type Checker struct {}

func (checker *Checker) CheckSpecialChars(dataToCheck string) error {

	_, err := regexp.MatchString("[[:word:]]", dataToCheck)
	return err
}

func (checker *Checker) HasFieldsRequired(fieldsRequired []string) bool{

	hasAllFields := true;

	for _, fieldRequired := range fieldsRequired {
		if len(fieldRequired) == 0 {
			hasAllFields = false
			break
		}
	}

	return hasAllFields
}