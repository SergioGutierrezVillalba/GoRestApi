package validation

import (
	// "fmt"
	"FirstProject/model"
	"regexp"
)

type Checker struct {}

// func (checker *Checker) UsernameAlreadyExists(username string) bool {

// 	userExists := false
// 	_ , err := userUsecase.GetUserByUsername(username)

// 	if err == nil {
// 		userExists = true
// 	}

// 	return userExists
// }

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

func (checker *Checker) HasPermissions(role string, httpRequestMethod string) bool {

	givePermission := false

	if role == "admin" {
		givePermission = true
	}

	return givePermission
}

func (checker *Checker) IsUpdatingItself(userRequesting model.User, userToUpdate model.User) bool {
	
	isUpdatingItself := false

	if userRequesting.Username == userToUpdate.Username {
		isUpdatingItself = true
	}

	return isUpdatingItself
}