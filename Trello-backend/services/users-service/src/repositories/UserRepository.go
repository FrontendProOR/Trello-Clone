package repositories

import (
	"context"
	"errors"
	"time"

	"users-service/src/config"
	"users-service/src/models"
	"users-service/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var usersCollection *mongo.Collection
var projectsCollection *mongo.Collection

// Init initializes the users collection.
func Init() {
	usersCollection = config.DB.Collection("users")
	projectsCollection = config.DB.Collection("projects")

}

// InsertUser inserts a new user into the database.
func InsertUser(user models.User) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	result, err := usersCollection.InsertOne(ctx, user)
	return result, err
}

// GetAllUsers retrieves all users from the database.
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserById retrieves a user by ID from the database.
func GetUserById(id primitive.ObjectID) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := usersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return user, err
}

func UpdateUser(id primitive.ObjectID, user models.User) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": user}
	result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return result, err
}

// DeleteUser deletes a user from the database.
func DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := usersCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no user found with that ID")
	}
	return nil
}

// AddUserToProject adds a user to the project by ID
func AddUserToProject(projectID primitive.ObjectID, userID primitive.ObjectID) error {
	_, err := projectsCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": projectID},
		bson.M{"$addToSet": bson.M{"members": userID.Hex()}},
	)
	return err
}

// RemoveUserFromProject removes a user from the project by ID
func RemoveUserFromProject(projectID, userID primitive.ObjectID) error {
	_, err := projectsCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": projectID},
		bson.M{"$pull": bson.M{"members": userID.Hex()}},
	)
	return err
}

// GetUserByEmail retrieves a user by their email from the database.
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username from the database.
func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := usersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}

// SaveVerificationCode updates the verification code for the user.
func SaveVerificationCode(email, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"verification_code": code}}
	_, err := usersCollection.UpdateOne(ctx, bson.M{"email": email}, update)
	return err
}

// ActivateUser activates a user by setting the 'IsActive' field to true if the verification code matches.
func ActivateUser(email, code string) error {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}

	if user.Code != code {
		return errors.New("incorrect verification code")
	}

	update := bson.M{"$set": bson.M{"is_active": true}}
	_, err = usersCollection.UpdateOne(ctx, bson.M{"email": email}, update)
	return err
}

// SaveUser saves a user to the database (useful when new users are created).
func SaveUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := usersCollection.InsertOne(ctx, user)
	return err
}

// LoginUser logs in a user by email and password.
func LoginUser(email, password string) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user not found")
		}
		return user, err
	}

	if !user.IsActive {
		return user, errors.New("user is not active")
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, errors.New("invalid password")
	}

	return user, nil
}

func ChangePassword(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := usersCollection.UpdateOne(ctx, bson.M{"email": user.Email}, bson.M{"$set": bson.M{"password": user.Password}})
	return err
}