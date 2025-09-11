package main

import (
	"log"
	"net/http"
	"os"

	"uber-like-system/server/api"
	database "uber-like-system/server/postgres"
	"uber-like-system/server/redis"
	ws "uber-like-system/server/ws"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {

	queries, dbConn := database.NewDatabase("database/schema.sql")
	defer dbConn.Close()

	redisAddr := os.Getenv("REDIS_ADDR")

	redisCli := redis.New(redisAddr, "")

	wsm := ws.NewWebSocketManager()

	srvr := &api.Server{
		DB:    queries,
		Redis: redisCli,
		Wsm:   wsm,
	}

	router := chi.NewRouter()
	srvr.RegisterRoute(router)

	log.Println("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
