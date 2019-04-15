package usecase

import (
	repo "FirstProject/Domains/user/entity"
	"FirstProject/Model"
	"FirstProject/Model/validation"
	"FirstProject/Model/Auth"

	"image"
	"image/jpeg"
	"image/png"
	"io"

	// "encoding/base64"
	// "strings"
	"mime/multipart"
	"errors"
	"bytes"
	// "log"
	// "fmt"
	"os"
)

type Usecase interface{
	GetAll() ([] model.User, error)
	GetById(string) (model.User, error)
	GetUserByJwt(string)(model.User, error)
	GetProfileImage(string)([]byte, error)
	GetUserByUsername(string)(model.User, error)
	GetUserByRecoverToken(string)(model.User, error)
	Create(model.User) error
	Update(model.User) error
	Delete(string) error
	Register(model.User) error
	SetProfileImage(string, multipart.File) error
}

var (
	checker 	validation.Checker
	crypter		Auth.Crypter
)

type UsersUsecase struct {
	repo 	repo.RepositoryInterface
}

func NewUsecase(r repo.RepositoryInterface) Usecase {
	return &UsersUsecase{
		repo: r,
	}
}

func (u *UsersUsecase) GetAll()(users [] model.User, err error){

	users, err = u.repo.GetAll()

	if err != nil {
		err = errors.New("UsersNotFoundError")
	}

	return
}

func (u *UsersUsecase) GetById(userId string)(user model.User, err error){

	user, err = u.repo.GetById(userId)

	if err != nil {
		err = errors.New("UserNotFoundError")
	}

	return
}

func (u *UsersUsecase) GetUserByUsername(username string)(user model.User, err error){
	user, err = u.repo.GetUserByUsername(username)

	if err != nil {
		err = errors.New("UserNotFoundError")
	}
	return
}

func (u *UsersUsecase) GetUserByJwt(jwt string)(user model.User, err error){
	user, err = u.repo.GetUserByJwt(jwt)

	if err != nil {
		err = errors.New("UserNotFoundError")
	}
	return
}

func (u *UsersUsecase) GetProfileImage(userId string) (imageBytes []byte, err error){

	var user model.User
	user, err = u.repo.GetById(userId)

	if err != nil {
		err = errors.New("UserNotFoundError")
		return
	}
	
	fileWithImgData, err2 := os.Open(user.RouteImg)
	defer fileWithImgData.Close()

	if err2 != nil {
		err = err2
		return
	}

	myImage, format, err3 := image.Decode(fileWithImgData)

	if err3 != nil {
		err = err3
		return
	}

	var buff bytes.Buffer

	switch format{
	case "jpg":
	case "jpeg":
		quality := jpeg.Options{Quality:100}
		jpeg.Encode(&buff, myImage, &quality)
	case "png":
		png.Encode(&buff, myImage)
	}

	imageBytes = buff.Bytes()
	return
}

func (u *UsersUsecase) GetUserByRecoverToken(token string)(user model.User, err error){
	user, err = u.repo.GetUserByRecoverToken(token)

	if err != nil {
		err = errors.New("UserNotFoundError")
	}
	return
}

func (u *UsersUsecase) Create(user model.User) (err error) {

	var fieldsRequired []string
	fieldsRequired = append(fieldsRequired, user.Username, user.Password, user.Email, user.Role)

	if !checker.HasFieldsRequired(fieldsRequired){
		err = errors.New("EmptyFieldsError")
		return
	}

	_, err = u.GetUserByUsername(user.Username)

	if err == nil {
		err = errors.New("UsernameAlreadyExistsError")
		return
	}

	err = u.repo.Create(user)

	if err != nil {
		err = errors.New("CreationUserError")
		return
	}
	return
}

func (u *UsersUsecase) Update(user model.User) (err error) {

	var fieldsRequired []string
	fieldsRequired = append(fieldsRequired, user.Username, user.Password, user.Email)

	if !checker.HasFieldsRequired(fieldsRequired){
		err = errors.New("EmptyFieldsError")
		return
	}

	_, err = u.repo.GetById(user.Id.Hex())

	if err != nil {
		err = errors.New("UserNotFoundError")
		return
	}

	err = u.repo.Update(user)

	if err != nil {
		err = errors.New("UpdateUserError")
		return
	}
	return
}

func (u *UsersUsecase) Delete(userId string) (err error) {
	err = u.repo.Delete(userId)

	if err != nil {
		err = errors.New("DeleteUserError")
	}
	return
}

func (u *UsersUsecase) Login(user model.User, userDb model.User) error {

	if user.NotExists(){
		return errors.New("UserNotExistsError")
	}

	err := crypter.PasswordCoincides(user.Password, userDb.Password)

	if err != nil {
		return err
	}
	
	return nil
}

func (u *UsersUsecase) Register(user model.User) (err error){
	var fieldsRequired []string
	fieldsRequired = append(fieldsRequired, user.Username, user.Password, user.Email)

	if !checker.HasFieldsRequired(fieldsRequired){
		err = errors.New("EmptyFieldsError")
		return
	}

	_, err = u.GetUserByUsername(user.Username)

	if err == nil {
		err = errors.New("UsernameAlreadyExistsError")
		return
	}

	user.Role = "user"
	err = u.repo.Create(user)

	if err != nil {
		err = errors.New("CreationUserError")
		return
	}
	return
}

func (u *UsersUsecase) SetProfileImage(userId string, file multipart.File) (error){

	buf := bytes.NewBuffer(nil)

	defer file.Close()

	if _, err := io.Copy(buf, file); err != nil {
		return err
	}
	
	// embeded in profile image struct
	myImage, format, err := image.Decode(bytes.NewReader(buf.Bytes()))

	if err != nil {
		return err
	}

	// SaveImageInServerFolder
	imgName := userId + "." + format
	creationRoute, err2 := SaveImageInServerFolder(myImage, imgName, 80, format)
	if err2 != nil {
		return err2
	}

	// outputFile, err2 := os.Create(imgRoute)
	// defer outputFile.Close()

	// if err2 != nil {
	// 	return err2
	// }
	
	// switch format {
	// 	case "jpg":
	// 	case "jpeg":
	// 		quality := jpeg.Options{Quality:80}
	// 		jpeg.Encode(outputFile, myImage, &quality)
	// 	case "png":
	// 		png.Encode(outputFile, myImage)
	// }

	userDb, _ := u.repo.GetById(userId)
	userDb.SetRouteImg(creationRoute)
	userDb.EmptyProfileImage()

	if err := u.repo.Update(userDb); err != nil {
		return err
	}
	return nil
}

// Must evolve into ImageDirectory Struct
func SaveImageInServerFolder(imageDecoded image.Image, imgName string, percentOfQuality int, format string)(string, error){

	imageFolderRoute := "./image/"
	creatingRoute := imageFolderRoute + imgName

	outputFile, err := os.Create(creatingRoute)
	defer outputFile.Close()

	if err != nil {
		return "", err
	}

	// other func??
	switch format {
	case "jpg":
	case "jpeg":
		quality := jpeg.Options{Quality:percentOfQuality}
		jpeg.Encode(outputFile, imageDecoded, &quality)
	case "png":
		png.Encode(outputFile, imageDecoded)
	}
	return creatingRoute, nil
}