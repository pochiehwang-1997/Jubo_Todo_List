package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Jubo_Todo_List/controllers"
	"Jubo_Todo_List/routes"
	"Jubo_Todo_List/utilities"
)

const (
	dbName         string = "todo_list"
	collectionName string = "todos"
	port           string = "8080"
)

func main() {

	// Connect to MongoDB
	config, err := utilities.LoadConfig("./")
	if err != nil {
		log.Fatal("Cannot load Config:", err)
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MONGO_URI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database(dbName).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to MongoDB!")

	// Initialize collection
	db := client.Database(dbName)
	controllers.InitTodoCollections(db)

	// Setup server
	setupServer()

}

func setupServer() {
	// Setup Server
	routeHandler := chi.NewRouter()
	routeHandler.Use(middleware.Logger)
	routeHandler.Get("/", controllers.RenderHome)
	routeHandler.Mount("/todos", routes.TodoHandlers())

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      routeHandler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Listening on port", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("listen:%s\n", err)
	}
	log.Println("Server down")
}
