package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Allenize/notes-api/internal/api"
	"github.com/Allenize/notes-api/internal/auth"
	"github.com/Allenize/notes-api/internal/middleware"
	"github.com/Allenize/notes-api/internal/store"
	"github.com/Allenize/notes-api/internal/web"
)

func loadOrGenerateSecret() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	b := make([]byte, 32)
	rand.Read(b)
	secret := hex.EncodeToString(b)
	log.Printf("WARNING: JWT_SECRET not set, using a random secret for this run only: %s", secret)
	log.Printf("set JWT_SECRET explicitly in production, or tokens won't survive a restart")
	return []byte(secret)
}

func main() {
	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "./data.json"
	}

	s, err := store.New(dataFile)
	if err != nil {
		log.Fatalf("could not open data store: %v", err)
	}

	secret := loadOrGenerateSecret()
	h := api.New(s, secret)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", web.Handler())
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("POST /auth/signup", h.Signup)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/refresh", h.Refresh)

	protected := http.NewServeMux()
	protected.HandleFunc("GET /tasks", h.ListTasks)
	protected.HandleFunc("POST /tasks", h.CreateTask)
	protected.HandleFunc("GET /tasks/{id}", h.GetTask)
	protected.HandleFunc("PUT /tasks/{id}", h.UpdateTask)
	protected.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)

	mux.Handle("/tasks", auth.RequireAuth(secret)(protected))
	mux.Handle("/tasks/", auth.RequireAuth(secret)(protected))

	limiter := auth.NewRateLimiter(120, 20)
	handler := middleware.Logging(limiter.Middleware(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		log.Printf("notes-api listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("shutdown complete")
}
