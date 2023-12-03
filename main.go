package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/northern-ai/mongo-golang/controllers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(getClient())
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	r.GET("/user", uc.GetAllUsers)
	fmt.Println("Starting server at port 8080")
	http.ListenAndServe("localhost:8080", r)

}

func getClient() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}

	return client
}
