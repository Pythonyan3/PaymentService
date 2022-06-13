package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Pythonyan3/payment-service/config"
	"github.com/Pythonyan3/payment-service/internal/database"
	"github.com/Pythonyan3/payment-service/internal/handlers"
	"github.com/Pythonyan3/payment-service/internal/middleware"
	"github.com/Pythonyan3/payment-service/internal/repositories"
	"github.com/Pythonyan3/payment-service/internal/server"
	"github.com/Pythonyan3/payment-service/internal/services"

	"github.com/gorilla/mux"
)

type Application struct{}

func (app *Application) Run() error {
	// application entrypoint method, initialize whole things.
	log.Println("Service is starting...")

	// declaring all of variables
	var err error
	var cfg *config.Config
	var router *mux.Router
	var postgresDB *database.PostgresDB
	var httpServer *server.Server
	// repositories
	var transactionRepository *repositories.TransactionPostgresRepository
	var userRepository *repositories.UserPostgresRepository
	// services
	var transactionService *services.TransactionService
	var userService *services.UserService
	// middlewares
	var authMiddleware *middleware.AuthMiddleware
	// handlers
	var userHandler *handlers.UserHandler
	var transactionHandler *handlers.TransactionHandler

	// parse config (env variables)
	cfg = config.GetConfig()

	// create postgres DB connection
	postgresDB, err = database.NewPostgresDB(cfg)
	if err != nil {
		return fmt.Errorf("NewPostgresDb failed: %w", err)
	}

	// create repositories
	transactionRepository = repositories.NewTransactionPostgresRepository(postgresDB)
	userRepository = repositories.NewUserPostgresRepository(postgresDB)

	// create services
	transactionService = services.NewTransactionService(transactionRepository)
	userService = services.NewUserService(userRepository)

	// create middleware
	authMiddleware = middleware.NewAuthMiddleware(cfg.JWTSignKey)

	// create handlers
	transactionHandler = handlers.NewTransactionHandler(transactionService, authMiddleware)
	userHandler = handlers.NewUserHandler(userService)

	router = mux.NewRouter().PathPrefix("/api").Subrouter()

	// init routes
	userHandler.InitRoutes(router)
	transactionHandler.InitRoutes(router)

	// create and starting server
	httpServer = server.NewServer(cfg.ServicePort, router)

	go func() {
		if err := httpServer.Run(); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	// waiting for Ctrl + C to exit application
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Disconnecting db...")
	postgresDB.Close()

	log.Println("Service is shutted down!")

	return nil
}
