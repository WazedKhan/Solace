package main

import (
	"log"
	"net/http"
	"os"

	"github.com/WazedKhan/Solace/db"
	"github.com/WazedKhan/Solace/internal/auth"
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

	mux.HandleFunc("/api/v1/login", auth.LoginHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := os.Getenv("SERVER")
	if server == "" {
		server = "localhost"
	}

	addr := ":" + port

	log.Printf("server starting at http://%s:%s", server, port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
