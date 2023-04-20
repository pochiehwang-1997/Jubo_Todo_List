package controllers

import (
	"Jubo_Todo_List/models"
	"Jubo_Todo_List/utilities"
	"context"
	"encoding/json"
	"log"
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

var rnd *renderer.Render
var collection *mongo.Collection

func TodoCollections(db *mongo.Database) {
	collection = db.Collection("todos")
}

func RenderHome(w http.ResponseWriter, r *http.Request) {
	rnd = renderer.New()
	err := rnd.Template(w, http.StatusOK, []string{"views/home.html"}, nil)
	utilities.CheckErr(err)
}

func FetchTodo(w http.ResponseWriter, r *http.Request) {
	rnd = renderer.New()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(chi.URLParam(r, "id")))
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Invalid ID",
			"error":   err,
		})
		return
	}
	t := models.TodoModel{}
	filter := bson.D{{Key: "_id", Value: id}}
	err = collection.FindOne(context.TODO(), filter).Decode(&t)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todos",
			"error":   err,
		})
		return
	}

	todoList := models.Todo{
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

func FetchAllTodos(w http.ResponseWriter, r *http.Request) {
	rnd = renderer.New()
	todos := []models.TodoModel{}

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

	todoList := []models.Todo{}

	for _, t := range todos {
		todoList = append(todoList, models.Todo{
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

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	rnd = renderer.New()
	var t models.Todo

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

	tm := models.TodoModel{
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

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	rnd = renderer.New()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(chi.URLParam(r, "id")))
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Invalid ID",
			"error":   err,
		})
		return
	}

	var t models.Todo

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
		"message": "Successfully update " + strconv.Itoa(int(result.ModifiedCount)) + " todo",
	})

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	rnd = renderer.New()
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
		"message": "Successfully delete " + strconv.Itoa(int(result.DeletedCount)) + " todo",
	})
}
