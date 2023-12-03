package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/northern-ai/mongo-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}

}

func (uc UserController) GetAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data []models.User

	collection := uc.client.Database("mongo-golang").Collection("users")
	cursor, err := collection.Find(r.Context(), bson.M{})
	if err != nil {
		panic(err)
	}

	for cursor.Next(r.Context()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			panic(err)
		}
		data = append(data, user)
	}
	uj, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", uj)
}
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	u := models.User{}

	collection := uc.client.Database("mongo-golang").Collection("users")
	if err := collection.FindOne(r.Context(), bson.M{"_id": oid}).Decode(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)
	collection := uc.client.Database("mongo-golang").Collection("users")

	result, err := collection.InsertOne(r.Context(), u)
	if err != nil {
		panic(err)
	}
	authToken := (result.InsertedID.(primitive.ObjectID)).Hex()
	filter := bson.M{"_id": result.InsertedID.(primitive.ObjectID)}
	update := bson.M{"$set": bson.M{"AuthToken": authToken}}
	_, er := collection.UpdateOne(r.Context(), filter, update)
	if er != nil {
		panic(er)
	}
	fmt.Println(authToken)

	uj, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	collection := uc.client.Database("mongo-golang").Collection("users")
	result, err := collection.DeleteOne(r.Context(), bson.M{"_id": oid})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user: %s\n", result)
}
