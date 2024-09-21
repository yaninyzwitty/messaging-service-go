package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/yaninyzwitty/messaging-service/configuration"
	"github.com/yaninyzwitty/messaging-service/controller"
	"github.com/yaninyzwitty/messaging-service/database"
	"github.com/yaninyzwitty/messaging-service/repository"
	"github.com/yaninyzwitty/messaging-service/router"
	"github.com/yaninyzwitty/messaging-service/service"
)

func main() {
	cfg, err := configuration.LoadConfig()
	if err != nil {
		slog.Error("Error loading configuration", "error", err)
	}

	session, err := database.NewDatabaseConnection(cfg.HOSTS)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)

	}

	defer session.Close()

	messageRepo := repository.NewMessagesRepository(session)
	messageService := service.NewMessagesService(messageRepo)
	messageController := controller.NewMessageController(messageService)

	mux := router.NewRouter(messageController)

	server := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server: ", "error", err)
		}

	}()

	slog.Info(fmt.Sprintf("Server is running on port %s", cfg.PORT))

	// Set up OS signal handling for graceful shutdown
	quitCH := make(chan os.Signal, 1)
	signal.Notify(quitCH, os.Interrupt)

	<-quitCH
	slog.Info("Received termination signal, shutting down server...")

	shutdownCTX, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(shutdownCTX); err != nil {
		slog.Error("Failed to gracefully shut down server", "error", err)

	}
	slog.Info("Server shutdown successful")

}
