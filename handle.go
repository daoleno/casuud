package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type GroupRequest struct {
	*Group
}

func (g *GroupRequest) Bind(r *http.Request) error {
	if g.Group == nil {
		return errors.New("missing required Group fields.")
	}
	return nil

}

// GroupResponse is the response payload for the Group data model.
type GroupResponse struct {
	*Group
}

func NewGroupResponse(group *Group) *GroupResponse {
	resp := &GroupResponse{Group: group}

	return resp
}

func (rd *GroupResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewGroupListResponse(groups []*Group) []render.Renderer {
	list := []render.Renderer{}
	for _, group := range groups {
		list = append(list, NewGroupResponse(group))
	}
	return list
}

func ListGroups(w http.ResponseWriter, r *http.Request) {
	var groups []*Group
	result := db.Find(&groups)
	if result.Error != nil {
		render.Render(w, r, ErrInternal(result.Error))
		return
	}

	if err := render.RenderList(w, r, NewGroupListResponse(groups)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// CreateGroup creates group
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	data := &GroupRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}

	db.Create(data.Group)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewGroupResponse(data.Group))
}

type CardRequest struct {
	*Card
}

func (g *CardRequest) Bind(r *http.Request) error {
	if g.Card == nil {
		return errors.New("missing required Card fields.")
	}
	return nil

}

// CardResponse is the response payload for the Card data model.
type CardResponse struct {
	*Card
}

func (rd *CardResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewCardResponse(card *Card) *CardResponse {
	resp := &CardResponse{Card: card}

	return resp
}

func NewCardListResponse(cards []*Card) []render.Renderer {
	list := []render.Renderer{}
	for _, card := range cards {
		list = append(list, NewCardResponse(card))
	}
	return list
}

// CardCtx middleware is used to load an Card object from
// the URL parameters passed through as the request. In case
// the Card could not be found, we stop here and return a 404.
func CardCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var card *Card

		if cardID := chi.URLParam(r, "cardID"); cardID != "" {
			result := db.First((&card))
			if result.Error != nil {
				render.Render(w, r, ErrNotFound)
				return
			}
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "card", card)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCard returns the specific Card. You'll notice it just
// fetches the Card right off the context, as its understood that
// if we made it this far, the Card must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetCard(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the card
	// context because this handler is a child of the CardCtx
	// middleware. The worst case, the recoverer middleware will save us.
	card := r.Context().Value("card").(*Card)

	if err := render.Render(w, r, NewCardResponse(card)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// UpdateCard updates an existing Card in our persistent store.
func UpdateCard(w http.ResponseWriter, r *http.Request) {
	card := r.Context().Value("card").(*Card)

	data := &CardRequest{Card: card}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	card = data.Card
	db.Model(&card).Updates(Card{Front: card.Front, Back: card.Back})

	render.Render(w, r, NewCardResponse(card))
}

// DeleteArticle removes an existing Card from our persistent store.
func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the card
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	card := r.Context().Value("card").(*Card)

	result := db.Delete(card)
	if result.Error != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
}
