package main

import (
	"log"
	"net/http"

	"github.com/Udehlee/Collab-playlist/internal/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	oauth := auth.NewOAuth()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", oauth.Index)
	r.Get("/login/spotify", oauth.LoginWithSpotify)
	r.Get("/callback/spotify", oauth.HandleCallback)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}