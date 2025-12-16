package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	httpadapter "github.com/yimsoijoi/7s-backend-challenge/internal/adapters/http"
	"github.com/yimsoijoi/7s-backend-challenge/internal/adapters/mongo"
	"github.com/yimsoijoi/7s-backend-challenge/internal/application"
	"github.com/yimsoijoi/7s-backend-challenge/internal/infrastructure"
)

func main() {
	// Root Context
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Infrastructure
	mongoCfg := infrastructure.MongoConfig{
		URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Database: getEnv("MONGO_DB", "users"),
		Timeout:  10 * time.Second,
	}

	mongoClient, mongoDB := infrastructure.NewMongoDatabase(mongoCfg)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()
		_ = mongoClient.Disconnect(shutdownCtx)
	}()

	err := infrastructure.EnsureMongoIndexes(ctx, mongoDB.Collection("users"))
	if err != nil {
		log.Println("!! MongoDB not indexes")
	}

	ttlMinutes, err := strconv.Atoi(
		getEnv("JWT_TTL_MINUTES", "15"),
	)
	if err != nil {
		log.Fatalf("config JWT_TTL_MINUTES failed: %s", err.Error())
	}

	jwtManager := infrastructure.NewJWTManager(
		getEnv("JWT_SECRET", "secret"),
		time.Duration(ttlMinutes)*time.Minute,
	)

	// Repositories
	userRepo := mongo.NewUserRepository(mongoDB)

	// Services
	userService := application.NewUserService(userRepo, jwtManager)

	// HTTP Handlers
	handler := httpadapter.NewHandler(userService)

	mux := http.NewServeMux()

	// Public
	mux.HandleFunc("/auth/login", handler.Login)
	mux.Handle("/users", httpadapter.Logging(http.HandlerFunc(handler.CreateUser)))

	// Protected
	mux.Handle(
		"/users/",
		httpadapter.Logging(
			httpadapter.Auth(jwtManager, http.HandlerFunc(handler.GetUser)),
		),
	)
	mux.Handle(
		"/users/",
		httpadapter.Logging(
			httpadapter.Auth(jwtManager, http.HandlerFunc(handler.ListUsers)),
		),
	)
	mux.Handle(
		"/users/",
		httpadapter.Logging(
			httpadapter.Auth(jwtManager, http.HandlerFunc(handler.UpdateUser)),
		),
	)
	mux.Handle(
		"/users/",
		httpadapter.Logging(
			httpadapter.Auth(jwtManager, http.HandlerFunc(handler.DeleteUser)),
		),
	)

	// HTTP Server
	server := &http.Server{
		Addr:         getEnv("REST_PORT", ":8080"),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Background Task
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				count, err := userRepo.Count(context.Background())
				if err == nil {
					log.Printf("Users in DB: %d", count)
				}
			case <-ctx.Done():
				log.Println("Stopping background worker")
				return
			}
		}
	}()

	// Start Server
	go func() {
		log.Println("HTTP server started on :8080")
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen failed: %v", err)
		}
	}()

	// Graceful Shutdown
	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	log.Println("Server exited cleanly")
}

// Helpers

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
