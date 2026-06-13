package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/WazedKhan/Solace/db"
	"github.com/WazedKhan/Solace/internal/auth"
	jwt_token "github.com/WazedKhan/Solace/internal/auth/token"
	"github.com/WazedKhan/Solace/middleware"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	mux := http.NewServeMux()

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Panicln("database connection string is missing")
	}

	ttlHours, err := strconv.Atoi(os.Getenv("TOKEN_VALID_PERIOD"))
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is missing")
	}

	pool, err := db.NewPool(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	log.Println("database connected")

	generator := jwt_token.NewGenerator(
		jwtSecret,
		time.Duration(ttlHours)*time.Hour,
	)

	repo := auth.NewRepository(pool)
	service := auth.NewService(repo, generator)
	authHandler := auth.NewHandler(service)

	mux.HandleFunc("POST /api/v1/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/login", authHandler.Login)
	mux.HandleFunc("/api/v1/users", authHandler.GetUsers)

	handler := middleware.RequestLog(mux)
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
		Handler: handler,
	}

	log.Printf("server starting at http://%s:%s", server, port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
