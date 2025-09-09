package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"server/api"
	"uber-like-system/server/redis"
	"uber-like-system/server/ws"

	"github.com/go-chi/chi/v5" // Chi v5 router
	_ "github.com/lib/pq"      // Postgres driver
)

func main() {
	// Load environment variables
	dbURL := os.Getenv("DB_URL")
	redisAddr := os.Getenv("REDIS_ADDR")

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	redisCli := redis.New(redisAddr, "")

	// Setup WebSocket Hub
	hub := ws.NewHub()

	// Initialize API Server
	srvr := &api.Server{
		DB:    db,       // Your database connection
		Redis: redisCli, // Redis client
		WSHub: hub,      // WebSocket hub
	}

	// Setup Chi Router
	route := chi.NewRouter()
	srvr.RegisterRoute(route)

	log.Println("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", route); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
