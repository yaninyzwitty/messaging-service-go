package router

import (
	"net/http"

	"github.com/yaninyzwitty/messaging-service/controller"
	"github.com/yaninyzwitty/messaging-service/middleware"
)

func NewRouter(controller *controller.MessageController) http.Handler {
	router := http.NewServeMux()

	// define middlewares
	loggingMiddleware := middleware.LoggingMiddleware
	corsMiddleware := middleware.CorsMiddleware

	// create a middleware chain
	middlewareChain := middleware.ChainMiddlewares(
		loggingMiddleware,
		corsMiddleware,
	)

	// Define routes and wrap them with the middleware stack
	router.HandleFunc("POST /messages", func(w http.ResponseWriter, r *http.Request) {
		middlewareChain(http.HandlerFunc(controller.CreateMessage)).ServeHTTP(w, r)
	})
	router.HandleFunc("PUT /messages/{id}", func(w http.ResponseWriter, r *http.Request) {
		middlewareChain(http.HandlerFunc(controller.UpdateMessage)).ServeHTTP(w, r)
	})
	router.HandleFunc("GET /messages", func(w http.ResponseWriter, r *http.Request) {
		middlewareChain(http.HandlerFunc(controller.GetMessages)).ServeHTTP(w, r)
	})
	router.HandleFunc("GET /messages/{id}", func(w http.ResponseWriter, r *http.Request) {
		middlewareChain(http.HandlerFunc(controller.GetMessage)).ServeHTTP(w, r)
	})
	router.HandleFunc("DELETE /messages/{id}", func(w http.ResponseWriter, r *http.Request) {
		middlewareChain(http.HandlerFunc(controller.DeleteMessage)).ServeHTTP(w, r)
	})
	return router

}
