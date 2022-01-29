package repositories

import (
	"context"
	"gihub.com/toufiq-austcse/todo-app-go/api/dto"
	"gihub.com/toufiq-austcse/todo-app-go/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return UserRepository{db: db}

}
func hashAndSalt(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash password")
	}
	return string(hash)
}
func (repo UserRepository) InsertUser(user dto.RegisterUserDto) (models.User, error) {
	user.Password = hashAndSalt([]byte(user.Password))
	result, err := repo.db.Collection("users").InsertOne(context.Background(), bson.M{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	})
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		ID:       result.InsertedID.(primitive.ObjectID),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil

}

func (repo UserRepository) FindUserByEmail(email string) (models.User, error) {
	var user models.User
	result := repo.db.Collection("users").FindOne(context.Background(), bson.M{
		"email": email,
	})
	if result.Err() != nil {
		if result.Err().Error() == "mongo: no documents in result" {
			return user, nil
		}
		return user, result.Err()
	}
	err := result.Decode(&user)
	if err != nil {
		return models.User{}, nil
	}

	return user, nil

}
func (repo UserRepository) FindUserByUserId(id string) (models.User, error) {
	log.Println("User id ", id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic("Error in ObjectId Conversion")
	}
	var user models.User
	result := repo.db.Collection("users").FindOne(context.Background(), bson.M{
		"_id": objectId,
	})
	if result.Err() != nil {
		if result.Err().Error() == "mongo: no documents in result" {
			return user, nil
		}
		return user, result.Err()
	}
	err = result.Decode(&user)
	if err != nil {
		return models.User{}, nil
	}
	return user, nil
}
