package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/WazedKhan/Solace/db"
	"github.com/WazedKhan/Solace/internal/auth"
	jwt_token "github.com/WazedKhan/Solace/internal/auth/token"
	"github.com/WazedKhan/Solace/internal/habit"
	"github.com/WazedKhan/Solace/internal/journal"
	"github.com/WazedKhan/Solace/internal/storage"
	"github.com/WazedKhan/Solace/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	mux := http.NewServeMux()
	ctx := context.Background()

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

	// bucket initialization
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("S3_BUCKET_NAME is missing")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatal("AWS_REGION is missing")
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(cfg)
	presingedClient := s3.NewPresignClient(s3Client)
	storage := storage.NewS3Storage(s3Client, presingedClient, bucketName)

	// auth
	repo := auth.NewRepository(pool)
	service := auth.NewService(repo, generator)
	authHandler := auth.NewHandler(service)

	// habit
	habitRepo := habit.NewRepository(pool)
	habitService := habit.NewService(habitRepo)
	habitHandler := habit.NewHandler(habitService)

	// journal
	journalRepo := journal.NewRepository(pool)
	journalService := journal.NewService(journalRepo, storage)
	journalHandler := journal.NewHandler(journalService)

	mux.HandleFunc("POST /api/v1/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/login", authHandler.Login)
	mux.Handle(
		"GET /api/v1/me",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(authHandler.Me),
		),
	)
	// habit routes
	mux.Handle(
		"POST /api/v1/habits",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(habitHandler.CreateHabit),
		),
	)
	mux.Handle(
		"GET /api/v1/habits",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(habitHandler.GetHabits),
		),
	)
	mux.Handle(
		"POST /api/v1/habits/{id}/check-in",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(habitHandler.CheckIn),
		),
	)

	// journal routes
	mux.Handle(
		"POST /api/v1/journals",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(journalHandler.CreateJournal),
		),
	)
	mux.Handle(
		"GET /api/v1/journals",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(journalHandler.GetJournals),
		),
	)
	mux.Handle(
		"GET /api/v1/journals/drafts",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(journalHandler.GetDrafts),
		),
	)
	mux.Handle(
		"POST /api/v1/journals/presign-upload",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(journalHandler.PresignUpload),
		),
	)
	mux.Handle(
		"POST /api/v1/journals/{id}/confirm-upload",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(journalHandler.ConfirmUpload),
		),
	)
	mux.Handle(
		"GET /api/v1/journals/{id}",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(journalHandler.GetJournalByID),
		),
	)
	mux.Handle(
		"PATCH /api/v1/journals/{id}",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(journalHandler.UpdateJournalByID),
		),
	)
	mux.Handle(
		"DELETE /api/v1/journals/{id}",
		middleware.AuthMiddleware(
			generator,
			http.HandlerFunc(journalHandler.SoftDeleteJournal),
		),
	)

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
