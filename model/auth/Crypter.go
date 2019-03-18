package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type Crypter struct {}

func (crypter *Crypter) Crypt(pwd string) (string, error) {

	passwordToBytes := []byte(pwd)
	cryptedPwd, err := bcrypt.GenerateFromPassword(passwordToBytes, bcrypt.DefaultCost)
	cryptedPwdToString := string(cryptedPwd)

	if err != nil {
		return "", err
	}

	return cryptedPwdToString, nil
}

func (crypter *Crypter) Uncrypt(){
}

func (crypter *Crypter) PasswordCoincides(DatabasePassword string, SentPassword string) error {

	err2 := bcrypt.CompareHashAndPassword([]byte(DatabasePassword), []byte(SentPassword))

	if err2 != nil { // if err = nil, correct password
		return err2
	}

	return nil
}