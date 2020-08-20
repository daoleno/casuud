package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Init databse
	dsn := "user=postgres dbname=casuu port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Group{}, &Card{})

	// Init router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Casuu!"))
	})

	r.Route("/group", func(r chi.Router) {
		r.Get("/", ListGroups)
		r.Post("/", CreateGroup)
		r.Post("/{groupName}", ListCards)
	})

	r.Route("/card", func(r chi.Router) {
		r.Post("/", CreateCard)
		r.Get("/", ListCards)
	})

	r.Route("/{cardID}", func(r chi.Router) {
		r.Use(CardCtx) // Load the *Card on the request context
		// r.Get("/", GetCard)       // GET /articles/123
		// r.Put("/", UpdateCard)    // PUT /articles/123
		// r.Delete("/", DeleteCard) // DELETE /articles/123
	})

	http.ListenAndServe(":8080", r)
}
