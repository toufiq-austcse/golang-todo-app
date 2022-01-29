package repositories

import (
	"context"
	"gihub.com/toufiq-austcse/todo-app-go/api/dto"
	"gihub.com/toufiq-austcse/todo-app-go/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type TodoRepository struct {
	db *mongo.Database
}

func NewTodoRepository(db *mongo.Database) TodoRepository {
	return TodoRepository{db: db}
}

func (repository TodoRepository) CreateTodo(dto dto.CreateTodoDto, user models.AuthUser) (models.Todo, error) {
	userObjId, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		panic("Error in User ObjectId conversion")
	}
	todo := models.Todo{
		ID:          primitive.NewObjectID(),
		Task:        dto.Task,
		IsCompleted: dto.IsCompleted,
		UserID:      userObjId,
	}
	result, err := repository.db.Collection("todos").InsertOne(context.Background(), todo)
	if err != nil {
		return models.Todo{}, err
	}
	return models.Todo{
		ID:          result.InsertedID.(primitive.ObjectID),
		Task:        dto.Task,
		IsCompleted: dto.IsCompleted,
	}, nil
}

func (repository TodoRepository) FindOneById(id string, user models.AuthUser) (models.Todo, error) {
	userObjId, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		panic("error in userId conversion")
	}
	todoObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {

		panic("error in userId conversion")
	}
	var todo models.Todo
	result := repository.db.Collection("todos").FindOne(context.Background(), bson.M{"_id": todoObjId, "user_id": userObjId})
	if result.Err() != nil {
		return models.Todo{}, err
	}
	err = result.Decode(&todo)
	if err != nil {
		return models.Todo{}, err
	}

	return todo, nil
}

func (repository TodoRepository) FindAll(user models.AuthUser) ([]models.Todo, error) {
	userObjId, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		panic("error in userId conversion")
	}
	var todos []models.Todo
	findResult, err := repository.db.Collection("todos").Find(context.Background(), bson.M{"user_id": userObjId})
	if err != nil {
		return todos, err
	}

	err = findResult.All(context.Background(), &todos)
	if err != nil {
		return todos, err
	}
	return todos, nil
}

func (repository TodoRepository) DeleteOneById(id string, user models.AuthUser) (bool, error) {
	userObjId, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		panic("error in userId conversion")
	}
	todoObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic("error in userId conversion")
	}
	deleteResult := repository.db.Collection("todos").FindOneAndDelete(context.Background(), bson.M{"_id": todoObjId, "user_id": userObjId})
	if deleteResult.Err() != nil {
		return false, deleteResult.Err()
	}
	return true, nil
}

func (repository TodoRepository) UpdateOneById(id string, user models.AuthUser, updateTodoDto dto.UpdateTodoDto) (models.Todo, error) {
	userObjId, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		panic("error in userId conversion")
	}
	todoObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic("error in userId conversion")
	}
	var updatedTodo models.Todo
	result := repository.db.Collection("todos").FindOneAndUpdate(context.Background(), bson.M{"_id": todoObjId, "user_id": userObjId}, bson.D{
		{"$set", bson.M{"task": updateTodoDto.Task, "is_completed": updateTodoDto.IsCompleted}},
	})
	if result.Err() != nil {
		log.Println("Err ", result.Err())
		return models.Todo{}, result.Err()
	}
	if err := result.Decode(&updatedTodo); err != nil {
		return models.Todo{}, err
	}
	return updatedTodo, nil
}
