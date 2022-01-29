package main

import (
	"gihub.com/toufiq-austcse/todo-app-go/api/controllers"
	"gihub.com/toufiq-austcse/todo-app-go/api/repositories"
	"gihub.com/toufiq-austcse/todo-app-go/api/services"
	"gihub.com/toufiq-austcse/todo-app-go/common/database"
	"gihub.com/toufiq-austcse/todo-app-go/common/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv(envFileName string) {
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatal("Error in loading env")
	}
	log.Println("Env Loaded")
}

func main() {
	LoadEnv(".env")
	var (
		db             = database.SetupDatabaseConnection()
		userRepository = repositories.NewUserRepository(db)
		todoRepository = repositories.NewTodoRepository(db)
		userService    = services.NewUserService(userRepository)
		todoService    = services.NewTodoService(todoRepository)
		authService    = services.NewAuthService(userService)
		authController = controllers.NewAuthController(authService)
		todoController = controllers.NewTodoController(todoService)
	)
	ginEngine := gin.Default()
	authRoutes := ginEngine.Group("/api/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.GET("/me", middlewares.JwtAuthMiddleware(authService), authController.Me)
	}
	todosRoutes := ginEngine.Group("/api/todos")
	{
		todosRoutes.POST("/", middlewares.JwtAuthMiddleware(authService), todoController.Create)
		todosRoutes.GET("/:todo_id", middlewares.JwtAuthMiddleware(authService), todoController.GetOne)
		todosRoutes.GET("/", middlewares.JwtAuthMiddleware(authService), todoController.GetAll)
		todosRoutes.PATCH("/:todo_id", middlewares.JwtAuthMiddleware(authService), todoController.Update)
		todosRoutes.DELETE("/:todo_id", middlewares.JwtAuthMiddleware(authService), todoController.Delete)
	}
	ginEngine.Run(os.Getenv("PORT"))

}
