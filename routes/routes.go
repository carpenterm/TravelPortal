package routes

import (
	"github.com/go-chi/chi/v5"
	handlers "github.com/skandpurohit/TravelPortal/internal/server"
	"github.com/skandpurohit/TravelPortal/internal/service"
	"github.com/skandpurohit/TravelPortal/middleware"
)

func RegisterRoutes(taskService *service.TaskService) *chi.Mux {
	r := chi.NewRouter()

	h := handlers.NewHandlers(taskService)

	r.Get("/api/nextdate", h.HandleNextDate)
	r.Post("/api/task", middleware.Auth(h.HandleAddTask))
	r.Get("/api/tasks", middleware.Auth(h.HandleGetTasks))
	r.Get("/api/task", middleware.Auth(h.HandleGetTask))
	r.Put("/api/task", middleware.Auth(h.HandlePutTask))
	r.Delete("/api/task", middleware.Auth(h.HandleDeleteTask))
	r.Post("/api/task/done", middleware.Auth(h.HandleDoneTask))
	r.Post("/api/signin", h.HandleSignIn)

	return r
}
