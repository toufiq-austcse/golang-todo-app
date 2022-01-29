package middlewares

import (
	"gihub.com/toufiq-austcse/todo-app-go/api/models"
	"gihub.com/toufiq-austcse/todo-app-go/api/services"
	"gihub.com/toufiq-austcse/todo-app-go/common/helper"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func JwtAuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.BuildErrorResponse("Failed to process the request", "No Token Found", helper.EmptyObj{})
			context.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		headerArr := strings.Split(authHeader, " ")
		if len(headerArr) != 2 {
			response := helper.BuildErrorResponse("Failed to process the request", "Token must start with Bearer", helper.EmptyObj{})
			context.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		user, err := authService.VerifyToken(headerArr[1])
		if err != nil {
			log.Println(err)
			response := helper.BuildErrorResponse("Token is not valid", err.Error(), nil)
			context.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		context.Set("user", models.AuthUser{
			ID:    user.ID.Hex(),
			Name:  user.Name,
			Email: user.Email,
		})
		context.Next()

	}
}
