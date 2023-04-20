package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rnd *renderer.Render

var db *mongo.Database

var collection *mongo.Collection

const (
	uri            string = "mongodb://localhost:27017"
	dbName         string = "todo_list"
	collectionName string = "todos"
	port           string = "9000"
)

type (
	todoModel struct {
		ID          primitive.ObjectID `bson:"_id,omitempty"`
		Title       string             `bson:"title"`
		Description string             `bson:"description"`
		Completed   bool               `bson:"completed"`
		CreatedAt   time.Time          `bson:"createdAt"`
	}

	todo struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Completed   bool      `json:"completed"`
		CreatedAt   time.Time `json:"createdAt"`
	}
)

func init() {
	rnd = renderer.New()
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
}

func fetchTodo(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(chi.URLParam(r, "id")))
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Invalid ID",
			"error":   err,
		})
		return
	}
	t := todoModel{}
	filter := bson.D{{Key: "_id", Value: id}}
	err = collection.FindOne(context.TODO(), filter).Decode(&t)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err,
		})
		return
	}

	todoList := todo{
		ID:          t.ID.Hex(),
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todoList,
	})

}

func fetchAllTodos(w http.ResponseWriter, r *http.Request) {
	todos := []todoModel{}

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err,
		})
		return
	}
	if err = cursor.All(context.TODO(), &todos); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err,
		})
		return
	}

	todoList := []todo{}

	for _, t := range todos {
		todoList = append(todoList, todo{
			ID:          t.ID.Hex(),
			Title:       t.Title,
			Description: t.Description,
			Completed:   t.Completed,
			CreatedAt:   t.CreatedAt,
		})
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todoList,
	})
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	// Validate if there is a title or not
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is required",
		})
		return
	}

	tm := todoModel{
		Title:       t.Title,
		Description: t.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	insertResult, err := collection.InsertOne(context.TODO(), tm)
	if err != nil {
		log.Fatal(err)
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Successfully create todo!",
		"todo_id": insertResult.InsertedID,
	})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(chi.URLParam(r, "id")))
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Invalid ID",
			"error":   err,
		})
		return
	}

	var t todo

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is required",
		})
		return
	}

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: t.Title}, {Key: "completed", Value: t.Completed}, {Key: "description", Value: t.Description}}}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to update todo",
			"error":   err,
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Successfully update" + strconv.Itoa(int(result.ModifiedCount)) + "todo",
	})

}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(chi.URLParam(r, "id")))
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Invalid ID",
			"error":   err,
		})
		return
	}

	filter := bson.D{{Key: "_id", Value: id}}
	result, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete todo",
			"error":   err,
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Successfully delete" + strconv.Itoa(int(result.DeletedCount)) + "todo",
	})
}

func main() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
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
	db = client.Database(dbName)
	collection = db.Collection(collectionName)
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todos", todoHandlers())

	srv := &http.Server{
		Addr:         "localhost:9000",
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Listening on port", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("listen:%s\n", err)
	}

}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/{id}", fetchTodo)
		r.Get("/", fetchAllTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return rg
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
