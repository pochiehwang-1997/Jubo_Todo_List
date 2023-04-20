package routes

import (
	controllers "Jubo_Todo_List/controllers"
	"net/http"

	"github.com/go-chi/chi"
)

// Routers
func TodoHandlers() http.Handler {
	routeGroup := chi.NewRouter()
	routeGroup.Group(func(r chi.Router) {
		r.Get("/{id}", controllers.FetchTodo)
		r.Get("/", controllers.FetchAllTodos)
		r.Post("/", controllers.CreateTodo)
		r.Put("/{id}", controllers.UpdateTodo)
		r.Delete("/{id}", controllers.DeleteTodo)
	})
	return routeGroup
}
