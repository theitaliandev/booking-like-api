package store

import (
	"context"

	"github.com/theitaliandev/booking-like-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type UserStore interface {
	GetUserByID(primitive.ObjectID) (*types.User, error)
	GetUsers() (*[]types.User, error)
	CreateUser(*types.User) (*types.User, error)
	DeleteUser(primitive.ObjectID) error
	UpdateUser(primitive.ObjectID, *types.UpdateUserParams) (*types.UpdateUserParams, error)
}

type MongoUserStore struct {
	client     *mongo.Client
	context    context.Context
	collection *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client, ctx context.Context) *MongoUserStore {
	return &MongoUserStore{
		client:     client,
		context:    ctx,
		collection: client.Database(dbname).Collection(userCollection),
	}
}

func (s *MongoUserStore) GetUserByID(objID primitive.ObjectID) (*types.User, error) {
	var user types.User
	if err := s.collection.FindOne(s.context, bson.M{"_id": objID}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers() (*[]types.User, error) {
	var users []types.User
	cur, err := s.collection.Find(s.context, bson.D{})
	if err != nil {
		return nil, err
	}
	if err := cur.All(s.context, &users); err != nil {
		return nil, err
	}
	return &users, nil
}

func (s *MongoUserStore) CreateUser(user *types.User) (*types.User, error) {
	res, err := s.collection.InsertOne(s.context, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) DeleteUser(objID primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(s.context, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(objID primitive.ObjectID, user *types.UpdateUserParams) (*types.UpdateUserParams, error) {
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "firstName", Value: user.FirtsName}, {Key: "lastName", Value: user.LastName}, {Key: "email", Value: user.Email}}}}
	_, err := s.collection.UpdateByID(s.context, objID, update)
	if err != nil {
		return nil, err
	}
	return user, nil
}
