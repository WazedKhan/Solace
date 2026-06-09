package main

import (
	"log"
	"net/http"
	"os"

	"github.com/WazedKhan/Solace/db"
	"github.com/WazedKhan/Solace/internal/auth"
	"github.com/WazedKhan/Solace/middleware"
)

func main() {
	mux := http.NewServeMux()
	dsn := "postgres://solace:strong-password@localhost:5432/solace"

	pool, err := db.NewPool(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	log.Println("database connected")

	repo := auth.NewRepository(pool)
	service := auth.NewService(repo)
	authHandler := auth.NewHandler(service)

	mux.HandleFunc("/api/v1/login", auth.LoginHandler)
	mux.HandleFunc("/api/v1/register", authHandler.Register)
	mux.HandleFunc("/api/v1/users", authHandler.GetUsers)

	middleware := middleware.RequestLog(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := os.Getenv("SERVER")
	if server == "" {
		server = "localhost"
	}

	addr := ":" + port

	srv := &http.Server{
		Addr:    addr,
		Handler: middleware,
	}

	log.Printf("server starting at http://%s:%s", server, port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
