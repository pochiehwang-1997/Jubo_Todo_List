package controllers

import (
	"Jubo_Todo_List/models"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var render *renderer.Render
var collection *mongo.Collection

// Initialze collection
func InitTodoCollections(db *mongo.Database) {
	collection = db.Collection("todos")
}

// Render home page
func RenderHome(w http.ResponseWriter, r *http.Request) {
	render = renderer.New()
	render.JSON(w, http.StatusOK, renderer.M{"message": "Hi!"})
}

// Get one specific todo
func FetchTodo(w http.ResponseWriter, r *http.Request) {
	render = renderer.New()

	// Get id from URL params and change to mongo object id form
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(chi.URLParam(r, "id")))
	if err != nil {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Invalid ID",
			"error":   err,
		})
		return
	}

	// Find one by id and save in todoInstance
	todoInstance := models.TodoModel{}
	filter := bson.D{{Key: "_id", Value: id}}
	err = collection.FindOne(context.TODO(), filter).Decode(&todoInstance)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err,
		})
		return
	}

	// Change to json type and send to client
	result := models.Todo{
		ID:          todoInstance.ID.Hex(),
		Title:       todoInstance.Title,
		Description: todoInstance.Description,
		Completed:   todoInstance.Completed,
		CreatedAt:   todoInstance.CreatedAt,
	}
	render.JSON(w, http.StatusOK, renderer.M{
		"data": result,
	})
}

// Get all todo list
func FetchAllTodos(w http.ResponseWriter, r *http.Request) {
	render = renderer.New()

	// Get all todos and save in todoInstances
	todoInstances := []models.TodoModel{}
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err,
		})
		return
	}
	if err = cursor.All(context.TODO(), &todoInstances); err != nil {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err,
		})
		return
	}

	// Change to json type and send to client
	result := []models.Todo{}
	for _, t := range todoInstances {
		result = append(result, models.Todo{
			ID:          t.ID.Hex(),
			Title:       t.Title,
			Description: t.Description,
			Completed:   t.Completed,
			CreatedAt:   t.CreatedAt,
		})
	}
	render.JSON(w, http.StatusOK, renderer.M{
		"data": result,
	})
}

// Create one todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	render = renderer.New()

	// Decode input from json type
	var input models.Todo
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.JSON(w, http.StatusBadRequest, err)
		return
	}

	// Validate if there is a title or not, if not send error msg
	if input.Title == "" {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is required",
		})
		return
	}

	// Create one todo and insert into collection
	todoInstance := models.TodoModel{
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	insertResult, err := collection.InsertOne(context.TODO(), todoInstance)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Failed to create todo",
			"error":   err,
		})
		return
	}

	// Return success msg
	render.JSON(w, http.StatusCreated, renderer.M{
		"message": "Successfully create todo!",
		"todo_id": insertResult.InsertedID,
	})
}

// Update one todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	render = renderer.New()

	// Get id from url params and change it to mongo object id type
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(chi.URLParam(r, "id")))
	if err != nil {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Invalid ID",
			"error":   err,
		})
		return
	}

	// Decode input from json type
	var input models.Todo
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.JSON(w, http.StatusBadRequest, err)
		return
	}

	// Validate if there is a title or not, if not send error msg
	if input.Title == "" {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is required",
		})
		return
	}

	// Update todo instance
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: input.Title}, {Key: "completed", Value: input.Completed}, {Key: "description", Value: input.Description}}}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Failed to update todo",
			"error":   err,
		})
		return
	}

	// Send success msg
	render.JSON(w, http.StatusOK, renderer.M{
		"message": "Successfully update " + strconv.Itoa(int(result.ModifiedCount)) + " todo",
	})
}

// Delete one todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	render = renderer.New()

	// Get id from url params and change it to mongo object id type
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(chi.URLParam(r, "id")))
	if err != nil {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Invalid ID",
			"error":   err,
		})
		return
	}

	// Delete one todo instance
	filter := bson.D{{Key: "_id", Value: id}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Failed to delete todo",
			"error":   err,
		})
		return
	}

	// Send success msg
	render.JSON(w, http.StatusOK, renderer.M{
		"message": "Successfully delete " + strconv.Itoa(int(result.DeletedCount)) + " todo",
	})
}
