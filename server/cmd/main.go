package main

import (
	"log"
	"net/http"
	"os"
	api "uber-like-system/cmd"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	redisAddr := os.Getenv("REDIS_ADDR")

	redisCli := redis.New(redisAddr, "")
	hub := ws.Newhub()
	srvr := &api.Server{}

	route := chi.NewRouter()
	srvr.RegisterRoute(r)

	log.Println("listening on")
	http.ListenAndServe(":8080", r)
}
