package usecase

import (
	repo "FirstProject/Domains/user/entity"
	"FirstProject/Model"
	"FirstProject/Model/validation"

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
	"fmt"
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
	UpdateSelf(model.User) error
	UpdateAdmin(model.User) error
	Delete(string) error
	Register(model.User) error
	SetProfileImage(string, multipart.File) error
}

var (
	checker 	validation.Checker
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

func (u *UsersUsecase) UpdateSelf(user model.User) (err error){

	var fieldsRequired []string
	fieldsRequired = append(fieldsRequired, user.Username, user.Password, user.Email, user.GroupId)

	if !checker.HasFieldsRequired(fieldsRequired){
		err = errors.New("EmptyFieldsError")
		return
	}

	if err = u.repo.Update(user); err != nil {
		err = errors.New("UpdateUserError")
		return
	}
	return
}

func (u *UsersUsecase) UpdateAdmin(user model.User) (err error){

	var fieldsRequired  []string
	fieldsRequired = append(fieldsRequired, user.Username, user.Role, user.Email, user.GroupId)

	if !checker.HasFieldsRequired(fieldsRequired){
		err = errors.New("EmptyFieldsError")
		return
	}

	if err = u.repo.Update(user); err != nil {
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
	
	myImage, format, err2 := image.Decode(bytes.NewReader(buf.Bytes()))

	if err2 != nil {
		fmt.Println("Error decoding")
		return err2
	}

	imgRoute := "./image/" + userId + "." + format
	outputFile, err3 := os.Create(imgRoute)
	defer outputFile.Close()

	if err3 != nil {
		return err3
	}
	
	switch format {
		case "jpg":
		case "jpeg":
			quality := jpeg.Options{Quality:80}
			jpeg.Encode(outputFile, myImage, &quality)
		case "png":
			png.Encode(outputFile, myImage)
	}

	fmt.Println("(UsersUsecase) Usecase says, this is the route: " + imgRoute)

	userDb, _ := u.repo.GetById(userId)
	userDb.RouteImg = imgRoute
	userDb.ProfileImage = ""

	if err := u.repo.Update(userDb); err != nil {
		return err
	}
	return nil
}

// func (u *UsersUsecase) SetProfileImage(user model.User) (error){

// 	// 1. Get bytes from Multipart

// 	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(user.ProfileImage))
// 	myImage, format, err := image.Decode(reader)

	// if err != nil {
	// 	fmt.Println("Error decoding")
	// 	return err
	// }

	// imgRoute := "./image/" + user.Id.Hex() + "." + format
	// outputFile, err := os.Create(imgRoute)
	// defer outputFile.Close()

	// if err != nil {
	// 	return err
	// }
	
	// switch format {
	// 	case "jpg":
	// 	case "jpeg":
	// 		quality := jpeg.Options{Quality:80}
	// 		jpeg.Encode(outputFile, myImage, &quality)
	// 	case "png":
	// 		png.Encode(outputFile, myImage)
	// }

	// fmt.Println("(UsersUsecase) Usecase says, this is the route: " + imgRoute)

	// userDb, _ := u.repo.GetById(user.Id.Hex())
	// userDb.RouteImg = imgRoute
	// userDb.ProfileImage = ""

	// if err := u.repo.Update(userDb); err != nil {
	// 	return err
	// }
	// return nil
// }