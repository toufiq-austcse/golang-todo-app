package controllers

import (
	"gihub.com/toufiq-austcse/todo-app-go/api/dto"
	"gihub.com/toufiq-austcse/todo-app-go/api/services"
	"gihub.com/toufiq-austcse/todo-app-go/common/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	authService services.AuthService
	jwtService  services.JwtService
}

func NewAuthController(authService services.AuthService) AuthController {
	return AuthController{
		authService: authService,
	}
}

func (controller AuthController) Register(ctx *gin.Context) {
	var body dto.RegisterUserDto
	if err := ctx.ShouldBind(&body); err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	isDuplicateEmail, err := controller.authService.CheckDuplicateEmail(body.Email)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
		return
	}
	if isDuplicateEmail {
		response := helper.BuildErrorResponse("Failed to process request", "duplicate email", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
		return
	}

	accessToken, err := controller.authService.UserRegistration(body)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BuildResponse(true, "Token Created", gin.H{
		"access_token": accessToken,
	})
	ctx.JSON(http.StatusCreated, response)
	return

}

func (controller AuthController) Login(ctx *gin.Context) {
	var body dto.LoginUserDto
	if err := ctx.ShouldBind(&body); err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	isVerified, user := controller.authService.VerifyCredentials(body)
	if !isVerified {
		response := helper.BuildErrorResponse("Please check again your credentials", "Invalid Credentials", helper.EmptyObj{})
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}
	accessToken := controller.jwtService.GenerateToken(user.ID.Hex())
	response := helper.BuildResponse(true, "Token Created", gin.H{
		"access_token": accessToken,
	})
	ctx.JSON(http.StatusOK, response)
	return
}
func (controller AuthController) Me(ctx *gin.Context) {
	user, isExist := ctx.Get("user")
	if isExist {
		response := helper.BuildResponse(true, "Authorized", user)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("Unauthenticated", "Invalid token", helper.EmptyObj{})
	ctx.JSON(http.StatusUnauthorized, response)
	return

}
