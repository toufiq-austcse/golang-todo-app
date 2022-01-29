package services

import (
	"gihub.com/toufiq-austcse/todo-app-go/api/dto"
	"gihub.com/toufiq-austcse/todo-app-go/api/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthService struct {
	userService UserService
	jwtService  JwtService
}

func NewAuthService(service UserService) AuthService {
	return AuthService{userService: service}
}

func (authService AuthService) UserRegistration(dto dto.RegisterUserDto) (string, error) {
	user, err := authService.userService.InsertUser(dto)
	if err != nil {
		return "", err
	}
	accessToken := authService.jwtService.GenerateToken(user.ID.Hex())
	return accessToken, nil
}

func (authService AuthService) CheckDuplicateEmail(email string) (bool, error) {
	user, err := authService.userService.findUserByEmail(email)
	if err != nil {
		log.Println("error in CheckDuplicateEmail", err.Error())
		return false, err
	}
	if user == (models.User{}) {
		return false, nil
	}
	return true, nil

}

func (authService AuthService) VerifyCredentials(dto dto.LoginUserDto) (bool, models.User) {
	user, err := authService.userService.findUserByEmail(dto.Email)
	if err != nil {
		return false, models.User{}
	}
	if user == (models.User{}) {
		return false, models.User{}
	}
	if comparePassword(user.Password, dto.Password) {
		return true, user
	}
	return false, models.User{}
}

func comparePassword(hashedPassword string, givenPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(givenPassword))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (authService AuthService) VerifyToken(token string) (models.User, error) {
	validateToken, err := authService.jwtService.ValidateToken(token)
	if err != nil {
		return models.User{}, err
	}
	userId, err := authService.jwtService.GetUserIdFromToken(validateToken)
	log.Println("userId ", userId)
	if err != nil {
		return models.User{}, err
	}
	user, err := authService.userService.FindUserById(userId)
	return user, err
}
