package first_tests

import (
	// helper 	"FirstProject/Model/Helper"
	"testing"
	"fmt"
	// "log"
)

type HelperChecker struct {
	Helper		Helper
}

func NewHelperChecker() HelperChecker {
	return HelperChecker {
		Helper: Helper{},
	}
}

func TestQuitBearer(t *testing.T){

	// Arrange : Create variables needed for test
	fakeJWT := "Bearer jwt"

	// Act: 
	result := QuitBearer(fakeJWT)

	// Assert
	if result != "jwt" {
		t.Error(fmt.Println("Expected output: 'jwt', but received: '" + result + "'"))
	}
}