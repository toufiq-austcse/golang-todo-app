package controllers

import (
	"gihub.com/toufiq-austcse/todo-app-go/api/dto"
	"gihub.com/toufiq-austcse/todo-app-go/api/models"
	"gihub.com/toufiq-austcse/todo-app-go/api/services"
	"gihub.com/toufiq-austcse/todo-app-go/common/helper"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TodoController struct {
	service services.TodoService
}

func NewTodoController(todoService services.TodoService) TodoController {
	return TodoController{service: todoService}
}

func (controller TodoController) Create(ctx *gin.Context) {
	var body dto.CreateTodoDto
	if err := ctx.ShouldBind(&body); err != nil {
		response := helper.BuildErrorResponse("Unable to process the request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	user, _ := ctx.Get("user")
	createdTodo, err := controller.service.Create(body, user.(models.AuthUser))
	if err != nil {
		response := helper.BuildErrorResponse("Unable to process the request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helper.BuildResponse(true, "Todo Created", gin.H{
		"_id":          createdTodo.ID,
		"task":         createdTodo.Task,
		"is_completed": createdTodo.IsCompleted,
	})
	ctx.JSON(http.StatusCreated, response)
	return
}

func (controller TodoController) GetOne(ctx *gin.Context) {
	todoId := ctx.Param("todo_id")
	user, _ := ctx.Get("user")
	todo, err := controller.service.FineOneTodoById(todoId, user.(models.AuthUser))
	if err != nil {
		response := helper.BuildErrorResponse("Unable to process the request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	if todo == (models.Todo{}) {
		response := helper.BuildErrorResponse("Unable to process the request", "Not Found", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	response := helper.BuildResponse(true, "", gin.H{
		"_id":          todo.ID,
		"task":         todo.Task,
		"is_completed": todo.IsCompleted,
	})
	ctx.JSON(http.StatusOK, response)
	return
}

func (controller TodoController) GetAll(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	todos, err := controller.service.FindAll(user.(models.AuthUser))
	log.Println("todos ", todos)
	if err != nil {
		response := helper.BuildErrorResponse("Unable to process the request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helper.BuildResponse(true, "", todos)
	ctx.JSON(http.StatusOK, response)
	return
}

func (controller TodoController) Update(ctx *gin.Context) {
	todoId := ctx.Param("todo_id")
	user, _ := ctx.Get("user")
	var body dto.UpdateTodoDto
	if err := ctx.ShouldBind(&body); err != nil {
		response := helper.BuildErrorResponse("Unable to process the request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	updatedTodo, err := controller.service.UpdateOneById(todoId, user.(models.AuthUser), body)
	if err != nil {
		response := helper.BuildErrorResponse("Unable to process the request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	if updatedTodo == (models.Todo{}) {
		response := helper.BuildErrorResponse("Unable to process the request", "Not Updated", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.BuildResponse(true, "Todo Updated", updatedTodo)
	ctx.JSON(http.StatusOK, response)
	return
}

func (controller TodoController) Delete(ctx *gin.Context) {
	todoId := ctx.Param("todo_id")
	user, _ := ctx.Get("user")
	isDeleted, err := controller.service.DeleteOneById(todoId, user.(models.AuthUser))
	if err != nil {
		response := helper.BuildErrorResponse("Unable to process the request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helper.BuildResponse(true, "Todo Deleted", gin.H{"is_deleted": isDeleted})
	ctx.JSON(http.StatusOK, response)
	return
}
