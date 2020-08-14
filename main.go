package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Casuu!"))
	})

	r.Route("/group", func(r chi.Router) {
		r.Get("/", ListGroups)
		r.Post("/", CreateGroup)
	})
	http.ListenAndServe(":8080", r)
}

type Group struct {
	Name string `json:"name"`
}

type GroupRequest struct {
	*Group
}

func (g *GroupRequest) Bind(r *http.Request) error {
	if g.Group == nil {
		return errors.New("missing required Group fields.")
	}
	return nil

}

type Card struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type CardRequest struct {
	*Card
}

func ListGroups(w http.ResponseWriter, r *http.Request) {

}

// CreateGroup creates group
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	data := &GroupRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}

	db.Create(data.Group)
}

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
