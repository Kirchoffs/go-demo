package services

import (
    "context"
    "errors"
    "mongodb-api-demo/models"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
    userCollection *mongo.Collection
    ctx            context.Context
}

func NewUserService(userCollection *mongo.Collection, ctx context.Context) UserService {
    return &UserServiceImpl{
        userCollection: userCollection,
        ctx:            ctx,
    }
}

func (userService *UserServiceImpl) CreateUser(user *models.User) error {
    _, err := userService.userCollection.InsertOne(userService.ctx, user)
    return err
}

func (userService *UserServiceImpl) GetUser(name *string) (*models.User, error) {
    var user *models.User
    query := bson.D{bson.E{Key: "user_name", Value: name}}
    err := userService.userCollection.FindOne(userService.ctx, query).Decode(&user)
    return user, err
}

func (userService *UserServiceImpl) GetAll() ([]*models.User, error) {
    var users []*models.User
    cursor, err := userService.userCollection.Find(userService.ctx, bson.D{{}})
    if err != nil {
        return nil, err
    }
    for cursor.Next(userService.ctx) {
        var user models.User
        err := cursor.Decode(&user)
        if err != nil {
            return nil, err
        }
        users = append(users, &user)
    }
    if err := cursor.Err(); err != nil {
        return nil, err
    }
    cursor.Close(userService.ctx)
    if len(users) == 0 {
        return nil, errors.New("documents not found")
    }
    return users, nil
}

func (userService *UserServiceImpl) UpdateUser(user *models.User) error {
    filter := bson.D{bson.E{Key: "user_name", Value: user.Name}}
    update := bson.D{bson.E{
        Key: "$set",
        Value: bson.D{
            bson.E{Key: "user_name", Value: user.Name},
            bson.E{Key: "user_age", Value: user.Age},
            bson.E{Key: "user_address", Value: user.Address},
        },
    }}

    result, _ := userService.userCollection.UpdateOne(userService.ctx, filter, update)
    if result.MatchedCount != 1 {
        return errors.New("no matched document found for update")
    }
    return nil
}

func (userService *UserServiceImpl) DeleteUser(name *string) error {
    filter := bson.D{bson.E{Key: "user_name", Value: name}}
    result, _ := userService.userCollection.DeleteOne(userService.ctx, filter)
    if result.DeletedCount != 1 {
        return errors.New("no matched document found for delete")
    }
    return nil
}
